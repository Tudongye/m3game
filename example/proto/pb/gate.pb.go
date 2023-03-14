// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.2.0
// source: gate.proto

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

type SendToCli struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendToCli) Reset() {
	*x = SendToCli{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gate_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendToCli) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendToCli) ProtoMessage() {}

func (x *SendToCli) ProtoReflect() protoreflect.Message {
	mi := &file_gate_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendToCli.ProtoReflect.Descriptor instead.
func (*SendToCli) Descriptor() ([]byte, []int) {
	return file_gate_proto_rawDescGZIP(), []int{0}
}

type SendToCli_Req struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerID string `protobuf:"bytes,1,opt,name=PlayerID,proto3" json:"PlayerID,omitempty"`
	Content  string `protobuf:"bytes,2,opt,name=Content,proto3" json:"Content,omitempty"`
}

func (x *SendToCli_Req) Reset() {
	*x = SendToCli_Req{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gate_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendToCli_Req) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendToCli_Req) ProtoMessage() {}

func (x *SendToCli_Req) ProtoReflect() protoreflect.Message {
	mi := &file_gate_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendToCli_Req.ProtoReflect.Descriptor instead.
func (*SendToCli_Req) Descriptor() ([]byte, []int) {
	return file_gate_proto_rawDescGZIP(), []int{0, 0}
}

func (x *SendToCli_Req) GetPlayerID() string {
	if x != nil {
		return x.PlayerID
	}
	return ""
}

func (x *SendToCli_Req) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type SendToCli_Rsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendToCli_Rsp) Reset() {
	*x = SendToCli_Rsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gate_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendToCli_Rsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendToCli_Rsp) ProtoMessage() {}

func (x *SendToCli_Rsp) ProtoReflect() protoreflect.Message {
	mi := &file_gate_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendToCli_Rsp.ProtoReflect.Descriptor instead.
func (*SendToCli_Rsp) Descriptor() ([]byte, []int) {
	return file_gate_proto_rawDescGZIP(), []int{0, 1}
}

var File_gate_proto protoreflect.FileDescriptor

var file_gate_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x57, 0x0a, 0x09, 0x53, 0x65, 0x6e, 0x64, 0x54, 0x6f, 0x43, 0x6c, 0x69, 0x1a,
	0x3b, 0x0a, 0x03, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x49, 0x44, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x1a, 0x05, 0x0a, 0x03,
	0x52, 0x73, 0x70, 0x3a, 0x06, 0xca, 0xf3, 0x18, 0x02, 0x0a, 0x00, 0x32, 0x42, 0x0a, 0x07, 0x47,
	0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x12, 0x37, 0x0a, 0x09, 0x53, 0x65, 0x6e, 0x64, 0x54, 0x6f,
	0x43, 0x6c, 0x69, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x6e, 0x64,
	0x54, 0x6f, 0x43, 0x6c, 0x69, 0x2e, 0x52, 0x65, 0x71, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x54, 0x6f, 0x43, 0x6c, 0x69, 0x2e, 0x52, 0x73, 0x70, 0x42,
	0x0a, 0x5a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_gate_proto_rawDescOnce sync.Once
	file_gate_proto_rawDescData = file_gate_proto_rawDesc
)

func file_gate_proto_rawDescGZIP() []byte {
	file_gate_proto_rawDescOnce.Do(func() {
		file_gate_proto_rawDescData = protoimpl.X.CompressGZIP(file_gate_proto_rawDescData)
	})
	return file_gate_proto_rawDescData
}

var file_gate_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_gate_proto_goTypes = []interface{}{
	(*SendToCli)(nil),     // 0: proto.SendToCli
	(*SendToCli_Req)(nil), // 1: proto.SendToCli.Req
	(*SendToCli_Rsp)(nil), // 2: proto.SendToCli.Rsp
}
var file_gate_proto_depIdxs = []int32{
	1, // 0: proto.GateSer.SendToCli:input_type -> proto.SendToCli.Req
	2, // 1: proto.GateSer.SendToCli:output_type -> proto.SendToCli.Rsp
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_gate_proto_init() }
func file_gate_proto_init() {
	if File_gate_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gate_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendToCli); i {
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
		file_gate_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendToCli_Req); i {
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
		file_gate_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendToCli_Rsp); i {
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
			RawDescriptor: file_gate_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gate_proto_goTypes,
		DependencyIndexes: file_gate_proto_depIdxs,
		MessageInfos:      file_gate_proto_msgTypes,
	}.Build()
	File_gate_proto = out.File
	file_gate_proto_rawDesc = nil
	file_gate_proto_goTypes = nil
	file_gate_proto_depIdxs = nil
}
