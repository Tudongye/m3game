package actor

import (
	"context"
	"errors"
	"fmt"
	"m3game/example/proto/pb"
	"m3game/plugins/db"
	"m3game/plugins/db/wraper"
	"m3game/plugins/log"
	"m3game/runtime/server/actor"
)

var (
	_err_actor_dbplugin = errors.New("_err_actor_dbplugin")
)

func ActorCreater(actorid string) actor.Actor {
	return &Actor{
		ActorBase: actor.ActorBaseCreator(actorid),
		ready:     false,
		logp:      log.LogPlus{"ActorID": actorid},
	}
}

func ConvertActor(ctx context.Context) *Actor {
	a := actor.ParseActor(ctx)
	if a == nil {
		return nil
	}
	return a.(*Actor)
}

type Actor struct {
	*actor.ActorBase
	wraper *wraper.Wraper[*pb.ActorDB]
	ready  bool
	logp   log.LogPlus
}

func (a *Actor) Ready() bool {
	return a.ready
}

func (a *Actor) ActorID() string {
	return a.ID()
}

func (a *Actor) OnInit() error {
	log.InfoP(a.logp, "OnInit")
	a.wraper = wraper.New(actormeta, a.ID())
	return nil
}

func (a *Actor) OnTick() error {
	return nil
}

func (a *Actor) OnExit() error {
	log.InfoP(a.logp, "OnExit")
	return nil
}

func (a *Actor) Save() error {
	log.DebugP(a.logp, "Save")
	if a.wraper.HasDirty() {
		log.Debug("Saving %s", a.ID())
		dbp := db.Get()
		if dbp == nil {
			log.Error(_err_actor_dbplugin.Error())
			return _err_actor_dbplugin
		}
		if err := a.wraper.Update(dbp); err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

func (a *Actor) DB() *pb.ActorDB {
	log.Debug("%v", a.wraper.TObj())
	return a.wraper.TObj()
}

func (a *Actor) ModifyName(name string) error {
	log.DebugP(a.logp, "ModifyName %s", name)
	if !a.ready {
		return fmt.Errorf("Actor not Ready")
	}
	if actorname, err := wraper.Getter[*pb.ActorName](a.wraper); err != nil {
		log.Error(err.Error())
		return err
	} else {
		actorname.Name = name
		return wraper.Setter(a.wraper, actorname)
	}
}

func (a *Actor) LvUp() error {
	log.DebugP(a.logp, "LvUp ")
	if !a.ready {
		return fmt.Errorf("Actor not Ready")
	}
	if actorinfo, err := wraper.Getter[*pb.ActorInfo](a.wraper); err != nil {
		log.Error(err.Error())
		return err
	} else {
		actorinfo.Level += 1
		return wraper.Setter(a.wraper, actorinfo)
	}
}
func (a *Actor) Name() string {
	if actorname, err := wraper.Getter[*pb.ActorName](a.wraper); err != nil {
		log.Error(err.Error())
		return ""
	} else {
		return actorname.Name
	}
}

func (a *Actor) Login() error {
	dbplugin := db.Get()
	if dbplugin == nil {
		return _err_actor_dbplugin
	}
	if err := a.wraper.Read(dbplugin); err != nil {
		log.Error(err.Error())
		return _err_actor_dbplugin
	}
	a.ready = true
	return nil
}
