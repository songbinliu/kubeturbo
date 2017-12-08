package metrics

import (
	"fmt"
	"github.com/golang/glog"

	"github.com/turbonomic/kubeturbo/pkg/cluster"
	"github.com/turbonomic/kubeturbo/pkg/discovery/monitoring/istio"
	"github.com/turbonomic/kubeturbo/pkg/discovery/monitoring/kubelet"

	api "k8s.io/client-go/pkg/api/v1"
)

type VClusterBuilder struct {
	Name string
	UUID string

	TopoGetter            *cluster.ClusterScraper
	ContainerMetricGetter *kubelet.KubeletClient
	AppMetricGetter       *istio.AppMetricClient

	cluster *VirtualCluster
}

func NewVCluterBuilder(name, uuid string) *VClusterBuilder {
	return &VClusterBuilder{
		Name: name,
		UUID: uuid,
	}
}

func (b *VClusterBuilder) SetTopoBuilder(kube *cluster.ClusterScraper) {
	b.TopoGetter = kube
}

func (b *VClusterBuilder) SetPodMetricBuilder(klet *kubelet.KubeletClient) {
	b.ContainerMetricGetter = klet
}

func (b *VClusterBuilder) SetAppMetricBuilder(app *istio.AppMetricClient) {
	b.AppMetricGetter = app
}

func (b *VClusterBuilder) BuildCluster() (*VirtualCluster, error) {
	cluster := NewCluster(b.Name, b.UUID)
	b.cluster = cluster

	if err := b.BuildTopology(); err != nil {
		err = fmt.Errorf("Failed to build cluster: %v", err)
		glog.Error(err.Error())
		return nil, err
	}

	if err := b.AddContainerMetric(); err != nil {
		err := fmt.Errorf("Failed to add Container metric to cluster: %v", err)
		glog.Warning(err.Error())
	}

	if err := b.AddAppMetric(); err != nil {
		err = fmt.Errorf("Failed to add App metric to cluster: %v", err)
		glog.Warning(err.Error())
	}

	return cluster, nil
}

func getNodeIP(node *api.Node) string {
	addrs := node.Status.Addresses

	result := ""
	for i := range addrs {
		if addrs[i].Type == "InternalIP" {
			result = addrs[i].Address
			return result
		} else if result == "" {
			result = addrs[i].Address
		}
	}

	return result
}

func (b *VClusterBuilder) AddNodes(nodes []*api.Node) error {
	for _, node := range nodes {
		vnode := NewVNode(node, b.UUID)
		b.cluster.AddVnode(vnode)
	}
	return nil
}

func (b *VClusterBuilder) AddPods(pods []*api.Pod) error {
	for _, pod := range pods {
		vpod := NewPod(pod)
		b.cluster.AddPod(vpod, pod.Spec.NodeName)
	}

	return nil
}

func (b *VClusterBuilder) AddServices(services []*api.Service, endpoints []*api.Endpoints) error {
	dict := make(map[string]*api.Endpoints)

	for _, ep := range endpoints {
		fullName := genFullName(ep.Namespace, ep.Name)
		dict[fullName] = ep
	}

	for _, svc := range services {
		vapp := NewVirtualApp(svc)
		b.cluster.AddVirtualApp(vapp)

		// connecting the pods with this VApp
		fullName := genFullName(svc.Namespace, svc.Name)
		ep, exist := dict[fullName]
		if !exist {
			continue
		}

		for i := range ep.Subsets {
			addrs := ep.Subsets[i].Addresses
			for _, addr := range addrs {
				pod := addr.TargetRef
				if pod == nil {
					continue
				}

				podName := genFullName(pod.Namespace, pod.Name)
				b.cluster.ConnectPodVapp(podName, vapp)
			}
		}
	}

	return nil
}

func (b *VClusterBuilder) BuildTopology() error {
	if b.TopoGetter == nil {
		err := fmt.Errorf("Topology builder is unset.")
		glog.Errorf(err.Error())
		return err
	}

	// make sure that adding nodes, pods, and svc in order.
	//1. get all the nodes;
	nodes, err := b.TopoGetter.GetAllNodes()
	if err != nil {
		err := fmt.Errorf("Topology Getter failed to get all nodes: %v", err)
		glog.Errorf(err.Error())
		return err
	}
	if err := b.AddNodes(nodes); err != nil {
		glog.Errorf("Failed to add nodes to cluster: %v", err)
		return err
	}

	//2. get all the pods;
	pods, err := b.TopoGetter.GetAllPods()
	if err != nil {
		err := fmt.Errorf("Topology Getter failed to get all podes: %v", err)
		glog.Errorf(err.Error())
		return err
	}
	b.AddPods(pods)

	//3. get all services;
	services, err := b.TopoGetter.GetAllServices()
	if err != nil {
		err := fmt.Errorf("Topology Getter failed to get all services: %v", err)
		glog.Errorf(err.Error())
		return err
	}

	//4. get all endpoints;
	endpoints, err := b.TopoGetter.GetAllEndpoints()
	if err != nil {
		err := fmt.Errorf("Topology Getter failed to get all endpoints: %v", err)
		glog.Errorf(err.Error())
		return err
	}
	b.AddServices(services, endpoints)

	//5. set Capacity and Reservation
	b.cluster.SetCapacity()
	return nil
}

func (b *VClusterBuilder) AddContainerMetric() error {
	if b.ContainerMetricGetter == nil {
		err := fmt.Errorf("PodMetricGetter is unset.")
		glog.Errorf(err.Error())
		return err
	}

	//TODO: get container & pod metrics from kubelet
	b.cluster.SetContainerMetric()
	return nil
}

func (b *VClusterBuilder) AddAppMetric() error {
	if b.AppMetricGetter == nil {
		err := fmt.Errorf("AppMetricGetter is unset.")
		glog.Errorf(err.Error())
		return err
	}

	//1. get the metrics
	pods, vapps, err := b.AppMetricGetter.GetPodAppMetrics()
	if err != nil {
		glog.Errorf("Failed to get pod and service metrics: %v", err)
		return err
	}

	//2. add these metrics into the cluster
	b.cluster.SetAppMetric(pods, vapps)
	return nil
}
