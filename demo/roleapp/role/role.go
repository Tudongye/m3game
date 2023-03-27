package role

import (
	"context"
	"fmt"
	"m3game/config"
	"m3game/demo/onlineapp/onlinecli"
	"m3game/demo/proto/pb"
	"m3game/demo/roleapp/rolecli"
	"m3game/meta"
	"m3game/plugins/db"
	"m3game/plugins/log"
	"m3game/runtime/server/actor"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var (
	_err_actor_dbplugin = errors.New("_err_actor_dbplugin")
)

func RoleCreater(actorid string) actor.Actor {
	return &Role{
		ActorBase: actor.ActorBaseCreator(actorid),
		ready:     false,
		logp:      log.LogPlus{"RoleID": actorid},
	}
}

func ConvertRole(ctx context.Context) *Role {
	a := actor.ParseActor(ctx)
	if a == nil {
		return nil
	}
	return a.(*Role)
}

type Role struct {
	*actor.ActorBase
	wraper *db.Wraper[*pb.RoleDB, pb.RFlag]
	ready  bool
	logp   log.LogPlus

	gateapp string
}

func (a *Role) SetGate(g string) {
	a.gateapp = g
}

func (a *Role) GetGate() string {
	return a.gateapp
}

func (a *Role) Ready() bool {
	return a.ready
}

func (a *Role) ActorID() string {
	return a.ID()
}

func (a *Role) OnInit() error {
	log.InfoP(a.logp, "OnInit")
	if roleid, err := strconv.ParseInt(a.ID(), 10, 64); err != nil {
		return err
	} else {
		a.wraper = _rolewrapermeta.New(roleid)
	}
	return nil
}

func (a *Role) OnTick() error {
	if !a.ready {
		a.Exit()
	}
	return nil
}

func (a *Role) OnExit() error {
	log.InfoP(a.logp, "OnExit")
	// 向Online反注册
	if a.ready {
		if err := onlinecli.OnlineCreate(context.Background(), a.wraper.Obj().RoleId, config.GetAppID().String()); err != nil {
			return err
		}
	}
	return nil
}

func (a *Role) OnSave() error {
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

func (a *Role) DB() *pb.RoleDB {
	return a.wraper.Obj()
}

func (a *Role) ModifyName(name string) error {
	log.DebugP(a.logp, "ModifyName %s", name)
	if !a.ready {
		return fmt.Errorf("Role not Ready")
	}
	a.wraper.Set(pb.RFlag_RName, name)
	return nil
}

func (a *Role) PowerUp(up int32) error {
	log.DebugP(a.logp, "PowerUp ")
	if !a.ready {
		return fmt.Errorf("Role not Ready")
	}
	power := a.wraper.Get(pb.RFlag_RPower).(int32)
	a.wraper.Set(pb.RFlag_RPower, power+up)
	return nil
}

func (a *Role) GetInfo() (*pb.RoleDB, *pb.ClubRoleDB, error) {
	log.DebugP(a.logp, "GetInfo ")
	if !a.ready {
		return nil, nil, fmt.Errorf("Role not Ready")
	}
	return a.wraper.Obj(), nil, nil
}

func (a *Role) Login(ctx context.Context) error {
	// 查询Online
	if appid, err := onlinecli.OnlineRead(ctx, a.DB().RoleId); err != nil {
		log.Error("%s", err.Error())
		return err
	} else if appid != "" && appid != config.GetAppID().String() {
		// 踢下线
		if err := rolecli.RoleKick(ctx, a.DB().RoleId, meta.RouteApp(appid)); err != nil {
			log.Error("%s", err.Error())
			return err
		}
		time.Sleep(3 * time.Second)
	}
	// 向Online注册
	if err := onlinecli.OnlineCreate(ctx, a.DB().RoleId, config.GetAppID().String()); err != nil {
		log.Error("%s", err.Error())
		return err
	}
	dbp := db.Instance()
	if err := a.wraper.Read(ctx, dbp); err != nil {
		if db.IsErrKeyNotFound(err) {
			// 未注册，
			a.wraper.Set(pb.RFlag_RName, fmt.Sprintf("Role%d", a.DB().RoleId))
			a.wraper.Set(pb.RFlag_RPower, int32(0))
			// DB写入失败
			if err := a.wraper.Create(ctx, dbp); err != nil {
				log.Error("%s", err.Error())
				return err
			}
		} else {
			log.Error("%s", err.Error())
			return err
		}
	}
	a.ready = true
	return nil
}

func (a *Role) Kick(ctx context.Context) error {
	log.DebugP(a.logp, "Kick ")
	a.Exit()
	return nil
}
