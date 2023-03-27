package actor

import (
	"context"
	"fmt"
	"m3game/plugins/log"
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

func (am *ActorMgr) CreateActor(actorid string) (*actorRuntime, error) {
	if a, ok := am.GetActor(actorid); ok {
		return a, nil
	}
	ar := newActorRuntime(am.actorcreater(actorid), am.cfg)
	ctx, cancel := context.WithCancel(context.Background())
	ar.ctx = ctx
	ar.cancel = cancel
	if a, loaded := am.actorruntimes.LoadOrStore(actorid, ar); loaded {
		return a.(*actorRuntime), nil
	} else {
		ar = a.(*actorRuntime)
		go am.runActor(ar, actorid)
		return ar, nil
	}
}

func (am *ActorMgr) runActor(ar *actorRuntime, actorid string) {
	if err := ar.run(); err != nil {
		log.Error("actor.run() err %s", err.Error())
	}
	ar.cancel()
	am.actorruntimes.Delete(actorid)
}

func (am *ActorMgr) CallFunc(actorid string, ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ar, ok := am.GetActor(actorid)
	if !ok {
		var err error
		if ar, err = am.CreateActor(actorid); err != nil {
			return nil, err
		}
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
		return nil, fmt.Errorf("Actor %s have exit", actorid)
	case rsp := <-actorreq.rspchan:
		return rsp.rsp, rsp.err
	case <-ctx.Done():
		return nil, fmt.Errorf("Wait Rsp ctx Done")
	}
}

func (am *ActorMgr) KickLease(leaseid string) error {
	actorid := parseLeaseId(am.cfg.LeasePrefix, leaseid)
	return am.KickOne(actorid)
}

func (am *ActorMgr) KickOne(actorid string) error {
	if actorruntime, ok := am.GetActor(actorid); !ok {
		return fmt.Errorf("Not find Actor %s", actorid)
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
