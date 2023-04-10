package asynccli

import (
	"context"
	"m3game/plugins/log"
	"m3game/plugins/transport"
	"m3game/runtime/client"
	"m3game/runtime/rpc"

	"m3game/example/proto"
	"m3game/example/proto/pb"

	"m3game/meta"

	"github.com/pkg/errors"

	"google.golang.org/grpc"
)

var (
	_client *Client
)

func init() {
	if err := rpc.InjectionRPC(pb.File_async_proto.Services().Get(0)); err != nil {
		log.Fatal("InjectionRPC Async %s", err.Error())
	}
}

func New(srcapp meta.RouteApp, opts ...grpc.DialOption) (*Client, error) {
	if _client != nil {
		return _client, nil
	}
	dstsvc := meta.GenDstRouteSvc(srcapp, proto.AsyncAppFuncID)
	_client = &Client{
		Client: client.New(srcapp, dstsvc),
	}

	var err error
	if _client.conn, err = transport.Instance().ClientConn(_client.DstSvc().String(), opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", _client.DstSvc().String())
	} else {
		_client.AsyncSerClient = pb.NewAsyncSerClient(_client.conn)
		return _client, nil
	}
}

type Client struct {
	client.Client
	pb.AsyncSerClient
	conn grpc.ClientConnInterface
}

func Conn() grpc.ClientConnInterface {
	return _client.conn
}

func TransChannel(ctx context.Context, msg *pb.ChannelMsg, opts ...grpc.CallOption) error {
	var in pb.TransChannel_Req
	in.Msg = msg
	_, err := client.RPCCallBroadCast(_client, _client.TransChannel, ctx, &in, opts...)
	return err
}

func SSPullChannel(ctx context.Context, opts ...grpc.CallOption) ([]*pb.ChannelMsg, error) {
	var in pb.SSPullChannel_Req
	out, err := client.RPCCallRandom(_client, _client.SSPullChannel, ctx, &in, opts...)
	if err != nil {
		return nil, err
	} else {
		return out.Msgs, nil
	}
}
