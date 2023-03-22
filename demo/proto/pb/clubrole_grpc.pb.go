// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.2.0
// source: clubrole.proto

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

// ClubRoleSerClient is the client API for ClubRoleSer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClubRoleSerClient interface {
	ClubRoleRead(ctx context.Context, in *ClubRoleRead_Req, opts ...grpc.CallOption) (*ClubRoleRead_Rsp, error)
	ClubRoleCreate(ctx context.Context, in *ClubRoleCreate_Req, opts ...grpc.CallOption) (*ClubRoleCreate_Rsp, error)
	ClubRoleDelete(ctx context.Context, in *ClubRoleDelete_Req, opts ...grpc.CallOption) (*ClubRoleDelete_Rsp, error)
}

type clubRoleSerClient struct {
	cc grpc.ClientConnInterface
}

func NewClubRoleSerClient(cc grpc.ClientConnInterface) ClubRoleSerClient {
	return &clubRoleSerClient{cc}
}

func (c *clubRoleSerClient) ClubRoleRead(ctx context.Context, in *ClubRoleRead_Req, opts ...grpc.CallOption) (*ClubRoleRead_Rsp, error) {
	out := new(ClubRoleRead_Rsp)
	err := c.cc.Invoke(ctx, "/proto.ClubRoleSer/ClubRoleRead", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clubRoleSerClient) ClubRoleCreate(ctx context.Context, in *ClubRoleCreate_Req, opts ...grpc.CallOption) (*ClubRoleCreate_Rsp, error) {
	out := new(ClubRoleCreate_Rsp)
	err := c.cc.Invoke(ctx, "/proto.ClubRoleSer/ClubRoleCreate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clubRoleSerClient) ClubRoleDelete(ctx context.Context, in *ClubRoleDelete_Req, opts ...grpc.CallOption) (*ClubRoleDelete_Rsp, error) {
	out := new(ClubRoleDelete_Rsp)
	err := c.cc.Invoke(ctx, "/proto.ClubRoleSer/ClubRoleDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClubRoleSerServer is the server API for ClubRoleSer service.
// All implementations must embed UnimplementedClubRoleSerServer
// for forward compatibility
type ClubRoleSerServer interface {
	ClubRoleRead(context.Context, *ClubRoleRead_Req) (*ClubRoleRead_Rsp, error)
	ClubRoleCreate(context.Context, *ClubRoleCreate_Req) (*ClubRoleCreate_Rsp, error)
	ClubRoleDelete(context.Context, *ClubRoleDelete_Req) (*ClubRoleDelete_Rsp, error)
	mustEmbedUnimplementedClubRoleSerServer()
}

// UnimplementedClubRoleSerServer must be embedded to have forward compatible implementations.
type UnimplementedClubRoleSerServer struct {
}

func (UnimplementedClubRoleSerServer) ClubRoleRead(context.Context, *ClubRoleRead_Req) (*ClubRoleRead_Rsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClubRoleRead not implemented")
}
func (UnimplementedClubRoleSerServer) ClubRoleCreate(context.Context, *ClubRoleCreate_Req) (*ClubRoleCreate_Rsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClubRoleCreate not implemented")
}
func (UnimplementedClubRoleSerServer) ClubRoleDelete(context.Context, *ClubRoleDelete_Req) (*ClubRoleDelete_Rsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClubRoleDelete not implemented")
}
func (UnimplementedClubRoleSerServer) mustEmbedUnimplementedClubRoleSerServer() {}

// UnsafeClubRoleSerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClubRoleSerServer will
// result in compilation errors.
type UnsafeClubRoleSerServer interface {
	mustEmbedUnimplementedClubRoleSerServer()
}

func RegisterClubRoleSerServer(s grpc.ServiceRegistrar, srv ClubRoleSerServer) {
	s.RegisterService(&ClubRoleSer_ServiceDesc, srv)
}

func _ClubRoleSer_ClubRoleRead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClubRoleRead_Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClubRoleSerServer).ClubRoleRead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ClubRoleSer/ClubRoleRead",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClubRoleSerServer).ClubRoleRead(ctx, req.(*ClubRoleRead_Req))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClubRoleSer_ClubRoleCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClubRoleCreate_Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClubRoleSerServer).ClubRoleCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ClubRoleSer/ClubRoleCreate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClubRoleSerServer).ClubRoleCreate(ctx, req.(*ClubRoleCreate_Req))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClubRoleSer_ClubRoleDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClubRoleDelete_Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClubRoleSerServer).ClubRoleDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ClubRoleSer/ClubRoleDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClubRoleSerServer).ClubRoleDelete(ctx, req.(*ClubRoleDelete_Req))
	}
	return interceptor(ctx, in, info, handler)
}

// ClubRoleSer_ServiceDesc is the grpc.ServiceDesc for ClubRoleSer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClubRoleSer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ClubRoleSer",
	HandlerType: (*ClubRoleSerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ClubRoleRead",
			Handler:    _ClubRoleSer_ClubRoleRead_Handler,
		},
		{
			MethodName: "ClubRoleCreate",
			Handler:    _ClubRoleSer_ClubRoleCreate_Handler,
		},
		{
			MethodName: "ClubRoleDelete",
			Handler:    _ClubRoleSer_ClubRoleDelete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "clubrole.proto",
}