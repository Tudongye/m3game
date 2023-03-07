package roleserver

import (
	"context"
	"fmt"
	"m3game/db/wraper"
	"m3game/demo/proto/pb"
	"m3game/runtime/plugin"
	"m3game/server/actor"
	"m3game/util/log"

	"google.golang.org/protobuf/proto"
)

func roleCreater(actorid string) actor.Actor {
	return &RoleActor{
		ActorBase: actor.ActorBaseCreator(actorid),
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
	wraper *wraper.Wraper[*pb.RoleDB]
	ready  bool
}

func (a *RoleActor) RoleID() string {
	return a.ID()
}

func (a *RoleActor) OnInit() error {
	log.InfoP(a.ID(), "OnInit")
	a.wraper = wraper.New(rolemeta, a.ID())
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
	if a.wraper.HasDirty() {
		log.Debug("Saving %s", a.ID())
		db := plugin.GetDBPlugin()
		if db == nil {
			return _err_actor_dberr
		}
		if err := a.wraper.Update(db); err != nil {
			log.Error(err.Error())
			return err
		}
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
	if rolename, err := wraper.Getter[*pb.RoleName](a.wraper); err != nil {
		log.Error(err.Error())
		return err
	} else {
		rolename.Name = name
		return wraper.Setter(a.wraper, rolename)
	}
}

func (a *RoleActor) ModifyLocation(locationname string, location int32) error {
	log.DebugP(a.ID(), "ModifyLocation %s %d", locationname, location)
	if !a.ready {
		return fmt.Errorf("Actor not Ready")
	}
	if locationinfo, err := wraper.Getter[*pb.LocationInfo](a.wraper); err != nil {
		log.Error(err.Error())
		return err
	} else {
		locationinfo.LocateName = locationname
		locationinfo.Location = location
		return wraper.Setter(a.wraper, locationinfo)
	}
}

func (a *RoleActor) Name() string {
	if rolename, err := wraper.Getter[*pb.RoleName](a.wraper); err != nil {
		log.Error(err.Error())
		return ""
	} else {
		return rolename.Name
	}
}

func (a *RoleActor) Register(name string) error {
	dbplugin := plugin.GetDBPlugin()
	if dbplugin == nil {
		return _err_actor_dberr
	}
	if err := wraper.KeySetter(a.wraper, a.ID()); err != nil {
		log.Error(err.Error())
		return _err_actor_dberr
	}
	if err := wraper.Setter(a.wraper, &pb.RoleName{
		Name: name,
	}); err != nil {
		log.Error(err.Error())
		return _err_actor_dberr
	}
	if err := wraper.Setter(a.wraper, &pb.LocationInfo{
		Location:   0,
		LocateName: "",
	}); err != nil {
		log.Error(err.Error())
		return _err_actor_dberr
	}
	return a.wraper.Create(dbplugin)
}

func (a *RoleActor) Login() error {
	dbplugin := plugin.GetDBPlugin()
	if dbplugin == nil {
		return _err_actor_dberr
	}
	if err := a.wraper.Read(dbplugin); err != nil {
		log.Error(err.Error())
		return _err_actor_dberr
	}
	a.ready = true
	return nil
}
