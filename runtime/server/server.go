package server

import (
	"context"
	"m3game/runtime/app"
	"m3game/runtime/transport"
	"sync"

	"google.golang.org/grpc"
)

type Type string

const (
	Mutil Type = "mutil" // mutilthread
	Actor Type = "actor" // actor model, sync for per actor
	Async Type = "async" // single thread async
)

const (
	_ctxkey = "serverctx"
)

type Context interface {
	Reciver() *transport.Reciver
	Server() Server
}

func WithContext(ctx context.Context, sctx Context) context.Context {
	return context.WithValue(ctx, _ctxkey, sctx)
}

func ParseContext(ctx context.Context) Context {
	if ctx.Value(_ctxkey) == nil {
		return nil
	}
	return ctx.Value(_ctxkey).(Context)
}

type Server interface {
	Init(map[string]interface{}, app.App) error                     // 初始化
	Type() Type                                                     // 类型
	Name() string                                                   // 服务名
	Start(wg *sync.WaitGroup) error                                 // 启动
	Stop() error                                                    // 停止
	Reload(map[string]interface{}) error                            // 重载
	RecvInterFunc(*transport.Reciver) (resp interface{}, err error) // RPC ServerInterceptor
	SendInterFunc(sender *transport.Sender) error                   // RPC ClientInterceptor
	CreateContext(*transport.Reciver) Context                       //
	TransportRegister() func(grpc.ServiceRegistrar) error           // register grpcser to conn
}
