package worker

import (
	"errors"
	"sync"
	"time"

	"github.com/turbonomic/kubeturbo/pkg/discovery/dtofactory"
	"github.com/turbonomic/kubeturbo/pkg/discovery/metrics"
	"github.com/turbonomic/kubeturbo/pkg/discovery/monitoring"
	"github.com/turbonomic/kubeturbo/pkg/discovery/monitoring/types"
	"github.com/turbonomic/kubeturbo/pkg/discovery/stitching"
	"github.com/turbonomic/kubeturbo/pkg/discovery/task"

	"github.com/turbonomic/turbo-go-sdk/pkg/proto"

	"github.com/golang/glog"
)

const (
	defaultMonitoringWorkerTimeout time.Duration = time.Minute * 5
)

type k8sDiscoveryWorkerConfig struct {
	//k8sClusterScraper *cluster.ClusterScraper

	// a collection of all configs for building different monitoring clients.
	// key: monitor type; value: monitor worker config.
	monitoringSourceConfigs map[types.MonitorType][]monitoring.MonitorWorkerConfig

	stitchingPropertyType stitching.StitchingPropertyType
}

func NewK8sDiscoveryWorkerConfig(sType stitching.StitchingPropertyType) *k8sDiscoveryWorkerConfig {
	return &k8sDiscoveryWorkerConfig{
		stitchingPropertyType:   sType,
		monitoringSourceConfigs: make(map[types.MonitorType][]monitoring.MonitorWorkerConfig),
	}
}

// Add new monitoring worker config to the discovery worker config.
func (c *k8sDiscoveryWorkerConfig) WithMonitoringWorkerConfig(config monitoring.MonitorWorkerConfig) *k8sDiscoveryWorkerConfig {
	monitorType := config.GetMonitorType()
	configs, exist := c.monitoringSourceConfigs[monitorType]
	if !exist {
		configs = []monitoring.MonitorWorkerConfig{}
	}

	configs = append(configs, config)
	c.monitoringSourceConfigs[monitorType] = configs

	return c
}

// k8sDiscoveryWorker receives a discovery task from dispatcher(DiscoveryClient). Then ask available monitoring workers
// to scrape metrics source and get topology information. Finally it builds entityDTOs and send back to DiscoveryClient.
type k8sDiscoveryWorker struct {
	// A UID of a discovery worker.
	id string

	// config
	config *k8sDiscoveryWorkerConfig

	// a collection of all monitoring worker of different types.
	// key: monitoring types; value: monitoring worker instance.
	monitoringWorker map[types.MonitorType][]monitoring.MonitoringWorker

	// sink is a central place to store all the monitored data.
	sink *metrics.EntityMetricSink

	taskChan chan *task.Task
}

func NewK8sDiscoveryWorker(config *k8sDiscoveryWorkerConfig, wid string) (*k8sDiscoveryWorker, error) {
	if len(config.monitoringSourceConfigs) == 0 {
		return nil, errors.New("No monitoring source config found in config.")
	}

	// Build all the monitoring worker based on configs.
	monitoringWorkerMap := make(map[types.MonitorType][]monitoring.MonitoringWorker)
	for monitorType, configList := range config.monitoringSourceConfigs {
		monitorList, exist := monitoringWorkerMap[monitorType]
		if !exist {
			monitorList = []monitoring.MonitoringWorker{}
		}
		for _, config := range configList {
			monitoringWorker, err := monitoring.BuildMonitorWorker(config.GetMonitoringSource(), config)
			if err != nil {
				// TODO return?
				glog.Errorf("Failed to build monitoring worker configuration: %v", err)
				continue
			}
			monitorList = append(monitorList, monitoringWorker)
		}
		monitoringWorkerMap[monitorType] = monitorList
	}

	return &k8sDiscoveryWorker{
		id:               wid,
		config:           config,
		monitoringWorker: monitoringWorkerMap,
		sink:             metrics.NewEntityMetricSink(),

		taskChan: make(chan *task.Task),
	}, nil
}

func (worker *k8sDiscoveryWorker) RegisterAndRun(dispatcher *Dispatcher, collector *ResultCollector) {
	dispatcher.RegisterWorker(worker)
	for {
		select {
		case currTask := <-worker.taskChan:
			glog.V(3).Infof("Worker %s has received a discovery task.", worker.id)
			result := worker.executeTask(currTask)
			collector.ResultPool() <- result
			glog.V(3).Infof("Worker %s has finished the discovery task.", worker.id)

			dispatcher.RegisterWorker(worker)
		}
	}
}

