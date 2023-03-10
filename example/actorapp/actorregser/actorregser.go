package actorregser

import (
	"context"
	"fmt"
	"m3game/example/actorapp/actor"
	"m3game/example/proto/pb"
	"m3game/runtime/rpc"
	"m3game/runtime/server/async"

	"google.golang.org/grpc"
)

func init() {
	if err := rpc.RegisterRPCSvc(pb.File_actor_proto.Services().Get(1)); err != nil {
		panic(fmt.Sprintf("RegisterRPCSvc ActorRegSer %s", err.Error()))
	}
}

func New() *ActorRegSer {
	return &ActorRegSer{
		Server: async.New("ActorRegSer"),
	}
}

type ActorRegSer struct {
	*async.Server
	pb.UnimplementedActorRegSerServer
}

func (d *ActorRegSer) Register(ctx context.Context, in *pb.Register_Req) (*pb.Register_Rsp, error) {
	out := new(pb.Register_Rsp)
	if actorid, err := actor.Register(in.Name); err != nil {
		return out, err
	} else {
		out.ActorID = actorid
	}
	return out, nil
}

func (s *ActorRegSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterActorRegSerServer(t, s)
		return nil
	}
}
