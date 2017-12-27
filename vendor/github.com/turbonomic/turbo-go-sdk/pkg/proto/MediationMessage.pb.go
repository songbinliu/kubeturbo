// Code generated by protoc-gen-go.
// source: MediationMessage.proto
// DO NOT EDIT!

package proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Messages, sent from client to server
type MediationClientMessage struct {
	// Types that are valid to be assigned to MediationClientMessage:
	//	*MediationClientMessage_ValidationResponse
	//	*MediationClientMessage_DiscoveryResponse
	//	*MediationClientMessage_KeepAlive
	//	*MediationClientMessage_ActionProgress
	//	*MediationClientMessage_ActionResponse
	MediationClientMessage isMediationClientMessage_MediationClientMessage `protobuf_oneof:"mediation_client_message"`
	// this is always required in reality. it's optional here because
	// we don't know if in the future, with embedded targets, we will
	// still use it or not
	MessageID        *int32 `protobuf:"varint,15,opt,name=messageID" json:"messageID,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *MediationClientMessage) Reset()                    { *m = MediationClientMessage{} }
func (m *MediationClientMessage) String() string            { return proto.CompactTextString(m) }
func (*MediationClientMessage) ProtoMessage()               {}
func (*MediationClientMessage) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{0} }

type isMediationClientMessage_MediationClientMessage interface {
	isMediationClientMessage_MediationClientMessage()
}

type MediationClientMessage_ValidationResponse struct {
	ValidationResponse *ValidationResponse `protobuf:"bytes,2,opt,name=validationResponse,oneof"`
}
type MediationClientMessage_DiscoveryResponse struct {
	DiscoveryResponse *DiscoveryResponse `protobuf:"bytes,3,opt,name=discoveryResponse,oneof"`
}
type MediationClientMessage_KeepAlive struct {
	KeepAlive *KeepAlive `protobuf:"bytes,4,opt,name=keepAlive,oneof"`
}
type MediationClientMessage_ActionProgress struct {
	ActionProgress *ActionProgress `protobuf:"bytes,5,opt,name=actionProgress,oneof"`
}
type MediationClientMessage_ActionResponse struct {
	ActionResponse *ActionResult `protobuf:"bytes,6,opt,name=actionResponse,oneof"`
}

func (*MediationClientMessage_ValidationResponse) isMediationClientMessage_MediationClientMessage() {}
func (*MediationClientMessage_DiscoveryResponse) isMediationClientMessage_MediationClientMessage()  {}
func (*MediationClientMessage_KeepAlive) isMediationClientMessage_MediationClientMessage()          {}
func (*MediationClientMessage_ActionProgress) isMediationClientMessage_MediationClientMessage()     {}
func (*MediationClientMessage_ActionResponse) isMediationClientMessage_MediationClientMessage()     {}

func (m *MediationClientMessage) GetMediationClientMessage() isMediationClientMessage_MediationClientMessage {
	if m != nil {
		return m.MediationClientMessage
	}
	return nil
}

func (m *MediationClientMessage) GetValidationResponse() *ValidationResponse {
	if x, ok := m.GetMediationClientMessage().(*MediationClientMessage_ValidationResponse); ok {
		return x.ValidationResponse
	}
	return nil
}

func (m *MediationClientMessage) GetDiscoveryResponse() *DiscoveryResponse {
	if x, ok := m.GetMediationClientMessage().(*MediationClientMessage_DiscoveryResponse); ok {
		return x.DiscoveryResponse
	}
	return nil
}

func (m *MediationClientMessage) GetKeepAlive() *KeepAlive {
	if x, ok := m.GetMediationClientMessage().(*MediationClientMessage_KeepAlive); ok {
		return x.KeepAlive
	}
	return nil
}

func (m *MediationClientMessage) GetActionProgress() *ActionProgress {
	if x, ok := m.GetMediationClientMessage().(*MediationClientMessage_ActionProgress); ok {
		return x.ActionProgress
	}
	return nil
}

func (m *MediationClientMessage) GetActionResponse() *ActionResult {
	if x, ok := m.GetMediationClientMessage().(*MediationClientMessage_ActionResponse); ok {
		return x.ActionResponse
	}
	return nil
}

func (m *MediationClientMessage) GetMessageID() int32 {
	if m != nil && m.MessageID != nil {
		return *m.MessageID
	}
	return 0
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*MediationClientMessage) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _MediationClientMessage_OneofMarshaler, _MediationClientMessage_OneofUnmarshaler, _MediationClientMessage_OneofSizer, []interface{}{
		(*MediationClientMessage_ValidationResponse)(nil),
		(*MediationClientMessage_DiscoveryResponse)(nil),
		(*MediationClientMessage_KeepAlive)(nil),
		(*MediationClientMessage_ActionProgress)(nil),
		(*MediationClientMessage_ActionResponse)(nil),
	}
}

func _MediationClientMessage_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*MediationClientMessage)
	// mediation_client_message
	switch x := m.MediationClientMessage.(type) {
	case *MediationClientMessage_ValidationResponse:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ValidationResponse); err != nil {
			return err
		}
	case *MediationClientMessage_DiscoveryResponse:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.DiscoveryResponse); err != nil {
			return err
		}
	case *MediationClientMessage_KeepAlive:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.KeepAlive); err != nil {
			return err
		}
	case *MediationClientMessage_ActionProgress:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ActionProgress); err != nil {
			return err
		}
	case *MediationClientMessage_ActionResponse:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ActionResponse); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("MediationClientMessage.MediationClientMessage has unexpected type %T", x)
	}
	return nil
}

func _MediationClientMessage_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*MediationClientMessage)
	switch tag {
	case 2: // mediation_client_message.validationResponse
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ValidationResponse)
		err := b.DecodeMessage(msg)
		m.MediationClientMessage = &MediationClientMessage_ValidationResponse{msg}
		return true, err
	case 3: // mediation_client_message.discoveryResponse
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(DiscoveryResponse)
		err := b.DecodeMessage(msg)
		m.MediationClientMessage = &MediationClientMessage_DiscoveryResponse{msg}
		return true, err
	case 4: // mediation_client_message.keepAlive
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(KeepAlive)
		err := b.DecodeMessage(msg)
		m.MediationClientMessage = &MediationClientMessage_KeepAlive{msg}
		return true, err
	case 5: // mediation_client_message.actionProgress
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ActionProgress)
		err := b.DecodeMessage(msg)
		m.MediationClientMessage = &MediationClientMessage_ActionProgress{msg}
		return true, err
	case 6: // mediation_client_message.actionResponse
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ActionResult)
		err := b.DecodeMessage(msg)
		m.MediationClientMessage = &MediationClientMessage_ActionResponse{msg}
		return true, err
	default:
		return false, nil
	}
}

func _MediationClientMessage_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*MediationClientMessage)
	// mediation_client_message
	switch x := m.MediationClientMessage.(type) {
	case *MediationClientMessage_ValidationResponse:
		s := proto.Size(x.ValidationResponse)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *MediationClientMessage_DiscoveryResponse:
		s := proto.Size(x.DiscoveryResponse)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *MediationClientMessage_KeepAlive:
		s := proto.Size(x.KeepAlive)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *MediationClientMessage_ActionProgress:
		s := proto.Size(x.ActionProgress)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *MediationClientMessage_ActionResponse:
		s := proto.Size(x.ActionResponse)
		n += proto.SizeVarint(6<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// Messages, sent from server to client
type MediationServerMessage struct {
	// Types that are valid to be assigned to MediationServerMessage:
	//	*MediationServerMessage_ValidationRequest
	//	*MediationServerMessage_DiscoveryRequest
	//	*MediationServerMessage_ActionRequest
	//	*MediationServerMessage_InterruptOperation
	MediationServerMessage isMediationServerMessage_MediationServerMessage `protobuf_oneof:"mediation_server_message"`
	// this is always required in reality. it's optional here because
	// we don't know if in the future, with embedded targets, we will
	// still use it or not
	MessageID        *int32 `protobuf:"varint,15,opt,name=messageID" json:"messageID,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *MediationServerMessage) Reset()                    { *m = MediationServerMessage{} }
func (m *MediationServerMessage) String() string            { return proto.CompactTextString(m) }
func (*MediationServerMessage) ProtoMessage()               {}
func (*MediationServerMessage) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{1} }

