package actorser

import (
	"context"
	"fmt"
	"m3game/example/actorapp/actor"
	"m3game/example/asyncapp/asynccli"
	"m3game/example/loader"
	"m3game/example/proto/pb"
	"m3game/meta"
	"m3game/plugins/log"
	"m3game/runtime/resource"
	"m3game/runtime/rpc"
	mactor "m3game/runtime/server/actor"
	"time"

	"github.com/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	_err_actor_parsefail = errors.New("_err_actor_parsefail")
	_err_actor_created   = errors.New("_err_actor_created")
	_err_actor_readyed   = errors.New("_err_actor_created")
	_err_actor_notcreate = errors.New("_err_actor_notcreate")
	_err_actor_notready  = errors.New("_err_actor_notready")
	_err_actor_dbnotfind = errors.New("_err_actor_dbnotfind")
	_err_actor_dberr     = errors.New("_err_actor_dberr")
)

var (
	_ser *ActorSer
)

func init() {
	if err := rpc.RegisterRPCSvc(pb.File_actor_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("RegisterRPCSvc Actor %s", err.Error()))
	}
}

func New() *ActorSer {
	if _ser != nil {
		return _ser
	}
	_ser = &ActorSer{
		Server: mactor.New("ActorSer", actor.ActorCreater),
	}
	return _ser
}

func Ser() *ActorSer {
	return _ser
}

type ActorSer struct {
	*mactor.Server
	pb.UnimplementedActorSerServer
}

func (s *ActorSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterActorSerServer(t, s)
		return nil
	}
}

func (d *ActorSer) Login(ctx context.Context, in *pb.Login_Req) (*pb.Login_Rsp, error) {
	out := new(pb.Login_Rsp)
	log.Info("Login")
	actor := actor.ConvertActor(ctx)
	if actor == nil {
		return out, _err_actor_parsefail
	}
	if actor.Ready() {
		return out, _err_actor_readyed
	}
	if err := actor.Login(ctx); err != nil {
		return out, err
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vlist, ok := md[string(meta.M3RouteSrcApp)]; ok {
			if len(vlist) > 0 {
				log.Debug("Actor %s at Gate %s", actor.ActorID(), vlist[0])
				actor.SetGate(vlist[0])
			}
		}
		if vlist, ok := md[string(meta.M3PlayerID)]; ok {
			if len(vlist) > 0 {
				log.Debug("Actor %s at PlayerID %s", actor.ActorID(), vlist[0])
				actor.SetPlayerId(vlist[0])
			}
		}
	}
	out.ActorDB = actor.DB()
	return out, nil
}

func (d *ActorSer) ModifyName(ctx context.Context, in *pb.ModifyName_Req) (*pb.ModifyName_Rsp, error) {
	out := new(pb.ModifyName_Rsp)
	actor := actor.ConvertActor(ctx)
	if actor == nil || !actor.Ready() {
		return out, _err_actor_notready
	}
	actor.ModifyName(in.NewName)
	out.ActorName = actor.DB().ActorName
	return out, nil
}

func (d *ActorSer) LvUp(ctx context.Context, in *pb.LvUp_Req) (*pb.LvUp_Rsp, error) {
	out := new(pb.LvUp_Rsp)
	actor := actor.ConvertActor(ctx)
	if actor == nil || !actor.Ready() {
		return out, _err_actor_notready
	}
	actor.LvUp()
	out.ActorInfo = actor.DB().ActorInfo
	return out, nil
}
func (d *ActorSer) GetInfo(ctx context.Context, in *pb.GetInfo_Req) (*pb.GetInfo_Rsp, error) {
	out := new(pb.GetInfo_Rsp)
	actor := actor.ConvertActor(ctx)
	if actor == nil || !actor.Ready() {
		return out, _err_actor_notready
	}
	out.Name = actor.DB().ActorName.Name

	titlecfgloader := resource.GetLoader[*loader.TitleCfgLoader](ctx)
	if titlecfgloader == nil {
		return out, fmt.Errorf("titlecfgloader Err")
	}
	out.Title = titlecfgloader.GetTitleByLv(actor.DB().ActorInfo.Level)
	return out, nil
}

func (d *ActorSer) PostChannel(ctx context.Context, in *pb.PostChannel_Req) (*pb.PostChannel_Rsp, error) {
	out := new(pb.PostChannel_Rsp)

	actor := actor.ConvertActor(ctx)
	if actor == nil || !actor.Ready() {
		return out, _err_actor_notready
	}
	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	ctx = metadata.NewOutgoingContext(ctx, md)
	if err := asynccli.TransChannel(ctx, &pb.ChannelMsg{Name: actor.Name(), Content: in.Content}); err != nil {
		return out, err
	} else {
		return out, nil
	}
}

func (d *ActorSer) PullChannel(ctx context.Context, in *pb.PullChannel_Req) (*pb.PullChannel_Rsp, error) {
	out := new(pb.PullChannel_Rsp)
	actor := actor.ConvertActor(ctx)
	if actor == nil || !actor.Ready() {
		return out, _err_actor_notready
	}
	if msg, err := asynccli.SSPullChannel(ctx); err != nil {
		log.Error(err.Error())
		return nil, err
	} else {
		out.Msgs = msg
	}
	return out, nil
}
