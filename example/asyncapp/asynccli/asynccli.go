package asynccli

import (
	"context"
	"fmt"
	"m3game/runtime/client"
	"m3game/runtime/rpc"

	"m3game/example/proto"
	"m3game/example/proto/pb"

	"m3game/meta"

	"github.com/pkg/errors"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var (
	_client *Client
)

func init() {
	if err := rpc.InjectionRPC(pb.File_async_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("InjectionRPC Async %s", err.Error()))
	}
}

func Init(srcapp meta.RouteApp, opts ...grpc.DialOption) error {
	if _client != nil {
		return nil
	}
	if env, world, _, _, err := srcapp.Parse(); err != nil {
		return nil
	} else {
		dstsvc := meta.GenRouteSvc(env, world, proto.AsyncAppFuncID)
		_client = &Client{
			Client: client.New(srcapp, dstsvc),
		}
	}
	var err error
	target := fmt.Sprintf("router://%s", _client.DstSvc().String())
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"Balance_m3g"}`))
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(client.ClientInterceptors()...)))
	if _client.conn, err = grpc.Dial(target, opts...); err != nil {
		return errors.Wrapf(err, "Dial Target %s", target)
	} else {
		_client.AsyncSerClient = pb.NewAsyncSerClient(_client.conn)
		return nil
	}
}

type Client struct {
	client.Client
	pb.AsyncSerClient
	conn *grpc.ClientConn
}

func Conn() *grpc.ClientConn {
	return _client.conn
}

func TransChannel(ctx context.Context, msg *pb.ChannelMsg, opts ...grpc.CallOption) error {
	var in pb.TransChannel_Req
	in.Msg = msg
	_, err := client.RPCCallBroadCast(_client, _client.TransChannel, ctx, &in, opts...)
	return err
}

func SSPullChannel(ctx context.Context, opts ...grpc.CallOption) ([]*pb.ChannelMsg, error) {
	var in pb.SSPullChannel_Req
	out, err := client.RPCCallRandom(_client, _client.SSPullChannel, ctx, &in, opts...)
	if err != nil {
		return nil, err
	} else {
		return out.Msgs, nil
	}
}
