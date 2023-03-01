package dirserver

import (
	"context"
	"fmt"
	dpb "m3game/demo/proto/pb"
	"m3game/proto"
	"m3game/server"
	"m3game/server/mutil"

	"google.golang.org/grpc"
)

func New() *DirSer {
	return &DirSer{
		Server: mutil.New("DirSer"),
	}
}

type DirSer struct {
	*mutil.Server
	dpb.UnimplementedDirSerServer
}

func (d *DirSer) Hello(ctx context.Context, in *dpb.Hello_Req) (*dpb.Hello_Rsp, error) {
	out := new(dpb.Hello_Rsp)
	sctx := server.ParseContext(ctx)
	out.Rsp = fmt.Sprintf("Hello , %s", in.Req)
	if sctx != nil {
		if v, ok := sctx.Reciver().Metas().Get(proto.META_CLIENT); ok && v == proto.META_FLAG_TRUE {
			out.Rsp = fmt.Sprintf("Hello Client , %s", in.Req)
		}
	}
	return out, nil
}

func (s *DirSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		dpb.RegisterDirSerServer(t, s)
		return nil
	}
}
