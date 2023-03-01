package roleserver

import (
	"context"
	"fmt"
	"m3game/db/cache"
	"m3game/demo/proto/pb"
	"m3game/runtime/plugin"
	"m3game/server/actor"
	"m3game/util/log"

	"google.golang.org/protobuf/proto"
)

func roleCreater(actorid string) actor.Actor {
	return &RoleActor{
		ActorBase: actor.ActorBaseCreator(actorid),
		dirty:     false,
		ready:     false,
	}
}

func ParseRoleActor(ctx context.Context) *RoleActor {
	a := actor.ParseActor(ctx)
	if a == nil {
		return nil
	}
	return a.(*RoleActor)
}

type RoleActor struct {
	*actor.ActorBase
	db    *pb.RoleDB
	dirty bool
	ready bool
}

func (a *RoleActor) RoleID() string {
	return a.ID()
}

func (a *RoleActor) OnInit() error {
	log.InfoP(a.ID(), "OnInit")
	return nil
}

func (a *RoleActor) OnTick() error {
	return nil
}

func (a *RoleActor) OnExit() error {
	log.InfoP(a.ID(), "OnExit")
	return nil
}

func (a *RoleActor) Save() error {
	log.DebugP(a.ID(), "Save")
	if a.dirty {
		log.Debug("Saving %s", a.ID())
		a.dirty = false
		db := plugin.GetDBPluginByName(cache.Name())
		if db == nil {
			return _err_actor_dberr
		}
		return db.Update(rolemeta, a.ID(), a.db)
	}
	return nil
}

func (a *RoleActor) ReBuild(proto.Message) error {
	log.DebugP(a.ID(), "ReBuild")
	return nil
}

func (a *RoleActor) Pack() (proto.Message, error) {
	log.DebugP(a.ID(), "Pack")
	return nil, nil
}

func (a *RoleActor) ModifyName(name string) error {
	log.DebugP(a.ID(), "ModifyName %s", name)
	if !a.ready {
		return fmt.Errorf("Actor not Ready")
	}
	a.db.Name = name
	a.dirty = true
	return nil
}
func (a *RoleActor) Name() string {
	return a.db.Name
}

func (a *RoleActor) Register(name string) error {
	db := plugin.GetDBPluginByName(cache.Name())
	if db == nil {
		return _err_actor_dberr
	}
	obj := roleDBCreater()
	obj.RoleID = a.ID()
	obj.Name = name
	return db.Insert(rolemeta, a.ID(), obj)
}

func (a *RoleActor) Login() error {
	db := plugin.GetDBPluginByName(cache.Name())
	if db == nil {
		return _err_actor_dberr
	}
	if obj, err := db.Get(rolemeta, a.ID()); err != nil {
		return err
	} else if roleobj, ok := obj.(*pb.RoleDB); !ok {
		return _err_actor_dberr
	} else {
		a.db = roleobj
		a.ready = true
		return nil
	}
}
