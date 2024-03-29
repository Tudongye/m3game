// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.2.0
// source: options.proto

package metapb

import (
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
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

type M3RpcOption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RouteKey string `protobuf:"bytes,1,opt,name=route_key,json=routeKey,proto3" json:"route_key,omitempty"`
	Ntf      bool   `protobuf:"varint,2,opt,name=ntf,proto3" json:"ntf,omitempty"`
	Trace    bool   `protobuf:"varint,3,opt,name=trace,proto3" json:"trace,omitempty"`
	Cs       bool   `protobuf:"varint,4,opt,name=cs,proto3" json:"cs,omitempty"`
}

func (x *M3RpcOption) Reset() {
	*x = M3RpcOption{}
	if protoimpl.UnsafeEnabled {
		mi := &file_options_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *M3RpcOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*M3RpcOption) ProtoMessage() {}

func (x *M3RpcOption) ProtoReflect() protoreflect.Message {
	mi := &file_options_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use M3RpcOption.ProtoReflect.Descriptor instead.
func (*M3RpcOption) Descriptor() ([]byte, []int) {
	return file_options_proto_rawDescGZIP(), []int{0}
}

func (x *M3RpcOption) GetRouteKey() string {
	if x != nil {
		return x.RouteKey
	}
	return ""
}

func (x *M3RpcOption) GetNtf() bool {
	if x != nil {
		return x.Ntf
	}
	return false
}

func (x *M3RpcOption) GetTrace() bool {
	if x != nil {
		return x.Trace
	}
	return false
}

func (x *M3RpcOption) GetCs() bool {
	if x != nil {
		return x.Cs
	}
	return false
}

type M3DBFieldOption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Flag    string `protobuf:"bytes,1,opt,name=flag,proto3" json:"flag,omitempty"`        // 字段标记
	Primary bool   `protobuf:"varint,2,opt,name=primary,proto3" json:"primary,omitempty"` // 主键
}

func (x *M3DBFieldOption) Reset() {
	*x = M3DBFieldOption{}
	if protoimpl.UnsafeEnabled {
		mi := &file_options_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *M3DBFieldOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*M3DBFieldOption) ProtoMessage() {}

func (x *M3DBFieldOption) ProtoReflect() protoreflect.Message {
	mi := &file_options_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use M3DBFieldOption.ProtoReflect.Descriptor instead.
func (*M3DBFieldOption) Descriptor() ([]byte, []int) {
	return file_options_proto_rawDescGZIP(), []int{1}
}

func (x *M3DBFieldOption) GetFlag() string {
	if x != nil {
		return x.Flag
	}
	return ""
}

func (x *M3DBFieldOption) GetPrimary() bool {
	if x != nil {
		return x.Primary
	}
	return false
}

type M3ViewFieldOption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Wflag string `protobuf:"bytes,1,opt,name=wflag,proto3" json:"wflag,omitempty"` // 视野白名单
	Bflag string `protobuf:"bytes,2,opt,name=bflag,proto3" json:"bflag,omitempty"` // 视野黑名单
}

func (x *M3ViewFieldOption) Reset() {
	*x = M3ViewFieldOption{}
	if protoimpl.UnsafeEnabled {
		mi := &file_options_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *M3ViewFieldOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*M3ViewFieldOption) ProtoMessage() {}

func (x *M3ViewFieldOption) ProtoReflect() protoreflect.Message {
	mi := &file_options_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use M3ViewFieldOption.ProtoReflect.Descriptor instead.
func (*M3ViewFieldOption) Descriptor() ([]byte, []int) {
	return file_options_proto_rawDescGZIP(), []int{2}
}

func (x *M3ViewFieldOption) GetWflag() string {
	if x != nil {
		return x.Wflag
	}
	return ""
}

func (x *M3ViewFieldOption) GetBflag() string {
	if x != nil {
		return x.Bflag
	}
	return ""
}

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
		mi := &file_options_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Meta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Meta) ProtoMessage() {}

func (x *Meta) ProtoReflect() protoreflect.Message {
	mi := &file_options_proto_msgTypes[3]
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
	return file_options_proto_rawDescGZIP(), []int{3}
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

var file_options_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptor.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         50001,
		Name:          "proto.db_primary_key",
		Tag:           "bytes,50001,opt,name=db_primary_key",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptor.MessageOptions)(nil),
		ExtensionType: (*M3RpcOption)(nil),
		Field:         51001,
		Name:          "proto.rpc_option",
		Tag:           "bytes,51001,opt,name=rpc_option",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptor.FieldOptions)(nil),
		ExtensionType: (*M3DBFieldOption)(nil),
		Field:         60001,
		Name:          "proto.dbfield_option",
		Tag:           "bytes,60001,opt,name=dbfield_option",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptor.FieldOptions)(nil),
		ExtensionType: (*M3ViewFieldOption)(nil),
		Field:         60002,
		Name:          "proto.viewfield_option",
		Tag:           "bytes,60002,opt,name=viewfield_option",
		Filename:      "options.proto",
	},
}

// Extension fields to descriptor.MessageOptions.
var (
	// optional string db_primary_key = 50001;
	E_DbPrimaryKey = &file_options_proto_extTypes[0]
	// optional proto.M3RpcOption rpc_option = 51001;
	E_RpcOption = &file_options_proto_extTypes[1]
)

