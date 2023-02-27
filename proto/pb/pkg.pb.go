// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.2
// source: pkg.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RouteType int32

const (
	RouteType_RT_MIN    RouteType = 0
	RouteType_RT_P2P    RouteType = 1
	RouteType_RT_RAND   RouteType = 2
	RouteType_RT_HASH   RouteType = 3
	RouteType_RT_BROAD  RouteType = 4
	RouteType_RT_MUTIL  RouteType = 5
	RouteType_RT_SINGLE RouteType = 6
)

// Enum value maps for RouteType.
var (
	RouteType_name = map[int32]string{
		0: "RT_MIN",
		1: "RT_P2P",
		2: "RT_RAND",
		3: "RT_HASH",
		4: "RT_BROAD",
		5: "RT_MUTIL",
		6: "RT_SINGLE",
	}
	RouteType_value = map[string]int32{
		"RT_MIN":    0,
		"RT_P2P":    1,
		"RT_RAND":   2,
		"RT_HASH":   3,
		"RT_BROAD":  4,
		"RT_MUTIL":  5,
		"RT_SINGLE": 6,
	}
)

func (x RouteType) Enum() *RouteType {
	p := new(RouteType)
	*p = x
	return p
}

func (x RouteType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RouteType) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_proto_enumTypes[0].Descriptor()
}

func (RouteType) Type() protoreflect.EnumType {
	return &file_pkg_proto_enumTypes[0]
}

func (x RouteType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RouteType.Descriptor instead.
func (RouteType) EnumDescriptor() ([]byte, []int) {
	return file_pkg_proto_rawDescGZIP(), []int{0}
}

//
type Meta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *Meta) Reset() {
	*x = Meta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Meta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Meta) ProtoMessage() {}

func (x *Meta) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Meta.ProtoReflect.Descriptor instead.
func (*Meta) Descriptor() ([]byte, []int) {
	return file_pkg_proto_rawDescGZIP(), []int{0}
}

func (x *Meta) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Meta) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type Metas struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metas []*Meta `protobuf:"bytes,1,rep,name=Metas,proto3" json:"Metas,omitempty"`
}

func (x *Metas) Reset() {
	*x = Metas{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metas) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metas) ProtoMessage() {}

func (x *Metas) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metas.ProtoReflect.Descriptor instead.
func (*Metas) Descriptor() ([]byte, []int) {
	return file_pkg_proto_rawDescGZIP(), []int{1}
}

func (x *Metas) GetMetas() []*Meta {
	if x != nil {
		return x.Metas
	}
	return nil
}

type RouteIns struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EnvID   string `protobuf:"bytes,1,opt,name=EnvID,proto3" json:"EnvID,omitempty"`
	WorldID string `protobuf:"bytes,2,opt,name=WorldID,proto3" json:"WorldID,omitempty"`
	FuncID  string `protobuf:"bytes,3,opt,name=FuncID,proto3" json:"FuncID,omitempty"`
	InsID   string `protobuf:"bytes,4,opt,name=InsID,proto3" json:"InsID,omitempty"`
	IDStr   string `protobuf:"bytes,5,opt,name=IDStr,proto3" json:"IDStr,omitempty"`
}

func (x *RouteIns) Reset() {
	*x = RouteIns{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteIns) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteIns) ProtoMessage() {}

func (x *RouteIns) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteIns.ProtoReflect.Descriptor instead.
func (*RouteIns) Descriptor() ([]byte, []int) {
	return file_pkg_proto_rawDescGZIP(), []int{2}
}

func (x *RouteIns) GetEnvID() string {
	if x != nil {
		return x.EnvID
	}
	return ""
}

func (x *RouteIns) GetWorldID() string {
	if x != nil {
		return x.WorldID
	}
	return ""
}

func (x *RouteIns) GetFuncID() string {
	if x != nil {
		return x.FuncID
	}
	return ""
}

