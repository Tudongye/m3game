/*
package client,According to pb's custom option rpoption, obtain the relevant parameters of rpc call through reflection and perform the response operation
client包，根据pb的自定义选项rpcoption，通过反射获得rpc调用的相关参数，并执行响应操作
*/
package client

import (
	"context"
	"fmt"
	"m3game/meta"
	"m3game/runtime/rpc"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"google.golang.org/grpc"
)

// RPC Random Route
func RPCCallRandom[T1, T2 proto.Message](c Client, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, opts ...grpc.CallOption) (T2, error) {
	t1fullname := t1.ProtoReflect().Descriptor().FullName()
	method := rpc.Meta(t1fullname)
	if method == nil {
		var t2 T2
		return t2, fmt.Errorf("Method not find %s", t1fullname)
	}
	ctx = FillRouteHeadRandom(ctx, c.SrcApp(), c.DstSvc(), meta.IsNtyFalse)
	return rpcCall(method, f, ctx, t1, opts...)
}

// RPC Hash Route
func RPCCallHash[T1, T2 proto.Message](c Client, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, opts ...grpc.CallOption) (T2, error) {
	t1fullname := t1.ProtoReflect().Descriptor().FullName()
	method := rpc.Meta(t1fullname)
	var t2 T2
	if method == nil {
		return t2, fmt.Errorf("Method not find %s", t1fullname)
	}
	hashkey, err := method.HashKey(t1)
	if err != nil {
		return t2, errors.Wrapf(err, "Get HashKey")
	}
	ctx = FillRouteHeadHash(ctx, c.SrcApp(), c.DstSvc(), hashkey, meta.IsNtyFalse)
	return rpcCall(method, f, ctx, t1, opts...)
}

// RPC P2P Route
func RPCCallP2P[T1, T2 proto.Message](c Client, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, dstapp meta.RouteApp, opts ...grpc.CallOption) (T2, error) {
	t1fullname := t1.ProtoReflect().Descriptor().FullName()
	method := rpc.Meta(t1fullname)
	var t2 T2
	if method == nil {
		return t2, fmt.Errorf("Method not find %s", t1fullname)
	}
	ctx = FillRouteHeadP2P(ctx, c.SrcApp(), c.DstSvc(), dstapp, meta.IsNtyFalse)
	return rpcCall(method, f, ctx, t1, opts...)
}

// RPC Single Route
func RPCCallSingle[T1, T2 proto.Message](c Client, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, opts ...grpc.CallOption) (T2, error) {
	t1fullname := t1.ProtoReflect().Descriptor().FullName()
	method := rpc.Meta(t1fullname)
	var t2 T2
	if method == nil {
		return t2, fmt.Errorf("Method not find %s", t1fullname)
	}
	ctx = FillRouteHeadSingle(ctx, c.SrcApp(), c.DstSvc(), meta.IsNtyFalse)
	return rpcCall(method, f, ctx, t1, opts...)
}

// RPC Multi Route
func RPCCallMulti[T1, T2 proto.Message](c Client, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, topicid string, opts ...grpc.CallOption) (T2, error) {
	t1fullname := t1.ProtoReflect().Descriptor().FullName()
	method := rpc.Meta(t1fullname)
	var t2 T2
	if method == nil {
		return t2, fmt.Errorf("Method not find %s", t1fullname)
	}
	ctx = FillRouteHeadMulti(ctx, c.SrcApp(), topicid)
	return rpcCall(method, f, ctx, t1, opts...)
}

// RPC BroadCast Route
func RPCCallBroadCast[T1, T2 proto.Message](c Client, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, opts ...grpc.CallOption) (T2, error) {
	t1fullname := t1.ProtoReflect().Descriptor().FullName()
	method := rpc.Meta(t1fullname)
	if method == nil {
		var t2 T2
		return t2, fmt.Errorf("Method not find %s", t1fullname)
	}
	ctx = FillRouteHeadBroad(ctx, c.SrcApp(), c.DstSvc())
	return rpcCall(method, f, ctx, t1, opts...)
}

func rpcCall[T1, T2 proto.Message](method *rpc.RPCMeta, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, opts ...grpc.CallOption) (T2, error) {
	t2, err := f(ctx, t1, opts...)
	if err != nil {
		return t2, errors.Wrapf(err, "Rpc %s", method.RpcName())
	}
	return t2, err
}
