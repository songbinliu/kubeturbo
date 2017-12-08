package metrics

import (
	"fmt"

	"github.com/turbonomic/kubeturbo/pkg/discovery/monitoring/istio"

	"github.com/golang/glog"
	api "k8s.io/client-go/pkg/api/v1"
)

const (
	defaultLatencyCapacity = 500 // 500 ms
	defaultTransactionCapacity = 200 // 200 tps
)

type ObjectMeta struct {
	Name       string
	Namespace  string
	UUID       string
	Kind       string
	ProviderID string
}

type Resource struct {
	Capacity float64
	Used     float64
	Reserved float64
}

//type KeyedResource struct {
//	Key string
//	Capacity float64
//}

// virtual machine
type VNode struct {
	ObjectMeta
	ClusterId string

	Memory Resource       // memory unit is KB
	CPU    Resource       // cpu unit is cpu time in milliseconds: 900m = 0.9 cpu core
	CoreFrequency float64 // single cpu-core frequency in Mhz

	Detail *api.Node

	// key=pod.FullName
	Pods map[string]*Pod
}

type Pod struct {
	ObjectMeta
	FullName string //FullName = "Namespace/Name"

	Memory Resource      // memory unit is KB
	CPU    Resource      // cpu unit is cpu time in milliseconds: 900m = 0.9 cpu core

	Transaction Resource
	Latency Resource

	Detail *api.Pod

	//TODO: There will be error if the Pod is associated with multiple services.
	// This is because the transactions monitored by Service does not know who is the provider.
	Service *VirtualApp
	MainContainerIdx int
}

type VirtualApp struct {
	ObjectMeta
	FullName string //FullName = "Namespace/Name"

	Transaction Resource
	Latency     Resource

	// indexed by Pod.FullName
	Pods map[string]*Pod

	// indexed by Port number
	Ports map[int]struct{}
}

type VirtualCluster struct {
	ObjectMeta

	// indexed by node.Name
	Nodes    map[string]*VNode

	// indexed by pod.fullName (namespace/name)
	Pods     map[string]*Pod

	// indexed by vapp.fullName (namespace/name)
	Services map[string]*VirtualApp
}

func NewCluster(name, uid string) *VirtualCluster {
	return &VirtualCluster {
		ObjectMeta.Name: name,
		ObjectMeta.UUID: uid,
		Nodes: make(map[string]*VNode),
		Pods: make(map[string]*Pod),
		Services: []*VirtualApp{},
	}
}

func (vc *VirtualCluster) AddVnode(vnode *VNode) error {
	uid := vnode.UUID

	if b, exist := vc.Nodes[uid]; exist {
		err := fmt.Errorf("Vnode with uid(%v) already exists : %++v", uid, b)
		glog.Errorf(err.Error())
		return err
	}

	vc.Nodes[uid] = vnode
	return nil
}

func (vc *VirtualCluster) AddPod(pod *Pod, hostName string) error {
	uid := pod.FullName
	if b, exist := vc.Pods[uid]; exist {
		err := fmt.Errorf("Pod with uid(%v) already exists : %++v", uid, b)
		glog.Errorf(err.Error())
		return err
	}

	node, exist := vc.Nodes[hostName]
	if !exist {
		err := fmt.Errorf("hosting node(%v) does not exist.", hostName)
		glog.Errorf(err.Error())
		return err
	}

	vc.Pods[uid] = pod
	pod.ProviderID = node.Name
	node.AddPod(pod)
	return nil
}

func (vc *VirtualCluster) AddVirtualApp(vapp *VirtualApp) error {
	fullName := vapp.FullName
	if b, exist := vc.Services[fullName]; exist {
		err := fmt.Errorf("Pod with uid(%v) already exists : %++v", fullName, b)
		glog.Errorf(err.Error())
		return err
	}

	vc.Services[fullName] = vapp
	return nil
}