func (x *RouteIns) GetInsID() string {
	if x != nil {
		return x.InsID
	}
	return ""
}

func (x *RouteIns) GetIDStr() string {
	if x != nil {
		return x.IDStr
	}
	return ""
}

type RouteSvc struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EnvID   string `protobuf:"bytes,1,opt,name=EnvID,proto3" json:"EnvID,omitempty"`
	WorldID string `protobuf:"bytes,2,opt,name=WorldID,proto3" json:"WorldID,omitempty"`
	FuncID  string `protobuf:"bytes,3,opt,name=FuncID,proto3" json:"FuncID,omitempty"`
	IDStr   string `protobuf:"bytes,4,opt,name=IDStr,proto3" json:"IDStr,omitempty"`
}

func (x *RouteSvc) Reset() {
	*x = RouteSvc{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteSvc) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteSvc) ProtoMessage() {}

func (x *RouteSvc) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteSvc.ProtoReflect.Descriptor instead.
func (*RouteSvc) Descriptor() ([]byte, []int) {
	return file_pkg_proto_rawDescGZIP(), []int{3}
}

func (x *RouteSvc) GetEnvID() string {
	if x != nil {
		return x.EnvID
	}
	return ""
}

func (x *RouteSvc) GetWorldID() string {
	if x != nil {
		return x.WorldID
	}
	return ""
}

func (x *RouteSvc) GetFuncID() string {
	if x != nil {
		return x.FuncID
	}
	return ""
}

func (x *RouteSvc) GetIDStr() string {
	if x != nil {
		return x.IDStr
	}
	return ""
}

type RouteWorld struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EnvID   string `protobuf:"bytes,1,opt,name=EnvID,proto3" json:"EnvID,omitempty"`
	WorldID string `protobuf:"bytes,2,opt,name=WorldID,proto3" json:"WorldID,omitempty"`
}

func (x *RouteWorld) Reset() {
	*x = RouteWorld{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteWorld) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteWorld) ProtoMessage() {}

func (x *RouteWorld) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteWorld.ProtoReflect.Descriptor instead.
func (*RouteWorld) Descriptor() ([]byte, []int) {
	return file_pkg_proto_rawDescGZIP(), []int{4}
}

func (x *RouteWorld) GetEnvID() string {
	if x != nil {
		return x.EnvID
	}
	return ""
}

func (x *RouteWorld) GetWorldID() string {
	if x != nil {
		return x.WorldID
	}
	return ""
}

type RouteP2PHead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DstIns *RouteIns `protobuf:"bytes,1,opt,name=DstIns,proto3" json:"DstIns,omitempty"`
}

func (x *RouteP2PHead) Reset() {
	*x = RouteP2PHead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteP2PHead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteP2PHead) ProtoMessage() {}

func (x *RouteP2PHead) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteP2PHead.ProtoReflect.Descriptor instead.
func (*RouteP2PHead) Descriptor() ([]byte, []int) {
	return file_pkg_proto_rawDescGZIP(), []int{5}
}

func (x *RouteP2PHead) GetDstIns() *RouteIns {
	if x != nil {
		return x.DstIns
	}
	return nil
}

type RouteRandHead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pass string `protobuf:"bytes,1,opt,name=Pass,proto3" json:"Pass,omitempty"`
}

func (x *RouteRandHead) Reset() {
	*x = RouteRandHead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteRandHead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteRandHead) ProtoMessage() {}

func (x *RouteRandHead) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteRandHead.ProtoReflect.Descriptor instead.
func (*RouteRandHead) Descriptor() ([]byte, []int) {
	return file_pkg_proto_rawDescGZIP(), []int{6}
}

func (x *RouteRandHead) GetPass() string {
	if x != nil {
		return x.Pass
	}
	return ""
}

type RouteHashHead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HashKey string `protobuf:"bytes,1,opt,name=HashKey,proto3" json:"HashKey,omitempty"`
}

