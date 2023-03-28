// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.2.0
// source: online.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "m3game/meta/metapb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ORFlag int32

const (
	ORFlag_OROnlineRoleMin ORFlag = 0
	ORFlag_ORRoleId        ORFlag = 1
	ORFlag_OROnlineApp     ORFlag = 2
)

// Enum value maps for ORFlag.
var (
	ORFlag_name = map[int32]string{
		0: "OROnlineRoleMin",
		1: "ORRoleId",
		2: "OROnlineApp",
	}
	ORFlag_value = map[string]int32{
		"OROnlineRoleMin": 0,
		"ORRoleId":        1,
		"OROnlineApp":     2,
	}
)

func (x ORFlag) Enum() *ORFlag {
	p := new(ORFlag)
	*p = x
	return p
}

func (x ORFlag) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ORFlag) Descriptor() protoreflect.EnumDescriptor {
	return file_online_proto_enumTypes[0].Descriptor()
}

func (ORFlag) Type() protoreflect.EnumType {
	return &file_online_proto_enumTypes[0]
}

func (x ORFlag) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ORFlag.Descriptor instead.
func (ORFlag) EnumDescriptor() ([]byte, []int) {
	return file_online_proto_rawDescGZIP(), []int{0}
}

type OnlineRoleDB struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoleId    int64      `protobuf:"varint,1,opt,name=RoleId,proto3" json:"RoleId,omitempty"`
	OnlineApp *OnlineApp `protobuf:"bytes,2,opt,name=OnlineApp,proto3" json:"OnlineApp,omitempty"`
}

func (x *OnlineRoleDB) Reset() {
	*x = OnlineRoleDB{}
	if protoimpl.UnsafeEnabled {
		mi := &file_online_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnlineRoleDB) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnlineRoleDB) ProtoMessage() {}

func (x *OnlineRoleDB) ProtoReflect() protoreflect.Message {
	mi := &file_online_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnlineRoleDB.ProtoReflect.Descriptor instead.
func (*OnlineRoleDB) Descriptor() ([]byte, []int) {
	return file_online_proto_rawDescGZIP(), []int{0}
}

func (x *OnlineRoleDB) GetRoleId() int64 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

func (x *OnlineRoleDB) GetOnlineApp() *OnlineApp {
	if x != nil {
		return x.OnlineApp
	}
	return nil
}

type OnlineApp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AppId string `protobuf:"bytes,1,opt,name=AppId,proto3" json:"AppId,omitempty"`
	Ver   string `protobuf:"bytes,2,opt,name=Ver,proto3" json:"Ver,omitempty"`
}

func (x *OnlineApp) Reset() {
	*x = OnlineApp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_online_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnlineApp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnlineApp) ProtoMessage() {}

func (x *OnlineApp) ProtoReflect() protoreflect.Message {
	mi := &file_online_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnlineApp.ProtoReflect.Descriptor instead.
func (*OnlineApp) Descriptor() ([]byte, []int) {
	return file_online_proto_rawDescGZIP(), []int{1}
}

func (x *OnlineApp) GetAppId() string {
	if x != nil {
		return x.AppId
	}
	return ""
}

func (x *OnlineApp) GetVer() string {
	if x != nil {
		return x.Ver
	}
	return ""
}

type OnlineCreate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *OnlineCreate) Reset() {
	*x = OnlineCreate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_online_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnlineCreate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnlineCreate) ProtoMessage() {}

func (x *OnlineCreate) ProtoReflect() protoreflect.Message {
	mi := &file_online_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnlineCreate.ProtoReflect.Descriptor instead.
func (*OnlineCreate) Descriptor() ([]byte, []int) {
	return file_online_proto_rawDescGZIP(), []int{2}
}

type OnlineRead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *OnlineRead) Reset() {
	*x = OnlineRead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_online_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnlineRead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnlineRead) ProtoMessage() {}

