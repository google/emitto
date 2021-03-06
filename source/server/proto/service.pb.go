// Code generated by protoc-gen-go. DO NOT EDIT.
// source: source/server/proto/service.proto

package emitto_service

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	status "google.golang.org/genproto/googleapis/rpc/status"
	field_mask "google.golang.org/genproto/protobuf/field_mask"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Location struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Zones                []string `protobuf:"bytes,2,rep,name=zones,proto3" json:"zones,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Location) Reset()         { *m = Location{} }
func (m *Location) String() string { return proto.CompactTextString(m) }
func (*Location) ProtoMessage()    {}
func (*Location) Descriptor() ([]byte, []int) {
	return fileDescriptor_6256c75c7d842e44, []int{0}
}

func (m *Location) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Location.Unmarshal(m, b)
}
func (m *Location) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Location.Marshal(b, m, deterministic)
}
func (m *Location) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Location.Merge(m, src)
}
func (m *Location) XXX_Size() int {
	return xxx_messageInfo_Location.Size(m)
}
func (m *Location) XXX_DiscardUnknown() {
	xxx_messageInfo_Location.DiscardUnknown(m)
}

var xxx_messageInfo_Location proto.InternalMessageInfo

func (m *Location) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Location) GetZones() []string {
	if m != nil {
		return m.Zones
	}
	return nil
}

type Rule struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Body                 string   `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	LocationZones        []string `protobuf:"bytes,3,rep,name=location_zones,json=locationZones,proto3" json:"location_zones,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Rule) Reset()         { *m = Rule{} }
func (m *Rule) String() string { return proto.CompactTextString(m) }
func (*Rule) ProtoMessage()    {}
func (*Rule) Descriptor() ([]byte, []int) {
	return fileDescriptor_6256c75c7d842e44, []int{1}
}

func (m *Rule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Rule.Unmarshal(m, b)
}
func (m *Rule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Rule.Marshal(b, m, deterministic)
}
func (m *Rule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Rule.Merge(m, src)
}
func (m *Rule) XXX_Size() int {
	return xxx_messageInfo_Rule.Size(m)
}
func (m *Rule) XXX_DiscardUnknown() {
	xxx_messageInfo_Rule.DiscardUnknown(m)
}

var xxx_messageInfo_Rule proto.InternalMessageInfo

func (m *Rule) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Rule) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

func (m *Rule) GetLocationZones() []string {
	if m != nil {
		return m.LocationZones
	}
	return nil
}

type DeployRulesRequest struct {
	Location             *Location `protobuf:"bytes,1,opt,name=location,proto3" json:"location,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *DeployRulesRequest) Reset()         { *m = DeployRulesRequest{} }
func (m *DeployRulesRequest) String() string { return proto.CompactTextString(m) }
func (*DeployRulesRequest) ProtoMessage()    {}
func (*DeployRulesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6256c75c7d842e44, []int{2}
}

