package server

import (
	"context"
	"m3game/runtime/app"

	"google.golang.org/grpc"
)

type Type string

const (
	Multi Type = "multi" // multithread
	Actor Type = "actor" // actor model, sync for per actor
	Async Type = "async" // single thread async
)
const (
	_serctxKey = "_serverctx"
)

func WithServer(ctx context.Context, sctx Server) context.Context {
	return context.WithValue(ctx, _serctxKey, sctx)
}

func ParseServer(ctx context.Context) Server {
	if value := ctx.Value(_serctxKey); value == nil {
		return nil
	} else {
		return value.(Server)
	}
}

type Server interface {
	Init(map[string]interface{}, app.App) error                                                                                                                 // 初始化
	Type() Type                                                                                                                                                 // 类型
	Name() string                                                                                                                                               // 服务名
	Prepare(context.Context) error                                                                                                                              // 预启动
	Start(context.Context)                                                                                                                                      // 启动
	Reload(map[string]interface{}) error                                                                                                                        // 重载
	ServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error)                // RPC ServerInterceptor
	ClientInterceptor(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error // RPC ClientInterceptor
	TransportRegister() func(grpc.ServiceRegistrar) error                                                                                                       // register grpcser to conn
}