func (x *RouteHashHead) Reset() {
	*x = RouteHashHead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteHashHead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteHashHead) ProtoMessage() {}

func (x *RouteHashHead) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteHashHead.ProtoReflect.Descriptor instead.
func (*RouteHashHead) Descriptor() ([]byte, []int) {
	return file_pkg_proto_rawDescGZIP(), []int{7}
}

func (x *RouteHashHead) GetHashKey() string {
	if x != nil {
		return x.HashKey
	}
	return ""
}

type RouteBroadHead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pass string `protobuf:"bytes,1,opt,name=Pass,proto3" json:"Pass,omitempty"`
}

func (x *RouteBroadHead) Reset() {
	*x = RouteBroadHead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteBroadHead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteBroadHead) ProtoMessage() {}

func (x *RouteBroadHead) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteBroadHead.ProtoReflect.Descriptor instead.
func (*RouteBroadHead) Descriptor() ([]byte, []int) {
	return file_pkg_proto_rawDescGZIP(), []int{8}
}

func (x *RouteBroadHead) GetPass() string {
	if x != nil {
		return x.Pass
	}
	return ""
}

type RouteMutilHead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic string `protobuf:"bytes,1,opt,name=Topic,proto3" json:"Topic,omitempty"`
}

func (x *RouteMutilHead) Reset() {
	*x = RouteMutilHead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteMutilHead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteMutilHead) ProtoMessage() {}

func (x *RouteMutilHead) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteMutilHead.ProtoReflect.Descriptor instead.
func (*RouteMutilHead) Descriptor() ([]byte, []int) {
	return file_pkg_proto_rawDescGZIP(), []int{9}
}

func (x *RouteMutilHead) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

type RouteSingleHead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pass string `protobuf:"bytes,1,opt,name=Pass,proto3" json:"Pass,omitempty"`
}

func (x *RouteSingleHead) Reset() {
	*x = RouteSingleHead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteSingleHead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteSingleHead) ProtoMessage() {}

func (x *RouteSingleHead) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteSingleHead.ProtoReflect.Descriptor instead.
func (*RouteSingleHead) Descriptor() ([]byte, []int) {
	return file_pkg_proto_rawDescGZIP(), []int{10}
}

func (x *RouteSingleHead) GetPass() string {
	if x != nil {
		return x.Pass
	}
	return ""
}

type RoutePara struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RouteP2PHead    []*RouteP2PHead    `protobuf:"bytes,1,rep,name=RouteP2PHead,proto3" json:"RouteP2PHead,omitempty"`
	RouteRandHead   []*RouteRandHead   `protobuf:"bytes,2,rep,name=RouteRandHead,proto3" json:"RouteRandHead,omitempty"`
	RouteHashHead   []*RouteHashHead   `protobuf:"bytes,3,rep,name=RouteHashHead,proto3" json:"RouteHashHead,omitempty"`
	RouteBroadHead  []*RouteBroadHead  `protobuf:"bytes,4,rep,name=RouteBroadHead,proto3" json:"RouteBroadHead,omitempty"`
	RouteMutilHead  []*RouteMutilHead  `protobuf:"bytes,5,rep,name=RouteMutilHead,proto3" json:"RouteMutilHead,omitempty"`
	RouteSingleHead []*RouteSingleHead `protobuf:"bytes,6,rep,name=RouteSingleHead,proto3" json:"RouteSingleHead,omitempty"`
}

func (x *RoutePara) Reset() {
	*x = RoutePara{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoutePara) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoutePara) ProtoMessage() {}

func (x *RoutePara) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoutePara.ProtoReflect.Descriptor instead.
func (*RoutePara) Descriptor() ([]byte, []int) {
	return file_pkg_proto_rawDescGZIP(), []int{11}
}

func (x *RoutePara) GetRouteP2PHead() []*RouteP2PHead {
	if x != nil {
		return x.RouteP2PHead
	}
	return nil
}

