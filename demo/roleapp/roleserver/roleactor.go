package roleserver

import (
	"context"
	"fmt"
	"log"
	"m3game/db/cache"
	"m3game/demo/proto/pb"
	"m3game/runtime/plugin"
	"m3game/server/actor"

	"google.golang.org/protobuf/proto"
)

func RoleCreater(actorid string) actor.Actor {
	return &RoleActor{
		roleid:    actorid,
		ActorBase: actor.ActorBaseCreator(),
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
	roleid string
	db     *pb.RoleDB
	dirty  bool
	ready  bool
}

func (a *RoleActor) OnInit() error {
	log.Printf("OnInit %s\n", a.roleid)
	return nil
}

func (a *RoleActor) OnTick() error {
	return nil
}

func (a *RoleActor) OnExit() error {
	log.Printf("OnExit %s\n", a.roleid)
	return nil
}

func (a *RoleActor) Save() error {
	log.Printf("Save %s\n", a.roleid)
	if a.dirty {
		log.Printf("Saving %s\n", a.roleid)
		a.dirty = false
		db := plugin.GetDBPluginByName(cache.Name())
		if db == nil {
			return _err_actor_dberr
		}
		return db.Update(rolemeta, a.roleid, a.db)
	}
	return nil
}

func (a *RoleActor) ReBuild(proto.Message) error {
	log.Printf("ReBuild %s\n", a.roleid)
	return nil
}

func (a *RoleActor) Pack() (proto.Message, error) {
	log.Printf("Pack %s\n", a.roleid)
	return nil, nil
}

func (a *RoleActor) ModifyName(name string) error {
	log.Printf("ModifyName %s %s\n", a.roleid, name)
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
	obj.RoleID = a.roleid
	obj.Name = name
	return db.Insert(rolemeta, a.roleid, obj)
}

func (a *RoleActor) Login() error {
	db := plugin.GetDBPluginByName(cache.Name())
	if db == nil {
		return _err_actor_dberr
	}
	if obj, err := db.Get(rolemeta, a.roleid); err != nil {
		return err
	} else if roleobj, ok := obj.(*pb.RoleDB); !ok {
		return _err_actor_dberr
	} else {
		a.db = roleobj
		a.ready = true
		return nil
	}
}
