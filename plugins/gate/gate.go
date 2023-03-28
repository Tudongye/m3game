package gate

import (
	"context"
	"m3game/meta/errs"
	"m3game/meta/metapb"
	"m3game/plugins/log"
	"m3game/runtime/plugin"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	_gatereciver GateReciver
	_gate        Gate
)

type Gate interface {
	plugin.PluginIns
	GetConn(connid string) CSConn
}

type CSConn interface {
	Send(ctx context.Context, msg *metapb.CSMsg) error
	Kick()
}

type GateReciver interface {
	AuthCall([]byte) (string, []byte, error)
	LogicCall(string, *metapb.CSMsg) (*metapb.CSMsg, error)
}

func New(g Gate) (Gate, error) {
	if _gate != nil {
		log.Fatal("Gate Only One")
		return nil, errs.GateInsHasNewed.New("gate is newed %s", _gate.Factory().Name())
	}
	_gate = g
	return _gate, nil
}

func Instance() Gate {
	if _gate == nil {
		log.Fatal("Gate not newd")
		return nil
	}
	return _gate
}

func GetConn(connid string) CSConn {
	return _gate.GetConn(connid)
}

func SetReciver(g GateReciver) {
	_gatereciver = g
}

func LogicCall(connid string, req *metapb.CSMsg) (*metapb.CSMsg, error) {
	return _gatereciver.LogicCall(connid, req)
}

func AuthCall(req []byte) (string, []byte, error) {
	return _gatereciver.AuthCall(req)
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
