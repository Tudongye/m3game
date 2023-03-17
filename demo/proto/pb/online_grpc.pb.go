// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.2.0
// source: online.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// OnlineSerClient is the client API for OnlineSer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OnlineSerClient interface {
	OnlineCreate(ctx context.Context, in *OnlineCreate_Req, opts ...grpc.CallOption) (*OnlineCreate_Rsp, error)
	OnlineRead(ctx context.Context, in *OnlineRead_Req, opts ...grpc.CallOption) (*OnlineRead_Rsp, error)
	OnlineDelete(ctx context.Context, in *OnlineDelete_Req, opts ...grpc.CallOption) (*OnlineDelete_Rsp, error)
}

type onlineSerClient struct {
	cc grpc.ClientConnInterface
}

func NewOnlineSerClient(cc grpc.ClientConnInterface) OnlineSerClient {
	return &onlineSerClient{cc}
}

func (c *onlineSerClient) OnlineCreate(ctx context.Context, in *OnlineCreate_Req, opts ...grpc.CallOption) (*OnlineCreate_Rsp, error) {
	out := new(OnlineCreate_Rsp)
	err := c.cc.Invoke(ctx, "/proto.OnlineSer/OnlineCreate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *onlineSerClient) OnlineRead(ctx context.Context, in *OnlineRead_Req, opts ...grpc.CallOption) (*OnlineRead_Rsp, error) {
	out := new(OnlineRead_Rsp)
	err := c.cc.Invoke(ctx, "/proto.OnlineSer/OnlineRead", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *onlineSerClient) OnlineDelete(ctx context.Context, in *OnlineDelete_Req, opts ...grpc.CallOption) (*OnlineDelete_Rsp, error) {
	out := new(OnlineDelete_Rsp)
	err := c.cc.Invoke(ctx, "/proto.OnlineSer/OnlineDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OnlineSerServer is the server API for OnlineSer service.
// All implementations must embed UnimplementedOnlineSerServer
// for forward compatibility
type OnlineSerServer interface {
	OnlineCreate(context.Context, *OnlineCreate_Req) (*OnlineCreate_Rsp, error)
	OnlineRead(context.Context, *OnlineRead_Req) (*OnlineRead_Rsp, error)
	OnlineDelete(context.Context, *OnlineDelete_Req) (*OnlineDelete_Rsp, error)
	mustEmbedUnimplementedOnlineSerServer()
}

// UnimplementedOnlineSerServer must be embedded to have forward compatible implementations.
type UnimplementedOnlineSerServer struct {
}

func (UnimplementedOnlineSerServer) OnlineCreate(context.Context, *OnlineCreate_Req) (*OnlineCreate_Rsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OnlineCreate not implemented")
}
func (UnimplementedOnlineSerServer) OnlineRead(context.Context, *OnlineRead_Req) (*OnlineRead_Rsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OnlineRead not implemented")
}
func (UnimplementedOnlineSerServer) OnlineDelete(context.Context, *OnlineDelete_Req) (*OnlineDelete_Rsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OnlineDelete not implemented")
}
func (UnimplementedOnlineSerServer) mustEmbedUnimplementedOnlineSerServer() {}

// UnsafeOnlineSerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OnlineSerServer will
// result in compilation errors.
type UnsafeOnlineSerServer interface {
	mustEmbedUnimplementedOnlineSerServer()
}

func RegisterOnlineSerServer(s grpc.ServiceRegistrar, srv OnlineSerServer) {
	s.RegisterService(&OnlineSer_ServiceDesc, srv)
}

func _OnlineSer_OnlineCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OnlineCreate_Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OnlineSerServer).OnlineCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.OnlineSer/OnlineCreate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OnlineSerServer).OnlineCreate(ctx, req.(*OnlineCreate_Req))
	}
	return interceptor(ctx, in, info, handler)
}

func _OnlineSer_OnlineRead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OnlineRead_Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OnlineSerServer).OnlineRead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.OnlineSer/OnlineRead",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OnlineSerServer).OnlineRead(ctx, req.(*OnlineRead_Req))
	}
	return interceptor(ctx, in, info, handler)
}

func _OnlineSer_OnlineDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OnlineDelete_Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OnlineSerServer).OnlineDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.OnlineSer/OnlineDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OnlineSerServer).OnlineDelete(ctx, req.(*OnlineDelete_Req))
	}
	return interceptor(ctx, in, info, handler)
}

// OnlineSer_ServiceDesc is the grpc.ServiceDesc for OnlineSer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OnlineSer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.OnlineSer",
	HandlerType: (*OnlineSerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OnlineCreate",
			Handler:    _OnlineSer_OnlineCreate_Handler,
		},
		{
			MethodName: "OnlineRead",
			Handler:    _OnlineSer_OnlineRead_Handler,
		},
		{
			MethodName: "OnlineDelete",
			Handler:    _OnlineSer_OnlineDelete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "online.proto",
}
