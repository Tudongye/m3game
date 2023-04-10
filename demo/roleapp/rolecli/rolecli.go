package rolecli

import (
	"context"
	"m3game/plugins/log"
	"m3game/plugins/transport"
	"m3game/runtime/client"
	"m3game/runtime/rpc"

	"m3game/demo/proto/pb"

	"m3game/demo/proto"

	"m3game/meta"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	_client *Client
)

func init() {
	if err := rpc.InjectionRPC(pb.File_role_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC Role %s", err.Error())
	}
}

func New(srcapp meta.RouteApp, opts ...grpc.DialOption) (*Client, error) {
	if _client != nil {
		return _client, nil
	}
	dstsvc := meta.GenDstRouteSvc(srcapp, proto.RoleFuncID)
	_client = &Client{
		Client: client.New(srcapp, dstsvc),
	}

	var err error
	if _client.conn, err = transport.Instance().ClientConn(_client.DstSvc().String(), opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", _client.DstSvc().String())
	} else {
		_client.RoleSerClient = pb.NewRoleSerClient(_client.conn)
		return _client, nil
	}
}

func Conn() grpc.ClientConnInterface {
	return _client.conn
}

type Client struct {
	client.Client
	pb.RoleSerClient
	conn grpc.ClientConnInterface
}

func RoleKick(ctx context.Context, roleid int64, dstapp meta.RouteApp, opts ...grpc.CallOption) error {
	var in pb.RoleKick_Req
	in.RoleId = roleid
	_, err := client.RPCCallP2P(_client, _client.RoleKick, ctx, &in, dstapp, opts...)
	if err != nil {
		return err
	} else {
		return nil
	}
}
