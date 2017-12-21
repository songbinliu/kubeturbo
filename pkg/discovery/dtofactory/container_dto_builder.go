package dtofactory

import (
	"fmt"
	"github.com/golang/glog"
	api "k8s.io/client-go/pkg/api/v1"

	"github.com/turbonomic/kubeturbo/pkg/discovery/dtofactory/property"
	"github.com/turbonomic/kubeturbo/pkg/discovery/metrics"
	"github.com/turbonomic/kubeturbo/pkg/discovery/util"
	sdkbuilder "github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

// commodity to sell: VCPU/VMemory/Application
// commodity to buy: VCPU/VMemory with reservation, and VMPMAccess to bind to pod

type containerDTOBuilder struct {
	generalBuilder
}

type containerInfo struct {
	fullName string
	uuid     string

	index     int
	name      string
	podId     string
	pod       *api.Pod
	container *api.Container

	cpuFreq   float64
	monitored bool

	// cpu commodity
	cpuUsed     float64
	cpuReserved float64
	cpuCap      float64

	//2. memory commodity
	memUsed     float64
	memReserved float64
	memCap      float64
}

func (c *containerInfo) genNameID() {
	c.fullName = util.ContainerNameFunc(c.pod, c.container)
	c.uuid = util.ContainerIdFunc(c.podId, c.index)
}

func (c *containerInfo) buildDTO() (*proto.EntityDTO, error) {

	ebuilder := sdkbuilder.NewEntityDTOBuilder(proto.EntityDTO_CONTAINER, c.uuid)
	ebuilder.DisplayName(c.fullName)

	//1. build commodities to sell
	commS, err := c.buildSelling()
	if err != nil {
		glog.Errorf("Failed to build container(%v) DTO: %v", c.fullName, err)
		return nil, err
	}
	ebuilder.SellsCommodities(commS)

	//2. build commodities to buy
	commB, err := c.buildBuying()
	if err != nil {
		glog.Errorf("Failed to build container(%v) DTO: %v", c.fullName, err)
		return nil, err
	}

	provider := sdkbuilder.CreateProvider(proto.EntityDTO_CONTAINER_POD, c.podId)
	ebuilder.Provider(provider).BuysCommodities(commB)

	//3. set properties
	properties := addContainerProperty(c.pod, c.index)
	ebuilder.WithProperties(properties)
	ebuilder.Monitored(c.monitored)

	//4. create it
	dto, err := ebuilder.Create()
	if err != nil {
		glog.Errorf("Failed to build container(%v) DTO: %v", c.fullName, err)
		return nil, err
	}

	return dto, nil
}

func (c *containerInfo) buildSelling() ([]*proto.CommodityDTO, error) {
	var result []*proto.CommodityDTO

	//1. cpu
	resize := true
	cpuComm := buildSellComm(c.cpuUsed, c.cpuCap, resize, proto.CommodityDTO_VCPU)
	if cpuComm == nil {
		err := fmt.Errorf("Failed to build selling cpu commodity for container: %v", c.fullName)
		glog.Errorf(err.Error())
		return result, err
	}
	result = append(result, cpuComm)

	//2. memory
	memComm := buildSellComm(c.memUsed, c.memCap, resize, proto.CommodityDTO_VMEM)
	if memComm == nil {
		err := fmt.Errorf("Failed to build selling memory commodity for container: %v", c.fullName)
		glog.Errorf(err.Error())
		return result, err
	}
	result = append(result, memComm)

	//3. Application
	appComm := buildSellKeyComm(c.uuid, applicationCommodityDefaultCapacity, proto.CommodityDTO_APPLICATION)
	if appComm == nil {
		err := fmt.Errorf("Failed to build selling app commodity for container: %v", c.fullName)
		glog.Errorf(err.Error())
	}
	result = append(result, appComm)

	return result, nil
}

// buy cpu/memory with used and reservation, and vmpmAccess
func (c *containerInfo) buildBuying() ([]*proto.CommodityDTO, error) {
	var result []*proto.CommodityDTO

	//1. cpu
	cpuComm := buildBuyComm(c.cpuUsed, c.cpuReserved, proto.CommodityDTO_VCPU)
	if cpuComm == nil {
		err := fmt.Errorf("Failed to build cpu commodity to buy: %v", c.fullName)
		glog.Error(err.Error())
		return result, err
	}
	result = append(result, cpuComm)

	//2. memory
	memComm := buildBuyComm(c.memUsed, c.memReserved, proto.CommodityDTO_VMEM)
	if memComm == nil {
		err := fmt.Errorf("Failed to build memory commodity to buy: %v", c.fullName)
		glog.Error(err.Error())
		return result, err
	}
	result = append(result, memComm)

	//3. vmpmAccess to bind container to pod
	podComm := buildBuyKeyComm(c.podId, proto.CommodityDTO_VMPM_ACCESS)
	if podComm == nil {
		err := fmt.Errorf("Failed to builld pod access commodity to buy: %v", c.fullName)
		glog.Errorf(err.Error())
		return result, err
	}
	result = append(result, podComm)

	return result, nil
}

func newContainerInfo() *containerInfo {
	return &containerInfo{}
}

func NewContainerDTOBuilder(sink *metrics.EntityMetricSink) *containerDTOBuilder {
	return &containerDTOBuilder{
		generalBuilder: newGeneralBuilder(sink),
	}
}

// get cpu frequency
func (builder *containerDTOBuilder) getNodeCPUFrequency(pod *api.Pod) (float64, error) {
	key := util.NodeKeyFromPodFunc(pod)
	cpuFrequencyUID := metrics.GenerateEntityStateMetricUID(metrics.NodeType, key, metrics.CpuFrequency)
	cpuFrequencyMetric, err := builder.metricsSink.GetMetric(cpuFrequencyUID)
	if err != nil {
		err := fmt.Errorf("Failed to get cpu frequency from sink for node %s: %v", key, err)
		glog.Error(err)
		return 0.0, err
	}

	cpuFrequency := cpuFrequencyMetric.GetValue().(float64)
	return cpuFrequency, nil
}

func (builder *containerDTOBuilder) BuildDTOs(pods []*api.Pod) ([]*proto.EntityDTO, error) {
	var result []*proto.EntityDTO

	cinfo := newContainerInfo()

	for _, pod := range pods {
		nodeCPUFrequency, err := builder.getNodeCPUFrequency(pod)
		if err != nil {
			glog.Errorf("failed to build ContainerDTOs for pod[%s]: %v", pod.Name, err)
			continue
		}

		cinfo.pod = pod
		cinfo.podId = string(pod.UID)
		cinfo.cpuFreq = nodeCPUFrequency
		cinfo.monitored = util.Monitored(pod)

		for i := range pod.Spec.Containers {

			//1. set up container Info
			container := &(pod.Spec.Containers[i])
			cinfo.name = container.Name
			cinfo.index = i
			cinfo.container = container
			cinfo.genNameID()

			//2. get used/capacity/reserved metrics
			builder.getCPUMemoryUsed(cinfo)
			builder.getCPUMemoryCapacity(cinfo)

			//3. build container DTO
			dto, err := cinfo.buildDTO()
			if err != nil {
				glog.Errorf("Failed to build container DTO: %v", err)
				continue
			}

			result = append(result, dto)
		}
	}

	return result, nil
}

func (builder *containerDTOBuilder) getUsed(key string, entityType metrics.DiscoveredEntityType, rtype metrics.ResourceType) (float64, error) {
	usedMetricUID := metrics.GenerateEntityResourceMetricUID(entityType, key, rtype, metrics.Used)
	usedMetric, err := builder.metricsSink.GetMetric(usedMetricUID)
	if err != nil {
		glog.Errorf("Failed to get %s used for %s %s: %v", rtype, entityType, key, err)
		return 0, err
	}
	usedValue := usedMetric.GetValue().(float64)

	return usedValue, nil
}

func (builder *containerDTOBuilder) getCPUMemoryUsed(cinfo *containerInfo) error {
	usedId := util.ContainerStatNameFunc(cinfo.podId, cinfo.name)

	cpuUsed, err := builder.getUsed(usedId, metrics.ContainerType, metrics.CPU)
	if err != nil {
		glog.Errorf("Failed to get used cpu for container: %v, set to 0", usedId)
		cpuUsed = 0
	}
	memUsed, err := builder.getUsed(usedId, metrics.ContainerType, metrics.Memory)
	if err != nil {
		glog.Errorf("Failed to get used memory for container: %v, set to 0", usedId)
		memUsed = 0
	}

	cinfo.cpuUsed = cpuUsed * cinfo.cpuFreq
	cinfo.memUsed = memUsed
	return nil
}

// Get capacity and reservation from container.Spec
func (builder *containerDTOBuilder) getCPUMemoryCapacity(cinfo *containerInfo) error {
	res := cinfo.container.Resources
	cpuCap := util.UnifyCPUTime(res.Limits.Cpu())
	memCap := util.UnifyMemory(res.Limits.Memory())

	cpuReserved := util.UnifyCPUTime(res.Requests.Cpu())
	memReserved := util.UnifyMemory(res.Requests.Memory())

	cinfo.cpuCap = cpuCap * cinfo.cpuFreq
	cinfo.memCap = memCap

	cinfo.cpuReserved = cpuReserved * cinfo.cpuFreq
	cinfo.memReserved = memReserved

	return nil
}

func addContainerProperty(pod *api.Pod, index int) []*proto.EntityDTO_EntityProperty {
	var properties []*proto.EntityDTO_EntityProperty
	podProperties := property.AddHostingPodProperties(pod.Namespace, pod.Name, index)
	properties = append(properties, podProperties...)

	return properties
}

func buildSellKeyComm(key string, cap float64, cType proto.CommodityDTO_CommodityType) *proto.CommodityDTO {
	cb := sdkbuilder.NewCommodityDTOBuilder(cType)
	comm, err := cb.Capacity(cap).Key(key).Create()
	if err != nil {
		glog.Errorf("Failed to create commodity: %v", err)
		return nil
	}

	return comm
}

func buildSellComm(used, cap float64, resize bool, cType proto.CommodityDTO_CommodityType) *proto.CommodityDTO {
	cb := sdkbuilder.NewCommodityDTOBuilder(cType)
	comm, err := cb.Capacity(cap).Used(used).Resizable(resize).Create()
	if err != nil {
		glog.Errorf("Failed to create commodity: %v", err)
		return nil
	}

	return comm
}

func buildBuyKeyComm(key string, cType proto.CommodityDTO_CommodityType) *proto.CommodityDTO {
	cb := sdkbuilder.NewCommodityDTOBuilder(cType)
	comm, err := cb.Key(key).Used(0.0).Create()
	if err != nil {
		glog.Errorf("Failed to create commodity: %v", err)
		return nil
	}

	return comm
}

func buildBuyComm(used, reserved float64, cType proto.CommodityDTO_CommodityType) *proto.CommodityDTO {
	cb := sdkbuilder.NewCommodityDTOBuilder(cType)
	comm, err := cb.Reservation(reserved).Used(used).Create()
	if err != nil {
		glog.Errorf("Failed to create commodity to buy: %v", err)
		return nil
	}

	return comm
}
