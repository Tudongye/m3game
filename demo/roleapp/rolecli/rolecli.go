package rolecli

import (
	"context"
	"fmt"
	"m3game/plugins/log"
	"m3game/plugins/transport"
	"m3game/runtime/client"
	"m3game/runtime/rpc"
	"time"

	"m3game/demo/proto/pb"

	"m3game/demo/proto"

	"m3game/meta"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
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
		_client.RoleSerClient = pb.NewRoleSerClient(_client.conn)
		return _client, nil
	}
}

func Conn() *grpc.ClientConn {
	return _client.conn
}

type Client struct {
	client.Client
	pb.RoleSerClient
	conn *grpc.ClientConn
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
