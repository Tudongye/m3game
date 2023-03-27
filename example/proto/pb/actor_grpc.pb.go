// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.2.0
// source: actor.proto

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

// ActorSerClient is the client API for ActorSer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ActorSerClient interface {
	Login(ctx context.Context, in *Login_Req, opts ...grpc.CallOption) (*Login_Rsp, error)
	GetInfo(ctx context.Context, in *GetInfo_Req, opts ...grpc.CallOption) (*GetInfo_Rsp, error)
	ModifyName(ctx context.Context, in *ModifyName_Req, opts ...grpc.CallOption) (*ModifyName_Rsp, error)
	LvUp(ctx context.Context, in *LvUp_Req, opts ...grpc.CallOption) (*LvUp_Rsp, error)
	PostChannel(ctx context.Context, in *PostChannel_Req, opts ...grpc.CallOption) (*PostChannel_Rsp, error)
	PullChannel(ctx context.Context, in *PullChannel_Req, opts ...grpc.CallOption) (*PullChannel_Rsp, error)
}

type actorSerClient struct {
	cc grpc.ClientConnInterface
}

func NewActorSerClient(cc grpc.ClientConnInterface) ActorSerClient {
	return &actorSerClient{cc}
}

func (c *actorSerClient) Login(ctx context.Context, in *Login_Req, opts ...grpc.CallOption) (*Login_Rsp, error) {
	out := new(Login_Rsp)
	err := c.cc.Invoke(ctx, "/proto.ActorSer/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *actorSerClient) GetInfo(ctx context.Context, in *GetInfo_Req, opts ...grpc.CallOption) (*GetInfo_Rsp, error) {
	out := new(GetInfo_Rsp)
	err := c.cc.Invoke(ctx, "/proto.ActorSer/GetInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *actorSerClient) ModifyName(ctx context.Context, in *ModifyName_Req, opts ...grpc.CallOption) (*ModifyName_Rsp, error) {
	out := new(ModifyName_Rsp)
	err := c.cc.Invoke(ctx, "/proto.ActorSer/ModifyName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *actorSerClient) LvUp(ctx context.Context, in *LvUp_Req, opts ...grpc.CallOption) (*LvUp_Rsp, error) {
	out := new(LvUp_Rsp)
	err := c.cc.Invoke(ctx, "/proto.ActorSer/LvUp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *actorSerClient) PostChannel(ctx context.Context, in *PostChannel_Req, opts ...grpc.CallOption) (*PostChannel_Rsp, error) {
	out := new(PostChannel_Rsp)
	err := c.cc.Invoke(ctx, "/proto.ActorSer/PostChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *actorSerClient) PullChannel(ctx context.Context, in *PullChannel_Req, opts ...grpc.CallOption) (*PullChannel_Rsp, error) {
	out := new(PullChannel_Rsp)
	err := c.cc.Invoke(ctx, "/proto.ActorSer/PullChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ActorSerServer is the server API for ActorSer service.
// All implementations must embed UnimplementedActorSerServer
// for forward compatibility
type ActorSerServer interface {
	Login(context.Context, *Login_Req) (*Login_Rsp, error)
	GetInfo(context.Context, *GetInfo_Req) (*GetInfo_Rsp, error)
	ModifyName(context.Context, *ModifyName_Req) (*ModifyName_Rsp, error)
	LvUp(context.Context, *LvUp_Req) (*LvUp_Rsp, error)
	PostChannel(context.Context, *PostChannel_Req) (*PostChannel_Rsp, error)
	PullChannel(context.Context, *PullChannel_Req) (*PullChannel_Rsp, error)
	mustEmbedUnimplementedActorSerServer()
}

// UnimplementedActorSerServer must be embedded to have forward compatible implementations.
type UnimplementedActorSerServer struct {
}

func (UnimplementedActorSerServer) Login(context.Context, *Login_Req) (*Login_Rsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedActorSerServer) GetInfo(context.Context, *GetInfo_Req) (*GetInfo_Rsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfo not implemented")
}
func (UnimplementedActorSerServer) ModifyName(context.Context, *ModifyName_Req) (*ModifyName_Rsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifyName not implemented")
}
func (UnimplementedActorSerServer) LvUp(context.Context, *LvUp_Req) (*LvUp_Rsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LvUp not implemented")
}
func (UnimplementedActorSerServer) PostChannel(context.Context, *PostChannel_Req) (*PostChannel_Rsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostChannel not implemented")
}
func (UnimplementedActorSerServer) PullChannel(context.Context, *PullChannel_Req) (*PullChannel_Rsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PullChannel not implemented")
}
func (UnimplementedActorSerServer) mustEmbedUnimplementedActorSerServer() {}

// UnsafeActorSerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ActorSerServer will
// result in compilation errors.
type UnsafeActorSerServer interface {
	mustEmbedUnimplementedActorSerServer()
}

func RegisterActorSerServer(s grpc.ServiceRegistrar, srv ActorSerServer) {
	s.RegisterService(&ActorSer_ServiceDesc, srv)
}

func _ActorSer_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Login_Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActorSerServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ActorSer/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActorSerServer).Login(ctx, req.(*Login_Req))
	}
	return interceptor(ctx, in, info, handler)
}

func _ActorSer_GetInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetInfo_Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActorSerServer).GetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ActorSer/GetInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActorSerServer).GetInfo(ctx, req.(*GetInfo_Req))
	}
	return interceptor(ctx, in, info, handler)
}

