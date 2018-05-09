package registration

import (
	"github.com/golang/glog"
	"github.com/turbonomic/kubeturbo/pkg/discovery/stitching"

	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

const (
	TargetIdentifierField string = "targetIdentifier"
	Username              string = "username"
	Password              string = "password"
)

type RegistrationConfig struct {
	// The property used for stitching.
	stitchingPropertyType stitching.StitchingPropertyType
	vmPriority            int32
	vmIsBase              bool
}

func NewRegistrationClientConfig(pType stitching.StitchingPropertyType, p int32, isbase bool) *RegistrationConfig {
	return &RegistrationConfig{
		stitchingPropertyType: pType,
		vmPriority:            p,
		vmIsBase:              isbase,
	}
}

type K8sRegistrationClient struct {
	config *RegistrationConfig
}

func NewK8sRegistrationClient(config *RegistrationConfig) *K8sRegistrationClient {
	return &K8sRegistrationClient{
		config: config,
	}
}

func (rClient *K8sRegistrationClient) GetSupplyChainDefinition() []*proto.TemplateDTO {
	supplyChainFactory := NewSupplyChainFactory(rClient.config.stitchingPropertyType, rClient.config.vmPriority, rClient.config.vmIsBase)
	supplyChain, err := supplyChainFactory.createSupplyChain()
	if err != nil {
		glog.Errorf("Failed to create supply chain: %v", err)
		// TODO error handling
	}
	return supplyChain
}

func (rClient *K8sRegistrationClient) GetAccountDefinition() []*proto.AccountDefEntry {
	var acctDefProps []*proto.AccountDefEntry

	// target ID
	targetIDAcctDefEntry := builder.NewAccountDefEntryBuilder(TargetIdentifierField, "Address",
		"IP of the Kubernetes master", ".*", false, false).Create()
	acctDefProps = append(acctDefProps, targetIDAcctDefEntry)

	// username
	usernameAcctDefEntry := builder.NewAccountDefEntryBuilder(Username, "Username",
		"Username of the Kubernetes master", ".*", false, false).Create()
	acctDefProps = append(acctDefProps, usernameAcctDefEntry)

	// password
	passwordAcctDefEntry := builder.NewAccountDefEntryBuilder(Password, "Password",
		"Password of the Kubernetes master", ".*", false, true).Create()
	acctDefProps = append(acctDefProps, passwordAcctDefEntry)

	return acctDefProps
}

func (rClient *K8sRegistrationClient) GetIdentifyingFields() string {
	return TargetIdentifierField
}

func (rClient *K8sRegistrationClient) GetActionPolicy() []*proto.ActionPolicyDTO {
	glog.V(3).Infof("Begin to build Action Policies")
	ab := builder.NewActionPolicyBuilder()
	supported := proto.ActionPolicyDTO_SUPPORTED
	recommend := proto.ActionPolicyDTO_NOT_EXECUTABLE
	notSupported := proto.ActionPolicyDTO_NOT_SUPPORTED

	//1. containerPod: move, provision; not resize;
	pod := proto.EntityDTO_CONTAINER_POD
	podPolicy := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	podPolicy[proto.ActionItemDTO_MOVE] = supported
	podPolicy[proto.ActionItemDTO_PROVISION] = supported
	podPolicy[proto.ActionItemDTO_RIGHT_SIZE] = notSupported
	addActionPolicy(ab, pod, podPolicy)

	//2. container: support resize; recommend provision; not move;
	container := proto.EntityDTO_CONTAINER
	containerPolicy := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	containerPolicy[proto.ActionItemDTO_RIGHT_SIZE] = supported
	containerPolicy[proto.ActionItemDTO_PROVISION] = recommend
	containerPolicy[proto.ActionItemDTO_MOVE] = notSupported
	addActionPolicy(ab, container, containerPolicy)

	//3. application: only recommend provision; all else are not supported
	app := proto.EntityDTO_APPLICATION
	appPolicy := make(map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability)
	appPolicy[proto.ActionItemDTO_PROVISION] = recommend
	appPolicy[proto.ActionItemDTO_RIGHT_SIZE] = notSupported
	appPolicy[proto.ActionItemDTO_MOVE] = notSupported
	addActionPolicy(ab, app, appPolicy)

	return ab.Create()
}

func addActionPolicy(ab *builder.ActionPolicyBuilder,
	entity proto.EntityDTO_EntityType,
	policies map[proto.ActionItemDTO_ActionType]proto.ActionPolicyDTO_ActionCapability) {

	for action, policy := range policies {
		ab.WithEntityActions(entity, action, policy)
	}
}