func (x *RoutePara) GetRouteRandHead() []*RouteRandHead {
	if x != nil {
		return x.RouteRandHead
	}
	return nil
}

func (x *RoutePara) GetRouteHashHead() []*RouteHashHead {
	if x != nil {
		return x.RouteHashHead
	}
	return nil
}

func (x *RoutePara) GetRouteBroadHead() []*RouteBroadHead {
	if x != nil {
		return x.RouteBroadHead
	}
	return nil
}

func (x *RoutePara) GetRouteMutilHead() []*RouteMutilHead {
	if x != nil {
		return x.RouteMutilHead
	}
	return nil
}

func (x *RoutePara) GetRouteSingleHead() []*RouteSingleHead {
	if x != nil {
		return x.RouteSingleHead
	}
	return nil
}

type RouteHead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SrcIns    *RouteIns  `protobuf:"bytes,1,opt,name=SrcIns,proto3" json:"SrcIns,omitempty"`
	DstSvc    *RouteSvc  `protobuf:"bytes,2,opt,name=DstSvc,proto3" json:"DstSvc,omitempty"`
	RouteType RouteType  `protobuf:"varint,3,opt,name=RouteType,proto3,enum=proto.RouteType" json:"RouteType,omitempty"`
	RoutePara *RoutePara `protobuf:"bytes,4,opt,name=RoutePara,proto3" json:"RoutePara,omitempty"`
	Metas     *Metas     `protobuf:"bytes,5,opt,name=Metas,proto3" json:"Metas,omitempty"`
}

func (x *RouteHead) Reset() {
	*x = RouteHead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteHead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteHead) ProtoMessage() {}

func (x *RouteHead) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteHead.ProtoReflect.Descriptor instead.
func (*RouteHead) Descriptor() ([]byte, []int) {
	return file_pkg_proto_rawDescGZIP(), []int{12}
}

func (x *RouteHead) GetSrcIns() *RouteIns {
	if x != nil {
		return x.SrcIns
	}
	return nil
}

func (x *RouteHead) GetDstSvc() *RouteSvc {
	if x != nil {
		return x.DstSvc
	}
	return nil
}

func (x *RouteHead) GetRouteType() RouteType {
	if x != nil {
		return x.RouteType
	}
	return RouteType_RT_MIN
}

func (x *RouteHead) GetRoutePara() *RoutePara {
	if x != nil {
		return x.RoutePara
	}
	return nil
}

func (x *RouteHead) GetMetas() *Metas {
	if x != nil {
		return x.Metas
	}
	return nil
}

type BrokerPkg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FullMethod string `protobuf:"bytes,1,opt,name=FullMethod,proto3" json:"FullMethod,omitempty"`
	Content    []byte `protobuf:"bytes,2,opt,name=Content,proto3" json:"Content,omitempty"`
}

func (x *BrokerPkg) Reset() {
	*x = BrokerPkg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BrokerPkg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BrokerPkg) ProtoMessage() {}

func (x *BrokerPkg) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BrokerPkg.ProtoReflect.Descriptor instead.
func (*BrokerPkg) Descriptor() ([]byte, []int) {
	return file_pkg_proto_rawDescGZIP(), []int{13}
}

func (x *BrokerPkg) GetFullMethod() string {
	if x != nil {
		return x.FullMethod
	}
	return ""
}

func (x *BrokerPkg) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

var File_pkg_proto protoreflect.FileDescriptor