func _ActorSer_ModifyName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifyName_Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActorSerServer).ModifyName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ActorSer/ModifyName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActorSerServer).ModifyName(ctx, req.(*ModifyName_Req))
	}
	return interceptor(ctx, in, info, handler)
}

func _ActorSer_LvUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LvUp_Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActorSerServer).LvUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ActorSer/LvUp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActorSerServer).LvUp(ctx, req.(*LvUp_Req))
	}
	return interceptor(ctx, in, info, handler)
}

func _ActorSer_PostChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostChannel_Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActorSerServer).PostChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ActorSer/PostChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActorSerServer).PostChannel(ctx, req.(*PostChannel_Req))
	}
	return interceptor(ctx, in, info, handler)
}

func _ActorSer_PullChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PullChannel_Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActorSerServer).PullChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ActorSer/PullChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActorSerServer).PullChannel(ctx, req.(*PullChannel_Req))
	}
	return interceptor(ctx, in, info, handler)
}

// ActorSer_ServiceDesc is the grpc.ServiceDesc for ActorSer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ActorSer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ActorSer",
	HandlerType: (*ActorSerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _ActorSer_Login_Handler,
		},
		{
			MethodName: "GetInfo",
			Handler:    _ActorSer_GetInfo_Handler,
		},
		{
			MethodName: "ModifyName",
			Handler:    _ActorSer_ModifyName_Handler,
		},
		{
			MethodName: "LvUp",
			Handler:    _ActorSer_LvUp_Handler,
		},
		{
			MethodName: "PostChannel",
			Handler:    _ActorSer_PostChannel_Handler,
		},
		{
			MethodName: "PullChannel",
			Handler:    _ActorSer_PullChannel_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "actor.proto",
}

// ActorRegSerClient is the client API for ActorRegSer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ActorRegSerClient interface {
	Register(ctx context.Context, in *Register_Req, opts ...grpc.CallOption) (*Register_Rsp, error)
	Kick(ctx context.Context, in *Kick_Req, opts ...grpc.CallOption) (*Kick_Rsp, error)
}

type actorRegSerClient struct {
	cc grpc.ClientConnInterface
}

func NewActorRegSerClient(cc grpc.ClientConnInterface) ActorRegSerClient {
	return &actorRegSerClient{cc}
}

func (c *actorRegSerClient) Register(ctx context.Context, in *Register_Req, opts ...grpc.CallOption) (*Register_Rsp, error) {
	out := new(Register_Rsp)
	err := c.cc.Invoke(ctx, "/proto.ActorRegSer/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *actorRegSerClient) Kick(ctx context.Context, in *Kick_Req, opts ...grpc.CallOption) (*Kick_Rsp, error) {
	out := new(Kick_Rsp)
	err := c.cc.Invoke(ctx, "/proto.ActorRegSer/Kick", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ActorRegSerServer is the server API for ActorRegSer service.
// All implementations must embed UnimplementedActorRegSerServer
// for forward compatibility
type ActorRegSerServer interface {
	Register(context.Context, *Register_Req) (*Register_Rsp, error)
	Kick(context.Context, *Kick_Req) (*Kick_Rsp, error)
	mustEmbedUnimplementedActorRegSerServer()
}

// UnimplementedActorRegSerServer must be embedded to have forward compatible implementations.
type UnimplementedActorRegSerServer struct {
}

func (UnimplementedActorRegSerServer) Register(context.Context, *Register_Req) (*Register_Rsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedActorRegSerServer) Kick(context.Context, *Kick_Req) (*Kick_Rsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Kick not implemented")
}
func (UnimplementedActorRegSerServer) mustEmbedUnimplementedActorRegSerServer() {}

// UnsafeActorRegSerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ActorRegSerServer will
// result in compilation errors.
type UnsafeActorRegSerServer interface {
	mustEmbedUnimplementedActorRegSerServer()
}

func RegisterActorRegSerServer(s grpc.ServiceRegistrar, srv ActorRegSerServer) {
	s.RegisterService(&ActorRegSer_ServiceDesc, srv)
}

func _ActorRegSer_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Register_Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActorRegSerServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ActorRegSer/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActorRegSerServer).Register(ctx, req.(*Register_Req))
	}
	return interceptor(ctx, in, info, handler)
}

func _ActorRegSer_Kick_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Kick_Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActorRegSerServer).Kick(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ActorRegSer/Kick",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActorRegSerServer).Kick(ctx, req.(*Kick_Req))
	}
	return interceptor(ctx, in, info, handler)
}

// ActorRegSer_ServiceDesc is the grpc.ServiceDesc for ActorRegSer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ActorRegSer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ActorRegSer",
	HandlerType: (*ActorRegSerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _ActorRegSer_Register_Handler,
		},
		{
			MethodName: "Kick",
			Handler:    _ActorRegSer_Kick_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "actor.proto",
}
