package kubeturbo

import (
	"github.com/turbonomic/kubeturbo/pkg/discovery/stitching"
	"github.com/turbonomic/kubeturbo/pkg/kubeclient"
	client "k8s.io/client-go/kubernetes"
)

// Configuration created using the parameters passed to the kubeturbo service container.
type Config struct {
	tapSpec *K8sTAPServiceSpec

	StitchingPropType stitching.StitchingPropertyType

	//VMPriority: priority of VM in supplyChain definition from kubeturbo, should be less than 0;
	VMPriority int32
	//VMIsBase: Is VM is the base template from kubeturbo, when stitching with other VM probes, should be false;
	VMIsBase bool

	Client        *client.Clientset
	KubeletClient *kubeclient.KubeletClient

	// Close this to stop all reflectors
	StopEverything chan struct{}

	DiscoveryIntervalSec int
}

func NewVMTConfig2() *Config {
	cfg := &Config{
		StopEverything: make(chan struct{}),
	}

	return cfg
}

func (c *Config) WithKubeClient(client *client.Clientset) *Config {
	c.Client = client
	return c
}

func (c *Config) WithKubeletClient(client *kubeclient.KubeletClient) *Config {
	c.KubeletClient = client
	return c
}

func (c *Config) WithTapSpec(spec *K8sTAPServiceSpec) *Config {
	c.tapSpec = spec
	return c
}

// Create the StitchingPropertyType for reconciling the kubernetes cluster nodes with the infrastructure VMs
func (c *Config) UsingUUIDStitch(useUUID bool) *Config {
	stitchingPropType := stitching.IP
	if useUUID {
		stitchingPropType = stitching.UUID
	}
	c.StitchingPropType = stitchingPropType
	return c
}

func (c *Config) WithVMPriority(p int32) *Config {
	c.VMPriority = p
	return c
}

func (c *Config) WithVMIsBase(isBase bool) *Config {
	c.VMIsBase = isBase
	return c
}

func (c *Config) WithDiscoveryInterval(di int) *Config {
	c.DiscoveryIntervalSec = di
	return c
}
