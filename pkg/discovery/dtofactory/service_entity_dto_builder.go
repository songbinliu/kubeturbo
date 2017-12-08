package dtofactory

import (
	"fmt"
	"github.com/golang/glog"
	api "k8s.io/client-go/pkg/api/v1"

	"github.com/turbonomic/kubeturbo/pkg/discovery/util"
	"github.com/turbonomic/kubeturbo/pkg/discovery/metrics"
	sdkbuilder "github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

const (
	vAppPrefix string = "vApp"
)

var (
	commodityTypeBetweenAppAndService map[proto.CommodityDTO_CommodityType]struct{} = map[proto.CommodityDTO_CommodityType]struct{}{
		proto.CommodityDTO_TRANSACTION: struct{}{},
	}
)

type ServiceEntityDTOBuilder struct{
	vcluster *metrics.VirtualCluster
}

type VAppEntityDTOBuilder struct {
	vcluster *metrics.VirtualCluster
}

func NewVAppEntityDTOBuilder(vc *metrics.VirtualCluster) *VAppEntityDTOBuilder {
	return &VAppEntityDTOBuilder{
		vcluster: vc,
	}
}

func (builder *VAppEntityDTOBuilder) BuildEntityDTO() ([]*proto.EntityDTO, error) {
	result := []*proto.EntityDTO{}

	for _, vapp := range builder.vcluster.Services {
		id := string(vapp.UUID)
		serviceName := vapp.FullName
		displayName := fmt.Sprintf("%s-%s", vAppPrefix, serviceName)

		ebuilder := sdkbuilder.NewEntityDTOBuilder(proto.EntityDTO_VIRTUAL_APPLICATION, id).
			DisplayName(displayName)

		//1. commodities it will buy
		if err := builder.getCommoditiesBought(ebuilder, vapp); err != nil {
			glog.Errorf("Failed to create Vapp(%) entityDTO: %v", serviceName, err)
			continue
		}

		//2. commodities it will sell (only for monitoring)
		commoditiesSold, err := builder.getCommoditiesSold(vapp)
		if err != nil {
			glog.Errorf("Failed to create VApp(%s) entityDTO: %v", serviceName, err)
			continue
		}
		ebuilder.SellsCommodities(commoditiesSold)

		//3. virtual application data.
		vAppData := &proto.EntityDTO_VirtualApplicationData{
			ServiceType: &serviceName,
		}
		ebuilder.VirtualApplicationData(vAppData)

		//4. create EntityDTO
		entityDTO, err := ebuilder.Create()
		if err != nil {
			glog.Errorf("Failed to create EntityDTO for service[%s]: %v", serviceName, err)
			continue
		}

		result = append(result, entityDTO)
	}

	return result, nil
}

func (builder *VAppEntityDTOBuilder) getCommoditiesBought(ebuilder *sdkbuilder.EntityDTOBuilder, vapp *metrics.VirtualApp) error {
	if len(vapp.Pods) < 1 {
		err := fmt.Errorf("Failed to create commoditiesBought for Service[%s]: it is not associated with any pod.", vapp.FullName)
		glog.Error(err.Error())
		return err
	}

	for _, pod := range vapp.Pods {
		if pod.MainContainerIdx < 0 {
			glog.Warningf("Pod[%s] mainContainerIndex is not valide.", pod.FullName)
			continue
		}

		containerId := util.ContainerIdFunc(pod.UUID, pod.MainContainerIdx)
		appId := util.ApplicationIdFunc(containerId)

		provider := sdkbuilder.CreateProvider(proto.EntityDTO_APPLICATION, appId)

		tranComm, err := sdkbuilder.NewCommodityDTOBuilder(proto.CommodityDTO_TRANSACTION).
									Key(appId).
									Used(pod.Transaction.Used).
									Create()
		if err != nil {
			glog.Errorf("Faild to create transaction commodity for Vapp(%v) Pod(%v): %v", vapp.FullName, pod.FullName, err)
			continue
		}

		latencyComm, err := sdkbuilder.NewCommodityDTOBuilder(proto.CommodityDTO_RESPONSE_TIME).
			                    Key(appId).
		                        Used(pod.Latency.Used).
		                        Create()
		if err != nil {
			glog.Errorf("Faild to create latency commodity for Vapp(%v) Pod(%v): %v", vapp.FullName, pod.FullName, err)
			continue
		}

		bought := []*proto.CommodityDTO{tranComm, latencyComm}
		ebuilder.Provider(provider).BuysCommodities(bought)
	}

	return nil
}

// it will sell Transaction and Latency, for monitoring only.
func (builder *VAppEntityDTOBuilder) getCommoditiesSold(vapp *metrics.VirtualApp) ([]*proto.CommodityDTO, error) {
	result := []*proto.CommodityDTO{}
	svcId := vapp.UUID

	//1. Transactions per second
	ebuilder := sdkbuilder.NewCommodityDTOBuilder(proto.CommodityDTO_TRANSACTION).
		                   Key(svcId).
	                       Capacity(vapp.Transaction.Capacity).
	                       Used(vapp.Transaction.Used)
	tranCommodity, err := ebuilder.Create()
	if err != nil {
		glog.Errorf("Failed to create Transaction commodity for service(%v) for selling: %v", vapp.FullName, err)
		return result, err
	}
	result = append(result, tranCommodity)

	//2. Latency
	ebuilder2 := sdkbuilder.NewCommodityDTOBuilder(proto.CommodityDTO_RESPONSE_TIME).
		                   Key(svcId).
	                       Capacity(vapp.Latency.Capacity).
	                       Used(vapp.Latency.Used)
	latencyCommodity, err := ebuilder2.Create()
	if err != nil {
		glog.Errorf("Failed to create Latency commodity for service(%v) for selling: %v", vapp.FullName, err)
		return result, err
	}
	result = append(result, latencyCommodity)


	return result, nil
}

// TODO: delete from here
func (builder *ServiceEntityDTOBuilder) BuildSvcEntityDTO(servicePodMap map[*api.Service][]*api.Pod, clusterID string, appDTOs map[string]*proto.EntityDTO) ([]*proto.EntityDTO, error) {
	result := []*proto.EntityDTO{}

	for service, pods := range servicePodMap {
		id := string(service.UID)
		serviceName := util.GetServiceClusterID(service)
		displayName := fmt.Sprintf("%s-%s", vAppPrefix, serviceName)

		ebuilder := sdkbuilder.NewEntityDTOBuilder(proto.EntityDTO_VIRTUAL_APPLICATION, id).
			DisplayName(displayName)

		//1. commodities bought
		if err := builder.createCommodityBought(ebuilder, pods, appDTOs); err != nil {
			glog.Errorf("failed to create server[%s] EntityDTO: %v", serviceName, err)
			continue
		}

		//2. virtual application data.
		vAppData := &proto.EntityDTO_VirtualApplicationData{
			ServiceType: &service.Name,
		}
		ebuilder.VirtualApplicationData(vAppData)

		//3. create it
		entityDto, err := ebuilder.Create()
		if err != nil {
			glog.Errorf("failed to create service[%s] EntityDTO: %v", serviceName, err)
			continue
		}
		result = append(result, entityDto)
	}

	return result, nil
}

func (builder *ServiceEntityDTOBuilder) createCommodityBought(ebuilder *sdkbuilder.EntityDTOBuilder, pods []*api.Pod, appDTOs map[string]*proto.EntityDTO) error {

	for _, pod := range pods {
		podId := string(pod.UID)
		for i := range pod.Spec.Containers {
			containerName := util.ContainerNameFunc(pod, &(pod.Spec.Containers[i]))
			containerId := util.ContainerIdFunc(podId, i)
			appId := util.ApplicationIdFunc(containerId)

			appDTO, exist := appDTOs[appId]
			if !exist {
				glog.Errorf("Cannot find app[%s] entityDTO of container[%s].", appId, containerName)
				continue
			}

			bought, err := builder.getCommoditiesBought(appDTO)
			if err != nil {
				glog.Errorf("failed to get commodity bought from container[%s]: %v", containerName, err)
				continue
			}
			provider := sdkbuilder.CreateProvider(proto.EntityDTO_APPLICATION, appId)
			ebuilder.Provider(provider).BuysCommodities(bought)
		}
	}

	return nil
}

func (svcEntityDTOBuilder *ServiceEntityDTOBuilder) getCommoditiesBought(appDTO *proto.EntityDTO) ([]*proto.CommodityDTO, error) {
	commoditiesSoldByApp := appDTO.GetCommoditiesSold()
	var commoditiesBoughtFromApp []*proto.CommodityDTO
	for _, commSold := range commoditiesSoldByApp {
		if _, exist := commodityTypeBetweenAppAndService[commSold.GetCommodityType()]; exist {
			commBoughtByService, err := sdkbuilder.NewCommodityDTOBuilder(commSold.GetCommodityType()).
				Key(commSold.GetKey()).
				Used(commSold.GetUsed()).
				Create()
			if err != nil {
				return nil, err
			}
			commoditiesBoughtFromApp = append(commoditiesBoughtFromApp, commBoughtByService)
		}
	}

	if len(commoditiesBoughtFromApp) < 1 {
		return nil, fmt.Errorf("no commodity found.")
	}

	return commoditiesBoughtFromApp, nil
}
