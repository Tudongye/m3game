package onlinecli

import (
	"context"
	"fmt"
	"m3game/runtime/client"
	"m3game/runtime/rpc"
	"m3game/runtime/transport"

	"m3game/demo/proto"
	"m3game/demo/proto/pb"

	"github.com/pkg/errors"

	"m3game/meta"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var (
	_client *Client
)

func init() {
	if err := rpc.RegisterRPCSvc(pb.File_online_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("RegisterRPCSvc online %s", err.Error()))
	}
}

func Init(srcapp meta.RouteApp, opts ...grpc.DialOption) error {
	if _client != nil {
		return nil
	}
	if env, world, _, _, err := srcapp.Parse(); err != nil {
		return nil
	} else {
		dstsvc := meta.GenRouteSvc(env, world, proto.OnlineFuncID)
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
		_client.OnlineSerClient = pb.NewOnlineSerClient(_client.conn)
		return nil
	}
}

type Client struct {
	client.Client
	pb.OnlineSerClient
	conn *grpc.ClientConn
}

func Conn() *grpc.ClientConn {
	return _client.conn
}

func OnlineCreate(ctx context.Context, roleid string, appid string, opts ...grpc.CallOption) error {
	var in pb.OnlineCreate_Req
	in.RoleId = roleid
	in.AppId = appid
	_, err := client.RPCCallRandom(_client, _client.OnlineCreate, ctx, &in, opts...)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func OnlineDelete(ctx context.Context, roleid string, appid string, opts ...grpc.CallOption) error {
	var in pb.OnlineDelete_Req
	in.RoleId = roleid
	in.AppId = appid
	_, err := client.RPCCallRandom(_client, _client.OnlineDelete, ctx, &in, opts...)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func OnlineRead(ctx context.Context, roleid string, opts ...grpc.CallOption) (string, error) {
	var in pb.OnlineRead_Req
	in.RoleId = roleid
	out, err := client.RPCCallRandom(_client, _client.OnlineRead, ctx, &in, opts...)
	if err != nil {
		return "", err
	} else {
		return out.AppId, nil
	}
}
