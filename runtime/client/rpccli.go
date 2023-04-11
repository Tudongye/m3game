/*
package client,According to pb's custom option rpoption, obtain the relevant parameters of rpc call through reflection and perform the response operation
client包，根据pb的自定义选项rpcoption，通过反射获得rpc调用的相关参数，并执行响应操作
*/
package client

import (
	"context"
	"m3game/meta/errs"
	"m3game/runtime/mesh"
	"m3game/runtime/rpc"

	"google.golang.org/protobuf/proto"

	"google.golang.org/grpc"
)

// RPC Random Route
func RPCCallRandom[T1, T2 proto.Message](c Client, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, opts ...grpc.CallOption) (T2, error) {
	t1fullname := t1.ProtoReflect().Descriptor().FullName()
	method := rpc.Meta(t1fullname)
	if method == nil {
		var t2 T2
		return t2, errs.RPCMethodNotRegister.New("Method not find %s", t1fullname)
	}
	ctx = FillRouteHeadRandom(ctx, c.SrcApp(), c.DstSvc(), "0")
	return rpcCall(method, f, ctx, t1, opts...)
}

// RPC Hash Route
func RPCCallHash[T1, T2 proto.Message](c Client, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, opts ...grpc.CallOption) (T2, error) {
	t1fullname := t1.ProtoReflect().Descriptor().FullName()
	method := rpc.Meta(t1fullname)
	var t2 T2
	if method == nil {
		return t2, errs.RPCMethodNotRegister.New("Method not find %s", t1fullname)
	}
	hashkey, err := method.HashKey(t1)
	if err != nil {
		return t2, errs.RPCCantFindHashKey.Wrap(err, "RPCCallHash Get HashKey For %s", t1fullname)
	}
	ctx = FillRouteHeadHash(ctx, c.SrcApp(), c.DstSvc(), hashkey, "0")
	return rpcCall(method, f, ctx, t1, opts...)
}

// RPC P2P Route
func RPCCallP2P[T1, T2 proto.Message](c Client, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, dstapp mesh.RouteApp, opts ...grpc.CallOption) (T2, error) {
	t1fullname := t1.ProtoReflect().Descriptor().FullName()
	method := rpc.Meta(t1fullname)
	var t2 T2
	if method == nil {
		return t2, errs.RPCMethodNotRegister.New("Method not find %s", t1fullname)
	}
	ctx = FillRouteHeadP2P(ctx, c.SrcApp(), c.DstSvc(), dstapp, "0")
	return rpcCall(method, f, ctx, t1, opts...)
}

// RPC Single Route
func RPCCallSingle[T1, T2 proto.Message](c Client, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, opts ...grpc.CallOption) (T2, error) {
	t1fullname := t1.ProtoReflect().Descriptor().FullName()
	method := rpc.Meta(t1fullname)
	var t2 T2
	if method == nil {
		return t2, errs.RPCMethodNotRegister.New("Method not find %s", t1fullname)
	}
	ctx = FillRouteHeadSingle(ctx, c.SrcApp(), c.DstSvc(), "0")
	return rpcCall(method, f, ctx, t1, opts...)
}

// RPC BroadCast Route
func RPCCallBroadCast[T1, T2 proto.Message](c Client, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, opts ...grpc.CallOption) (T2, error) {
	t1fullname := t1.ProtoReflect().Descriptor().FullName()
	method := rpc.Meta(t1fullname)
	if method == nil {
		var t2 T2
		return t2, errs.RPCMethodNotRegister.New("Method not find %s", t1fullname)
	}
	ctx = FillRouteHeadBroad(ctx, c.SrcApp(), c.DstSvc())
	return rpcCall(method, f, ctx, t1, opts...)
}

func rpcCall[T1, T2 proto.Message](method *rpc.RPCMeta, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, opts ...grpc.CallOption) (T2, error) {
	t2, err := f(ctx, t1, opts...)
	if err != nil {
		return t2, errs.RPCCallFuncFail.Wrap(err, "Call Rpc Func %s Fail", method.RpcName())
	}
	return t2, err
}
