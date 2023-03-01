package transport

import (
	"context"
	"m3game/broker"
	"m3game/proto"
	"m3game/proto/pb"

	"google.golang.org/grpc"
)

// Contain ClientInterPara

const (
	_senderkey = "_senderkey"
)

func NewSender(
	ctx context.Context,
	method string,
	req, resp interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts []grpc.CallOption,
) (*Sender, error) {
	s := &Sender{
		ctx:     ctx,
		method:  method,
		req:     req,
		resp:    resp,
		cc:      cc,
		invoker: invoker,
		opts:    opts,
		metas:   proto.NewM3Metas(),
	}
	if m3pkg, ok := req.(proto.M3Pkg); !ok {
		return nil, _err_msgisnotm3pkg
	} else {
		s.routehead = m3pkg.GetRouteHead()
	}
	return s, nil
}

func WithSender(ctx context.Context, s *Sender) context.Context {
	return context.WithValue(ctx, _senderkey, s)
}

func ParseSender(ctx context.Context) *Sender {
	return ctx.Value(_senderkey).(*Sender)
}

type Sender struct {
	ctx       context.Context
	method    string
	req, resp interface{}
	cc        *grpc.ClientConn
	invoker   grpc.UnaryInvoker
	opts      []grpc.CallOption
	metas     *proto.M3Metas
	routehead *pb.RouteHead
}

func (s *Sender) Ctx() context.Context {
	return s.ctx
}

func (s *Sender) Method() string {
	return s.method
}

func (s *Sender) Req() interface{} {
	return s.req
}

func (s *Sender) Resp() interface{} {
	return s.resp
}

func (s *Sender) Cc() *grpc.ClientConn {
	return s.cc
}

func (s *Sender) Invoker() grpc.UnaryInvoker {
	return s.invoker
}

func (s *Sender) Opts() []grpc.CallOption {
	return s.opts
}

func (s *Sender) Metas() *proto.M3Metas {
	return s.metas
}

func (s *Sender) RouteHead() *pb.RouteHead {
	return s.routehead
}

func (s *Sender) sendMsg() error {
	s.RouteHead().Metas = s.Metas().Encode()
	if s.RouteHead().RouteType == pb.RouteType_RT_BROAD {
		return sendToBrokerSer(s, broker.GenTopic(s.RouteHead().DstSvc.IDStr))
	} else if s.RouteHead().RouteType == pb.RouteType_RT_MUTIL {
		return sendToBrokerSer(s, s.RouteHead().RoutePara.RouteMutilHead[0].Topic)
	}

	ctx := WithSender(s.ctx, s)
	return s.invoker(ctx, s.method, s.req, s.resp, s.cc, s.opts...)
}
