package actorregser

import (
	"context"
	"fmt"
	"m3game/example/actorapp/actor"
	"m3game/example/proto/pb"
	"m3game/plugins/lease"
	"m3game/plugins/log"
	"m3game/runtime/rpc"
	"m3game/runtime/server/async"

	"google.golang.org/grpc"
)

func init() {
	if err := rpc.InjectionRPC(pb.File_actor_proto.Services().Get(1)); err != nil {
		panic(fmt.Sprintf("InjectionRPC ActorRegSer %s", err.Error()))
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
	log.Info("Register")
	if _, err := actor.Register(ctx, in.PlayerID, in.Name); err != nil {
		return out, err
	}
	return out, nil
}

func (d *ActorRegSer) Kick(ctx context.Context, in *pb.Kick_Req) (*pb.Kick_Rsp, error) {
	out := new(pb.Kick_Rsp)
	log.Info("Kick")
	if _, err := lease.RecvKickLease(ctx, in.Leaseid); err != nil {
		return out, err
	}
	return out, nil
}

func (s *ActorRegSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterActorRegSerServer(t, s)
		return nil
	}
}
