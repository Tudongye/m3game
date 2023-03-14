package actor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
)

func newActorMgr(creater ActorCreater) *ActorMgr {
	return &ActorMgr{
		actormap:     make(map[string]*actorRuntime),
		actorcreater: creater,
	}
}

type ActorMgr struct {
	actormap     map[string]*actorRuntime
	lock         sync.RWMutex
	actorcreater ActorCreater
}

func (am *ActorMgr) getActor(actorid string) (*actorRuntime, bool) {
	am.lock.RLock()
	defer am.lock.RUnlock()
	a, ok := am.actormap[actorid]
	return a, ok
}

func (am *ActorMgr) createActor(actorid string) (*actorRuntime, error) {
	am.lock.Lock()
	defer am.lock.Unlock()
	if a, ok := am.actormap[actorid]; ok {
		return a, nil
	}
	actor := am.actorcreater(actorid)
	if err := actor.OnInit(); err != nil {
		return nil, err
	}
	am.actormap[actorid] = newActorRuntime(actor)
	go func() {
		am.actormap[actorid].run()
		am.lock.Lock()
		defer am.lock.Unlock()
		delete(am.actormap, actorid)
	}()
	return am.actormap[actorid], nil
}

func (am *ActorMgr) CallFunc(actorid string, ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var ar *actorRuntime
	var ok bool
	if ar, ok = am.getActor(actorid); !ok {
		var err error
		if ar, err = am.createActor(actorid); err != nil {
			return nil, err
		}
	}
	actorreq := &actorReq{
		ctx:     ctx,
		req:     req,
		info:    info,
		handler: handler,
		rspchan: make(chan *actorRsp),
	}
	if err := ar.pushreq(actorreq); err != nil {
		return nil, err
	}
	t := time.NewTimer(time.Duration(_cfg.MaxReqWaitTime) * time.Second)
	select {
	case rsp := <-actorreq.rspchan:
		return rsp.rsp, rsp.err
	case <-t.C:
		return nil, fmt.Errorf("Wait Rsp TimeOut")
	}
}

func (am *ActorMgr) KickOne(actorid string) error {
	if actorruntime, ok := am.getActor(actorid); !ok {
		return fmt.Errorf("Not find Actor %s", actorid)
	} else {
		actorruntime.kick()
		return nil
	}
}

func (am *ActorMgr) KickAll() {
	am.lock.RLock()
	defer am.lock.RUnlock()
	for _, actorrumntime := range am.actormap {
		actorrumntime.kick()
	}
}