func (vc *VirtualCluster) ConnectPodVapp(podName string, vapp *VirtualApp) error {
	//1. find the pod in the cluster
	pod, exist := vc.Pods[podName]
	if !exist {
		err := fmt.Errorf("Failed to connect pod[%v] to vapp[%v], pod does not exist.", podName, vapp.FullName)
		glog.Error(err.Error())
		return err
	}

	//2. check whether this pod is already associcated with a service
	// TODO: This will be problem when the pod is associated with multiple services.
	if pod.Service != nil {
		glog.Errorf("Potential bug: pod[%v] is associated with multiple services.", pod.FullName)

		if len(pod.Service.Pods) == 1 {
			glog.Errorf("Won't associate pod[%v] to service[%v]", pod.FullName, pod.Service.FullName)
			return
		} else {
			pod.Service.DeletePod(pod)
		}
	}

	pod.Service = vapp
	vapp.AddPod(pod)

	//3. find the main container of the pod
	if pod.Detail == nil {
		glog.Warningf("Pod's detail is empty.")
		return nil
	}

	containers := pod.Detail.Spec.Containers
	pod.MainContainerIdx = findContainerIndex(containers, vapp.Ports)
	if pod.MainContainerIdx < 0 {
		glog.Warningf("Failed to find the main container of pod[%v], use the first one")
		pod.MainContainerIdx = 0
	}
	return nil
}

// set capacity and reservation
func (vc *VirtualCluster) SetCapacity() {
	//1. set Capacity for nodes and pods
	if len(vc.Nodes) == 0 {
		glog.Warningf("no node in cluster.")
		return
	}

	if len(vc.Pods) == 0 {
		glog.Warningf("no pod in cluster.")
	}

	for _, node := range vc.Nodes {
		node.SetCapacity()
		if len (node.Pods) == 0 {
			glog.V(2).Infof("Node(%s) has no pod.", node.Name)
			continue
		}

		for _, pod := range node.Pods {
			pod.SetCapacity(node)
		}
	}

	//2. set Capacity for services
	for _, vapp := range vc.Services {
		vapp.SetCapacity()
	}

	return nil
}

func (vc *VirtualCluster) SetAppMetric(podMetrics, svcMetrics istio.MetricSet) error {
	glog.V(2).Infof("Got %d Pod metrics, %d Service metrics", len(podMetrics), len(svcMetrics))

	//1. set metrics for services
	counter := 0
	if len(vc.Services) == 0 {
		glog.Warningf("no service in cluster")
	} else {
		for k, v := range svcMetrics {
			if vapp, exist := vc.Services[k]; exist {
				vapp.Transaction.Used = v.RequestPerSecond
				vapp.Latency.Used = v.Latency
				counter ++
			}
		}
	}
	glog.V(2).Infof("Add app metrics for [%d/%d] Services.", counter, len(vc.Services))

	//2. set metrics for Pods
	counter = 0
	if len(vc.Pods) == 0 {
		glog.Warningf("no pod in cluster")
		return nil
	}

	for k, v := range podMetrics {
		if pod, exist := vc.Pods[k]; exist {
			pod.Transaction.Used = v.RequestPerSecond
			pod.Latency.Used = v.Latency
			counter ++
		}
	}
	glog.V(2).Infof("Add app metrics for [%d/%d] Pods.", counter, len(vc.Pods))
	return nil
}

func (vc *VirtualCluster) SetContainerMetric() error {
	if len(vc.Nodes) == 0 {
		glog.Warningf("There is no nodes in cluster.")
		return nil
	}

	for _, node := range vc.Nodes {
		//TODO: get the real value from kubelet
		node.CoreFrequency = 2200.0 // 2200.0 Mhz
	}

	return nil
}

func NewVNode(node *api.Node, clusterId string) *VNode {
	return &VNode{
		ObjectMeta.Name: node.Name,
		ObjectMeta.UUID: node.UID,
		ObjectMeta.Kind: "vnode",
		ClusterId: clusterId,
		Pods: make(map[string]*Pod),
		Detail: node,
	}
}

