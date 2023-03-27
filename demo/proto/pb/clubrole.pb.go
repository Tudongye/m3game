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

type ClubRoleRead_Req struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoleId int64 `protobuf:"varint,1,opt,name=RoleId,proto3" json:"RoleId,omitempty"`
}

func (x *ClubRoleRead_Req) Reset() {
	*x = ClubRoleRead_Req{}
	if protoimpl.UnsafeEnabled {
		mi := &file_clubrole_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClubRoleRead_Req) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClubRoleRead_Req) ProtoMessage() {}

func (x *ClubRoleRead_Req) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use ClubRoleRead_Req.ProtoReflect.Descriptor instead.
func (*ClubRoleRead_Req) Descriptor() ([]byte, []int) {
	return file_clubrole_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ClubRoleRead_Req) GetRoleId() int64 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

type ClubRoleRead_Rsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClubId int64 `protobuf:"varint,1,opt,name=ClubId,proto3" json:"ClubId,omitempty"`
}

func (x *ClubRoleRead_Rsp) Reset() {
	*x = ClubRoleRead_Rsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_clubrole_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClubRoleRead_Rsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClubRoleRead_Rsp) ProtoMessage() {}

func (x *ClubRoleRead_Rsp) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use ClubRoleRead_Rsp.ProtoReflect.Descriptor instead.
func (*ClubRoleRead_Rsp) Descriptor() ([]byte, []int) {
	return file_clubrole_proto_rawDescGZIP(), []int{0, 1}
}

func (x *ClubRoleRead_Rsp) GetClubId() int64 {
	if x != nil {
		return x.ClubId
	}
	return 0
}

var File_clubrole_proto protoreflect.FileDescriptor

var file_clubrole_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x63, 0x6c, 0x75, 0x62, 0x72, 0x6f, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5a, 0x0a, 0x0c, 0x43, 0x6c, 0x75, 0x62, 0x52, 0x6f,
	0x6c, 0x65, 0x52, 0x65, 0x61, 0x64, 0x1a, 0x1d, 0x0a, 0x03, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a,
	0x06, 0x52, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x52,
	0x6f, 0x6c, 0x65, 0x49, 0x64, 0x1a, 0x1d, 0x0a, 0x03, 0x52, 0x73, 0x70, 0x12, 0x16, 0x0a, 0x06,
	0x43, 0x6c, 0x75, 0x62, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x43, 0x6c,
	0x75, 0x62, 0x49, 0x64, 0x3a, 0x0c, 0xca, 0xf3, 0x18, 0x02, 0x0a, 0x00, 0xca, 0xf3, 0x18, 0x02,
	0x10, 0x01, 0x32, 0x4f, 0x0a, 0x0b, 0x43, 0x6c, 0x75, 0x62, 0x52, 0x6f, 0x6c, 0x65, 0x53, 0x65,
	0x72, 0x12, 0x40, 0x0a, 0x0c, 0x43, 0x6c, 0x75, 0x62, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x61,
	0x64, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6c, 0x75, 0x62, 0x52, 0x6f,
	0x6c, 0x65, 0x52, 0x65, 0x61, 0x64, 0x2e, 0x52, 0x65, 0x71, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x43, 0x6c, 0x75, 0x62, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x61, 0x64, 0x2e,
	0x52, 0x73, 0x70, 0x42, 0x0a, 0x5a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_clubrole_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_clubrole_proto_goTypes = []interface{}{
	(*ClubRoleRead)(nil),     // 0: proto.ClubRoleRead
	(*ClubRoleRead_Req)(nil), // 1: proto.ClubRoleRead.Req
	(*ClubRoleRead_Rsp)(nil), // 2: proto.ClubRoleRead.Rsp
}
var file_clubrole_proto_depIdxs = []int32{
	1, // 0: proto.ClubRoleSer.ClubRoleRead:input_type -> proto.ClubRoleRead.Req
	2, // 1: proto.ClubRoleSer.ClubRoleRead:output_type -> proto.ClubRoleRead.Rsp
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_clubrole_proto_init() }
func file_clubrole_proto_init() {
	if File_clubrole_proto != nil {
		return
	}
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
		file_clubrole_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_clubrole_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
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
