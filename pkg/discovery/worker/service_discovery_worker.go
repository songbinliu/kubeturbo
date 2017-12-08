package worker

import (
	"errors"
	"github.com/golang/glog"

	"github.com/turbonomic/kubeturbo/pkg/discovery/dtofactory"
	"github.com/turbonomic/kubeturbo/pkg/discovery/metrics"
	"github.com/turbonomic/kubeturbo/pkg/discovery/task"

	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

const (
	k8sSvcDiscWorkerID string = "ServiceDiscoveryWorker"
)

type k8sServiceDiscoveryWorker struct {
	id       string
	vCluster *metrics.VirtualCluster
}

func NewK8sServiceDiscoveryWorker(vc *metrics.VirtualCluster) (*k8sServiceDiscoveryWorker, error) {
	return &k8sServiceDiscoveryWorker{
		id:       k8sSvcDiscWorkerID,
		vCluster: vc,
	}, nil
}

// post-process the entityDTOs and create service entityDTOs.
func (svcDiscWorker *k8sServiceDiscoveryWorker) Do(entityDTOs []*proto.EntityDTO) *task.TaskResult {

	builder := dtofactory.NewVAppEntityDTOBuilder(svcDiscWorker.vCluster)

	result, err := builder.BuildEntityDTO()
	if err != nil {
		return task.NewTaskResult(svcDiscWorker.id, task.TaskFailed).WithErr(errors.New("Failed to build VApp EntityDTOs"))
	}

	glog.V(4).Infof("Service discovery result is: %++v", result)
	glog.V(3).Infof("There are %d virtualApp entityDTOs", len(result))
	return task.NewTaskResult(svcDiscWorker.id, task.TaskSucceeded).WithContent(result)
}
