package gate

import (
	"context"
	"m3game/meta/metapb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	_gatereciver GateReciver
	_gate        Gate
)

type Gate interface {
	GetConn(playerid string) CSConn
}

type CSConn interface {
	Send(ctx context.Context, msg *metapb.CSMsg) error
	Kick()
}

type GateReciver interface {
	AuthCall(*metapb.AuthReq) (*metapb.AuthRsp, error)
	LogicCall(*metapb.CSMsg) (*metapb.CSMsg, error)
}

func Set(g Gate) {
	_gate = g
}

func GetConn(playerid string) CSConn {
	return _gate.GetConn(playerid)
}

func SetReciver(g GateReciver) {
	_gatereciver = g
}

func LogicCall(r *metapb.CSMsg) (*metapb.CSMsg, error) {
	return _gatereciver.LogicCall(r)
}
func AuthCall(r *metapb.AuthReq) (*metapb.AuthRsp, error) {
	return _gatereciver.AuthCall(r)
}

func CallGrpcCli(ctx context.Context, c grpc.ClientConnInterface, in *metapb.CSMsg, opts ...grpc.CallOption) (*metapb.CSMsg, error) {
	inmsg := &gateBuff{
		buff: in.Content,
	}
	outmsg := &gateBuff{}

	metas := make(map[string]string)
	for _, meta := range in.Metas {
		metas[meta.Key] = meta.Value
	}
	md := metadata.New(metas)
	ctx = metadata.NewOutgoingContext(ctx, md)
	if err := c.Invoke(ctx, in.Method, inmsg, outmsg, opts...); err != nil {
		return nil, err
	}
	out := &metapb.CSMsg{
		Method:  in.Method,
		Content: outmsg.buff,
	}
	return out, nil
}

type gateBuff struct {
	buff []byte
}

type GateCodec struct {
}

func (*GateCodec) Marshal(v interface{}) ([]byte, error) {
	c := v.(*gateBuff)
	return c.buff, nil
}

func (*GateCodec) Unmarshal(data []byte, v interface{}) error {
	c := v.(*gateBuff)
	c.buff = data
	return nil
}

func (*GateCodec) String() string {
	return "GateCodec"
}
