package actor

import (
	"context"
	"fmt"
	"m3game/meta/errs"
	"m3game/meta/monitor"
	"m3game/plugins/log"
	"m3game/plugins/metric"
	"regexp"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type ActorCreater func(string) Actor

func newActorMgr(creater ActorCreater, cfg *Config) *ActorMgr {
	return &ActorMgr{
		cfg:          cfg,
		actorcreater: creater,
		actorreqpool: sync.Pool{
			New: func() any {
				return &actorReq{
					rspchan: make(chan *actorRsp),
				}
			},
		},
	}
}

type ActorMgr struct {
	cfg           *Config
	actorcreater  ActorCreater
	actorruntimes sync.Map
	actorreqpool  sync.Pool
}

func (am *ActorMgr) GetActor(actorid string) (*actorRuntime, bool) {
	a, ok := am.actorruntimes.Load(actorid)
	if !ok {
		return nil, false
	}
	return a.(*actorRuntime), true
}

func (am *ActorMgr) NewActor(actorid string) *actorRuntime {
	if a, ok := am.GetActor(actorid); ok {
		return a
	}
	ar := newActorRuntime(am.actorcreater(actorid), am.cfg)
	ctx, cancel := context.WithCancel(context.Background())
	ar.ctx = ctx
	ar.cancel = cancel
	if a, loaded := am.actorruntimes.LoadOrStore(actorid, ar); loaded {
		return a.(*actorRuntime)
	} else {
		metric.Gauge(monitor.ActorRuntimeTotal).Inc()
		ar = a.(*actorRuntime)
		go am.runActor(ar, actorid)
		return ar
	}
}

func (am *ActorMgr) runActor(ar *actorRuntime, actorid string) {
	if err := ar.run(); err != nil {
		log.Error("actor.run() err %s", err.Error())
	}
	ar.cancel()
	metric.Gauge(monitor.ActorRuntimeTotal).Dec()
	am.actorruntimes.Delete(actorid)
}

func (am *ActorMgr) CallFunc(actorid string, ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ar, ok := am.GetActor(actorid)
	if !ok {
		if am.cfg.AutoCreate == 0 {
			return nil, errs.ActorCantAutoCreate.New("%s", actorid)
		}
		ar = am.NewActor(actorid)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(am.cfg.MaxReqWaitTime)*time.Second)
	defer cancel()
	value := am.actorreqpool.Get()
	actorreq := value.(*actorReq)
	actorreq.ctx = ctx
	actorreq.req = req
	actorreq.info = info
	actorreq.handler = handler
	if err := ar.pushreq(actorreq); err != nil {
		return nil, err
	}
	select {
	case <-ar.ctx.Done():
		return nil, errs.ActorRuntimeCallHandleActorDone.New("Actor %s have exit", actorid)
	case rsp := <-actorreq.rspchan:
		return rsp.rsp, rsp.err
	case <-ctx.Done():
		return nil, errs.ActorRuntimeCallHandleRPCDone.New("Wait Rsp ctx Done")
	}
}

func (am *ActorMgr) KickOne(actorid string) error {
	if actorruntime, ok := am.GetActor(actorid); !ok {
		return errs.ActorKickNoFindActor.New("Not find Actor %s", actorid)
	} else {
		actorruntime.kick()
		return nil
	}
}

func (am *ActorMgr) KickAll() {
	am.actorruntimes.Range(func(key, value interface{}) bool {
		ar := value.(*actorRuntime)
		ar.kick()
		return true
	})
}

func genLeaseId(prefix string, actorid string) string {
	return fmt.Sprintf("%s/%s", prefix, actorid)
}

func parseLeaseId(prefix string, leaseid string) string {
	regexLeaseId, _ := regexp.Compile(fmt.Sprintf("^/%s/(.+)$", prefix))
	groups := regexLeaseId.FindStringSubmatch(leaseid)
	if len(groups) == 0 {
		return ""
	}
	return groups[0]
}