// Worker start to working on task.
func (worker *k8sDiscoveryWorker) executeTask(currTask *task.Task) *task.TaskResult {
	if currTask == nil {
		err := errors.New("No task has been assigned to current worker.")
		glog.Errorf("%s", err)
		return task.NewTaskResult(worker.id, task.TaskFailed).WithErr(err)
	}

	// wait group to make sure metrics scraping finishes.
	var wg sync.WaitGroup

	// Resource monitoring
	resourceMonitorTask := currTask
	//if resourceMonitoringWorkers, exist := worker.monitoringWorker[types.ResourceMonitor]; exist {
	for _, resourceMonitoringWorkers := range worker.monitoringWorker {
		for _, rmWorker := range resourceMonitoringWorkers {
			wg.Add(1)
			go func(w monitoring.MonitoringWorker) {
				finishCh := make(chan struct{})
				stopCh := make(chan struct{}, 1)
				defer close(finishCh)
				defer close(stopCh)
				defer wg.Done()

				w.ReceiveTask(resourceMonitorTask)
				t := time.NewTimer(defaultMonitoringWorkerTimeout)
				go func() {
					glog.V(2).Infof("A %s monitoring worker is invoked.", w.GetMonitoringSource())
					// Assign task to monitoring worker.
					monitoringSink := w.Do()
					select {
					case <-stopCh:
						//glog.Infof("Calling thread: %s monitoring worker timeout!", w.GetMonitoringSource())
						return
					default:
					}
					//glog.Infof("%s has finished", w.GetMonitoringSource())
					t.Stop()
					// Don't do any filtering
					worker.sink.MergeSink(monitoringSink, nil)
					//glog.Infof("send to finish channel %p", finishCh)
					finishCh <- struct{}{}
				}()

				// either finish as expected or timeout.
				select {
				case <-finishCh:
					//glog.Infof("Worker %s finished as expected.", w.GetMonitoringSource())
					return
				case <-t.C:
					glog.Errorf("%s monitoring worker exceeds the max time limit for "+
						"completing the task.", w.GetMonitoringSource())
					stopCh <- struct{}{}
					//glog.Infof("%s stop", w.GetMonitoringSource())
					w.Stop()
					return
				}
			}(rmWorker)
		}
	}

	wg.Wait()

	var discoveryResult []*proto.EntityDTO
	// Build EntityDTO
	// node
	stitchingManager := stitching.NewStitchingManager(worker.config.stitchingPropertyType)

	nodes := currTask.NodeList()
	nodeNameUIDMap := make(map[string]string)
	for _, node := range nodes {
		nodeNameUIDMap[node.Name] = string(node.UID)
		stitchingManager.StoreStitchingValue(node)
	}

	nodeEntityDTOBuilder := dtofactory.NewNodeEntityDTOBuilder(worker.sink, stitchingManager)
	nodeEntityDTOs, err := nodeEntityDTOBuilder.BuildEntityDTOs(nodes)
	if err != nil {
		glog.Errorf("Error while creating node entityDTOs: %v", err)
		// TODO Node discovery fails, directly return?
	}
	glog.V(2).Infof("Worker %s builds %d node entityDTOs.", worker.id, len(nodeEntityDTOs))
	discoveryResult = append(discoveryResult, nodeEntityDTOs...)

	// pod
	pods := currTask.PodList()
	podEntityDTOBuilder := dtofactory.NewPodEntityDTOBuilder(worker.sink, stitchingManager, nodeNameUIDMap)
	podEntityDTOs, err := podEntityDTOBuilder.BuildEntityDTOs(pods)
	if err != nil {
		glog.Errorf("Error while creating pod entityDTOs: %v", err)
		// TODO Pod discovery fails, directly return?
	}
	discoveryResult = append(discoveryResult, podEntityDTOs...)
	glog.V(2).Infof("Worker %s builds %d pod entityDTOs.", worker.id, len(podEntityDTOs))

	// application
	applicationEntityDTOBuilder := dtofactory.NewApplicationEntityDTOBuilder(worker.sink)
	appEntityDTOs, err := applicationEntityDTOBuilder.BuildEntityDTOs(pods)
	if err != nil {
		glog.Errorf("Error while creating application entityDTOs: %v", err)
		// TODO Application discovery fails, return?
	}
	discoveryResult = append(discoveryResult, appEntityDTOs...)
	glog.V(2).Infof("Worker %s builds %d application entityDTOs.", worker.id, len(appEntityDTOs))

	// Send result
	glog.V(3).Infof("Discovery result of worker %s is: %++v", worker.id, discoveryResult)

	result := task.NewTaskResult(worker.id, task.TaskSucceeded).WithContent(discoveryResult)
	return result
}
