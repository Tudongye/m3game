package clubrolecli

import (
	"context"
	"fmt"
	"m3game/demo/proto"
	"m3game/demo/proto/pb"
	"m3game/meta"
	"m3game/plugins/log"
	"m3game/plugins/transport"
	"m3game/runtime/client"
	"m3game/runtime/rpc"
	"time"

	"github.com/pkg/errors"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
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
	target := fmt.Sprintf("router://%s", _client.DstSvc().String())
	opts = append(opts,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"Balance_m3g"}`),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(transport.Instance().ClientInterceptors()...)),
		grpc.WithTimeout(time.Second*10),
	)
	if _client.conn, err = grpc.Dial(target, opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", target)
	} else {
		_client.ClubRoleSerClient = pb.NewClubRoleSerClient(_client.conn)
		return _client, nil
	}
}

type Client struct {
	client.Client
	pb.ClubRoleSerClient
	conn *grpc.ClientConn
}

func Conn() *grpc.ClientConn {
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