type isMediationServerMessage_MediationServerMessage interface {
	isMediationServerMessage_MediationServerMessage()
}

type MediationServerMessage_ValidationRequest struct {
	ValidationRequest *ValidationRequest `protobuf:"bytes,2,opt,name=validationRequest,oneof"`
}
type MediationServerMessage_DiscoveryRequest struct {
	DiscoveryRequest *DiscoveryRequest `protobuf:"bytes,3,opt,name=discoveryRequest,oneof"`
}
type MediationServerMessage_ActionRequest struct {
	ActionRequest *ActionRequest `protobuf:"bytes,4,opt,name=actionRequest,oneof"`
}
type MediationServerMessage_InterruptOperation struct {
	InterruptOperation int32 `protobuf:"varint,5,opt,name=interruptOperation,oneof"`
}

func (*MediationServerMessage_ValidationRequest) isMediationServerMessage_MediationServerMessage()  {}
func (*MediationServerMessage_DiscoveryRequest) isMediationServerMessage_MediationServerMessage()   {}
func (*MediationServerMessage_ActionRequest) isMediationServerMessage_MediationServerMessage()      {}
func (*MediationServerMessage_InterruptOperation) isMediationServerMessage_MediationServerMessage() {}

func (m *MediationServerMessage) GetMediationServerMessage() isMediationServerMessage_MediationServerMessage {
	if m != nil {
		return m.MediationServerMessage
	}
	return nil
}

func (m *MediationServerMessage) GetValidationRequest() *ValidationRequest {
	if x, ok := m.GetMediationServerMessage().(*MediationServerMessage_ValidationRequest); ok {
		return x.ValidationRequest
	}
	return nil
}

