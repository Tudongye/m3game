package server

import (
	"context"
	"m3game/runtime/app"
	"sync"

	"google.golang.org/grpc"
)

type Type string

const (
	Multi Type = "multi" // multithread
	Actor Type = "actor" // actor model, sync for per actor
	Async Type = "async" // single thread async
)

const (
	_ctxkey = "serverctx"
)

type Context struct {
	ser Server
	kv  map[string]interface{}
}

func (c *Context) Server() Server {
	return c.ser
}

func (c *Context) Set(key string, value interface{}) {
	c.kv[key] = value
}

func (c *Context) Get(key string) interface{} {
	return c.kv[key]
}

func GenContext(ser Server) *Context {
	return &Context{
		ser: ser,
		kv:  make(map[string]interface{}),
	}
}

func WithContext(ctx context.Context, sctx *Context) context.Context {
	return context.WithValue(ctx, _ctxkey, sctx)
}

func ParseContext(ctx context.Context) *Context {
	if ctx.Value(_ctxkey) == nil {
		return nil
	}
	return ctx.Value(_ctxkey).(*Context)
}

type Server interface {
	Init(map[string]interface{}, app.App) error                                                                                                                 // 初始化
	Type() Type                                                                                                                                                 // 类型
	Name() string                                                                                                                                               // 服务名
	Start(wg *sync.WaitGroup) error                                                                                                                             // 启动
	Stop() error                                                                                                                                                // 停止
	Reload(map[string]interface{}) error                                                                                                                        // 重载
	ServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error)                // RPC ServerInterceptor
	ClientInterceptor(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error // RPC ClientInterceptor
	TransportRegister() func(grpc.ServiceRegistrar) error                                                                                                       // register grpcser to conn
}
