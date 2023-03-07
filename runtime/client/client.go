/*
package client,According to pb's custom option rpoption, obtain the relevant parameters of rpc call through reflection and perform the response operation
client包，根据pb的自定义选项rpcoption，通过反射获得rpc调用的相关参数，并执行响应操作
*/
package client

import (
	"context"
	"errors"
	"m3game/log"
	"m3game/proto"
	"m3game/proto/pb"
	"m3game/runtime"
	"m3game/runtime/server"
	"m3game/runtime/transport"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	err_client_methodnotfind  = errors.New("err_client_methodnotfind")
	err_client_hashkeynotfind = errors.New("err_client_hashkeynotfind")
)

func init() {
	transport.RegisterClientInterceptor(SendInteror())
}

// RPC Random Route
func RPCCallRandom[T1, T2 proto.M3Pkg](m Meta, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, opts ...grpc.CallOption) (T2, error) {
	t1fullname := t1.ProtoReflect().Descriptor().FullName()
	method := m.Method(t1fullname)
	if method == nil {
		var t2 T2
		return t2, err_client_methodnotfind
	}
	routehead := NewRouteHeadRandom(m.SrcIns(), m.DstSvc())
	routehead.Ntf = method.grpcoption.Ntf

	return rpcCall(method, routehead, f, ctx, t1, opts...)
}

// RPC Hash Route
func RPCCallHash[T1, T2 proto.M3Pkg](m Meta, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, opts ...grpc.CallOption) (T2, error) {
	t1fullname := t1.ProtoReflect().Descriptor().FullName()
	method := m.Method(t1fullname)
	var t2 T2
	if method == nil {
		return t2, err_client_methodnotfind
	}
	if method.hashkeyd == nil {
		return t2, err_client_hashkeynotfind
	}
	hashkey, ok := t1.ProtoReflect().Get(method.hashkeyd).Interface().(string)
	if !ok {
		return t2, err_client_hashkeynotfind
	}
	routehead := NewRouteHeadHash(m.SrcIns(), m.DstSvc(), hashkey)
	routehead.Ntf = method.grpcoption.Ntf

	return rpcCall(method, routehead, f, ctx, t1, opts...)
}

// RPC P2P Route
func RPCCallP2P[T1, T2 proto.M3Pkg](m Meta, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, dstins *pb.RouteIns, opts ...grpc.CallOption) (T2, error) {
	t1fullname := t1.ProtoReflect().Descriptor().FullName()
	method := m.Method(t1fullname)
	var t2 T2
	if method == nil {
		return t2, err_client_methodnotfind
	}
	routehead := NewRouteHeadP2P(m.SrcIns(), m.DstSvc(), dstins)
	routehead.Ntf = method.grpcoption.Ntf
	return rpcCall(method, routehead, f, ctx, t1, opts...)
}

// RPC Single Route
func RPCCallSingle[T1, T2 proto.M3Pkg](m Meta, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, opts ...grpc.CallOption) (T2, error) {
	t1fullname := t1.ProtoReflect().Descriptor().FullName()
	method := m.Method(t1fullname)
	var t2 T2
	if method == nil {
		return t2, err_client_methodnotfind
	}
	routehead := NewRouteHeadSingle(m.SrcIns(), m.DstSvc())
	routehead.Ntf = method.grpcoption.Ntf
	return rpcCall(method, routehead, f, ctx, t1, opts...)
}

// RPC Mutil Route
func RPCCallMutil[T1, T2 proto.M3Pkg](m Meta, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, dsttopicid string, opts ...grpc.CallOption) (T2, error) {
	t1fullname := t1.ProtoReflect().Descriptor().FullName()
	method := m.Method(t1fullname)
	var t2 T2
	if method == nil {
		return t2, err_client_methodnotfind
	}
	routehead := NewRouteHeadMutil(m.SrcIns(), dsttopicid)
	routehead.Ntf = method.grpcoption.Ntf
	return rpcCall(method, routehead, f, ctx, t1, opts...)
}

// RPC BroadCast Route
func RPCCallBroadCast[T1, T2 proto.M3Pkg](m Meta, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, opts ...grpc.CallOption) (T2, error) {
	t1fullname := t1.ProtoReflect().Descriptor().FullName()
	method := m.Method(t1fullname)
	if method == nil {
		var t2 T2
		return t2, err_client_methodnotfind
	}
	routehead := NewRouteHeadBroad(m.SrcIns(), m.DstSvc())
	routehead.Ntf = method.grpcoption.Ntf

	return rpcCall(method, routehead, f, ctx, t1, opts...)
}

func rpcCall[T1, T2 proto.M3Pkg](method *m3Method, routehead *pb.RouteHead, f func(context.Context, T1, ...grpc.CallOption) (T2, error), ctx context.Context, t1 T1, opts ...grpc.CallOption) (T2, error) {
	t1.ProtoReflect().Set(method.routeheadd, protoreflect.ValueOfMessage(routehead.ProtoReflect()))
	return f(ctx, t1, opts...)
}

// grpc.ClientInterceptor
func SendInteror() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		s, err := transport.NewSender(ctx, method, req, resp, cc, invoker, opts)
		if err != nil {
			log.Error("NewSender err %s", err.Error())
			return err
		}
		for _, opt := range opts {
			if mopt, ok := opt.(M3CallOption); ok {
				if err := mopt.Filter()(s); err != nil {
					return err
				}
			}
		}
		sctx := server.ParseContext(ctx)
		if sctx == nil {
			// if not server call, direct call transport
			return runtime.SendInterFunc(s)
		}
		return sctx.Server().SendInterFunc(s)
	}
}