// Extension fields to descriptor.FieldOptions.
var (
	// optional proto.M3DBFieldOption dbfield_option = 60001;
	E_DbfieldOption = &file_options_proto_extTypes[2]
	// optional proto.M3ViewFieldOption viewfield_option = 60002;
	E_ViewfieldOption = &file_options_proto_extTypes[3]
)

var File_options_proto protoreflect.FileDescriptor

var file_options_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x62, 0x0a, 0x0b, 0x4d, 0x33, 0x52, 0x70,
	0x63, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x6f, 0x75, 0x74, 0x65,
	0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x6f, 0x75, 0x74,
	0x65, 0x4b, 0x65, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6e, 0x74, 0x66, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x03, 0x6e, 0x74, 0x66, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x72, 0x61, 0x63, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x74, 0x72, 0x61, 0x63, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x63, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x63, 0x73, 0x22, 0x3f, 0x0a, 0x0f,
	0x4d, 0x33, 0x44, 0x42, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x12, 0x0a, 0x04, 0x66, 0x6c, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66,
	0x6c, 0x61, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x22, 0x3f, 0x0a,
	0x11, 0x4d, 0x33, 0x56, 0x69, 0x65, 0x77, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x66, 0x6c, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x77, 0x66, 0x6c, 0x61, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x66, 0x6c, 0x61,
	0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x62, 0x66, 0x6c, 0x61, 0x67, 0x22, 0x2e,
	0x0a, 0x04, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x47,
	0x0a, 0x0e, 0x64, 0x62, 0x5f, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x5f, 0x6b, 0x65, 0x79,
	0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0xd1, 0x86, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x64, 0x62, 0x50, 0x72, 0x69,
	0x6d, 0x61, 0x72, 0x79, 0x4b, 0x65, 0x79, 0x3a, 0x54, 0x0a, 0x0a, 0x72, 0x70, 0x63, 0x5f, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xb9, 0x8e, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x33, 0x52, 0x70, 0x63, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x09, 0x72, 0x70, 0x63, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x5e, 0x0a,
	0x0e, 0x64, 0x62, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xe1,
	0xd4, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d,
	0x33, 0x44, 0x42, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d,
	0x64, 0x62, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x64, 0x0a,
	0x10, 0x76, 0x69, 0x65, 0x77, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0xe2, 0xd4, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4d, 0x33, 0x56, 0x69, 0x65, 0x77, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x0f, 0x76, 0x69, 0x65, 0x77, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x42, 0x14, 0x5a, 0x12, 0x6d, 0x33, 0x67, 0x61, 0x6d, 0x65, 0x2f, 0x6d, 0x65,
	0x74, 0x61, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_options_proto_rawDescOnce sync.Once
	file_options_proto_rawDescData = file_options_proto_rawDesc
)

func file_options_proto_rawDescGZIP() []byte {
	file_options_proto_rawDescOnce.Do(func() {
		file_options_proto_rawDescData = protoimpl.X.CompressGZIP(file_options_proto_rawDescData)
	})
	return file_options_proto_rawDescData
}

var file_options_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_options_proto_goTypes = []interface{}{
	(*M3RpcOption)(nil),               // 0: proto.M3RpcOption
	(*M3DBFieldOption)(nil),           // 1: proto.M3DBFieldOption
	(*M3ViewFieldOption)(nil),         // 2: proto.M3ViewFieldOption
	(*Meta)(nil),                      // 3: proto.Meta
	(*descriptor.MessageOptions)(nil), // 4: google.protobuf.MessageOptions
	(*descriptor.FieldOptions)(nil),   // 5: google.protobuf.FieldOptions
}
var file_options_proto_depIdxs = []int32{
	4, // 0: proto.db_primary_key:extendee -> google.protobuf.MessageOptions
	4, // 1: proto.rpc_option:extendee -> google.protobuf.MessageOptions
	5, // 2: proto.dbfield_option:extendee -> google.protobuf.FieldOptions
	5, // 3: proto.viewfield_option:extendee -> google.protobuf.FieldOptions
	0, // 4: proto.rpc_option:type_name -> proto.M3RpcOption
	1, // 5: proto.dbfield_option:type_name -> proto.M3DBFieldOption
	2, // 6: proto.viewfield_option:type_name -> proto.M3ViewFieldOption
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	4, // [4:7] is the sub-list for extension type_name
	0, // [0:4] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_options_proto_init() }
func file_options_proto_init() {
	if File_options_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_options_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*M3RpcOption); i {
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
		file_options_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*M3DBFieldOption); i {
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
		file_options_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*M3ViewFieldOption); i {
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
		file_options_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_options_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 4,
			NumServices:   0,
		},
		GoTypes:           file_options_proto_goTypes,
		DependencyIndexes: file_options_proto_depIdxs,
		MessageInfos:      file_options_proto_msgTypes,
		ExtensionInfos:    file_options_proto_extTypes,
	}.Build()
	File_options_proto = out.File
	file_options_proto_rawDesc = nil
	file_options_proto_goTypes = nil
	file_options_proto_depIdxs = nil
}
