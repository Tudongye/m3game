package actor

import (
	"context"
	"fmt"
	"m3game/log"
	"m3game/runtime/server"
	"time"

	"google.golang.org/protobuf/proto"
)

type ActorCreater func(string) Actor

type Actor interface {
	ID() string
	OnInit() error                // 加载后触发
	OnTick() error                // 触发定时任务时
	OnExit() error                // 退出时触发
	Save() error                  // 写回
	ReBuild(proto.Message) error  // 重建,服务迁移
	Pack() (proto.Message, error) // 打包,服务迁移
}

func ActorBaseCreator(actorid string) *ActorBase {
	return &ActorBase{
		actorid: actorid,
	}
}

type ActorBase struct {
	actorid string
}

func (a *ActorBase) ID() string {
	return a.actorid
}

func (a *ActorBase) OnInit() error {
	return nil
}

func (a *ActorBase) OnTick() error {
	return nil
}

func (a *ActorBase) OnExit() error {
	return nil
}

func (a *ActorBase) Save() error {
	return nil
}

func (a *ActorBase) ReBuild(proto.Message) error {
	return nil
}

func (a *ActorBase) Pack() (proto.Message, error) {
	return nil, nil
}

func newActorRuntime(actor Actor) *actorRuntime {
	return &actorRuntime{
		actor:   actor,
		reqchan: make(chan *actorReq, _cfg.MaxReqChanSize),
	}
}

type actorRuntime struct {
	actor      Actor
	reqchan    chan *actorReq
	cancel     context.CancelFunc
	activetime int64 // 激活时间
	savetime   int64 // 回写时间
}

func (ar *actorRuntime) run() {
	now := time.Now().Unix()
	ar.activetime = now
	ar.savetime = now
	ctx, cancel := context.WithCancel(context.Background())
	ar.cancel = cancel
	t := time.NewTicker(time.Duration(_cfg.TickTimeInter) * time.Millisecond)
	defer t.Stop()
	for {
		if !ar.loop(ctx, t) {
			break
		}
	}
}

func (ar *actorRuntime) kick() {
	ar.cancel()
}

func (ar *actorRuntime) exit() {
	if err := ar.actor.Save(); err != nil {
		log.ErrorP(log.LogPlus{"ActorID": ar.actor.ID()}, "Save err %s", err.Error())
	}
	if err := ar.actor.OnExit(); err != nil {
		log.ErrorP(log.LogPlus{"ActorID": ar.actor.ID()}, "OnExit err %s", err.Error())
	}
}

func (ar *actorRuntime) ontick() {
	now := time.Now().Unix()
	if err := ar.actor.OnTick(); err != nil {
		log.ErrorP(log.LogPlus{"ActorID": ar.actor.ID()}, "OnTick err %s", err.Error())
	}
	if now-ar.savetime > int64(_cfg.SaveTimeInter) {
		ar.savetime = now
		if err := ar.actor.Save(); err != nil {
			log.ErrorP(log.LogPlus{"ActorID": ar.actor.ID()}, "Save err %s", err.Error())
		}
	}
	if now-ar.activetime > int64(_cfg.ActiveTimeOut) {
		ar.cancel()
		return
	}
}

func (ar *actorRuntime) pushreq(req *actorReq) error {
	select {
	case ar.reqchan <- req:
		return nil
	default:
		return fmt.Errorf("ReqChan Full")
	}
}

func (ar *actorRuntime) loop(ctx context.Context, t *time.Ticker) bool {
	select {
	case <-ctx.Done():
		ar.exit()
		return false
	case <-t.C:
		ar.ontick()
	case req := <-ar.reqchan:
		now := time.Now().Unix()
		ar.activetime = now
		req.ctx.SetActor(ar.actor)
		ctx := server.WithContext(req.ctx.Reciver().Ctx(), req.ctx)
		rsp, err := req.ctx.Reciver().HandleMsg(ctx)
		req.rspchan <- &actorRsp{
			rsp: rsp,
			err: err,
		}
	}
	return true
}