var file_pkg_proto_rawDesc = []byte{
	0x0a, 0x09, 0x70, 0x6b, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x2e, 0x0a, 0x04, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x2a, 0x0a, 0x05, 0x4d, 0x65, 0x74, 0x61, 0x73, 0x12, 0x21, 0x0a, 0x05, 0x4d,
	0x65, 0x74, 0x61, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x05, 0x4d, 0x65, 0x74, 0x61, 0x73, 0x22, 0x7e,
	0x0a, 0x08, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x49, 0x6e, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x6e,
	0x76, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x6e, 0x76, 0x49, 0x44,
	0x12, 0x18, 0x0a, 0x07, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x75,
	0x6e, 0x63, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x46, 0x75, 0x6e, 0x63,
	0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x49, 0x6e, 0x73, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x49, 0x6e, 0x73, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x49, 0x44, 0x53, 0x74,
	0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x49, 0x44, 0x53, 0x74, 0x72, 0x22, 0x68,
	0x0a, 0x08, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x53, 0x76, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x6e,
	0x76, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x6e, 0x76, 0x49, 0x44,
	0x12, 0x18, 0x0a, 0x07, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x75,
	0x6e, 0x63, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x46, 0x75, 0x6e, 0x63,
	0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x49, 0x44, 0x53, 0x74, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x49, 0x44, 0x53, 0x74, 0x72, 0x22, 0x3c, 0x0a, 0x0a, 0x52, 0x6f, 0x75, 0x74,
	0x65, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x6e, 0x76, 0x49, 0x44, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x6e, 0x76, 0x49, 0x44, 0x12, 0x18, 0x0a, 0x07,
	0x57, 0x6f, 0x72, 0x6c, 0x64, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x57,
	0x6f, 0x72, 0x6c, 0x64, 0x49, 0x44, 0x22, 0x37, 0x0a, 0x0c, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x50,
	0x32, 0x50, 0x48, 0x65, 0x61, 0x64, 0x12, 0x27, 0x0a, 0x06, 0x44, 0x73, 0x74, 0x49, 0x6e, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52,
	0x6f, 0x75, 0x74, 0x65, 0x49, 0x6e, 0x73, 0x52, 0x06, 0x44, 0x73, 0x74, 0x49, 0x6e, 0x73, 0x22,
	0x23, 0x0a, 0x0d, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x61, 0x6e, 0x64, 0x48, 0x65, 0x61, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x50, 0x61, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x50, 0x61, 0x73, 0x73, 0x22, 0x29, 0x0a, 0x0d, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x48, 0x61, 0x73,
	0x68, 0x48, 0x65, 0x61, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x48, 0x61, 0x73, 0x68, 0x4b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x48, 0x61, 0x73, 0x68, 0x4b, 0x65, 0x79, 0x22,
	0x24, 0x0a, 0x0e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x48, 0x65, 0x61,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x50, 0x61, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x50, 0x61, 0x73, 0x73, 0x22, 0x26, 0x0a, 0x0e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x4d, 0x75,
	0x74, 0x69, 0x6c, 0x48, 0x65, 0x61, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x6f, 0x70, 0x69, 0x63,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x22, 0x25, 0x0a,
	0x0f, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x48, 0x65, 0x61, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x50, 0x61, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x50, 0x61, 0x73, 0x73, 0x22, 0xfc, 0x02, 0x0a, 0x09, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x50, 0x61,
	0x72, 0x61, 0x12, 0x37, 0x0a, 0x0c, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x50, 0x32, 0x50, 0x48, 0x65,
	0x61, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x50, 0x32, 0x50, 0x48, 0x65, 0x61, 0x64, 0x52, 0x0c, 0x52,
	0x6f, 0x75, 0x74, 0x65, 0x50, 0x32, 0x50, 0x48, 0x65, 0x61, 0x64, 0x12, 0x3a, 0x0a, 0x0d, 0x52,
	0x6f, 0x75, 0x74, 0x65, 0x52, 0x61, 0x6e, 0x64, 0x48, 0x65, 0x61, 0x64, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65,
	0x52, 0x61, 0x6e, 0x64, 0x48, 0x65, 0x61, 0x64, 0x52, 0x0d, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52,
	0x61, 0x6e, 0x64, 0x48, 0x65, 0x61, 0x64, 0x12, 0x3a, 0x0a, 0x0d, 0x52, 0x6f, 0x75, 0x74, 0x65,
	0x48, 0x61, 0x73, 0x68, 0x48, 0x65, 0x61, 0x64, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x48, 0x61, 0x73, 0x68,
	0x48, 0x65, 0x61, 0x64, 0x52, 0x0d, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x48, 0x61, 0x73, 0x68, 0x48,
	0x65, 0x61, 0x64, 0x12, 0x3d, 0x0a, 0x0e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x42, 0x72, 0x6f, 0x61,
	0x64, 0x48, 0x65, 0x61, 0x64, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x48, 0x65,
	0x61, 0x64, 0x52, 0x0e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x48, 0x65,
	0x61, 0x64, 0x12, 0x3d, 0x0a, 0x0e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x4d, 0x75, 0x74, 0x69, 0x6c,
	0x48, 0x65, 0x61, 0x64, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x4d, 0x75, 0x74, 0x69, 0x6c, 0x48, 0x65, 0x61,
	0x64, 0x52, 0x0e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x4d, 0x75, 0x74, 0x69, 0x6c, 0x48, 0x65, 0x61,
	0x64, 0x12, 0x40, 0x0a, 0x0f, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65,
	0x48, 0x65, 0x61, 0x64, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x48, 0x65,
	0x61, 0x64, 0x52, 0x0f, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x48,
	0x65, 0x61, 0x64, 0x22, 0xe1, 0x01, 0x0a, 0x09, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x48, 0x65, 0x61,
	0x64, 0x12, 0x27, 0x0a, 0x06, 0x53, 0x72, 0x63, 0x49, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x49,
	0x6e, 0x73, 0x52, 0x06, 0x53, 0x72, 0x63, 0x49, 0x6e, 0x73, 0x12, 0x27, 0x0a, 0x06, 0x44, 0x73,
	0x74, 0x53, 0x76, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x53, 0x76, 0x63, 0x52, 0x06, 0x44, 0x73, 0x74,
	0x53, 0x76, 0x63, 0x12, 0x2e, 0x0a, 0x09, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52,
	0x6f, 0x75, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x2e, 0x0a, 0x09, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x50, 0x61, 0x72, 0x61,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52,
	0x6f, 0x75, 0x74, 0x65, 0x50, 0x61, 0x72, 0x61, 0x52, 0x09, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x50,
	0x61, 0x72, 0x61, 0x12, 0x22, 0x0a, 0x05, 0x4d, 0x65, 0x74, 0x61, 0x73, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x73,
	0x52, 0x05, 0x4d, 0x65, 0x74, 0x61, 0x73, 0x22, 0x45, 0x0a, 0x09, 0x42, 0x72, 0x6f, 0x6b, 0x65,
	0x72, 0x50, 0x6b, 0x67, 0x12, 0x1e, 0x0a, 0x0a, 0x46, 0x75, 0x6c, 0x6c, 0x4d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x46, 0x75, 0x6c, 0x6c, 0x4d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2a, 0x68,
	0x0a, 0x09, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x52,
	0x54, 0x5f, 0x4d, 0x49, 0x4e, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x52, 0x54, 0x5f, 0x50, 0x32,
	0x50, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x54, 0x5f, 0x52, 0x41, 0x4e, 0x44, 0x10, 0x02,
	0x12, 0x0b, 0x0a, 0x07, 0x52, 0x54, 0x5f, 0x48, 0x41, 0x53, 0x48, 0x10, 0x03, 0x12, 0x0c, 0x0a,
	0x08, 0x52, 0x54, 0x5f, 0x42, 0x52, 0x4f, 0x41, 0x44, 0x10, 0x04, 0x12, 0x0c, 0x0a, 0x08, 0x52,
	0x54, 0x5f, 0x4d, 0x55, 0x54, 0x49, 0x4c, 0x10, 0x05, 0x12, 0x0d, 0x0a, 0x09, 0x52, 0x54, 0x5f,
	0x53, 0x49, 0x4e, 0x47, 0x4c, 0x45, 0x10, 0x06, 0x42, 0x11, 0x5a, 0x0f, 0x6d, 0x33, 0x67, 0x61,
	0x6d, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_pkg_proto_rawDescOnce sync.Once
	file_pkg_proto_rawDescData = file_pkg_proto_rawDesc
)

