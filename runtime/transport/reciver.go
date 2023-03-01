package transport

import (
	"context"
	"m3game/proto"
	"m3game/proto/pb"

	"google.golang.org/grpc"
)

// Contain ServerInterPara

func newReciver(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (*Reciver, error) {
	rec := &Reciver{
		ctx:     ctx,
		req:     req,
		info:    info,
		handler: handler,
		metas:   proto.NewM3Metas(),
	}
	if mspkg, ok := req.(proto.M3Pkg); !ok {
		return nil, _err_msgisnotm3pkg
	} else {
		rec.metas.Decode(mspkg.GetRouteHead().Metas)
		rec.routehead = mspkg.GetRouteHead()
	}
	return rec, nil
}

type Reciver struct {
	ctx       context.Context
	req       interface{}
	info      *grpc.UnaryServerInfo
	handler   grpc.UnaryHandler
	metas     *proto.M3Metas
	routehead *pb.RouteHead
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

func (s *Reciver) Metas() *proto.M3Metas {
	return s.metas
}

func (s *Reciver) RouteHead() *pb.RouteHead {
	return s.routehead
}

func (r *Reciver) HandleMsg(ctx context.Context) (resp interface{}, err error) {
	return r.handler(ctx, r.req)
}