func (m *MediationServerMessage) GetDiscoveryRequest() *DiscoveryRequest {
	if x, ok := m.GetMediationServerMessage().(*MediationServerMessage_DiscoveryRequest); ok {
		return x.DiscoveryRequest
	}
	return nil
}

func (m *MediationServerMessage) GetActionRequest() *ActionRequest {
	if x, ok := m.GetMediationServerMessage().(*MediationServerMessage_ActionRequest); ok {
		return x.ActionRequest
	}
	return nil
}

func (m *MediationServerMessage) GetInterruptOperation() int32 {
	if x, ok := m.GetMediationServerMessage().(*MediationServerMessage_InterruptOperation); ok {
		return x.InterruptOperation
	}
	return 0
}

func (m *MediationServerMessage) GetMessageID() int32 {
	if m != nil && m.MessageID != nil {
		return *m.MessageID
	}
	return 0
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*MediationServerMessage) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _MediationServerMessage_OneofMarshaler, _MediationServerMessage_OneofUnmarshaler, _MediationServerMessage_OneofSizer, []interface{}{
		(*MediationServerMessage_ValidationRequest)(nil),
		(*MediationServerMessage_DiscoveryRequest)(nil),
		(*MediationServerMessage_ActionRequest)(nil),
		(*MediationServerMessage_InterruptOperation)(nil),
	}
}

func _MediationServerMessage_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*MediationServerMessage)
	// mediation_server_message
	switch x := m.MediationServerMessage.(type) {
	case *MediationServerMessage_ValidationRequest:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ValidationRequest); err != nil {
			return err
		}
	case *MediationServerMessage_DiscoveryRequest:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.DiscoveryRequest); err != nil {
			return err
		}
	case *MediationServerMessage_ActionRequest:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ActionRequest); err != nil {
			return err
		}
	case *MediationServerMessage_InterruptOperation:
		b.EncodeVarint(5<<3 | proto.WireVarint)
		b.EncodeVarint(uint64(x.InterruptOperation))
	case nil:
	default:
		return fmt.Errorf("MediationServerMessage.MediationServerMessage has unexpected type %T", x)
	}
	return nil
}

func _MediationServerMessage_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*MediationServerMessage)
	switch tag {
	case 2: // mediation_server_message.validationRequest
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ValidationRequest)
		err := b.DecodeMessage(msg)
		m.MediationServerMessage = &MediationServerMessage_ValidationRequest{msg}
		return true, err
	case 3: // mediation_server_message.discoveryRequest
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(DiscoveryRequest)
		err := b.DecodeMessage(msg)
		m.MediationServerMessage = &MediationServerMessage_DiscoveryRequest{msg}
		return true, err
	case 4: // mediation_server_message.actionRequest
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ActionRequest)
		err := b.DecodeMessage(msg)
		m.MediationServerMessage = &MediationServerMessage_ActionRequest{msg}
		return true, err
	case 5: // mediation_server_message.interruptOperation
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.MediationServerMessage = &MediationServerMessage_InterruptOperation{int32(x)}
		return true, err
	default:
		return false, nil
	}
}

func _MediationServerMessage_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*MediationServerMessage)
	// mediation_server_message
	switch x := m.MediationServerMessage.(type) {
	case *MediationServerMessage_ValidationRequest:
		s := proto.Size(x.ValidationRequest)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *MediationServerMessage_DiscoveryRequest:
		s := proto.Size(x.DiscoveryRequest)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *MediationServerMessage_ActionRequest:
		s := proto.Size(x.ActionRequest)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *MediationServerMessage_InterruptOperation:
		n += proto.SizeVarint(5<<3 | proto.WireVarint)
		n += proto.SizeVarint(uint64(x.InterruptOperation))
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// Request for action to be performed in probe
type ActionRequest struct {
	ProbeType *string `protobuf:"bytes,1,req,name=probeType" json:"probeType,omitempty"`
	// Account values provide data to allow the probe to allow it to connect
	// to the probe
	AccountValue []*AccountValue `protobuf:"bytes,2,rep,name=accountValue" json:"accountValue,omitempty"`
	// An action execution DTO contains one or more action items
	ActionExecutionDTO *ActionExecutionDTO `protobuf:"bytes,3,req,name=actionExecutionDTO" json:"actionExecutionDTO,omitempty"`
	// For Cross Destination actions (from one target to another) 2 sets of account
	// values are needed
	SecondaryAccountValue []*AccountValue `protobuf:"bytes,4,rep,name=secondaryAccountValue" json:"secondaryAccountValue,omitempty"`
	XXX_unrecognized      []byte          `json:"-"`
}

func (m *ActionRequest) Reset()                    { *m = ActionRequest{} }
func (m *ActionRequest) String() string            { return proto.CompactTextString(m) }
func (*ActionRequest) ProtoMessage()               {}
func (*ActionRequest) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{2} }

func (m *ActionRequest) GetProbeType() string {
	if m != nil && m.ProbeType != nil {
		return *m.ProbeType
	}
	return ""
}

func (m *ActionRequest) GetAccountValue() []*AccountValue {
	if m != nil {
		return m.AccountValue
	}
	return nil
}

