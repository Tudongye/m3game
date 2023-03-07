package dirclient

import (
	"context"
	"fmt"
	"m3game/proto/pb"
	"m3game/runtime/client"
	"m3game/runtime/rpc"
	"m3game/runtime/transport"
	"m3game/util"

	dproto "m3game/demo/proto"
	dpb "m3game/demo/proto/pb"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"google.golang.org/grpc"
)

var (
	_client *Client
)

func init() {
	if err := rpc.RegisterRPCSvc(dpb.File_dir_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("RegisterRPCSvc Dir %s", err.Error()))
	}
}

func Init(srcins *pb.RouteIns, opts ...grpc.CallOption) error {
	_client = &Client{
		Meta: client.NewMeta(
			srcins,
			&pb.RouteSvc{
				EnvID:   srcins.EnvID,
				WorldID: srcins.WorldID,
				FuncID:  srcins.FuncID,
				IDStr:   util.SvcID2Str(srcins.EnvID, srcins.WorldID, dproto.DirAppFuncID),
			},
		),
		opts: opts,
	}
	var err error
	if _client.conn, err = grpc.Dial(
		fmt.Sprintf("router://%s", util.SvcID2Str(srcins.EnvID, srcins.WorldID, dproto.DirAppFuncID)),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(transport.ClientInterceptors()...)),
	); err != nil {
		return err
	} else {
		_client.DirSerClient = dpb.NewDirSerClient(_client.conn)
		return nil
	}
}

type Client struct {
	client.Meta
	dpb.DirSerClient
	conn *grpc.ClientConn
	opts []grpc.CallOption
}

func Hello(ctx context.Context, hellostr string, opts ...grpc.CallOption) (string, error) {
	var in dpb.Hello_Req
	in.Req = hellostr
	out, err := client.RPCCallRandom(_client, _client.Hello, ctx, &in, append(opts, _client.opts...)...)
	if err != nil {
		return "", err
	} else {
		return out.Rsp, nil
	}
}

func TraceHello(ctx context.Context, hellostr string, opts ...grpc.CallOption) (string, error) {
	var in dpb.TraceHello_Req
	in.Req = hellostr
	out, err := client.RPCCallRandom(_client, _client.TraceHello, ctx, &in, append(opts, _client.opts...)...)
	if err != nil {
		return "", err
	} else {
		return out.Rsp, nil
	}
}

func BreakHello(ctx context.Context, hellostr string, opts ...grpc.CallOption) (string, error) {
	var in dpb.BreakHello_Req
	in.Req = hellostr
	out, err := client.RPCCallRandom(_client, _client.BreakHello, ctx, &in, append(opts, _client.opts...)...)
	if err != nil {
		return "", err
	} else {
		return out.Rsp, nil
	}
}