func file_pkg_proto_rawDescGZIP() []byte {
	file_pkg_proto_rawDescOnce.Do(func() {
		file_pkg_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_proto_rawDescData)
	})
	return file_pkg_proto_rawDescData
}

var file_pkg_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pkg_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_pkg_proto_goTypes = []interface{}{
	(RouteType)(0),          // 0: proto.RouteType
	(*Meta)(nil),            // 1: proto.Meta
	(*Metas)(nil),           // 2: proto.Metas
	(*RouteIns)(nil),        // 3: proto.RouteIns
	(*RouteSvc)(nil),        // 4: proto.RouteSvc
	(*RouteWorld)(nil),      // 5: proto.RouteWorld
	(*RouteP2PHead)(nil),    // 6: proto.RouteP2PHead
	(*RouteRandHead)(nil),   // 7: proto.RouteRandHead
	(*RouteHashHead)(nil),   // 8: proto.RouteHashHead
	(*RouteBroadHead)(nil),  // 9: proto.RouteBroadHead
	(*RouteMutilHead)(nil),  // 10: proto.RouteMutilHead
	(*RouteSingleHead)(nil), // 11: proto.RouteSingleHead
	(*RoutePara)(nil),       // 12: proto.RoutePara
	(*RouteHead)(nil),       // 13: proto.RouteHead
	(*BrokerPkg)(nil),       // 14: proto.BrokerPkg
}
var file_pkg_proto_depIdxs = []int32{
	1,  // 0: proto.Metas.Metas:type_name -> proto.Meta
	3,  // 1: proto.RouteP2PHead.DstIns:type_name -> proto.RouteIns
	6,  // 2: proto.RoutePara.RouteP2PHead:type_name -> proto.RouteP2PHead
	7,  // 3: proto.RoutePara.RouteRandHead:type_name -> proto.RouteRandHead
	8,  // 4: proto.RoutePara.RouteHashHead:type_name -> proto.RouteHashHead
	9,  // 5: proto.RoutePara.RouteBroadHead:type_name -> proto.RouteBroadHead
	10, // 6: proto.RoutePara.RouteMutilHead:type_name -> proto.RouteMutilHead
	11, // 7: proto.RoutePara.RouteSingleHead:type_name -> proto.RouteSingleHead
	3,  // 8: proto.RouteHead.SrcIns:type_name -> proto.RouteIns
	4,  // 9: proto.RouteHead.DstSvc:type_name -> proto.RouteSvc
	0,  // 10: proto.RouteHead.RouteType:type_name -> proto.RouteType
	12, // 11: proto.RouteHead.RoutePara:type_name -> proto.RoutePara
	2,  // 12: proto.RouteHead.Metas:type_name -> proto.Metas
	13, // [13:13] is the sub-list for method output_type
	13, // [13:13] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_pkg_proto_init() }
func file_pkg_proto_init() {
	if File_pkg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Meta); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metas); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteIns); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteSvc); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteWorld); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteP2PHead); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteRandHead); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteHashHead); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteBroadHead); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteMutilHead); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteSingleHead); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoutePara); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteHead); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BrokerPkg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_proto_goTypes,
		DependencyIndexes: file_pkg_proto_depIdxs,
		EnumInfos:         file_pkg_proto_enumTypes,
		MessageInfos:      file_pkg_proto_msgTypes,
	}.Build()
	File_pkg_proto = out.File
	file_pkg_proto_rawDesc = nil
	file_pkg_proto_goTypes = nil
	file_pkg_proto_depIdxs = nil
}
