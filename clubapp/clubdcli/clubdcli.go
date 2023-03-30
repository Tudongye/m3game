package clubdcli

import (
	"context"
	"fmt"
	"m3game/demo/proto"
	"m3game/demo/proto/pb"
	"m3game/meta"
	"m3game/plugins/log"
	"m3game/runtime/client"
	"m3game/runtime/rpc"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
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

func New(srcapp meta.RouteApp, opts ...grpc.DialOption) (*Client, error) {
	if _client != nil {
		return _client, nil
	}
	dstsvc := meta.GenDstRouteSvc(srcapp, proto.ClubFuncID)
	_client = &Client{
		Client: client.New(srcapp, dstsvc),
	}

	var err error
	target := fmt.Sprintf("router://%s", _client.DstSvc().String())
	opts = append(opts,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"Balance_m3g"}`),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(client.ClientInterceptors()...)),
		grpc.WithTimeout(time.Second*10),
	)
	if _client.conn, err = grpc.Dial(target, opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", target)
	} else {
		_client.ClubDaemonSerClient = pb.NewClubDaemonSerClient(_client.conn)
		return _client, nil
	}
}

type Client struct {
	client.Client
	pb.ClubDaemonSerClient
	conn *grpc.ClientConn
}

func Conn() *grpc.ClientConn {
	return _client.conn
}

func Kick(ctx context.Context, id string, app string, opts ...grpc.CallOption) ([]byte, error) {
	var in pb.ClubKick_Req
	in.LeaseId = id
	_, err := client.RPCCallP2P(_client, _client.ClubKick, ctx, &in, meta.RouteApp(app), opts...)
	return nil, err
}