func (m *ActionRequest) GetActionExecutionDTO() *ActionExecutionDTO {
	if m != nil {
		return m.ActionExecutionDTO
	}
	return nil
}

func (m *ActionRequest) GetSecondaryAccountValue() []*AccountValue {
	if m != nil {
		return m.SecondaryAccountValue
	}
	return nil
}

// Result of the action execution. It is translated only once
// after action execution is either completed or failed
type ActionResult struct {
	Response         *ActionResponse `protobuf:"bytes,1,req,name=response" json:"response,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *ActionResult) Reset()                    { *m = ActionResult{} }
func (m *ActionResult) String() string            { return proto.CompactTextString(m) }
func (*ActionResult) ProtoMessage()               {}
func (*ActionResult) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{3} }

func (m *ActionResult) GetResponse() *ActionResponse {
	if m != nil {
		return m.Response
	}
	return nil
}

// Progress of the currently executed action. Can be send multiple times
// for each action
type ActionProgress struct {
	Response         *ActionResponse `protobuf:"bytes,1,req,name=response" json:"response,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *ActionProgress) Reset()                    { *m = ActionProgress{} }
func (m *ActionProgress) String() string            { return proto.CompactTextString(m) }
func (*ActionProgress) ProtoMessage()               {}
func (*ActionProgress) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{4} }

func (m *ActionProgress) GetResponse() *ActionResponse {
	if m != nil {
		return m.Response
	}
	return nil
}

// This class holds response information about executing action. It contains: response: the
// ActionResponseState code representing the state of executing action responseDescription: the
// description message notifying detailed information about current status of executing action
type ActionResponse struct {
	// current action state
	ActionResponseState *ActionResponseState `protobuf:"varint,1,req,name=actionResponseState,enum=common_dto.ActionResponseState" json:"actionResponseState,omitempty"`
	// current action progress (0..100)
	Progress *int32 `protobuf:"varint,2,req,name=progress" json:"progress,omitempty"`
	// action state description, for example ("Moving VM...")
	ResponseDescription *string `protobuf:"bytes,3,req,name=responseDescription" json:"responseDescription,omitempty"`
	XXX_unrecognized    []byte  `json:"-"`
}

func (m *ActionResponse) Reset()                    { *m = ActionResponse{} }
func (m *ActionResponse) String() string            { return proto.CompactTextString(m) }
func (*ActionResponse) ProtoMessage()               {}
func (*ActionResponse) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{5} }

func (m *ActionResponse) GetActionResponseState() ActionResponseState {
	if m != nil && m.ActionResponseState != nil {
		return *m.ActionResponseState
	}
	return ActionResponseState_PENDING_ACCEPT
}

func (m *ActionResponse) GetProgress() int32 {
	if m != nil && m.Progress != nil {
		return *m.Progress
	}
	return 0
}

func (m *ActionResponse) GetResponseDescription() string {
	if m != nil && m.ResponseDescription != nil {
		return *m.ResponseDescription
	}
	return ""
}

// ContainerInfo message to the Operations Manager server.
// This message passes probe descriptions to the server.
type ContainerInfo struct {
	// Set of ProbeInfo objects, each one will carry information about one of the probe
	// that the container has loaded internally.
	Probes           []*ProbeInfo `protobuf:"bytes,1,rep,name=probes" json:"probes,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *ContainerInfo) Reset()                    { *m = ContainerInfo{} }
func (m *ContainerInfo) String() string            { return proto.CompactTextString(m) }
func (*ContainerInfo) ProtoMessage()               {}
func (*ContainerInfo) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{6} }

func (m *ContainerInfo) GetProbes() []*ProbeInfo {
	if m != nil {
		return m.Probes
	}
	return nil
}

type KeepAlive struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *KeepAlive) Reset()                    { *m = KeepAlive{} }
func (m *KeepAlive) String() string            { return proto.CompactTextString(m) }
func (*KeepAlive) ProtoMessage()               {}
func (*KeepAlive) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{7} }

type Ack struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *Ack) Reset()                    { *m = Ack{} }
func (m *Ack) String() string            { return proto.CompactTextString(m) }
func (*Ack) ProtoMessage()               {}
func (*Ack) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{8} }

type ValidationRequest struct {
	ProbeType        *string         `protobuf:"bytes,1,req,name=probeType" json:"probeType,omitempty"`
	AccountValue     []*AccountValue `protobuf:"bytes,2,rep,name=accountValue" json:"accountValue,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *ValidationRequest) Reset()                    { *m = ValidationRequest{} }
func (m *ValidationRequest) String() string            { return proto.CompactTextString(m) }
func (*ValidationRequest) ProtoMessage()               {}
func (*ValidationRequest) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{9} }

func (m *ValidationRequest) GetProbeType() string {
	if m != nil && m.ProbeType != nil {
		return *m.ProbeType
	}
	return ""
}

func (m *ValidationRequest) GetAccountValue() []*AccountValue {
	if m != nil {
		return m.AccountValue
	}
	return nil
}

