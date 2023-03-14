package gatecli

import (
	"context"
	"fmt"
	"m3game/runtime/client"
	"m3game/runtime/rpc"
	"m3game/runtime/transport"

	"m3game/example/proto"
	"m3game/example/proto/pb"

	"github.com/pkg/errors"

	"m3game/meta"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var (
	_client *Client
)

func init() {
	if err := rpc.RegisterRPCSvc(pb.File_gate_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("RegisterRPCSvc Gate %s", err.Error()))
	}
}

func Init(srcapp meta.RouteApp, opts ...grpc.DialOption) error {
	if _client != nil {
		return nil
	}
	if env, world, _, _, err := srcapp.Parse(); err != nil {
		return nil
	} else {
		dstsvc := meta.GenRouteSvc(env, world, proto.GateAppFuncID)
		_client = &Client{
			Client: client.New(srcapp, dstsvc),
		}
	}
	var err error
	target := fmt.Sprintf("router://%s", _client.DstSvc().String())
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"Balance_m3g"}`))
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(transport.ClientInterceptors()...)))
	if _client.conn, err = grpc.Dial(target, opts...); err != nil {
		return errors.Wrapf(err, "Dial Target %s", target)
	} else {
		_client.GateSerClient = pb.NewGateSerClient(_client.conn)
		return nil
	}
}

type Client struct {
	client.Client
	pb.GateSerClient
	conn *grpc.ClientConn
}

func Conn() *grpc.ClientConn {
	return _client.conn
}

func SendToCli(ctx context.Context, playerid string, content string, dstapp meta.RouteApp, opts ...grpc.CallOption) error {
	var in pb.SendToCli_Req
	in.PlayerID = playerid
	in.Content = content
	_, err := client.RPCCallP2P(_client, _client.SendToCli, ctx, &in, dstapp, opts...)
	if err != nil {
		return err
	} else {
		return nil
	}
}
