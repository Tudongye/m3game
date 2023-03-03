// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.2.0
// source: map.proto

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

// MapSerClient is the client API for MapSer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MapSerClient interface {
	Move(ctx context.Context, in *Move_Req, opts ...grpc.CallOption) (*Move_Rsp, error)
}

type mapSerClient struct {
	cc grpc.ClientConnInterface
}

func NewMapSerClient(cc grpc.ClientConnInterface) MapSerClient {
	return &mapSerClient{cc}
}

func (c *mapSerClient) Move(ctx context.Context, in *Move_Req, opts ...grpc.CallOption) (*Move_Rsp, error) {
	out := new(Move_Rsp)
	err := c.cc.Invoke(ctx, "/proto.MapSer/Move", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MapSerServer is the server API for MapSer service.
// All implementations must embed UnimplementedMapSerServer
// for forward compatibility
type MapSerServer interface {
	Move(context.Context, *Move_Req) (*Move_Rsp, error)
	mustEmbedUnimplementedMapSerServer()
}

// UnimplementedMapSerServer must be embedded to have forward compatible implementations.
type UnimplementedMapSerServer struct {
}

func (UnimplementedMapSerServer) Move(context.Context, *Move_Req) (*Move_Rsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Move not implemented")
}
func (UnimplementedMapSerServer) mustEmbedUnimplementedMapSerServer() {}

// UnsafeMapSerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MapSerServer will
// result in compilation errors.
type UnsafeMapSerServer interface {
	mustEmbedUnimplementedMapSerServer()
}

func RegisterMapSerServer(s grpc.ServiceRegistrar, srv MapSerServer) {
	s.RegisterService(&MapSer_ServiceDesc, srv)
}

func _MapSer_Move_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Move_Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MapSerServer).Move(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.MapSer/Move",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MapSerServer).Move(ctx, req.(*Move_Req))
	}
	return interceptor(ctx, in, info, handler)
}

// MapSer_ServiceDesc is the grpc.ServiceDesc for MapSer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MapSer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.MapSer",
	HandlerType: (*MapSerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Move",
			Handler:    _MapSer_Move_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "map.proto",
}
