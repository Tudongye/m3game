package actorregcli

import (
	"context"
	"fmt"
	"m3game/plugins/log"
	"m3game/runtime/client"
	"m3game/runtime/rpc"
	"time"

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
	if err := rpc.InjectionRPC(pb.File_actor_proto.Services().Get(1)); err != nil {
		log.Fatal("InjectionRPC ActorRegSer %s", err.Error())
	}
}

func Init(srcapp meta.RouteApp, opts ...grpc.DialOption) error {
	if _client != nil {
		return nil
	}
	dstsvc := meta.GenDstRouteSvc(srcapp, proto.ActorAppFuncID)
	_client = &Client{
		Client: client.New(srcapp, dstsvc),
	}

	var err error
	target := fmt.Sprintf("router://%s", _client.DstSvc().String())
	opts = append(opts,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"Balance_m3g"}`),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(client.ClientInterceptors()...)),
		grpc.WithTimeout(time.Second*10),
	)
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

func Kick(ctx context.Context, id string, app string, opts ...grpc.CallOption) ([]byte, error) {
	var in pb.Kick_Req
	in.Leaseid = id
	_, err := client.RPCCallP2P(_client, _client.Kick, ctx, &in, meta.RouteApp(app), opts...)
	return nil, err
}
