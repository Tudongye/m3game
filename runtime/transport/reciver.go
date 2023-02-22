package transport

import (
	"context"
	"m3game/proto"
	"m3game/proto/pb"

	"google.golang.org/grpc"
)

func CreateReciver(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (*Reciver, error) {
	rec := &Reciver{
		ctx:     ctx,
		req:     req,
		info:    info,
		handler: handler,
		metas:   proto.CreateM3Metas(),
	}
	if mspkg, ok := req.(proto.M3Pkg); ok {
		rec.metas.Decode(mspkg.GetRouteHead().Metas)
	}
	return rec, nil
}

type Reciver struct {
	ctx     context.Context
	req     interface{}
	info    *grpc.UnaryServerInfo
	handler grpc.UnaryHandler
	metas   *proto.M3Metas
}

func (r *Reciver) Ctx() context.Context {
	return r.ctx
}

func (r *Reciver) Req() interface{} {
	return r.req
}

func (r *Reciver) Info() *grpc.UnaryServerInfo {
	return r.info
}

func (r *Reciver) Handler() grpc.UnaryHandler {
	return r.handler
}

func (r *Reciver) HandleMsg(ctx context.Context) (resp interface{}, err error) {
	return r.handler(ctx, r.req)
}

func (s *Reciver) Metas() *proto.M3Metas {
	return s.metas
}

func (s *Reciver) RouteHead() *pb.RouteHead {
	if m3pkg, ok := s.req.(proto.M3Pkg); ok {
		return m3pkg.GetRouteHead()
	}
	return nil
}
