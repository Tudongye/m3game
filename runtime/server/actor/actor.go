package actor

import (
	"context"
	"fmt"
	"m3game/plugins/log"
	"time"

	"google.golang.org/grpc"
)

const (
	_actorkey = "_actorkey"
)

type ActorCreater func(string) Actor

type Actor interface {
	ID() string
	OnInit() error // 加载后触发
	OnTick() error // 触发定时任务时
	OnExit() error // 退出时触发
	Save() error   // 写回
}

func WithActor(ctx context.Context, actor Actor) context.Context {
	return context.WithValue(ctx, _actorkey, actor)
}

func ParseActor(ctx context.Context) Actor {
	var a Actor
	if t := ctx.Value(_actorkey); t == nil {
		return a
	} else {
		return t.(Actor)
	}
}

type actorReq struct {
	ctx     context.Context
	req     interface{}
	info    *grpc.UnaryServerInfo
	handler grpc.UnaryHandler
	rspchan chan *actorRsp
}

type actorRsp struct {
	rsp interface{}
	err error
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
		ctx = WithActor(req.ctx, ar.actor)
		rsp, err := req.handler(ctx, req.req)
		req.rspchan <- &actorRsp{
			rsp: rsp,
			err: err,
		}
	}
	return true
}
