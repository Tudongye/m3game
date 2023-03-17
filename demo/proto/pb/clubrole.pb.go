// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.2.0
// source: clubrole.proto

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

type ClubRoleRead struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ClubRoleRead) Reset() {
	*x = ClubRoleRead{}
	if protoimpl.UnsafeEnabled {
		mi := &file_clubrole_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClubRoleRead) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClubRoleRead) ProtoMessage() {}

func (x *ClubRoleRead) ProtoReflect() protoreflect.Message {
	mi := &file_clubrole_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClubRoleRead.ProtoReflect.Descriptor instead.
func (*ClubRoleRead) Descriptor() ([]byte, []int) {
	return file_clubrole_proto_rawDescGZIP(), []int{0}
}

type ClubRoleCreate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ClubRoleCreate) Reset() {
	*x = ClubRoleCreate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_clubrole_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClubRoleCreate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClubRoleCreate) ProtoMessage() {}

func (x *ClubRoleCreate) ProtoReflect() protoreflect.Message {
	mi := &file_clubrole_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClubRoleCreate.ProtoReflect.Descriptor instead.
func (*ClubRoleCreate) Descriptor() ([]byte, []int) {
	return file_clubrole_proto_rawDescGZIP(), []int{1}
}

type ClubRoleDelete struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ClubRoleDelete) Reset() {
	*x = ClubRoleDelete{}
	if protoimpl.UnsafeEnabled {
		mi := &file_clubrole_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClubRoleDelete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClubRoleDelete) ProtoMessage() {}

func (x *ClubRoleDelete) ProtoReflect() protoreflect.Message {
	mi := &file_clubrole_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClubRoleDelete.ProtoReflect.Descriptor instead.
func (*ClubRoleDelete) Descriptor() ([]byte, []int) {
	return file_clubrole_proto_rawDescGZIP(), []int{2}
}

type ClubRoleRead_Req struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoleId string `protobuf:"bytes,1,opt,name=RoleId,proto3" json:"RoleId,omitempty"`
}

func (x *ClubRoleRead_Req) Reset() {
	*x = ClubRoleRead_Req{}
	if protoimpl.UnsafeEnabled {
		mi := &file_clubrole_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClubRoleRead_Req) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClubRoleRead_Req) ProtoMessage() {}

func (x *ClubRoleRead_Req) ProtoReflect() protoreflect.Message {
	mi := &file_clubrole_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClubRoleRead_Req.ProtoReflect.Descriptor instead.
func (*ClubRoleRead_Req) Descriptor() ([]byte, []int) {
	return file_clubrole_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ClubRoleRead_Req) GetRoleId() string {
	if x != nil {
		return x.RoleId
	}
	return ""
}

type ClubRoleRead_Rsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClubRoleRelationDB *ClubRoleRelationDB `protobuf:"bytes,1,opt,name=ClubRoleRelationDB,proto3" json:"ClubRoleRelationDB,omitempty"`
}

func (x *ClubRoleRead_Rsp) Reset() {
	*x = ClubRoleRead_Rsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_clubrole_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClubRoleRead_Rsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClubRoleRead_Rsp) ProtoMessage() {}

func (x *ClubRoleRead_Rsp) ProtoReflect() protoreflect.Message {
	mi := &file_clubrole_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClubRoleRead_Rsp.ProtoReflect.Descriptor instead.
func (*ClubRoleRead_Rsp) Descriptor() ([]byte, []int) {
	return file_clubrole_proto_rawDescGZIP(), []int{0, 1}
}

func (x *ClubRoleRead_Rsp) GetClubRoleRelationDB() *ClubRoleRelationDB {
	if x != nil {
		return x.ClubRoleRelationDB
	}
	return nil
}

type ClubRoleCreate_Req struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClubId string `protobuf:"bytes,1,opt,name=ClubId,proto3" json:"ClubId,omitempty"`
	RoleId string `protobuf:"bytes,2,opt,name=RoleId,proto3" json:"RoleId,omitempty"`
}

func (x *ClubRoleCreate_Req) Reset() {
	*x = ClubRoleCreate_Req{}
	if protoimpl.UnsafeEnabled {
		mi := &file_clubrole_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClubRoleCreate_Req) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClubRoleCreate_Req) ProtoMessage() {}

func (x *ClubRoleCreate_Req) ProtoReflect() protoreflect.Message {
	mi := &file_clubrole_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClubRoleCreate_Req.ProtoReflect.Descriptor instead.
func (*ClubRoleCreate_Req) Descriptor() ([]byte, []int) {
	return file_clubrole_proto_rawDescGZIP(), []int{1, 0}
}

func (x *ClubRoleCreate_Req) GetClubId() string {
	if x != nil {
		return x.ClubId
	}
	return ""
}

func (x *ClubRoleCreate_Req) GetRoleId() string {
	if x != nil {
		return x.RoleId
	}
	return ""
}

type ClubRoleCreate_Rsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ClubRoleCreate_Rsp) Reset() {
	*x = ClubRoleCreate_Rsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_clubrole_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClubRoleCreate_Rsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClubRoleCreate_Rsp) ProtoMessage() {}

func (x *ClubRoleCreate_Rsp) ProtoReflect() protoreflect.Message {
	mi := &file_clubrole_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClubRoleCreate_Rsp.ProtoReflect.Descriptor instead.
func (*ClubRoleCreate_Rsp) Descriptor() ([]byte, []int) {
	return file_clubrole_proto_rawDescGZIP(), []int{1, 1}
}

type ClubRoleDelete_Req struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClubId string `protobuf:"bytes,1,opt,name=ClubId,proto3" json:"ClubId,omitempty"`
	RoleId string `protobuf:"bytes,2,opt,name=RoleId,proto3" json:"RoleId,omitempty"`
}

func (x *ClubRoleDelete_Req) Reset() {
	*x = ClubRoleDelete_Req{}
	if protoimpl.UnsafeEnabled {
		mi := &file_clubrole_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClubRoleDelete_Req) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClubRoleDelete_Req) ProtoMessage() {}

func (x *ClubRoleDelete_Req) ProtoReflect() protoreflect.Message {
	mi := &file_clubrole_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClubRoleDelete_Req.ProtoReflect.Descriptor instead.
func (*ClubRoleDelete_Req) Descriptor() ([]byte, []int) {
	return file_clubrole_proto_rawDescGZIP(), []int{2, 0}
}

func (x *ClubRoleDelete_Req) GetClubId() string {
	if x != nil {
		return x.ClubId
	}
	return ""
}

func (x *ClubRoleDelete_Req) GetRoleId() string {
	if x != nil {
		return x.RoleId
	}
	return ""
}

type ClubRoleDelete_Rsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ClubRoleDelete_Rsp) Reset() {
	*x = ClubRoleDelete_Rsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_clubrole_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClubRoleDelete_Rsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClubRoleDelete_Rsp) ProtoMessage() {}

func (x *ClubRoleDelete_Rsp) ProtoReflect() protoreflect.Message {
	mi := &file_clubrole_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClubRoleDelete_Rsp.ProtoReflect.Descriptor instead.
func (*ClubRoleDelete_Rsp) Descriptor() ([]byte, []int) {
	return file_clubrole_proto_rawDescGZIP(), []int{2, 1}
}

var File_clubrole_proto protoreflect.FileDescriptor

var file_clubrole_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x63, 0x6c, 0x75, 0x62, 0x72, 0x6f, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x09, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x8d, 0x01, 0x0a, 0x0c, 0x43, 0x6c, 0x75, 0x62, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65,
	0x61, 0x64, 0x1a, 0x1d, 0x0a, 0x03, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x6f, 0x6c,
	0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49,
	0x64, 0x1a, 0x50, 0x0a, 0x03, 0x52, 0x73, 0x70, 0x12, 0x49, 0x0a, 0x12, 0x43, 0x6c, 0x75, 0x62,
	0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x42, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6c, 0x75,
	0x62, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x42, 0x52,
	0x12, 0x43, 0x6c, 0x75, 0x62, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x44, 0x42, 0x3a, 0x0c, 0xca, 0xf3, 0x18, 0x02, 0x0a, 0x00, 0xca, 0xf3, 0x18, 0x02, 0x10,
	0x01, 0x22, 0x5c, 0x0a, 0x0e, 0x43, 0x6c, 0x75, 0x62, 0x52, 0x6f, 0x6c, 0x65, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x1a, 0x35, 0x0a, 0x03, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x43, 0x6c,
	0x75, 0x62, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x43, 0x6c, 0x75, 0x62,
	0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x1a, 0x05, 0x0a, 0x03, 0x52, 0x73,
	0x70, 0x3a, 0x0c, 0xca, 0xf3, 0x18, 0x02, 0x0a, 0x00, 0xca, 0xf3, 0x18, 0x02, 0x10, 0x01, 0x22,
	0x5c, 0x0a, 0x0e, 0x43, 0x6c, 0x75, 0x62, 0x52, 0x6f, 0x6c, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x1a, 0x35, 0x0a, 0x03, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x43, 0x6c, 0x75, 0x62,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x43, 0x6c, 0x75, 0x62, 0x49, 0x64,
	0x12, 0x16, 0x0a, 0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x1a, 0x05, 0x0a, 0x03, 0x52, 0x73, 0x70, 0x3a,
	0x0c, 0xca, 0xf3, 0x18, 0x02, 0x0a, 0x00, 0xca, 0xf3, 0x18, 0x02, 0x10, 0x01, 0x32, 0xdf, 0x01,
	0x0a, 0x0b, 0x43, 0x6c, 0x75, 0x62, 0x52, 0x6f, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x12, 0x40, 0x0a,
	0x0c, 0x43, 0x6c, 0x75, 0x62, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x61, 0x64, 0x12, 0x17, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6c, 0x75, 0x62, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65,
	0x61, 0x64, 0x2e, 0x52, 0x65, 0x71, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43,
	0x6c, 0x75, 0x62, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x61, 0x64, 0x2e, 0x52, 0x73, 0x70, 0x12,
	0x46, 0x0a, 0x0e, 0x43, 0x6c, 0x75, 0x62, 0x52, 0x6f, 0x6c, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x12, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6c, 0x75, 0x62, 0x52, 0x6f,
	0x6c, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x1a, 0x19, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6c, 0x75, 0x62, 0x52, 0x6f, 0x6c, 0x65, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x2e, 0x52, 0x73, 0x70, 0x12, 0x46, 0x0a, 0x0e, 0x43, 0x6c, 0x75, 0x62, 0x52,
	0x6f, 0x6c, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x43, 0x6c, 0x75, 0x62, 0x52, 0x6f, 0x6c, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x2e, 0x52, 0x65, 0x71, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6c, 0x75,
	0x62, 0x52, 0x6f, 0x6c, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x2e, 0x52, 0x73, 0x70, 0x42,
	0x0a, 0x5a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_clubrole_proto_rawDescOnce sync.Once
	file_clubrole_proto_rawDescData = file_clubrole_proto_rawDesc
)

func file_clubrole_proto_rawDescGZIP() []byte {
	file_clubrole_proto_rawDescOnce.Do(func() {
		file_clubrole_proto_rawDescData = protoimpl.X.CompressGZIP(file_clubrole_proto_rawDescData)
	})
	return file_clubrole_proto_rawDescData
}

var file_clubrole_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_clubrole_proto_goTypes = []interface{}{
	(*ClubRoleRead)(nil),       // 0: proto.ClubRoleRead
	(*ClubRoleCreate)(nil),     // 1: proto.ClubRoleCreate
	(*ClubRoleDelete)(nil),     // 2: proto.ClubRoleDelete
	(*ClubRoleRead_Req)(nil),   // 3: proto.ClubRoleRead.Req
	(*ClubRoleRead_Rsp)(nil),   // 4: proto.ClubRoleRead.Rsp
	(*ClubRoleCreate_Req)(nil), // 5: proto.ClubRoleCreate.Req
	(*ClubRoleCreate_Rsp)(nil), // 6: proto.ClubRoleCreate.Rsp
	(*ClubRoleDelete_Req)(nil), // 7: proto.ClubRoleDelete.Req
	(*ClubRoleDelete_Rsp)(nil), // 8: proto.ClubRoleDelete.Rsp
	(*ClubRoleRelationDB)(nil), // 9: proto.ClubRoleRelationDB
}
var file_clubrole_proto_depIdxs = []int32{
	9, // 0: proto.ClubRoleRead.Rsp.ClubRoleRelationDB:type_name -> proto.ClubRoleRelationDB
	3, // 1: proto.ClubRoleSer.ClubRoleRead:input_type -> proto.ClubRoleRead.Req
	5, // 2: proto.ClubRoleSer.ClubRoleCreate:input_type -> proto.ClubRoleCreate.Req
	7, // 3: proto.ClubRoleSer.ClubRoleDelete:input_type -> proto.ClubRoleDelete.Req
	4, // 4: proto.ClubRoleSer.ClubRoleRead:output_type -> proto.ClubRoleRead.Rsp
	6, // 5: proto.ClubRoleSer.ClubRoleCreate:output_type -> proto.ClubRoleCreate.Rsp
	8, // 6: proto.ClubRoleSer.ClubRoleDelete:output_type -> proto.ClubRoleDelete.Rsp
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_clubrole_proto_init() }
func file_clubrole_proto_init() {
	if File_clubrole_proto != nil {
		return
	}
	file_com_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_clubrole_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClubRoleRead); i {
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
		file_clubrole_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClubRoleCreate); i {
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
		file_clubrole_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClubRoleDelete); i {
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
		file_clubrole_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClubRoleRead_Req); i {
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
		file_clubrole_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClubRoleRead_Rsp); i {
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
		file_clubrole_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClubRoleCreate_Req); i {
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
		file_clubrole_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClubRoleCreate_Rsp); i {
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
		file_clubrole_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClubRoleDelete_Req); i {
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
		file_clubrole_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClubRoleDelete_Rsp); i {
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
			RawDescriptor: file_clubrole_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_clubrole_proto_goTypes,
		DependencyIndexes: file_clubrole_proto_depIdxs,
		MessageInfos:      file_clubrole_proto_msgTypes,
	}.Build()
	File_clubrole_proto = out.File
	file_clubrole_proto_rawDesc = nil
	file_clubrole_proto_goTypes = nil
	file_clubrole_proto_depIdxs = nil
}
