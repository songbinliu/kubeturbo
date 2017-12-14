package discovery

import (
	"fmt"
	"sync"
	"time"

	kubeClient "k8s.io/client-go/kubernetes"

	"github.com/turbonomic/kubeturbo/pkg/cluster"
	"github.com/turbonomic/kubeturbo/pkg/discovery/configs"
	"github.com/turbonomic/kubeturbo/pkg/discovery/vcluster"
	"github.com/turbonomic/kubeturbo/pkg/discovery/worker"
	"github.com/turbonomic/kubeturbo/pkg/discovery/worker/compliance"
	"github.com/turbonomic/kubeturbo/pkg/registration"

	sdkprobe "github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"

	"github.com/golang/glog"
)

const (
	// TODO make this number programmatically.
	workerCount int = 4
)

type DiscoveryClientConfig struct {
	k8sClusterScraper *cluster.ClusterScraper

	probeConfig *configs.ProbeConfig

	targetConfig *configs.K8sTargetConfig
}

func NewDiscoveryConfig(kubeClient *kubeClient.Clientset, probeConfig *configs.ProbeConfig, targetConfig *configs.K8sTargetConfig) *DiscoveryClientConfig {
	return &DiscoveryClientConfig{
		k8sClusterScraper: cluster.NewClusterScraper(kubeClient),
		probeConfig:       probeConfig,
		targetConfig:      targetConfig,
	}
}

type K8sDiscoveryClient struct {
	config          *DiscoveryClientConfig
	vcBuilderConfig *vcluster.VClusterBuilderConfig

	dispatcher      *worker.Dispatcher
	resultCollector *worker.ResultCollector

	wg sync.WaitGroup
}

func NewK8sDiscoveryClient(config *DiscoveryClientConfig, vcBuilderConfig *vcluster.VClusterBuilderConfig) *K8sDiscoveryClient {
	// make maxWorkerCount of result collector twice the worker count.
	resultCollector := worker.NewResultCollector(workerCount * 2)

	dispatcherConfig := worker.NewDispatcherConfig(config.k8sClusterScraper, config.probeConfig, workerCount)
	dispatcher := worker.NewDispatcher(dispatcherConfig)
	dispatcher.Init(resultCollector)

	dc := &K8sDiscoveryClient{
		config:          config,
		vcBuilderConfig: vcBuilderConfig,
		dispatcher:      dispatcher,
		resultCollector: resultCollector,
	}
	return dc
}

func (dc *K8sDiscoveryClient) GetAccountValues() *sdkprobe.TurboTargetInfo {
	var accountValues []*proto.AccountValue
	targetConf := dc.config.targetConfig
	// Convert all parameters in clientConf to AccountValue list
	targetID := registration.TargetIdentifierField
	accVal := &proto.AccountValue{
		Key:         &targetID,
		StringValue: &targetConf.TargetIdentifier,
	}
	accountValues = append(accountValues, accVal)

	username := registration.Username
	accVal = &proto.AccountValue{
		Key:         &username,
		StringValue: &targetConf.TargetUsername,
	}
	accountValues = append(accountValues, accVal)

	password := registration.Password
	accVal = &proto.AccountValue{
		Key:         &password,
		StringValue: &targetConf.TargetPassword,
	}
	accountValues = append(accountValues, accVal)

	targetInfo := sdkprobe.NewTurboTargetInfoBuilder(targetConf.ProbeCategory, targetConf.TargetType, targetID, accountValues).Create()
	return targetInfo
}

// Validate the Target
func (dc *K8sDiscoveryClient) Validate(accountValues []*proto.AccountValue) (*proto.ValidationResponse, error) {
	glog.V(2).Infof("Validating Kubernetes target...")

	// TODO: connect to the client and get validation response
	validationResponse := &proto.ValidationResponse{}

	return validationResponse, nil
}

// DiscoverTopology receives a discovery request from server and start probing the k8s.
func (dc *K8sDiscoveryClient) Discover(accountValues []*proto.AccountValue) (*proto.DiscoveryResponse, error) {
	currentTime := time.Now()
	newDiscoveryResultDTOs, err := dc.discoverWithNewFramework()
	if err != nil {
		glog.Errorf("Failed to use the new framework to discover current Kubernetes cluster: %s", err)
	}

	discoveryResponse := &proto.DiscoveryResponse{
		EntityDTO: newDiscoveryResultDTOs,
	}

	newFrameworkDiscTime := time.Now().Sub(currentTime).Seconds()
	glog.V(2).Infof("New framework discovery time: %.3f seconds", newFrameworkDiscTime)

	return discoveryResponse, nil
}

func forDebug(dtos []*proto.EntityDTO) {
	for i := range dtos {
		if dtos[i].GetEntityType() == proto.EntityDTO_VIRTUAL_MACHINE {
			glog.V(2).Infof("[%s]:\n %++v", dtos[i].GetDisplayName(), dtos[i])
			continue
		}
	}
}

func (dc *K8sDiscoveryClient) discoverWithNewFramework() ([]*proto.EntityDTO, error) {
	result := []*proto.EntityDTO{}

	//1. build the virtual cluster with metrics
	vclusterBuilder := vcluster.NewVCluterBuilder(dc.vcBuilderConfig)
	vclusterBuilder.SetNameId("unknown-cluster", "unknown-Id")
	vcluster, err := vclusterBuilder.BuildCluster()
	if err != nil {
		glog.Errorf("Failed to build virtual cluster: %v", err)
		return result, nil
	}

	//2. build DTOs based on the virtualCluster
	nodes, err := dc.config.k8sClusterScraper.GetAllNodes()
	if err != nil {
		return nil, fmt.Errorf("Failed to get all nodes in the cluster: %s", err)
	}

	//2.1 build node, pod, container, and app Entity-DTOs
	workerCount := dc.dispatcher.Dispatch(nodes, vcluster)
	entityDTOs := dc.resultCollector.Collect(workerCount)
	glog.V(2).Infof("Discovery workers have finished discovery work with %d entityDTOs built. Now performing service discovery...", len(entityDTOs))

	forDebug(entityDTOs)

	//2.2 add affinity/anti-affinity commodities
	glog.V(2).Infof("begin to process affinity.")
	affinityProcessorConfig := compliance.NewAffinityProcessorConfig(dc.config.k8sClusterScraper)
	affinityProcessor, err := compliance.NewAffinityProcessor(affinityProcessorConfig)
	if err != nil {
		glog.Errorf("Failed during process affinity rules: %s", err)
	} else {
		entityDTOs = affinityProcessor.ProcessAffinityRules(entityDTOs)
	}

	//2.3 build vApp Entity-DTOs
	glog.V(2).Infof("begin to generate service EntityDTOs.")
	svcDiscWorker, err := worker.NewK8sServiceDiscoveryWorker(vcluster)
	svcDiscResult := svcDiscWorker.Do(entityDTOs)
	if svcDiscResult.Err() != nil {
		glog.Errorf("Failed to discover services from current Kubernetes cluster with the new discovery framework: %s", svcDiscResult.Err())
	} else {
		entityDTOs = append(entityDTOs, svcDiscResult.Content()...)
	}

	glog.V(2).Infof("There are %d entityDTOs.", len(entityDTOs))

	return entityDTOs, nil
}
