package gatecli

import (
	"context"
	"m3game/runtime/client"
	"m3game/runtime/rpc"

	"m3game/example/proto"
	"m3game/example/proto/pb"
	"m3game/plugins/log"
	"m3game/plugins/transport"

	"github.com/pkg/errors"

	"m3game/meta"

	"google.golang.org/grpc"
)

var (
	_client *Client
)

func init() {
	if err := rpc.InjectionRPC(pb.File_gate_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC Gate %s", err.Error())
	}
}

func New(srcapp meta.RouteApp, opts ...grpc.DialOption) (*Client, error) {
	if _client != nil {
		return _client, nil
	}
	dstsvc := meta.GenDstRouteSvc(srcapp, proto.GateAppFuncID)
	_client = &Client{
		Client: client.New(srcapp, dstsvc),
	}

	var err error
	if _client.conn, err = transport.Instance().ClientConn(_client.DstSvc().String(), opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", _client.DstSvc().String())
	} else {
		_client.GateSerClient = pb.NewGateSerClient(_client.conn)
		return _client, nil
	}
}

type Client struct {
	client.Client
	pb.GateSerClient
	conn grpc.ClientConnInterface
}

func Conn() grpc.ClientConnInterface {
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
