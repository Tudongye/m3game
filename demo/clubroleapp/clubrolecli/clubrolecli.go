package clubrolecli

import (
	"context"
	"m3game/demo/proto"
	"m3game/demo/proto/pb"
	"m3game/meta"
	"m3game/plugins/log"
	"m3game/plugins/transport"
	"m3game/runtime/client"
	"m3game/runtime/rpc"

	"github.com/pkg/errors"

	"google.golang.org/grpc"
)

var (
	_client *Client
)

func init() {
	if err := rpc.InjectionRPC(pb.File_clubrole_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC ClubRole %s", err.Error())
	}
}

func New(srcapp meta.RouteApp, opts ...grpc.DialOption) (*Client, error) {
	if _client != nil {
		return _client, nil
	}
	dstsvc := meta.GenDstRouteSvc(srcapp, proto.ClubRoleFuncID)
	_client = &Client{
		Client: client.New(srcapp, dstsvc),
	}

	var err error
	if _client.conn, err = transport.Instance().ClientConn(_client.DstSvc().String(), opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", _client.DstSvc().String())
	} else {
		_client.ClubRoleSerClient = pb.NewClubRoleSerClient(_client.conn)
		return _client, nil
	}
}

type Client struct {
	client.Client
	pb.ClubRoleSerClient
	conn grpc.ClientConnInterface
}

func Conn() grpc.ClientConnInterface {
	return _client.conn
}

func ClubRoleRead(ctx context.Context, roleid int64, opts ...grpc.CallOption) (int64, error) {
	var in pb.ClubRoleRead_Req
	in.RoleId = roleid
	out, err := client.RPCCallRandom(_client, _client.ClubRoleRead, ctx, &in, opts...)
	if err != nil {
		return 0, err
	} else {
		return out.ClubId, nil
	}
}
