package worldcli

import (
	"m3game/plugins/log"
	"m3game/plugins/transport"
	"m3game/runtime/client"
	"m3game/runtime/mesh"
	"m3game/runtime/rpc"

	"m3game/example/proto/pb"

	"m3game/example/proto"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	_client *Client
)

func init() {
	if err := rpc.InjectionRPC(pb.File_world_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC World %s", err.Error())
	}
}
func New(srcapp mesh.RouteApp, opts ...grpc.DialOption) (*Client, error) {
	if _client != nil {
		return _client, nil
	}
	dstsvc := mesh.GenDstRouteSvc(srcapp, proto.WorldAppFuncID)
	_client = &Client{
		Client: client.New(srcapp, dstsvc),
	}

	var err error
	if _client.conn, err = transport.Instance().ClientConn(_client.DstSvc().String(), opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", _client.DstSvc().String())
	} else {
		_client.WorldSerClient = pb.NewWorldSerClient(_client.conn)
		return _client, nil
	}
}

type Client struct {
	client.Client
	pb.WorldSerClient
	conn grpc.ClientConnInterface
}

func Conn() grpc.ClientConnInterface {
	return _client.conn
}
