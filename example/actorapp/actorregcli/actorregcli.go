package actorregcli

import (
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
	if err := rpc.RegisterRPCSvc(pb.File_actor_proto.Services().Get(1)); err != nil {
		panic(fmt.Sprintf("RegisterRPCSvc ActorRegSer %s", err.Error()))
	}
}

func Init(srcapp meta.RouteApp, opts ...grpc.DialOption) error {
	if _client != nil {
		return nil
	}
	if env, world, _, _, err := srcapp.Parse(); err != nil {
		return nil
	} else {
		dstsvc := meta.GenRouteSvc(env, world, proto.ActorAppFuncID)
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
		_client.ActorRegSerClient = pb.NewActorRegSerClient(_client.conn)
		return nil
	}
}

type Client struct {
	client.Client
	pb.ActorRegSerClient
	conn *grpc.ClientConn
}

func Conn() *grpc.ClientConn {
	return _client.conn
}
