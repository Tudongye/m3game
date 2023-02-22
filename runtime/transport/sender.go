package transport

import (
	"context"
	"fmt"
	"m3game/proto"
	"m3game/proto/pb"

	"google.golang.org/grpc"
)

const (
	_senderkey = "_senderkey"
)

func CreateSender(
	ctx context.Context,
	method string,
	req, resp interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts []grpc.CallOption,
) *Sender {
	return &Sender{
		ctx:     ctx,
		method:  method,
		req:     req,
		resp:    resp,
		cc:      cc,
		invoker: invoker,
		opts:    opts,
		metas:   proto.CreateM3Metas(),
	}
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

func (s *Sender) SendMsg() error {
	if m3pkg, ok := s.Req().(proto.M3Pkg); ok {
		m3pkg.GetRouteHead().Metas = s.Metas().Encode()
	} else {
		return fmt.Errorf("Req cant trans to M3Pkg")
	}
	ctx := WithSender(s.ctx, s)
	return s.invoker(ctx, s.method, s.req, s.resp, s.cc, s.opts...)
}

func (s *Sender) Metas() *proto.M3Metas {
	return s.metas
}

func (s *Sender) RouteHead() *pb.RouteHead {
	if m3pkg, ok := s.req.(proto.M3Pkg); ok {
		return m3pkg.GetRouteHead()
	}
	return nil
}
