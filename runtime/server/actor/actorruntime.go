package actor

import (
	"context"
	"m3game/meta/errs"
	"m3game/plugins/lease"
	"m3game/plugins/log"
	"time"

	"google.golang.org/grpc"
)

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

func newActorRuntime(actor Actor, cfg *Config) *actorRuntime {
	return &actorRuntime{
		cfg:     cfg,
		actor:   actor,
		reqchan: make(chan *actorReq, cfg.MaxReqChanSize),
	}
}

type actorRuntime struct {
	actor       Actor
	reqchan     chan *actorReq
	ctx         context.Context
	cancel      context.CancelFunc
	activetime  int64 // 激活时间
	savetime    int64 // 回写时间
	moveoutchch chan chan []byte
	cfg         *Config
}

func (ar *actorRuntime) run() error {
	ar.moveoutchch = make(chan chan []byte, 1)
	if ar.cfg.LeaseMode == 1 {
		// 获取租约
		var movebytes []byte
		var err error
		if movebytes, err = ar.allocLease(ar.ctx); err != nil {
			return err
		}
		if err := ar.actor.OnMoveIn(movebytes); err != nil {
			return errs.ActorMoveOutFail.Wrap(err, "actorRuntime run OnMoveIn")
		}
	}
	if err := ar.actor.OnInit(); err != nil {
		return errs.ActorOnInitFail.Wrap(err, "actorRuntime run OnInit")
	}

	now := time.Now().Unix()
	ar.activetime = now
	ar.savetime = now
	t := time.NewTicker(time.Duration(ar.cfg.TickTimeInter) * time.Millisecond)
	defer t.Stop()
	for {
		if !ar.loop(t) {
			break
		}
	}
	ar.exit()
	if ar.cfg.LeaseMode == 1 {
		// 释放租约
		leaseid := genLeaseId(ar.cfg.LeasePrefix, ar.actor.ID())
		if err := lease.FreeLease(ar.ctx, leaseid); err != nil {
			return errs.ActorRuntimeFreeLeaseFail.Wrap(err, "FreeLease %s", leaseid)
		}
	}
	return nil
}

func (ar *actorRuntime) kick() {
	ar.cancel()
}

func (ar *actorRuntime) exit() {
	if err := ar.actor.OnSave(); err != nil {
		log.ErrorP(log.LogPlus{"ActorID": ar.actor.ID()}, "OnSave err %s", err.Error())
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
	if now-ar.savetime > int64(ar.cfg.SaveTimeInter) && ar.cfg.SaveTimeInter != 0 {
		ar.savetime = now
		if err := ar.actor.OnSave(); err != nil {
			log.ErrorP(log.LogPlus{"ActorID": ar.actor.ID()}, "OnSave err %s", err.Error())
		}
	}
	if now-ar.activetime > int64(ar.cfg.ActiveTimeOut) && ar.cfg.ActiveTimeOut != 0 {
		ar.cancel()
		return
	}
}

func (ar *actorRuntime) pushreq(req *actorReq) error {
	select {
	case <-ar.ctx.Done():
		return errs.ActorRuntimePushReqFailActorDone.New("Actor %s have exit", ar.actor.ID())
	case ar.reqchan <- req:
		return nil
	default:
		return errs.ActorRuntimePushReqFailChanFull.New("ReqChan Full")
	}
}

func (ar *actorRuntime) loop(t *time.Ticker) bool {
	select {
	case <-ar.ctx.Done():
		return false
	case <-ar.actor.ExitCh():
		return false
	case moveoutch := <-ar.moveoutchch:
		moveoutch <- ar.actor.OnMoveOut()
		return false
	case <-t.C:
		ar.ontick()
	case req := <-ar.reqchan:
		now := time.Now().Unix()
		ar.activetime = now
		ctx := WithActor(req.ctx, ar.actor)
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
		return nil, errs.ActorRuntimeKickFailActorDone.New("Actor %s has MoveOut", ar.actor.ID())
	}

	select {
	case <-ctx.Done():
		return nil, errs.ActorRuntimeKickFailRPCDone.New("Actor %s MoveOut Timeout", ar.actor.ID())
	case movebytes := <-moveoutch:
		return movebytes, nil
	}
}

func (ar *actorRuntime) allocLease(ctx context.Context) ([]byte, error) {
	log.Info("allocLease")
	var movebytes []byte
	ctx, cancel := context.WithTimeout(ctx, time.Duration(ar.cfg.AllocLeaseTimeOut)*time.Second)
	defer cancel()
	leaseid := genLeaseId(ar.cfg.LeasePrefix, ar.actor.ID())
	if v, err := lease.GetLease(ctx, leaseid); err != nil {
		return nil, errs.ActorRuntimeAllocLeaseFail.Wrap(err, " GetLease %s ", leaseid)
	} else if v != nil {
		// 踢出租约
		if movebytes, err = lease.KickLease(ctx, leaseid); err != nil {
			return nil, errs.ActorRuntimeAllocLeaseFail.Wrap(err, " KickLease %s ", leaseid)
		}
		time.Sleep(time.Duration(ar.cfg.WaitFreeLeaseTimeOut) * time.Second)
	}
	// 申请租约
	if err := lease.AllocLease(ctx, leaseid, ar.kickLease); err != nil {
		// 强制释放一次
		if err := lease.FreeLease(ctx, leaseid); err != nil {
			log.Error("FreeLease %s Fail %s", leaseid, err.Error())
		}
		return nil, errs.ActorRuntimeAllocLeaseFail.Wrap(err, " AllocLease %s ", leaseid)
	}
	return movebytes, nil
}
