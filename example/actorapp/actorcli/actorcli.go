package actorcli

import (
	"fmt"
	"m3game/runtime/client"
	"m3game/runtime/rpc"
	"time"

	"m3game/example/proto/pb"

	"m3game/example/proto"

	"m3game/meta"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	_client *Client
)

func init() {
	if err := rpc.InjectionRPC(pb.File_actor_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("InjectionRPC Actor %s", err.Error()))
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
		_client.ActorSerClient = pb.NewActorSerClient(_client.conn)
		return nil
	}
}

type Client struct {
	client.Client
	pb.ActorSerClient
	conn *grpc.ClientConn
}

func Conn() *grpc.ClientConn {
	return _client.conn
}
