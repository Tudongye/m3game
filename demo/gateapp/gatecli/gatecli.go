package gatecli

import (
	"context"
	"m3game/plugins/log"
	"m3game/plugins/transport"
	"m3game/runtime/client"
	"m3game/runtime/mesh"
	"m3game/runtime/rpc"

	"m3game/demo/proto"
	"m3game/demo/proto/pb"

	"github.com/pkg/errors"

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

func New(srcapp mesh.RouteApp, opts ...grpc.DialOption) (*Client, error) {
	if _client != nil {
		return _client, nil
	}
	dstsvc := mesh.GenDstRouteSvc(srcapp, proto.GateFuncID)
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

func Conn() grpc.ClientConnInterface {
	return _client.conn
}

type Client struct {
	client.Client
	pb.GateSerClient
	conn grpc.ClientConnInterface
}

func SendToCli(ctx context.Context, roleid int64, ntymsg *pb.NtyMsg, dstapp mesh.RouteApp, opts ...grpc.CallOption) error {
	var in pb.SendToCli_Req
	in.RoleId = roleid
	in.NtyMsg = ntymsg
	_, err := client.RPCCallP2P(_client, _client.SendToCli, ctx, &in, dstapp, opts...)
	if err != nil {
		return err
	} else {
		return nil
	}
}
