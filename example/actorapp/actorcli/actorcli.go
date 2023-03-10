package actorcli

import (
	"fmt"
	mpb "m3game/proto/pb"
	"m3game/runtime/client"
	"m3game/runtime/rpc"
	"m3game/util"

	"m3game/example/proto/pb"

	"m3game/example/proto"

	"google.golang.org/grpc"
)

var (
	_client *Client
)

func init() {
	if err := rpc.RegisterRPCSvc(pb.File_actor_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("RegisterRPCSvc Actor %s", err.Error()))
	}
}
func Init(srcins *mpb.RouteIns, opts ...grpc.CallOption) error {
	if _client != nil {
		return nil
	}
	dstsvc := util.RouteIns2Svc(srcins, proto.MutilAppFuncID)
	_client = &Client{
		Meta: client.NewMeta(
			srcins,
			dstsvc,
		),
		opts: opts,
	}

	var err error
	if _client.conn, err = grpc.Dial(
		fmt.Sprintf("router://%s", util.SvcID2Str(srcins.EnvID, srcins.WorldID, proto.ActorAppFuncID)),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(client.SendInteror()),
	); err != nil {
		return err
	} else {
		_client.ActorSerClient = pb.NewActorSerClient(_client.conn)
		return nil
	}
}

type Client struct {
	client.Meta
	pb.ActorSerClient
	conn *grpc.ClientConn
	opts []grpc.CallOption
}

func Conn() *grpc.ClientConn {
	return _client.conn
}
