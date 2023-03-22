package trace

import (
	"context"
	"m3game/runtime/rpc"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var (
	_traceprovider        *trace.TracerProvider
	_defaulttraceprovider = trace.NewTracerProvider()
)

func Get() *trace.TracerProvider {
	if _traceprovider == nil {
		return _defaulttraceprovider
	}
	return _traceprovider
}

func Set(tp *trace.TracerProvider) {
	if _traceprovider != nil {
		panic("Trace only one")
	}
	_traceprovider = tp
}

func ClientInterceptor() grpc.UnaryClientInterceptor {
	if _traceprovider == nil {
		return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	}
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if pbmsg, ok := req.(proto.Message); !ok {
			return invoker(ctx, method, req, reply, cc, opts...)
		} else if m := rpc.Meta(pbmsg.ProtoReflect().Descriptor().FullName()); m == nil {
			return invoker(ctx, method, req, reply, cc, opts...)
		} else if !m.GrpcOption().Trace {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		return otelgrpc.UnaryClientInterceptor()(ctx, method, req, reply, cc, invoker, opts...)
	}
}

func ServerInterceptor() grpc.UnaryServerInterceptor {
	if _traceprovider == nil {
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			return handler(ctx, req)
		}
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if pbmsg, ok := req.(proto.Message); !ok {
			return handler(ctx, req)
		} else if m := rpc.Meta(pbmsg.ProtoReflect().Descriptor().FullName()); m == nil {
			return handler(ctx, req)
		} else if !m.GrpcOption().Trace {
			return handler(ctx, req)
		}
		return otelgrpc.UnaryServerInterceptor()(ctx, req, info, handler)
	}
}
