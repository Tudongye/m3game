package clubdcli

import (
	"context"
	"m3game/demo/proto"
	"m3game/demo/proto/pb"
	"m3game/plugins/log"
	"m3game/plugins/transport"
	"m3game/runtime/client"
	"m3game/runtime/mesh"
	"m3game/runtime/rpc"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	_client *Client
)

func init() {
	if err := rpc.InjectionRPC(pb.File_club_proto.Services().Get(1)); err != nil {
		log.Fatal("InjectionRPC ClubDSer %s", err.Error())
	}
}

func New(srcapp mesh.RouteApp, opts ...grpc.DialOption) (*Client, error) {
	if _client != nil {
		return _client, nil
	}
	dstsvc := mesh.GenDstRouteSvc(srcapp, proto.ClubFuncID)
	_client = &Client{
		Client: client.New(srcapp, dstsvc),
	}

	var err error
	if _client.conn, err = transport.Instance().ClientConn(_client.DstSvc().String(), opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", _client.DstSvc().String())
	} else {
		_client.ClubDaemonSerClient = pb.NewClubDaemonSerClient(_client.conn)
		return _client, nil
	}
}

type Client struct {
	client.Client
	pb.ClubDaemonSerClient
	conn grpc.ClientConnInterface
}

func Conn() grpc.ClientConnInterface {
	return _client.conn
}

func Kick(ctx context.Context, id string, app string, opts ...grpc.CallOption) ([]byte, error) {
	var in pb.ClubKick_Req
	in.LeaseId = id
	_, err := client.RPCCallP2P(_client, _client.ClubKick, ctx, &in, mesh.RouteApp(app), opts...)
	return nil, err
}
