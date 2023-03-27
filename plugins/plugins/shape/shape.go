package shape

import (
	"context"
	"fmt"
	"m3game/plugins/log"
	"m3game/runtime/plugin"

	"google.golang.org/grpc"
)

type Shape interface {
	plugin.PluginIns
	RegisterRule([]Rule) error
	ClientInterceptor() grpc.UnaryClientInterceptor
	ServerInterceptor() grpc.UnaryServerInterceptor
}

var (
	_shape Shape
)

func New(me Shape) (Shape, error) {
	if _shape != nil {
		log.Fatal("Shape Only One")
		return nil, fmt.Errorf("Shape is newed %s", me.Factory().Name())
	}
	_shape = me
	return _shape, nil
}

func Instance() Shape {
	if _shape == nil {
		return nil
	}
	return _shape
}

func ClientInterceptor() grpc.UnaryClientInterceptor {
	if Instance() == nil {
		return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	}
	return Instance().ClientInterceptor()
}

func ServerInterceptor() grpc.UnaryServerInterceptor {
	if Instance() == nil {
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			return handler(ctx, req)
		}
	}
	return Instance().ServerInterceptor()
}