func (x *OnlineRead) ProtoReflect() protoreflect.Message {
	mi := &file_online_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnlineRead.ProtoReflect.Descriptor instead.
func (*OnlineRead) Descriptor() ([]byte, []int) {
	return file_online_proto_rawDescGZIP(), []int{3}
}

type OnlineDelete struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *OnlineDelete) Reset() {
	*x = OnlineDelete{}
	if protoimpl.UnsafeEnabled {
		mi := &file_online_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnlineDelete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnlineDelete) ProtoMessage() {}

func (x *OnlineDelete) ProtoReflect() protoreflect.Message {
	mi := &file_online_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnlineDelete.ProtoReflect.Descriptor instead.
func (*OnlineDelete) Descriptor() ([]byte, []int) {
	return file_online_proto_rawDescGZIP(), []int{4}
}

type OnlineCreate_Req struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoleId int64  `protobuf:"varint,1,opt,name=RoleId,proto3" json:"RoleId,omitempty"`
	AppId  string `protobuf:"bytes,2,opt,name=AppId,proto3" json:"AppId,omitempty"`
}

func (x *OnlineCreate_Req) Reset() {
	*x = OnlineCreate_Req{}
	if protoimpl.UnsafeEnabled {
		mi := &file_online_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnlineCreate_Req) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnlineCreate_Req) ProtoMessage() {}

func (x *OnlineCreate_Req) ProtoReflect() protoreflect.Message {
	mi := &file_online_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnlineCreate_Req.ProtoReflect.Descriptor instead.
func (*OnlineCreate_Req) Descriptor() ([]byte, []int) {
	return file_online_proto_rawDescGZIP(), []int{2, 0}
}

func (x *OnlineCreate_Req) GetRoleId() int64 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

func (x *OnlineCreate_Req) GetAppId() string {
	if x != nil {
		return x.AppId
	}
	return ""
}

type OnlineCreate_Rsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *OnlineCreate_Rsp) Reset() {
	*x = OnlineCreate_Rsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_online_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnlineCreate_Rsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnlineCreate_Rsp) ProtoMessage() {}

func (x *OnlineCreate_Rsp) ProtoReflect() protoreflect.Message {
	mi := &file_online_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnlineCreate_Rsp.ProtoReflect.Descriptor instead.
func (*OnlineCreate_Rsp) Descriptor() ([]byte, []int) {
	return file_online_proto_rawDescGZIP(), []int{2, 1}
}

type OnlineRead_Req struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoleId int64 `protobuf:"varint,1,opt,name=RoleId,proto3" json:"RoleId,omitempty"`
}

func (x *OnlineRead_Req) Reset() {
	*x = OnlineRead_Req{}
	if protoimpl.UnsafeEnabled {
		mi := &file_online_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnlineRead_Req) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnlineRead_Req) ProtoMessage() {}

func (x *OnlineRead_Req) ProtoReflect() protoreflect.Message {
	mi := &file_online_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnlineRead_Req.ProtoReflect.Descriptor instead.
func (*OnlineRead_Req) Descriptor() ([]byte, []int) {
	return file_online_proto_rawDescGZIP(), []int{3, 0}
}

func (x *OnlineRead_Req) GetRoleId() int64 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

type OnlineRead_Rsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AppId string `protobuf:"bytes,1,opt,name=AppId,proto3" json:"AppId,omitempty"`
}

func (x *OnlineRead_Rsp) Reset() {
	*x = OnlineRead_Rsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_online_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnlineRead_Rsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnlineRead_Rsp) ProtoMessage() {}

func (x *OnlineRead_Rsp) ProtoReflect() protoreflect.Message {
	mi := &file_online_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnlineRead_Rsp.ProtoReflect.Descriptor instead.
func (*OnlineRead_Rsp) Descriptor() ([]byte, []int) {
	return file_online_proto_rawDescGZIP(), []int{3, 1}
}

func (x *OnlineRead_Rsp) GetAppId() string {
	if x != nil {
		return x.AppId
	}
	return ""
}

type OnlineDelete_Req struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoleId int64  `protobuf:"varint,1,opt,name=RoleId,proto3" json:"RoleId,omitempty"`
	AppId  string `protobuf:"bytes,2,opt,name=AppId,proto3" json:"AppId,omitempty"`
}

func (x *OnlineDelete_Req) Reset() {
	*x = OnlineDelete_Req{}
	if protoimpl.UnsafeEnabled {
		mi := &file_online_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnlineDelete_Req) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnlineDelete_Req) ProtoMessage() {}

func (x *OnlineDelete_Req) ProtoReflect() protoreflect.Message {
	mi := &file_online_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnlineDelete_Req.ProtoReflect.Descriptor instead.
func (*OnlineDelete_Req) Descriptor() ([]byte, []int) {
	return file_online_proto_rawDescGZIP(), []int{4, 0}
}

func (x *OnlineDelete_Req) GetRoleId() int64 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

func (x *OnlineDelete_Req) GetAppId() string {
	if x != nil {
		return x.AppId
	}
	return ""
}

type OnlineDelete_Rsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *OnlineDelete_Rsp) Reset() {
	*x = OnlineDelete_Rsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_online_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnlineDelete_Rsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnlineDelete_Rsp) ProtoMessage() {}

func (x *OnlineDelete_Rsp) ProtoReflect() protoreflect.Message {
	mi := &file_online_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnlineDelete_Rsp.ProtoReflect.Descriptor instead.
func (*OnlineDelete_Rsp) Descriptor() ([]byte, []int) {
	return file_online_proto_rawDescGZIP(), []int{4, 1}
}

var File_online_proto protoreflect.FileDescriptor

var file_online_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x6f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7b, 0x0a, 0x0c, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x6f,
	0x6c, 0x65, 0x44, 0x42, 0x12, 0x28, 0x0a, 0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x42, 0x10, 0x8a, 0xa6, 0x1d, 0x0c, 0x0a, 0x08, 0x4f, 0x52, 0x52, 0x6f,
	0x6c, 0x65, 0x49, 0x64, 0x10, 0x01, 0x52, 0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x41,
	0x0a, 0x09, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x70, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65,
	0x41, 0x70, 0x70, 0x42, 0x11, 0x8a, 0xa6, 0x1d, 0x0d, 0x0a, 0x0b, 0x4f, 0x52, 0x4f, 0x6e, 0x6c,
	0x69, 0x6e, 0x65, 0x41, 0x70, 0x70, 0x52, 0x09, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x70,
	0x70, 0x22, 0x33, 0x0a, 0x09, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x70, 0x70, 0x12, 0x14,
	0x0a, 0x05, 0x41, 0x70, 0x70, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x41,
	0x70, 0x70, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x56, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x56, 0x65, 0x72, 0x22, 0x52, 0x0a, 0x0c, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x1a, 0x33, 0x0a, 0x03, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a,
	0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x52,
	0x6f, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x41, 0x70, 0x70, 0x49, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x41, 0x70, 0x70, 0x49, 0x64, 0x1a, 0x05, 0x0a, 0x03, 0x52,
	0x73, 0x70, 0x3a, 0x06, 0xca, 0xf3, 0x18, 0x02, 0x0a, 0x00, 0x22, 0x50, 0x0a, 0x0a, 0x4f, 0x6e,
	0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x61, 0x64, 0x1a, 0x1d, 0x0a, 0x03, 0x52, 0x65, 0x71, 0x12,
	0x16, 0x0a, 0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x1a, 0x1b, 0x0a, 0x03, 0x52, 0x73, 0x70, 0x12, 0x14,
	0x0a, 0x05, 0x41, 0x70, 0x70, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x41,
	0x70, 0x70, 0x49, 0x64, 0x3a, 0x06, 0xca, 0xf3, 0x18, 0x02, 0x0a, 0x00, 0x22, 0x52, 0x0a, 0x0c,
	0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x1a, 0x33, 0x0a, 0x03,
	0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x41,
	0x70, 0x70, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x41, 0x70, 0x70, 0x49,
	0x64, 0x1a, 0x05, 0x0a, 0x03, 0x52, 0x73, 0x70, 0x3a, 0x06, 0xca, 0xf3, 0x18, 0x02, 0x0a, 0x00,
	0x2a, 0x3c, 0x0a, 0x06, 0x4f, 0x52, 0x46, 0x6c, 0x61, 0x67, 0x12, 0x13, 0x0a, 0x0f, 0x4f, 0x52,
	0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x4d, 0x69, 0x6e, 0x10, 0x00, 0x12,
	0x0c, 0x0a, 0x08, 0x4f, 0x52, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x10, 0x01, 0x12, 0x0f, 0x0a,
	0x0b, 0x4f, 0x52, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x70, 0x70, 0x10, 0x02, 0x32, 0xcb,
	0x01, 0x0a, 0x09, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x53, 0x65, 0x72, 0x12, 0x40, 0x0a, 0x0c,
	0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x17, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x2e, 0x52, 0x65, 0x71, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x6e,
	0x6c, 0x69, 0x6e, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x2e, 0x52, 0x73, 0x70, 0x12, 0x3a,
	0x0a, 0x0a, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x61, 0x64, 0x12, 0x15, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x61, 0x64, 0x2e,
	0x52, 0x65, 0x71, 0x1a, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x6e, 0x6c, 0x69,
	0x6e, 0x65, 0x52, 0x65, 0x61, 0x64, 0x2e, 0x52, 0x73, 0x70, 0x12, 0x40, 0x0a, 0x0c, 0x4f, 0x6e,
	0x6c, 0x69, 0x6e, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x2e,
	0x52, 0x65, 0x71, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x6e, 0x6c, 0x69,
	0x6e, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x2e, 0x52, 0x73, 0x70, 0x42, 0x0a, 0x5a, 0x08,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_online_proto_rawDescOnce sync.Once
	file_online_proto_rawDescData = file_online_proto_rawDesc
)

func file_online_proto_rawDescGZIP() []byte {
	file_online_proto_rawDescOnce.Do(func() {
		file_online_proto_rawDescData = protoimpl.X.CompressGZIP(file_online_proto_rawDescData)
	})
	return file_online_proto_rawDescData
}

var file_online_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_online_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_online_proto_goTypes = []interface{}{
	(ORFlag)(0),              // 0: proto.ORFlag
	(*OnlineRoleDB)(nil),     // 1: proto.OnlineRoleDB
	(*OnlineApp)(nil),        // 2: proto.OnlineApp
	(*OnlineCreate)(nil),     // 3: proto.OnlineCreate
	(*OnlineRead)(nil),       // 4: proto.OnlineRead
	(*OnlineDelete)(nil),     // 5: proto.OnlineDelete
	(*OnlineCreate_Req)(nil), // 6: proto.OnlineCreate.Req
	(*OnlineCreate_Rsp)(nil), // 7: proto.OnlineCreate.Rsp
	(*OnlineRead_Req)(nil),   // 8: proto.OnlineRead.Req
	(*OnlineRead_Rsp)(nil),   // 9: proto.OnlineRead.Rsp
	(*OnlineDelete_Req)(nil), // 10: proto.OnlineDelete.Req
	(*OnlineDelete_Rsp)(nil), // 11: proto.OnlineDelete.Rsp
}
var file_online_proto_depIdxs = []int32{
	2,  // 0: proto.OnlineRoleDB.OnlineApp:type_name -> proto.OnlineApp
	6,  // 1: proto.OnlineSer.OnlineCreate:input_type -> proto.OnlineCreate.Req
	8,  // 2: proto.OnlineSer.OnlineRead:input_type -> proto.OnlineRead.Req
	10, // 3: proto.OnlineSer.OnlineDelete:input_type -> proto.OnlineDelete.Req
	7,  // 4: proto.OnlineSer.OnlineCreate:output_type -> proto.OnlineCreate.Rsp
	9,  // 5: proto.OnlineSer.OnlineRead:output_type -> proto.OnlineRead.Rsp
	11, // 6: proto.OnlineSer.OnlineDelete:output_type -> proto.OnlineDelete.Rsp
	4,  // [4:7] is the sub-list for method output_type
	1,  // [1:4] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_online_proto_init() }
func file_online_proto_init() {
	if File_online_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_online_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnlineRoleDB); i {
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
		file_online_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnlineApp); i {
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
		file_online_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnlineCreate); i {
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
		file_online_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnlineRead); i {
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
		file_online_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnlineDelete); i {
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
		file_online_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnlineCreate_Req); i {
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
		file_online_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnlineCreate_Rsp); i {
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
		file_online_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnlineRead_Req); i {
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
		file_online_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnlineRead_Rsp); i {
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
		file_online_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnlineDelete_Req); i {
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
		file_online_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnlineDelete_Rsp); i {
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
			RawDescriptor: file_online_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_online_proto_goTypes,
		DependencyIndexes: file_online_proto_depIdxs,
		EnumInfos:         file_online_proto_enumTypes,
		MessageInfos:      file_online_proto_msgTypes,
	}.Build()
	File_online_proto = out.File
	file_online_proto_rawDesc = nil
	file_online_proto_goTypes = nil
	file_online_proto_depIdxs = nil
}
