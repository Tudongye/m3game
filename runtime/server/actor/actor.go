package actor

import (
	"context"
	"fmt"
	"m3game/plugins/lease"
	"m3game/plugins/log"
	"time"

	"github.com/pkg/errors"

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
	OnMoveIn([]byte) error
	OnMoveOut() []byte
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

func (a *ActorBase) OnMoveIn([]byte) error {
	return nil
}

func (a *ActorBase) OnMoveOut() []byte {
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
	actor       Actor
	reqchan     chan *actorReq
	cancel      context.CancelFunc
	ctx         context.Context
	activetime  int64 // 激活时间
	savetime    int64 // 回写时间
	moveoutchch chan chan []byte
}

func (ar *actorRuntime) run() error {
	ar.moveoutchch = make(chan chan []byte, 1)
	if _cfg.LeaseMode == 1 {
		// 获取租约
		var movebytes []byte
		var err error
		if movebytes, err = ar.allocLease(); err != nil {
			return err
		}
		if err := ar.actor.OnMoveIn(movebytes); err != nil {
			return err
		}
	}
	if err := ar.actor.OnInit(); err != nil {
		return err
	}

	now := time.Now().Unix()
	ar.activetime = now
	ar.savetime = now
	t := time.NewTicker(time.Duration(_cfg.TickTimeInter) * time.Millisecond)
	defer t.Stop()
	for {
		if !ar.loop(ar.ctx, t) {
			break
		}
	}
	ar.exit()
	if _cfg.LeaseMode == 1 {
		// 释放租约
		leaseid := fmt.Sprintf("%s/%s", _cfg.LeasePrefix, ar.actor.ID())
		if err := lease.FreeLease(context.Background(), leaseid); err != nil {
			return errors.Wrapf(err, "FreeLease %s", leaseid)
		}
	}
	return nil
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
		return false
	case moveoutch := <-ar.moveoutchch:
		moveoutch <- ar.actor.OnMoveOut()
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

func (ar *actorRuntime) kickLease(ctx context.Context) ([]byte, error) {
	log.Info("kickLease")
	var moveoutch = make(chan []byte)
	select {
	case ar.moveoutchch <- moveoutch:
		break
	default:
		return nil, fmt.Errorf("Actor %s has MoveOut", ar.actor.ID())
	}

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("Actor %s MoveOut Timeout", ar.actor.ID())
	case movebytes := <-moveoutch:
		return movebytes, nil
	}
}

func (ar *actorRuntime) allocLease() ([]byte, error) {
	log.Info("allocLease")
	var movebytes []byte
	ctx, cancel := context.WithTimeout(ar.ctx, time.Duration(_cfg.AllocLeaseTimeOut)*time.Second)
	defer cancel()
	leaseid := fmt.Sprintf("%s/%s", _cfg.LeasePrefix, ar.actor.ID())
	if v, err := lease.GetLease(ctx, leaseid); err != nil {
		return nil, errors.Wrapf(err, " GetLease %s ", leaseid)
	} else if v != nil {
		// 踢出租约
		if movebytes, err = lease.KickLease(ctx, leaseid); err != nil {
			return nil, errors.Wrapf(err, " KickLease %s ", leaseid)
		}
		time.Sleep(time.Duration(_cfg.WaitFreeLeaseTimeOut) * time.Second)
	}
	// 申请租约
	if err := lease.AllocLease(ctx, leaseid, ar.kickLease); err != nil {
		return nil, errors.Wrapf(err, " AllocLease %s ", leaseid)
	}
	return movebytes, nil
}
