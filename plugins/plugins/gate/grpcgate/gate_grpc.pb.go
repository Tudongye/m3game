// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.2.0
// source: gate.proto

package grpcgate

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	metapb "m3game/meta/metapb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GateSerClient is the client API for GateSer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GateSerClient interface {
	CSTransport(ctx context.Context, opts ...grpc.CallOption) (GateSer_CSTransportClient, error)
}

type gateSerClient struct {
	cc grpc.ClientConnInterface
}

func NewGateSerClient(cc grpc.ClientConnInterface) GateSerClient {
	return &gateSerClient{cc}
}

func (c *gateSerClient) CSTransport(ctx context.Context, opts ...grpc.CallOption) (GateSer_CSTransportClient, error) {
	stream, err := c.cc.NewStream(ctx, &GateSer_ServiceDesc.Streams[0], "/proto.GateSer/CSTransport", opts...)
	if err != nil {
		return nil, err
	}
	x := &gateSerCSTransportClient{stream}
	return x, nil
}

type GateSer_CSTransportClient interface {
	Send(*metapb.CSMsg) error
	Recv() (*metapb.CSMsg, error)
	grpc.ClientStream
}

type gateSerCSTransportClient struct {
	grpc.ClientStream
}

func (x *gateSerCSTransportClient) Send(m *metapb.CSMsg) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gateSerCSTransportClient) Recv() (*metapb.CSMsg, error) {
	m := new(metapb.CSMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GateSerServer is the server API for GateSer service.
// All implementations must embed UnimplementedGateSerServer
// for forward compatibility
type GateSerServer interface {
	CSTransport(GateSer_CSTransportServer) error
	mustEmbedUnimplementedGateSerServer()
}

// UnimplementedGateSerServer must be embedded to have forward compatible implementations.
type UnimplementedGateSerServer struct {
}

func (UnimplementedGateSerServer) CSTransport(GateSer_CSTransportServer) error {
	return status.Errorf(codes.Unimplemented, "method CSTransport not implemented")
}
func (UnimplementedGateSerServer) mustEmbedUnimplementedGateSerServer() {}

// UnsafeGateSerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GateSerServer will
// result in compilation errors.
type UnsafeGateSerServer interface {
	mustEmbedUnimplementedGateSerServer()
}

func RegisterGateSerServer(s grpc.ServiceRegistrar, srv GateSerServer) {
	s.RegisterService(&GateSer_ServiceDesc, srv)
}

func _GateSer_CSTransport_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GateSerServer).CSTransport(&gateSerCSTransportServer{stream})
}

type GateSer_CSTransportServer interface {
	Send(*metapb.CSMsg) error
	Recv() (*metapb.CSMsg, error)
	grpc.ServerStream
}

type gateSerCSTransportServer struct {
	grpc.ServerStream
}

func (x *gateSerCSTransportServer) Send(m *metapb.CSMsg) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gateSerCSTransportServer) Recv() (*metapb.CSMsg, error) {
	m := new(metapb.CSMsg)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GateSer_ServiceDesc is the grpc.ServiceDesc for GateSer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GateSer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.GateSer",
	HandlerType: (*GateSerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CSTransport",
			Handler:       _GateSer_CSTransport_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "gate.proto",
}
