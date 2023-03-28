package actor

import (
	"context"
	"fmt"
	"m3game/example/gateapp/gatecli"
	"m3game/example/proto/pb"
	"m3game/meta"
	"m3game/plugins/db"
	"m3game/plugins/log"
	"m3game/runtime/server/actor"

	"github.com/pkg/errors"
)

const (
	ActorIdMetaKey = "actorid"
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
	wraper   *db.Wraper[*pb.ActorDB, pb.AcFlag]
	ready    bool
	logp     log.LogPlus
	playerid string
	gateapp  string
}

func (a *Actor) SetPlayerId(s string) {
	a.playerid = s
}

func (a *Actor) GetPlayerId() string {
	return a.playerid
}
func (a *Actor) SetGate(s string) {
	a.gateapp = s
}

func (a *Actor) GetGate() string {
	return a.gateapp
}
func (a *Actor) Ready() bool {
	return a.ready
}

func (a *Actor) ActorID() string {
	return a.ID()
}

func (a *Actor) OnInit() error {
	log.InfoP(a.logp, "OnInit")
	a.wraper = actorwrapermeata.New(a.ID())
	return nil
}

func (a *Actor) OnTick() error {
	return nil
}

func (a *Actor) OnExit() error {
	log.InfoP(a.logp, "OnExit")
	if err := gatecli.SendToCli(context.Background(), a.playerid, "Exited", meta.RouteApp(a.gateapp)); err != nil {
		log.ErrorP(a.logp, "SendCli fail %s", err.Error())
	}
	return nil
}

func (a *Actor) OnSave() error {
	log.DebugP(a.logp, "Save")
	if a.wraper.IsDirty() {
		log.DebugP(a.logp, "Saving")
		dbp := db.Instance()
		if dbp == nil {
			log.Error(_err_actor_dbplugin.Error())
			return _err_actor_dbplugin
		}
		if err := a.wraper.Update(context.TODO(), dbp); err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

func (a *Actor) DB() *pb.ActorDB {
	return a.wraper.Obj()
}

func (a *Actor) ModifyName(name string) error {
	log.DebugP(a.logp, "ModifyName %s", name)
	if !a.ready {
		return fmt.Errorf("Actor not Ready")
	}
	a.wraper.Set(pb.AcFlag_FActorName, name)
	return nil
}

func (a *Actor) LvUp() error {
	log.DebugP(a.logp, "LvUp ")
	if !a.ready {
		return fmt.Errorf("Actor not Ready")
	}
	lv := a.wraper.Get(pb.AcFlag_FActorLevel).(int32)
	a.wraper.Set(pb.AcFlag_FActorLevel, lv+1)
	return nil
}

func (a *Actor) Name() string {
	return a.wraper.Get(pb.AcFlag_FActorName).(string)
}

func (a *Actor) Login(ctx context.Context) error {
	dbplugin := db.Instance()
	if dbplugin == nil {
		return _err_actor_dbplugin
	}
	if err := a.wraper.Read(ctx, dbplugin); err != nil {
		log.Error("%s %s", err.Error(), a.ActorID())
		return _err_actor_dbplugin
	}
	a.ready = true
	return nil
}
