package dtofactory

import (
	"fmt"
	"github.com/golang/glog"

	api "k8s.io/client-go/pkg/api/v1"

	"github.com/turbonomic/kubeturbo/pkg/discovery/dtofactory/property"
	"github.com/turbonomic/kubeturbo/pkg/discovery/metrics"
	"github.com/turbonomic/kubeturbo/pkg/discovery/task"
	"github.com/turbonomic/kubeturbo/pkg/discovery/util"

	sdkbuilder "github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

/*
*  Application will sell:
        (1) Transactions and ResponseTime(Latency);

   Application will buy:
        (1) VMem, VCPU, and Application(binding it to container)
*/

var (
	applicationResourceCommodityBought = []metrics.ResourceType{
		metrics.CPU,
		metrics.Memory,
	}
)

type applicationEntityDTOBuilder struct {
	generalBuilder
	vcluster *metrics.VirtualCluster
}

func NewApplicationEntityDTOBuilder(sink *metrics.EntityMetricSink) *applicationEntityDTOBuilder {
	return &applicationEntityDTOBuilder{
		generalBuilder: newGeneralBuilder(sink),
	}
}

// get hosting node cpu frequency
func (builder *applicationEntityDTOBuilder) getNodeCPUFrequency(pod *api.Pod) (float64, error) {
	key := util.NodeKeyFromPodFunc(pod)
	cpuFrequencyUID := metrics.GenerateEntityStateMetricUID(task.NodeType, key, metrics.CpuFrequency)
	cpuFrequencyMetric, err := builder.metricsSink.GetMetric(cpuFrequencyUID)
	if err != nil {
		err := fmt.Errorf("Failed to get cpu frequency from sink for node %s: %v", key, err)
		glog.Error(err)
		return 0.0, err
	}

	cpuFrequency := cpuFrequencyMetric.GetValue().(float64)
	return cpuFrequency, nil
}

func (builder *applicationEntityDTOBuilder) BuildEntityDTOs(pods []*api.Pod) ([]*proto.EntityDTO, error) {
	var result []*proto.EntityDTO

	for _, pod := range pods {
		podFullName := util.GetPodClusterID(pod)
		nodeCPUFrequency, err := builder.getNodeCPUFrequency(pod)
		if err != nil {
			glog.Errorf("failed to build ContainerDTOs for pod[%s]: %v", podFullName, err)
			continue
		}

		vpod, exist := builder.vcluster.Pods[podFullName];
		if !exist {
			err := fmt.Errorf("Potential bug: cannot find pod[%v] from vcluster.")
			glog.Errorf(err.Error())
			glog.Errorf("Failed to build EntityDTOs for pod[%v]'s containers.", podFullName)
			//vpod = nil
			continue
		}

		podId := string(pod.UID)
		for i := range pod.Spec.Containers {
			//1. Id and Name
			//container := &(pod.Spec.Containers[i])
			containerId := util.ContainerIdFunc(podId, i)
			appId := util.ApplicationIdFunc(containerId)
			displayName := util.ApplicationDisplayName(podFullName)

			ebuilder := sdkbuilder.NewEntityDTOBuilder(proto.EntityDTO_APPLICATION, appId).
				DisplayName(displayName)

			//2. sold commodities: transaction & latency
			commoditiesSold, err := builder.getCommoditiesSold(appId, i, vpod)
			if err != nil {
				glog.Errorf("Failed to create Application(%s) entityDTO: %v", displayName, err)
				continue
			}
			ebuilder.SellsCommodities(commoditiesSold)

			//3. bought commodities: vcpu/vmem/application
			commoditiesBought, err := builder.getApplicationCommoditiesBought(appId, podFullName, containerId, nodeCPUFrequency)
			if err != nil {
				glog.Errorf("Failed to create Application(%s) entityDTO: %v", displayName, err)
				continue
			}
			provider := sdkbuilder.CreateProvider(proto.EntityDTO_CONTAINER, containerId)
			ebuilder.Provider(provider).BuysCommodities(commoditiesBought)

			//4. set properties
			properties := builder.getApplicationProperties(pod, i)
			ebuilder.WithProperties(properties)

			if !util.Monitored(pod) {
				ebuilder.Monitored(false)
			}

			appType := util.GetAppType(pod)
			ebuilder.ApplicationData(&proto.EntityDTO_ApplicationData{
				Type: &appType,
			})

			//5. build the entityDTO
			entityDTO, err := ebuilder.Create()
			if err != nil {
				glog.Errorf("Failed to build Application entityDTO based on application %s: %s", displayName, err)
				continue
			}
			result = append(result, entityDTO)
		}
	}

	return result, nil
}

func (builder *applicationEntityDTOBuilder) getCommoditiesSold(appId string, index int, vpod *metrics.Pod) ([]*proto.CommodityDTO, error) {
	var result []*proto.CommodityDTO
	usedTrans := 0
	usedLatency := 0
	// TODO: If container is not the container, should it still sell transactions and latency?
	if index == vpod.MainContainerIdx { // && vpod.Service != nil
		usedTrans = vpod.Transaction.Used
		usedLatency = vpod.Latency.Used
	}

	//1. transaction per second
	ebuilder := sdkbuilder.NewCommodityDTOBuilder(proto.CommodityDTO_TRANSACTION).Key(appId).
		Capacity(vpod.Transaction.Capacity).
		Used(usedTrans)
	tranCommodity, err := ebuilder.Create()
	if err != nil {
		glog.Errorf("Failed to create Transaction commodity for App(%v-%d) for selling: %v", vpod.FullName, index, err)
		return result, err
	}
	result = append(result, tranCommodity)

	//2. latency
	ebuilder2 := sdkbuilder.NewCommodityDTOBuilder(proto.CommodityDTO_RESPONSE_TIME).Key(appId).
		Capacity(vpod.Latency.Capacity).
		Used(usedLatency)
	latencyCommodity, err := ebuilder2.Create()
	if err != nil {
		glog.Errorf("Failed to create Latency commodity for App(%v-%d) for selling: %v", vpod.FullName, index, err)
		return result, err
	}
	result = append(result, latencyCommodity)

	return result, nil
}

// Build the bought commodities by each application.
// An application buys vCPU, vMem and Application commodity from a container.
func (builder *applicationEntityDTOBuilder) getApplicationCommoditiesBought(appId, podName, containerId string, cpuFrequency float64) ([]*proto.CommodityDTO, error) {
	var commoditiesBought []*proto.CommodityDTO

	converter := NewConverter().Set(func(input float64) float64 { return input * cpuFrequency }, metrics.CPU)

	// Resource commodities.
	resourceCommoditiesBought, err := builder.getResourceCommoditiesBought(task.ApplicationType, appId, applicationResourceCommodityBought, converter, nil)
	if err != nil {
		return nil, err
	}
	if len(resourceCommoditiesBought) != len(applicationResourceCommodityBought) {
		err = fmt.Errorf("mismatch num of commidities (%d Vs. %d) for application:%s, %s", len(resourceCommoditiesBought), len(applicationResourceCommodityBought), podName, appId)
		glog.Error(err)
		//return nil, err
	}
	commoditiesBought = append(commoditiesBought, resourceCommoditiesBought...)

	// Application commodity
	applicationCommBought, err := sdkbuilder.NewCommodityDTOBuilder(proto.CommodityDTO_APPLICATION).
		Key(containerId).
		Create()
	if err != nil {
		return nil, err
	}
	commoditiesBought = append(commoditiesBought, applicationCommBought)
	return commoditiesBought, nil
}

// Get the properties of the pod. This includes property related to application cluster property.
func (builder *applicationEntityDTOBuilder) getApplicationProperties(pod *api.Pod, index int) []*proto.EntityDTO_EntityProperty {
	var properties []*proto.EntityDTO_EntityProperty
	// additional node cluster info property.
	appProperties := property.AddHostingPodProperties(pod.Namespace, pod.Name, index)
	properties = append(properties, appProperties...)

	return properties
}