func (m *DeployRulesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeployRulesRequest.Unmarshal(m, b)
}
func (m *DeployRulesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeployRulesRequest.Marshal(b, m, deterministic)
}
func (m *DeployRulesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeployRulesRequest.Merge(m, src)
}
func (m *DeployRulesRequest) XXX_Size() int {
	return xxx_messageInfo_DeployRulesRequest.Size(m)
}
func (m *DeployRulesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeployRulesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeployRulesRequest proto.InternalMessageInfo

func (m *DeployRulesRequest) GetLocation() *Location {
	if m != nil {
		return m.Location
	}
	return nil
}

type DeployRulesResponse struct {
	ClientId             string         `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	Status               *status.Status `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *DeployRulesResponse) Reset()         { *m = DeployRulesResponse{} }
func (m *DeployRulesResponse) String() string { return proto.CompactTextString(m) }
func (*DeployRulesResponse) ProtoMessage()    {}
func (*DeployRulesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6256c75c7d842e44, []int{3}
}

func (m *DeployRulesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeployRulesResponse.Unmarshal(m, b)
}
func (m *DeployRulesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeployRulesResponse.Marshal(b, m, deterministic)
}
func (m *DeployRulesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeployRulesResponse.Merge(m, src)
}
func (m *DeployRulesResponse) XXX_Size() int {
	return xxx_messageInfo_DeployRulesResponse.Size(m)
}
func (m *DeployRulesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeployRulesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeployRulesResponse proto.InternalMessageInfo

func (m *DeployRulesResponse) GetClientId() string {
	if m != nil {
		return m.ClientId
	}
	return ""
}

func (m *DeployRulesResponse) GetStatus() *status.Status {
	if m != nil {
		return m.Status
	}
	return nil
}

type AddRuleRequest struct {
	Rule                 *Rule    `protobuf:"bytes,1,opt,name=rule,proto3" json:"rule,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddRuleRequest) Reset()         { *m = AddRuleRequest{} }
func (m *AddRuleRequest) String() string { return proto.CompactTextString(m) }
func (*AddRuleRequest) ProtoMessage()    {}
func (*AddRuleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6256c75c7d842e44, []int{4}
}

func (m *AddRuleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddRuleRequest.Unmarshal(m, b)
}
func (m *AddRuleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddRuleRequest.Marshal(b, m, deterministic)
}
func (m *AddRuleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddRuleRequest.Merge(m, src)
}
func (m *AddRuleRequest) XXX_Size() int {
	return xxx_messageInfo_AddRuleRequest.Size(m)
}
func (m *AddRuleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddRuleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddRuleRequest proto.InternalMessageInfo

func (m *AddRuleRequest) GetRule() *Rule {
	if m != nil {
		return m.Rule
	}
	return nil
}

type ModifyRuleRequest struct {
	Rule                 *Rule                 `protobuf:"bytes,1,opt,name=rule,proto3" json:"rule,omitempty"`
	FieldMask            *field_mask.FieldMask `protobuf:"bytes,2,opt,name=field_mask,json=fieldMask,proto3" json:"field_mask,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *ModifyRuleRequest) Reset()         { *m = ModifyRuleRequest{} }
func (m *ModifyRuleRequest) String() string { return proto.CompactTextString(m) }
func (*ModifyRuleRequest) ProtoMessage()    {}
func (*ModifyRuleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6256c75c7d842e44, []int{5}
}

func (m *ModifyRuleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ModifyRuleRequest.Unmarshal(m, b)
}
func (m *ModifyRuleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ModifyRuleRequest.Marshal(b, m, deterministic)
}
func (m *ModifyRuleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ModifyRuleRequest.Merge(m, src)
}
func (m *ModifyRuleRequest) XXX_Size() int {
	return xxx_messageInfo_ModifyRuleRequest.Size(m)
}
func (m *ModifyRuleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ModifyRuleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ModifyRuleRequest proto.InternalMessageInfo

func (m *ModifyRuleRequest) GetRule() *Rule {
	if m != nil {
		return m.Rule
	}
	return nil
}

func (m *ModifyRuleRequest) GetFieldMask() *field_mask.FieldMask {
	if m != nil {
		return m.FieldMask
	}
	return nil
}

type DeleteRuleRequest struct {
	RuleId               int64    `protobuf:"varint,1,opt,name=rule_id,json=ruleId,proto3" json:"rule_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteRuleRequest) Reset()         { *m = DeleteRuleRequest{} }
func (m *DeleteRuleRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteRuleRequest) ProtoMessage()    {}
func (*DeleteRuleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6256c75c7d842e44, []int{6}
}

func (m *DeleteRuleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteRuleRequest.Unmarshal(m, b)
}
func (m *DeleteRuleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteRuleRequest.Marshal(b, m, deterministic)
}
func (m *DeleteRuleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteRuleRequest.Merge(m, src)
}
func (m *DeleteRuleRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteRuleRequest.Size(m)
}
func (m *DeleteRuleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteRuleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteRuleRequest proto.InternalMessageInfo

func (m *DeleteRuleRequest) GetRuleId() int64 {
	if m != nil {
		return m.RuleId
	}
	return 0
}

type ListRulesRequest struct {
	RuleIds              []int64  `protobuf:"varint,1,rep,packed,name=rule_ids,json=ruleIds,proto3" json:"rule_ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListRulesRequest) Reset()         { *m = ListRulesRequest{} }
func (m *ListRulesRequest) String() string { return proto.CompactTextString(m) }
func (*ListRulesRequest) ProtoMessage()    {}
func (*ListRulesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6256c75c7d842e44, []int{7}
}

func (m *ListRulesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRulesRequest.Unmarshal(m, b)
}
func (m *ListRulesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRulesRequest.Marshal(b, m, deterministic)
}
func (m *ListRulesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRulesRequest.Merge(m, src)
}
func (m *ListRulesRequest) XXX_Size() int {
	return xxx_messageInfo_ListRulesRequest.Size(m)
}
func (m *ListRulesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRulesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListRulesRequest proto.InternalMessageInfo

func (m *ListRulesRequest) GetRuleIds() []int64 {
	if m != nil {
		return m.RuleIds
	}
	return nil
}

type ListRulesResponse struct {
	Rules                []*Rule  `protobuf:"bytes,1,rep,name=rules,proto3" json:"rules,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListRulesResponse) Reset()         { *m = ListRulesResponse{} }
func (m *ListRulesResponse) String() string { return proto.CompactTextString(m) }
func (*ListRulesResponse) ProtoMessage()    {}
func (*ListRulesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6256c75c7d842e44, []int{8}
}

func (m *ListRulesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRulesResponse.Unmarshal(m, b)
}
func (m *ListRulesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRulesResponse.Marshal(b, m, deterministic)
}
func (m *ListRulesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRulesResponse.Merge(m, src)
}
func (m *ListRulesResponse) XXX_Size() int {
	return xxx_messageInfo_ListRulesResponse.Size(m)
}
func (m *ListRulesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRulesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListRulesResponse proto.InternalMessageInfo

func (m *ListRulesResponse) GetRules() []*Rule {
	if m != nil {
		return m.Rules
	}
	return nil
}

type AddLocationRequest struct {
	Location             *Location `protobuf:"bytes,1,opt,name=location,proto3" json:"location,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *AddLocationRequest) Reset()         { *m = AddLocationRequest{} }
func (m *AddLocationRequest) String() string { return proto.CompactTextString(m) }
func (*AddLocationRequest) ProtoMessage()    {}
func (*AddLocationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6256c75c7d842e44, []int{9}
}

func (m *AddLocationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddLocationRequest.Unmarshal(m, b)
}
func (m *AddLocationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddLocationRequest.Marshal(b, m, deterministic)
}
func (m *AddLocationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddLocationRequest.Merge(m, src)
}
func (m *AddLocationRequest) XXX_Size() int {
	return xxx_messageInfo_AddLocationRequest.Size(m)
}
func (m *AddLocationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddLocationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddLocationRequest proto.InternalMessageInfo

func (m *AddLocationRequest) GetLocation() *Location {
	if m != nil {
		return m.Location
	}
	return nil
}

type ModifyLocationRequest struct {
	Location             *Location             `protobuf:"bytes,1,opt,name=location,proto3" json:"location,omitempty"`
	FieldMask            *field_mask.FieldMask `protobuf:"bytes,2,opt,name=field_mask,json=fieldMask,proto3" json:"field_mask,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *ModifyLocationRequest) Reset()         { *m = ModifyLocationRequest{} }
func (m *ModifyLocationRequest) String() string { return proto.CompactTextString(m) }
func (*ModifyLocationRequest) ProtoMessage()    {}
func (*ModifyLocationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6256c75c7d842e44, []int{10}
}

func (m *ModifyLocationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ModifyLocationRequest.Unmarshal(m, b)
}
func (m *ModifyLocationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ModifyLocationRequest.Marshal(b, m, deterministic)
}
func (m *ModifyLocationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ModifyLocationRequest.Merge(m, src)
}
func (m *ModifyLocationRequest) XXX_Size() int {
	return xxx_messageInfo_ModifyLocationRequest.Size(m)
}
func (m *ModifyLocationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ModifyLocationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ModifyLocationRequest proto.InternalMessageInfo

func (m *ModifyLocationRequest) GetLocation() *Location {
	if m != nil {
		return m.Location
	}
	return nil
}

func (m *ModifyLocationRequest) GetFieldMask() *field_mask.FieldMask {
	if m != nil {
		return m.FieldMask
	}
	return nil
}

type DeleteLocationRequest struct {
	LocationName         string   `protobuf:"bytes,1,opt,name=location_name,json=locationName,proto3" json:"location_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteLocationRequest) Reset()         { *m = DeleteLocationRequest{} }
func (m *DeleteLocationRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteLocationRequest) ProtoMessage()    {}
func (*DeleteLocationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6256c75c7d842e44, []int{11}
}

func (m *DeleteLocationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteLocationRequest.Unmarshal(m, b)
}
func (m *DeleteLocationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteLocationRequest.Marshal(b, m, deterministic)
}
func (m *DeleteLocationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteLocationRequest.Merge(m, src)
}
func (m *DeleteLocationRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteLocationRequest.Size(m)
}
func (m *DeleteLocationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteLocationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteLocationRequest proto.InternalMessageInfo

func (m *DeleteLocationRequest) GetLocationName() string {
	if m != nil {
		return m.LocationName
	}
	return ""
}

type ListLocationsRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListLocationsRequest) Reset()         { *m = ListLocationsRequest{} }
func (m *ListLocationsRequest) String() string { return proto.CompactTextString(m) }
func (*ListLocationsRequest) ProtoMessage()    {}
func (*ListLocationsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6256c75c7d842e44, []int{12}
}

func (m *ListLocationsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListLocationsRequest.Unmarshal(m, b)
}
func (m *ListLocationsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListLocationsRequest.Marshal(b, m, deterministic)
}
func (m *ListLocationsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListLocationsRequest.Merge(m, src)
}
func (m *ListLocationsRequest) XXX_Size() int {
	return xxx_messageInfo_ListLocationsRequest.Size(m)
}
func (m *ListLocationsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListLocationsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListLocationsRequest proto.InternalMessageInfo

type ListLocationsResponse struct {
	Locations            []*Location `protobuf:"bytes,1,rep,name=locations,proto3" json:"locations,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *ListLocationsResponse) Reset()         { *m = ListLocationsResponse{} }
func (m *ListLocationsResponse) String() string { return proto.CompactTextString(m) }
func (*ListLocationsResponse) ProtoMessage()    {}
func (*ListLocationsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6256c75c7d842e44, []int{13}
}

func (m *ListLocationsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListLocationsResponse.Unmarshal(m, b)
}
func (m *ListLocationsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListLocationsResponse.Marshal(b, m, deterministic)
}
func (m *ListLocationsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListLocationsResponse.Merge(m, src)
}
func (m *ListLocationsResponse) XXX_Size() int {
	return xxx_messageInfo_ListLocationsResponse.Size(m)
}
func (m *ListLocationsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListLocationsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListLocationsResponse proto.InternalMessageInfo

func (m *ListLocationsResponse) GetLocations() []*Location {
	if m != nil {
		return m.Locations
	}
	return nil
}

func init() {
	proto.RegisterType((*Location)(nil), "emitto.service.Location")
	proto.RegisterType((*Rule)(nil), "emitto.service.Rule")
	proto.RegisterType((*DeployRulesRequest)(nil), "emitto.service.DeployRulesRequest")
	proto.RegisterType((*DeployRulesResponse)(nil), "emitto.service.DeployRulesResponse")
	proto.RegisterType((*AddRuleRequest)(nil), "emitto.service.AddRuleRequest")
	proto.RegisterType((*ModifyRuleRequest)(nil), "emitto.service.ModifyRuleRequest")
	proto.RegisterType((*DeleteRuleRequest)(nil), "emitto.service.DeleteRuleRequest")
	proto.RegisterType((*ListRulesRequest)(nil), "emitto.service.ListRulesRequest")
	proto.RegisterType((*ListRulesResponse)(nil), "emitto.service.ListRulesResponse")
	proto.RegisterType((*AddLocationRequest)(nil), "emitto.service.AddLocationRequest")
	proto.RegisterType((*ModifyLocationRequest)(nil), "emitto.service.ModifyLocationRequest")
	proto.RegisterType((*DeleteLocationRequest)(nil), "emitto.service.DeleteLocationRequest")
	proto.RegisterType((*ListLocationsRequest)(nil), "emitto.service.ListLocationsRequest")
	proto.RegisterType((*ListLocationsResponse)(nil), "emitto.service.ListLocationsResponse")
}

func init() { proto.RegisterFile("source/server/proto/service.proto", fileDescriptor_6256c75c7d842e44) }

var fileDescriptor_6256c75c7d842e44 = []byte{
	// 641 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0xdf, 0x6e, 0xd3, 0x3c,
	0x18, 0xc6, 0xd7, 0x66, 0x6b, 0x9b, 0xb7, 0xdf, 0xaa, 0xaf, 0xa6, 0xdd, 0x42, 0x26, 0xa1, 0xd6,
	0xa3, 0x52, 0x35, 0x41, 0x8a, 0xca, 0x84, 0x04, 0x42, 0x42, 0x13, 0x1b, 0x52, 0x61, 0x03, 0x2d,
	0x9c, 0xed, 0x60, 0x55, 0x9a, 0xb8, 0x53, 0xb4, 0xb4, 0x0e, 0x71, 0x82, 0x28, 0x57, 0xc0, 0x95,
	0x72, 0x1d, 0x28, 0x76, 0xdc, 0x36, 0x49, 0x17, 0x04, 0xe3, 0xcc, 0x7f, 0x7e, 0x7e, 0xfc, 0xe6,
	0x79, 0xfd, 0x04, 0xba, 0x8c, 0x46, 0x81, 0x4d, 0x06, 0x8c, 0x04, 0x5f, 0x49, 0x30, 0xf0, 0x03,
	0x1a, 0x52, 0x3e, 0x71, 0x6d, 0x62, 0xf0, 0x19, 0x6a, 0x90, 0x99, 0x1b, 0x86, 0xd4, 0x48, 0x56,
	0xf5, 0x83, 0x1b, 0x4a, 0x6f, 0x3c, 0x22, 0xd8, 0x49, 0x34, 0x1d, 0x90, 0x99, 0x1f, 0x2e, 0x04,
	0xac, 0xef, 0x27, 0x9b, 0x81, 0x6f, 0x0f, 0x58, 0x68, 0x85, 0x11, 0x4b, 0x36, 0x3a, 0xd9, 0x53,
	0x53, 0x97, 0x78, 0xce, 0x78, 0x66, 0xb1, 0x5b, 0x41, 0xe0, 0x63, 0xa8, 0x9d, 0x53, 0xdb, 0x0a,
	0x5d, 0x3a, 0x47, 0x08, 0xb6, 0xe7, 0xd6, 0x8c, 0x68, 0xa5, 0x4e, 0xa9, 0xaf, 0x9a, 0x7c, 0x8c,
	0x5a, 0xb0, 0xf3, 0x9d, 0xce, 0x09, 0xd3, 0xca, 0x1d, 0xa5, 0xaf, 0x9a, 0x62, 0x82, 0x2f, 0x61,
	0xdb, 0x8c, 0x3c, 0x82, 0x1a, 0x50, 0x76, 0x1d, 0xce, 0x2b, 0x66, 0xd9, 0x75, 0x62, 0x85, 0x09,
	0x75, 0x16, 0x5a, 0x59, 0x28, 0xc4, 0x63, 0xd4, 0x83, 0x86, 0x97, 0xdc, 0x30, 0x16, 0x52, 0x0a,
	0x97, 0xda, 0x95, 0xab, 0x57, 0x5c, 0xf2, 0x3d, 0xa0, 0x53, 0xe2, 0x7b, 0x74, 0x11, 0x0b, 0x33,
	0x93, 0x7c, 0x89, 0x08, 0x0b, 0xd1, 0x31, 0xd4, 0x24, 0xc6, 0xaf, 0xa9, 0x0f, 0x35, 0x23, 0xed,
	0x8c, 0x21, 0xcb, 0x37, 0x97, 0x24, 0xbe, 0x86, 0x07, 0x29, 0x2d, 0xe6, 0xd3, 0x39, 0x23, 0xe8,
	0x00, 0x54, 0xdb, 0x73, 0xc9, 0x3c, 0x1c, 0x27, 0x45, 0xab, 0x66, 0x4d, 0x2c, 0x8c, 0x1c, 0x74,
	0x04, 0x15, 0x61, 0x9d, 0xa6, 0xf0, 0x7b, 0x90, 0x21, 0xbc, 0x33, 0x02, 0xdf, 0x36, 0x3e, 0xf3,
	0x1d, 0x33, 0x21, 0xf0, 0x2b, 0x68, 0x9c, 0x38, 0x4e, 0x2c, 0x2e, 0xeb, 0xec, 0xc3, 0x76, 0x10,
	0x79, 0x24, 0xa9, 0xb1, 0x95, 0xad, 0x91, 0xa3, 0x9c, 0xc0, 0xdf, 0xa0, 0x79, 0x41, 0x1d, 0x77,
	0xba, 0xf8, 0xab, 0xe3, 0xe8, 0x25, 0xc0, 0xaa, 0x87, 0xdc, 0xe7, 0xfa, 0x50, 0x97, 0xa5, 0xca,
	0x36, 0x1b, 0xef, 0x62, 0xe4, 0xc2, 0x62, 0xb7, 0xa6, 0x3a, 0x95, 0x43, 0xfc, 0x04, 0x9a, 0xa7,
	0xc4, 0x23, 0x21, 0x59, 0xbf, 0x79, 0x1f, 0xaa, 0xb1, 0xee, 0x78, 0xd9, 0xc6, 0x4a, 0x3c, 0x1d,
	0x39, 0xf8, 0x29, 0xfc, 0x7f, 0xee, 0xb2, 0x30, 0xd5, 0x8d, 0x87, 0x50, 0x4b, 0x60, 0xa6, 0x95,
	0x3a, 0x4a, 0x5f, 0x31, 0xab, 0x82, 0x66, 0xf8, 0x0d, 0x34, 0xd7, 0xf0, 0xc4, 0xf0, 0x23, 0xd8,
	0x89, 0xf7, 0x05, 0x7c, 0xd7, 0x77, 0x09, 0x24, 0xee, 0xff, 0x89, 0xe3, 0x2c, 0x9b, 0x79, 0xaf,
	0xfe, 0xff, 0x28, 0x41, 0x5b, 0x98, 0xfc, 0x4f, 0xf4, 0xee, 0x63, 0xfa, 0x6b, 0x68, 0x0b, 0xd3,
	0xb3, 0x95, 0x1c, 0xc2, 0x32, 0x00, 0xe3, 0xb5, 0xd4, 0xfd, 0x27, 0x17, 0x3f, 0x5a, 0x33, 0x82,
	0xf7, 0xa0, 0x15, 0xbb, 0x2a, 0xcf, 0xca, 0x46, 0xe0, 0x4f, 0xd0, 0xce, 0xac, 0x27, 0x8e, 0xbf,
	0x00, 0x55, 0x0a, 0x48, 0xd7, 0xef, 0xfe, 0xc0, 0x15, 0x3a, 0xfc, 0xb9, 0x03, 0x95, 0x33, 0x8e,
	0xa1, 0x2b, 0xa8, 0xaf, 0x85, 0x07, 0xe1, 0xec, 0xf1, 0x7c, 0x4a, 0xf5, 0xc3, 0x42, 0x46, 0x94,
	0x86, 0xb7, 0x9e, 0x95, 0xd0, 0x5b, 0xa8, 0x26, 0xc1, 0x41, 0x8f, 0xb2, 0x67, 0xd2, 0x89, 0xd2,
	0xf7, 0x72, 0xfe, 0x9e, 0xc5, 0x7f, 0x3c, 0xbc, 0x85, 0x46, 0x00, 0xab, 0x04, 0xa1, 0x6e, 0x56,
	0x27, 0x97, 0xae, 0x62, 0xa9, 0x55, 0x24, 0xf2, 0x52, 0xb9, 0xb8, 0x14, 0x48, 0x99, 0xa0, 0x2e,
	0x03, 0x80, 0x3a, 0x39, 0xcf, 0x33, 0x51, 0xd2, 0xbb, 0x05, 0x84, 0x34, 0x0c, 0x7d, 0x80, 0xfa,
	0x5a, 0x26, 0xf2, 0xad, 0xc8, 0x07, 0xa6, 0xa0, 0xc0, 0x4b, 0x68, 0xa4, 0x33, 0x81, 0x7a, 0x9b,
	0xad, 0xfb, 0x23, 0xc9, 0xf4, 0xe3, 0xce, 0x4b, 0x6e, 0x7c, 0xfc, 0x05, 0x92, 0xd7, 0xb0, 0x9b,
	0x7a, 0xd9, 0xe8, 0xf1, 0x26, 0xa3, 0xb2, 0x81, 0xd0, 0x7b, 0xbf, 0xa1, 0xa4, 0xa5, 0x93, 0x0a,
	0xbf, 0xf1, 0xf9, 0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x46, 0x5f, 0xec, 0x02, 0x83, 0x07, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EmittoClient is the client API for Emitto service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EmittoClient interface {
	DeployRules(ctx context.Context, in *DeployRulesRequest, opts ...grpc.CallOption) (Emitto_DeployRulesClient, error)
	AddRule(ctx context.Context, in *AddRuleRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	ModifyRule(ctx context.Context, in *ModifyRuleRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	DeleteRule(ctx context.Context, in *DeleteRuleRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	ListRules(ctx context.Context, in *ListRulesRequest, opts ...grpc.CallOption) (*ListRulesResponse, error)
	AddLocation(ctx context.Context, in *AddLocationRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	ModifyLocation(ctx context.Context, in *ModifyLocationRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	DeleteLocation(ctx context.Context, in *DeleteLocationRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	ListLocations(ctx context.Context, in *ListLocationsRequest, opts ...grpc.CallOption) (*ListLocationsResponse, error)
}

type emittoClient struct {
	cc *grpc.ClientConn
}

func NewEmittoClient(cc *grpc.ClientConn) EmittoClient {
	return &emittoClient{cc}
}

func (c *emittoClient) DeployRules(ctx context.Context, in *DeployRulesRequest, opts ...grpc.CallOption) (Emitto_DeployRulesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Emitto_serviceDesc.Streams[0], "/emitto.service.Emitto/DeployRules", opts...)
	if err != nil {
		return nil, err
	}
	x := &emittoDeployRulesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Emitto_DeployRulesClient interface {
	Recv() (*DeployRulesResponse, error)
	grpc.ClientStream
}

type emittoDeployRulesClient struct {
	grpc.ClientStream
}

func (x *emittoDeployRulesClient) Recv() (*DeployRulesResponse, error) {
	m := new(DeployRulesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *emittoClient) AddRule(ctx context.Context, in *AddRuleRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/emitto.service.Emitto/AddRule", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emittoClient) ModifyRule(ctx context.Context, in *ModifyRuleRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/emitto.service.Emitto/ModifyRule", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emittoClient) DeleteRule(ctx context.Context, in *DeleteRuleRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/emitto.service.Emitto/DeleteRule", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emittoClient) ListRules(ctx context.Context, in *ListRulesRequest, opts ...grpc.CallOption) (*ListRulesResponse, error) {
	out := new(ListRulesResponse)
	err := c.cc.Invoke(ctx, "/emitto.service.Emitto/ListRules", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emittoClient) AddLocation(ctx context.Context, in *AddLocationRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/emitto.service.Emitto/AddLocation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emittoClient) ModifyLocation(ctx context.Context, in *ModifyLocationRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/emitto.service.Emitto/ModifyLocation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emittoClient) DeleteLocation(ctx context.Context, in *DeleteLocationRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/emitto.service.Emitto/DeleteLocation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emittoClient) ListLocations(ctx context.Context, in *ListLocationsRequest, opts ...grpc.CallOption) (*ListLocationsResponse, error) {
	out := new(ListLocationsResponse)
	err := c.cc.Invoke(ctx, "/emitto.service.Emitto/ListLocations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmittoServer is the server API for Emitto service.
type EmittoServer interface {
	DeployRules(*DeployRulesRequest, Emitto_DeployRulesServer) error
	AddRule(context.Context, *AddRuleRequest) (*empty.Empty, error)
	ModifyRule(context.Context, *ModifyRuleRequest) (*empty.Empty, error)
	DeleteRule(context.Context, *DeleteRuleRequest) (*empty.Empty, error)
	ListRules(context.Context, *ListRulesRequest) (*ListRulesResponse, error)
	AddLocation(context.Context, *AddLocationRequest) (*empty.Empty, error)
	ModifyLocation(context.Context, *ModifyLocationRequest) (*empty.Empty, error)
	DeleteLocation(context.Context, *DeleteLocationRequest) (*empty.Empty, error)
	ListLocations(context.Context, *ListLocationsRequest) (*ListLocationsResponse, error)
}

func RegisterEmittoServer(s *grpc.Server, srv EmittoServer) {
	s.RegisterService(&_Emitto_serviceDesc, srv)
}

func _Emitto_DeployRules_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DeployRulesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EmittoServer).DeployRules(m, &emittoDeployRulesServer{stream})
}

type Emitto_DeployRulesServer interface {
	Send(*DeployRulesResponse) error
	grpc.ServerStream
}

type emittoDeployRulesServer struct {
	grpc.ServerStream
}

func (x *emittoDeployRulesServer) Send(m *DeployRulesResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Emitto_AddRule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmittoServer).AddRule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/emitto.service.Emitto/AddRule",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmittoServer).AddRule(ctx, req.(*AddRuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Emitto_ModifyRule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifyRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmittoServer).ModifyRule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/emitto.service.Emitto/ModifyRule",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmittoServer).ModifyRule(ctx, req.(*ModifyRuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Emitto_DeleteRule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmittoServer).DeleteRule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/emitto.service.Emitto/DeleteRule",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmittoServer).DeleteRule(ctx, req.(*DeleteRuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Emitto_ListRules_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRulesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmittoServer).ListRules(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/emitto.service.Emitto/ListRules",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmittoServer).ListRules(ctx, req.(*ListRulesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Emitto_AddLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddLocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmittoServer).AddLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/emitto.service.Emitto/AddLocation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmittoServer).AddLocation(ctx, req.(*AddLocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Emitto_ModifyLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifyLocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmittoServer).ModifyLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/emitto.service.Emitto/ModifyLocation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmittoServer).ModifyLocation(ctx, req.(*ModifyLocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Emitto_DeleteLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteLocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmittoServer).DeleteLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/emitto.service.Emitto/DeleteLocation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmittoServer).DeleteLocation(ctx, req.(*DeleteLocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Emitto_ListLocations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLocationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmittoServer).ListLocations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/emitto.service.Emitto/ListLocations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmittoServer).ListLocations(ctx, req.(*ListLocationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Emitto_serviceDesc = grpc.ServiceDesc{
	ServiceName: "emitto.service.Emitto",
	HandlerType: (*EmittoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddRule",
			Handler:    _Emitto_AddRule_Handler,
		},
		{
			MethodName: "ModifyRule",
			Handler:    _Emitto_ModifyRule_Handler,
		},
		{
			MethodName: "DeleteRule",
			Handler:    _Emitto_DeleteRule_Handler,
		},
		{
			MethodName: "ListRules",
			Handler:    _Emitto_ListRules_Handler,
		},
		{
			MethodName: "AddLocation",
			Handler:    _Emitto_AddLocation_Handler,
		},
		{
			MethodName: "ModifyLocation",
			Handler:    _Emitto_ModifyLocation_Handler,
		},
		{
			MethodName: "DeleteLocation",
			Handler:    _Emitto_DeleteLocation_Handler,
		},
		{
			MethodName: "ListLocations",
			Handler:    _Emitto_ListLocations_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DeployRules",
			Handler:       _Emitto_DeployRules_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "source/server/proto/service.proto",
}