func (vn *VNode) AddPod(pod *Pod) {
	vn.Pods[pod.FullName] = pod
}

func (vn *VNode) SetCapacity() error {
	if vn.Detail == nil {
		err := fmt.Errorf("No detail info for node[%v]", vn.Name)
		glog.Error(err.Error())
		return err
	}

	res := vn.Detail.Status.Allocatable
	vn.Memory.Capacity = res.Memory().Value()
	vn.CPU.Capacity = res.Cpu().MilliValue()
	return nil
}

func NewPod(d *api.Pod) *Pod {
	return &Pod {
		ObjectMeta.Name: d.Name,
		ObjectMeta.Namespace: d.Namespace,
		ObjectMeta.UUID: d.UID,
		ObjectMeta.Kind: "pod",
		FullName: genFullName(d.Namespace, d.Name),
		Detail: d,
		MainContainerIdx: 0,
	}
}

// set resource reservation and capacity
func (pod *Pod) SetCapacity(node *VNode) error {
	if pod.Detail == nil {
		err := fmt.Errorf("No detail info for node[%v]", pod.FullName)
		glog.Error(err.Error())
		return err
	}

	// TODO:
	pod.Latency.Capacity = defaultLatencyCapacity
	pod.Transaction.Capacity = defaultTransactionCapacity

	pod.CPU.Capacity = node.CPU.Capacity
	pod.Memory.Capacity = node.Memory.Capacity

	// reserved = sum(container.Requests)
	pod.Memory.Reserved = 0
	pod.CPU.Reserved = 0

	containers := pod.Detail.Spec.Containers
	for i := range containers {
		container := &containers[i]

		req := container.Resources.Requests
		pod.Memory.Reserved += req.Memory().Value()
		pod.CPU.Reserved += req.Cpu().MilliValue()
	}

	return nil
}

func NewVirtualApp(svc *api.Service) *VirtualApp {
	ports := make(map[int]struct{})
	for _, port := range svc.Spec.Ports {
		n := port.TargetPort.IntValue()
		if n == 0 {
			n = int(port.Port)
		}

		ports[n] = struct{}{}
	}

	return &VirtualApp{
		ObjectMeta.Name: svc.Name,
		ObjectMeta.Namespace: svc.Namespace,
		ObjectMeta.UUID: svc.UID,
		ObjectMeta.Kind: "service",
		FullName: genFullName(svc.Namespace, svc.Name),
		Ports: ports,
		Pods: make(map[string]*Pod),
	}
}

func (vapp *VirtualApp) SetCapacity() {
	if len(vapp.Pods) == 0 {
		glog.Warningf("Failed to set Capacity for vapp[%v]: no pods associated with it.", vapp.FullName)
		return
	}

	totalTrans := 0
	totalLatency := 0
	for _, pod := range vapp.Pods {
		totalTrans += pod.Transaction.Capacity
		if totalLatency < 0.0001 || totalLatency > pod.Latency.Capacity {
			totalLatency = pod.Latency.Capacity
		}
	}
	vapp.Transaction.Capacity = totalTrans
	vapp.Latency.Capacity = totalLatency
	return
}

func (vapp *VirtualApp) AddPod(pod *Pod) {
	vapp.Pods[pod.FullName] = pod
}

func (vapp *VirtualApp) DeletePod(pod *Pod) {
	delete(vapp.Pods, pod.FullName)
}

func genFullName(namespace, name string) string {
	return fmt.Sprintf("%v/%v", namespace, name)
}

// find the index of the containers of a Pod, which exposes web service on the ports
// by port
func findContainerIndex(containers []api.Container, ports map[int]struct{}) int {
	found := false
	idx := - 1
	for i := range containers {
		container := &containers[i]

		for _, port := range container.Ports {
			rport := int(port.ContainerPort)
			if _, exist := ports[rport]; exist {
				found = true
				break
			}
		}

		if found {
			idx = i
			break
		}
	}

	return idx
}
