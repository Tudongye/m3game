package actorregcli

import (
	"context"
	"m3game/plugins/log"
	"m3game/plugins/transport"
	"m3game/runtime/client"
	"m3game/runtime/rpc"

	"m3game/example/proto"
	"m3game/example/proto/pb"

	"github.com/pkg/errors"

	"m3game/meta"

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

func New(srcapp meta.RouteApp, opts ...grpc.DialOption) (*Client, error) {
	if _client != nil {
		return _client, nil
	}
	dstsvc := meta.GenDstRouteSvc(srcapp, proto.ActorAppFuncID)
	_client = &Client{
		Client: client.New(srcapp, dstsvc),
	}

	var err error
	if _client.conn, err = transport.Instance().ClientConn(_client.DstSvc().String(), opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", _client.DstSvc().String())
	} else {
		_client.ActorRegSerClient = pb.NewActorRegSerClient(_client.conn)
		return _client, nil
	}
}

type Client struct {
	client.Client
	pb.ActorRegSerClient
	conn grpc.ClientConnInterface
}

func Conn() grpc.ClientConnInterface {
	return _client.conn
}

func Kick(ctx context.Context, id string, app string, opts ...grpc.CallOption) ([]byte, error) {
	var in pb.Kick_Req
	in.Leaseid = id
	_, err := client.RPCCallP2P(_client, _client.Kick, ctx, &in, meta.RouteApp(app), opts...)
	return nil, err
}
