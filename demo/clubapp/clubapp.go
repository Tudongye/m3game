package clubapp

import (
	"context"
	"fmt"
	"m3game/config"
	"m3game/demo/clubapp/club"
	"m3game/demo/clubapp/clubdcli"
	"m3game/demo/clubapp/clubdser"
	"m3game/demo/clubapp/clubser"
	"m3game/demo/proto"
	_ "m3game/plugins/broker/nats"
	_ "m3game/plugins/db/mongo"
	"m3game/plugins/lease"
	_ "m3game/plugins/lease/etcd"
	"m3game/plugins/log"
	_ "m3game/plugins/log/zap"
	_ "m3game/plugins/metric/prometheus"
	"m3game/plugins/router"
	_ "m3game/plugins/router/consul"
	_ "m3game/plugins/shape/sentinel"
	_ "m3game/plugins/trace/jaeger"
	_ "m3game/plugins/transport/http2trans"
	_ "m3game/plugins/transport/natstrans"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/mesh"
	"m3game/runtime/server"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func newApp() *ClubApp {
	return &ClubApp{
		App: app.New(proto.ClubFuncID),
	}
}

var (
	_cfg AppCfg
)

type ClubApp struct {
	app.App
	clubser *clubser.ClubSer
}

type AppCfg struct {
	PrePareTime int `mapstructure:"PrePareTime" validate:"gt=0"`
}

func (a *ClubApp) Init(c map[string]interface{}) error {
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return errors.Wrap(err, "App Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&_cfg); err != nil {
		return err
	}
	return nil
}

func (d *ClubApp) Prepare(ctx context.Context) error {
	lease.SetReciver(d)
	if _, err := clubdcli.New(config.GetAppID()); err != nil {
		return err
	}
	return nil
}

func (d *ClubApp) Start(ctx context.Context) {
	log.Info("ClubApp PrepareTime %d", _cfg.PrePareTime)
	time.Sleep(time.Duration(_cfg.PrePareTime) * time.Second)
	log.Info("ClubApp Ready")
	t := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			// 插件检查
			if router.Instance().Factory().CanUnload(router.Instance()) {
				t.Stop()
				runtime.ShutDown("Router Delete")
			}
			if lease.Instance().Factory().CanUnload(lease.Instance()) {
				t.Stop()
				runtime.ShutDown("Lease Delete")
			}
			// Slot检查
			if slots, err := d.VoteSlots(ctx); err != nil {
				log.Error("%s", err.Error())
			} else {
				actormgr := d.clubser.ActorMgr()
				for _, slotid := range slots {
					actormgr.NewActor(slotid)
				}
			}
			continue
		}
	}
}

// 获取本地Slot列表
func (d *ClubApp) VoteSlots(ctx context.Context) ([]string, error) {
	appid := config.GetAppID().String()
	svcid := config.GetSvcID().String()
	var slots []string
	if routeinss, err := router.Instance().GetAllInstances(svcid); err != nil || len(routeinss) == 0 {
		return slots, errors.Wrapf(err, "Vote GetInss %s", svcid)
	} else {
		routehelper := mesh.NewRouteHelper()
		for _, ins := range routeinss {
			routehelper.Add(ins.GetAppID())
		}
		routehelper.Compress()
		// 遍历
		for i := 0; i < club.ClubSlotNum; i++ {
			slotid := fmt.Sprintf("%d", i)
			if logicappid, err := routehelper.RouteHash(slotid); err != nil {
				return slots, err
			} else if logicappid == appid {
				slots = append(slots, slotid)
			}
		}
	}
	return slots, nil
}

func (d *ClubApp) Alive(app string, svc string) bool {
	return true
}

func (d *ClubApp) PreExitLease(ctx context.Context, id string) ([]byte, error) {
	return nil, clubser.Ser().ActorMgr().KickLease(id)
}

func (d *ClubApp) SendKickLease(ctx context.Context, id string, app string) ([]byte, error) {
	return clubdcli.Kick(ctx, id, app)
}

func Run(ctx context.Context) error {
	app := newApp()
	app.clubser = clubser.New()
	runtime.New().Run(ctx, app, []server.Server{app.clubser, clubdser.New()})
	return nil
}
