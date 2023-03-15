package shape

import (
	"context"
	"m3game/runtime/plugin"

	"google.golang.org/grpc"
)

var (
	_shape Shape
)

func Set(s Shape) {
	if _shape != nil {
		panic("Shape Only One")
	}
	_shape = s
}

func Get() Shape {
	return _shape
}

type Shape interface {
	plugin.PluginIns
	RegisterRule([]Rule) error
	ClientInterceptor() grpc.UnaryClientInterceptor
	ServerInterceptor() grpc.UnaryServerInterceptor
}

func ClientInterceptor() grpc.UnaryClientInterceptor {
	if Get() == nil {
		return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	}
	return Get().ClientInterceptor()
}

func ServerInterceptor() grpc.UnaryServerInterceptor {
	if Get() == nil {
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			return handler(ctx, req)
		}
	}
	return Get().ServerInterceptor()
}