type DiscoveryRequest struct {
	ProbeType        *string         `protobuf:"bytes,1,req,name=probeType" json:"probeType,omitempty"`
	AccountValue     []*AccountValue `protobuf:"bytes,2,rep,name=accountValue" json:"accountValue,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *DiscoveryRequest) Reset()                    { *m = DiscoveryRequest{} }
func (m *DiscoveryRequest) String() string            { return proto.CompactTextString(m) }
func (*DiscoveryRequest) ProtoMessage()               {}
func (*DiscoveryRequest) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{10} }

func (m *DiscoveryRequest) GetProbeType() string {
	if m != nil && m.ProbeType != nil {
		return *m.ProbeType
	}
	return ""
}

func (m *DiscoveryRequest) GetAccountValue() []*AccountValue {
	if m != nil {
		return m.AccountValue
	}
	return nil
}

// The ProbeInfo class provides a description of the probe that enables users to
// attach Operations Manager to a target, and enables the probe to add entities to the
// Operations Manager market as valid members of the supply chain.
//
// To enable users to use this probe, the ProbeInfo includes a probe type and
// the set of fields a user must give to provide credentials and other data necessary to
// attach to a target. The probe type is an arbitrary string, but REST API calls that
// invoke this probe must refer to it by the same type.
//
// To enable adding entities to the Operations Manager market, the ProbeInfo includes a
// set of {@link TemplateDTO} objects called the supplyChainDefinitionSet.
// Each template object describes an entity type that
// the probe can discover and add to the market. This description includes an EntityDTO object
// and its corresponding lists of bought and sold {@code CommodityDTO} objects. As the probe
// discovers entities, it must create instances that map to members of the supplyChainDefinitionSet.
type ProbeInfo struct {
	// probeType is a string identifier to define the type of the probe. You specify the value for this string in the
	// {@code probe-conf.xml} file for your probe. Note that for a given instance of Operations Manager,
	// every probe communicating with the server must have a unique type.
	//
	// The probe you create must include a type
	// to display in the user interface. The string you provide for the probe type appears
	// in the Target Configuration form as one of the choices for the category that you set for the probe.
	//
	// For example, in the standard targets, Hypervisor is a category. If your probe
	// category is also {@code Hypervisor} and you specify a type of 'MyProbe', then MyProbe
	// will appear in the user interface as an additional Hypervisor choice.
	// On the other hand, if the category you provide does not match one of the standard categories,
	// MyProbe will appear as a choice in the CUSTOM category.
	ProbeType *string `protobuf:"bytes,1,req,name=probeType" json:"probeType,omitempty"`
	// String identifier to define the category of the probe. You specify the value for this string in the
	// 'probe-conf.xml' file for your probe.
	//
	// The probe you create must include a category.
	// If the category you provide matches one of the standard categories, then your probe will appear
	// as a choice in the Target Configuration form alongside the other members of that category.
	// For example, in the standard targets, 'Hypervisor' is a category. If your probe
	// category is also 'Hypervisor' and you specify a type of 'MyProbe', then MyProbe
	// will appear in the user interface as an additional 'Hypervisor' choice.
	// On the other hand, if the category you provide does not match one of the standard categories,
	// MyProbe will appear as a choice in the 'CUSTOM' category.
	//
	// The set of standard categories is defined in the 'ProbeCategory' enumeration.
	ProbeCategory *string `protobuf:"bytes,2,req,name=probeCategory" json:"probeCategory,omitempty"`
	// Set of TemplateDTO objects that defines the types of entities the probe discovers, and
	// what their bought and sold commodities are. Any entity instances the probe creates must match
	// members of this set.
	SupplyChainDefinitionSet []*TemplateDTO `protobuf:"bytes,3,rep,name=supplyChainDefinitionSet" json:"supplyChainDefinitionSet,omitempty"`
	// List of AccountDefEntry objects that describe the fields users provide as
	// input (i.e. ip, user, pass, ...). These fields appear in the Operations Manager user interface
	// when users add targets of this probe's type. REST API calls to add targets also provide data
	// for these fields (i.e. ip, user, password, ...).
	//
	// Order of elements in the list specifyes the order they appear in the UI.
	// List must not contain entries with equal "name" field. This is up to client to ensure this.
	AccountDefinition []*AccountDefEntry `protobuf:"bytes,4,rep,name=accountDefinition" json:"accountDefinition,omitempty"`
	// The field name, that should be treated as target identifier
	TargetIdentifierField *string `protobuf:"bytes,5,req,name=targetIdentifierField" json:"targetIdentifierField,omitempty"`
	// Specifies the interval at which discoveries will be executed for this probe.
	// The value is specified in seconds. If no value is provided for rediscoveryIntervalSeconds
	// a default of 600 seconds (10 minutes) will be used. The minimum value allowed for this
	// field is 60 seconds (1 minute).
	RediscoveryIntervalSeconds *int32 `protobuf:"varint,6,opt,name=rediscoveryIntervalSeconds" json:"rediscoveryIntervalSeconds,omitempty"`
	FullRediscoveryIntervalSeconds *int32 `protobuf:"varint,6,opt,name=fullRediscoveryIntervalSeconds" json:"fullRediscoveryIntervalSeconds,omitempty"`
	// EntityIdentityMetadata supplies meta information describing the properties used to identify
	// all of the entities that a probe may discover. There should be one EntityIdentityMetadata
	// for each type of entity that a probe may discover.
	EntityMetadata []*EntityIdentityMetadata `protobuf:"bytes,7,rep,name=entityMetadata" json:"entityMetadata,omitempty"`
	// Action policy provides data about entity types and actions that can be applied to them
	ActionPolicy     []*ActionPolicyDTO `protobuf:"bytes,8,rep,name=actionPolicy" json:"actionPolicy,omitempty"`
	XXX_unrecognized []byte             `json:"-"`
}

func (m *ProbeInfo) Reset()                    { *m = ProbeInfo{} }
func (m *ProbeInfo) String() string            { return proto.CompactTextString(m) }
func (*ProbeInfo) ProtoMessage()               {}
func (*ProbeInfo) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{11} }

func (m *ProbeInfo) GetProbeType() string {
	if m != nil && m.ProbeType != nil {
		return *m.ProbeType
	}
	return ""
}

func (m *ProbeInfo) GetProbeCategory() string {
	if m != nil && m.ProbeCategory != nil {
		return *m.ProbeCategory
	}
	return ""
}

func (m *ProbeInfo) GetSupplyChainDefinitionSet() []*TemplateDTO {
	if m != nil {
		return m.SupplyChainDefinitionSet
	}
	return nil
}

func (m *ProbeInfo) GetAccountDefinition() []*AccountDefEntry {
	if m != nil {
		return m.AccountDefinition
	}
	return nil
}

func (m *ProbeInfo) GetTargetIdentifierField() string {
	if m != nil && m.TargetIdentifierField != nil {
		return *m.TargetIdentifierField
	}
	return ""
}

func (m *ProbeInfo) GetRediscoveryIntervalSeconds() int32 {
	if m != nil && m.RediscoveryIntervalSeconds != nil {
		return *m.RediscoveryIntervalSeconds
	}
	return 0
}

func (m *ProbeInfo) GetEntityMetadata() []*EntityIdentityMetadata {
	if m != nil {
		return m.EntityMetadata
	}
	return nil
}

func (m *ProbeInfo) GetActionPolicy() []*ActionPolicyDTO {
	if m != nil {
		return m.ActionPolicy
	}
	return nil
}

func init() {
	proto.RegisterType((*MediationClientMessage)(nil), "common_dto.MediationClientMessage")
	proto.RegisterType((*MediationServerMessage)(nil), "common_dto.MediationServerMessage")
	proto.RegisterType((*ActionRequest)(nil), "common_dto.ActionRequest")
	proto.RegisterType((*ActionResult)(nil), "common_dto.ActionResult")
	proto.RegisterType((*ActionProgress)(nil), "common_dto.ActionProgress")
	proto.RegisterType((*ActionResponse)(nil), "common_dto.ActionResponse")
	proto.RegisterType((*ContainerInfo)(nil), "common_dto.ContainerInfo")
	proto.RegisterType((*KeepAlive)(nil), "common_dto.KeepAlive")
	proto.RegisterType((*Ack)(nil), "common_dto.Ack")
	proto.RegisterType((*ValidationRequest)(nil), "common_dto.ValidationRequest")
	proto.RegisterType((*DiscoveryRequest)(nil), "common_dto.DiscoveryRequest")
	proto.RegisterType((*ProbeInfo)(nil), "common_dto.ProbeInfo")
}

func init() { proto.RegisterFile("MediationMessage.proto", fileDescriptor4) }

var fileDescriptor4 = []byte{
	// 831 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xb4, 0x56, 0x5b, 0x6e, 0xdb, 0x46,
	0x14, 0xb5, 0xa4, 0x28, 0xb5, 0x6e, 0x62, 0xc7, 0x9a, 0xc0, 0x29, 0xab, 0x26, 0x8d, 0x40, 0xf4,
	0xc3, 0x3f, 0x15, 0x82, 0xf4, 0xf1, 0x55, 0xa4, 0x90, 0x4d, 0x07, 0x52, 0x0a, 0x25, 0xee, 0xc8,
	0xc8, 0xaf, 0x31, 0x21, 0xaf, 0xdc, 0x81, 0x49, 0x0e, 0x3b, 0x33, 0x14, 0xca, 0x35, 0x74, 0x35,
	0x05, 0xba, 0xa2, 0x2e, 0xa1, 0x2b, 0x28, 0x38, 0x7c, 0x8b, 0x74, 0x83, 0x06, 0xf0, 0x9f, 0x34,
	0xf7, 0x9c, 0x33, 0x97, 0xe7, 0x3e, 0x48, 0x78, 0xb2, 0x42, 0x8f, 0x33, 0xcd, 0x45, 0xb8, 0x42,
	0xa5, 0xd8, 0x35, 0xce, 0x22, 0x29, 0xb4, 0x20, 0xe0, 0x8a, 0x20, 0x10, 0xe1, 0x95, 0xa7, 0xc5,
	0xe4, 0x78, 0xee, 0xa6, 0x80, 0xf3, 0xdf, 0xd1, 0x8d, 0xd3, 0x1f, 0x19, 0x64, 0xf2, 0xc8, 0xe1,
	0xca, 0x15, 0x5b, 0x94, 0x49, 0x7e, 0x30, 0x5e, 0xc7, 0x51, 0xe4, 0x27, 0x67, 0xbf, 0x32, 0x5e,
	0x60, 0x9e, 0x2c, 0x3d, 0x0c, 0x35, 0xd7, 0xc9, 0x0a, 0x35, 0xf3, 0x98, 0x66, 0xd9, 0xb9, 0xfd,
	0xd7, 0xa0, 0x76, 0xf3, 0x99, 0xcf, 0x31, 0xd4, 0xf9, 0xfd, 0xe4, 0x02, 0xc8, 0x96, 0xf9, 0xdc,
	0x33, 0x21, 0x8a, 0x2a, 0x12, 0xa1, 0x42, 0xab, 0x3f, 0xed, 0x9d, 0x3c, 0x78, 0xf9, 0xd5, 0xac,
	0x4a, 0x6b, 0xf6, 0xbe, 0x85, 0x5a, 0xec, 0xd1, 0x0e, 0x2e, 0x59, 0xc1, 0xd8, 0x2b, 0x52, 0x2d,
	0x05, 0x07, 0x46, 0xf0, 0x59, 0x5d, 0xd0, 0xd9, 0x05, 0x2d, 0xf6, 0x68, 0x9b, 0x49, 0xbe, 0x87,
	0xd1, 0x0d, 0x62, 0x34, 0xf7, 0xf9, 0x16, 0xad, 0x7b, 0x46, 0xe6, 0xb8, 0x2e, 0xf3, 0x73, 0x11,
	0x5c, 0xec, 0xd1, 0x0a, 0x49, 0x1c, 0x38, 0x64, 0xc6, 0xc7, 0x0b, 0x29, 0xae, 0x25, 0x2a, 0x65,
	0x0d, 0x0d, 0x77, 0x52, 0xe7, 0xce, 0x1b, 0x88, 0xc5, 0x1e, 0xdd, 0xe1, 0x90, 0xd3, 0x42, 0xa5,
	0x7c, 0x90, 0xfb, 0x46, 0xc5, 0x6a, 0xab, 0x50, 0x54, 0xb1, 0xaf, 0x2b, 0x8d, 0xf2, 0x01, 0x9e,
	0xc2, 0x28, 0xc8, 0xcc, 0x5e, 0x3a, 0xd6, 0xa3, 0x69, 0xef, 0x64, 0x48, 0xab, 0x83, 0xd3, 0x09,
	0x58, 0x41, 0x51, 0x99, 0x2b, 0xd7, 0x94, 0xe6, 0x2a, 0x8f, 0xda, 0x7f, 0xf7, 0x6b, 0x65, 0x5b,
	0xa3, 0xdc, 0xa2, 0x2c, 0xca, 0xb6, 0x82, 0x71, 0xdd, 0xfa, 0xdf, 0x62, 0x54, 0x3a, 0xaf, 0xda,
	0xb3, 0xdb, 0xaa, 0x66, 0x40, 0xa9, 0xc9, 0x2d, 0x26, 0x79, 0x03, 0x47, 0x35, 0xe7, 0x33, 0xb5,
	0xac, 0x64, 0x4f, 0x6f, 0x29, 0x59, 0x21, 0xd6, 0xe2, 0x91, 0x39, 0x1c, 0x14, 0x0e, 0x64, 0x42,
	0x59, 0xd1, 0xbe, 0xe8, 0xb2, 0xac, 0x50, 0x69, 0x32, 0xc8, 0x0b, 0x20, 0x3c, 0xd4, 0x28, 0x65,
	0x1c, 0xe9, 0x77, 0x11, 0x4a, 0x93, 0xab, 0x29, 0xe0, 0x30, 0x6d, 0xba, 0x76, 0xec, 0xff, 0x98,
	0xac, 0x8c, 0x91, 0xa5, 0xc9, 0x7f, 0xf4, 0xe1, 0xa0, 0x91, 0x4e, 0xaa, 0x15, 0x49, 0xf1, 0x01,
	0x2f, 0x93, 0x08, 0xad, 0xde, 0xb4, 0x7f, 0x32, 0xa2, 0xd5, 0x01, 0xf9, 0x11, 0x1e, 0x32, 0xd7,
	0x15, 0x71, 0xa8, 0xdf, 0x33, 0x3f, 0x4e, 0x47, 0x65, 0xd0, 0x6e, 0x88, 0x2a, 0x4e, 0x1b, 0x68,
	0xf2, 0x16, 0x08, 0x6b, 0x8e, 0xb7, 0x73, 0xf9, 0xce, 0x1a, 0x4c, 0xfb, 0xbb, 0xe3, 0x36, 0x6f,
	0xa1, 0x68, 0x07, 0x93, 0xbc, 0x85, 0x63, 0x85, 0xae, 0x08, 0x3d, 0x26, 0x93, 0xfa, 0xb5, 0xd6,
	0xbd, 0x8f, 0xa4, 0xd5, 0x4d, 0xb3, 0x5f, 0xc3, 0xc3, 0x7a, 0x3b, 0x93, 0x1f, 0x60, 0x5f, 0x16,
	0xad, 0xdf, 0x33, 0x59, 0x4e, 0x3a, 0x5b, 0xdf, 0x20, 0x68, 0x89, 0xb5, 0x17, 0x70, 0xd8, 0x1c,
	0xae, 0x4f, 0x56, 0xfa, 0xb3, 0x57, 0x48, 0x95, 0x13, 0xf5, 0x0b, 0x3c, 0x6e, 0xce, 0xd8, 0x5a,
	0x33, 0x9d, 0xa9, 0x1e, 0xbe, 0x7c, 0x7e, 0xbb, 0xaa, 0x81, 0xd1, 0x2e, 0x2e, 0x99, 0xc0, 0x7e,
	0x54, 0x2c, 0x8a, 0xfe, 0xb4, 0x7f, 0x32, 0xa4, 0xe5, 0x7f, 0xf2, 0x02, 0x1e, 0x17, 0xd9, 0x38,
	0xa8, 0x5c, 0xc9, 0x23, 0xd3, 0x8e, 0x03, 0xd3, 0x19, 0x5d, 0x21, 0xfb, 0x15, 0x1c, 0x9c, 0x89,
	0x50, 0x33, 0x1e, 0xa2, 0x5c, 0x86, 0x1b, 0x41, 0xbe, 0x81, 0xfb, 0xa6, 0x83, 0x94, 0xd5, 0x33,
	0x75, 0x69, 0x6c, 0xb0, 0x8b, 0x34, 0x92, 0xc2, 0x68, 0x0e, 0xb2, 0x1f, 0xc0, 0xa8, 0x5c, 0x6b,
	0xf6, 0x10, 0x06, 0x73, 0xf7, 0xc6, 0x16, 0x30, 0x6e, 0x0d, 0xf3, 0x5d, 0xb6, 0xaa, 0x1d, 0xc2,
	0xd1, 0xee, 0xbc, 0xdf, 0xe9, 0x7d, 0xff, 0x0c, 0x60, 0x54, 0x5a, 0xf1, 0x91, 0x9b, 0xbe, 0x86,
	0x03, 0xf3, 0xe7, 0x8c, 0x69, 0xbc, 0x16, 0x32, 0x31, 0x35, 0x1b, 0xd1, 0xe6, 0x21, 0x59, 0x83,
	0xa5, 0xaa, 0x77, 0xa4, 0x83, 0x1b, 0x1e, 0xf2, 0x6c, 0x95, 0xa6, 0xdb, 0x2d, 0xcd, 0xed, 0xf3,
	0x7a, 0x6e, 0x97, 0x18, 0x44, 0x3e, 0xd3, 0x98, 0xce, 0xda, 0xad, 0x44, 0xb2, 0x84, 0x71, 0x9e,
	0x76, 0x75, 0x9e, 0x4f, 0xdb, 0x97, 0x1d, 0x4f, 0xea, 0xe0, 0xe6, 0x3c, 0xd4, 0x32, 0xa1, 0x6d,
	0x16, 0xf9, 0x0e, 0x8e, 0x35, 0x93, 0xd7, 0xa8, 0xb3, 0xd7, 0xf6, 0x86, 0xa3, 0x7c, 0xcd, 0xd1,
	0xf7, 0xac, 0xa1, 0x79, 0x9a, 0xee, 0x20, 0x79, 0x05, 0x13, 0x89, 0xe5, 0xd6, 0x5d, 0xa6, 0xbb,
	0x70, 0xcb, 0xfc, 0xb5, 0x19, 0x67, 0x65, 0xde, 0x4f, 0x43, 0xfa, 0x1f, 0x08, 0xf2, 0x06, 0x0e,
	0x9b, 0x1f, 0x09, 0xd6, 0x67, 0x26, 0x7b, 0xbb, 0x9e, 0xfd, 0xb9, 0x41, 0xec, 0x7e, 0x4e, 0xd0,
	0x1d, 0x26, 0xf9, 0x29, 0xad, 0xb8, 0x19, 0x73, 0xe1, 0x73, 0x37, 0xb1, 0xf6, 0xbb, 0x7c, 0xa8,
	0xe2, 0xa9, 0xb3, 0x0d, 0xc2, 0xe9, 0xb7, 0xf0, 0xdc, 0x15, 0xc1, 0x6c, 0x1b, 0xe8, 0x58, 0x7e,
	0x10, 0xb3, 0xd4, 0xff, 0x8d, 0x90, 0xc1, 0x4c, 0x79, 0x37, 0xb9, 0xc8, 0xe9, 0xd1, 0xee, 0x37,
	0xd3, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x04, 0x0d, 0x96, 0x2f, 0x46, 0x09, 0x00, 0x00,
}
