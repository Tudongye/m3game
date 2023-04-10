package onlineapp

import (
	"context"
	"m3game/config"
	"m3game/demo/onlineapp/onlineser"
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
	_ "m3game/plugins/transport/tcptrans"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/mesh"
	"m3game/runtime/server"
	"time"

	_ "m3game/plugins/trace/jaeger"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func newApp() *OnlineApp {
	return &OnlineApp{
		App: app.New(proto.OnlineFuncID),
	}
}

type OnlineApp struct {
	app.App
	cfg AppCfg
}

type AppCfg struct {
	PrePareTime int    `mapstructure:"PrePareTime" validate:"gt=0"`
	VoteLease   string `mapstructure:"VoteLease" validate:"required"` // RoleApp是有状态服务，给他10s重连的机会
}

func (a *OnlineApp) Init(c map[string]interface{}) error {
	if err := mapstructure.Decode(c, &a.cfg); err != nil {
		return errors.Wrap(err, "App Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&a.cfg); err != nil {
		return err
	}
	if err := onlineser.Init(c); err != nil {
		return err
	}
	return nil
}

func (a *OnlineApp) Start(ctx context.Context) {
	log.Info("OnlineApp PrepareTime %d", a.cfg.PrePareTime)
	time.Sleep(time.Duration(a.cfg.PrePareTime) * time.Second)
	log.Info("OnlineApp Ready")
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			// 插件检查
			if router.Instance().Factory().CanUnload(router.Instance()) {
				runtime.ShutDown("Router Delete")
				return
			}
			// 选主逻辑
			if voted, err := a.VoteMain(ctx); err != nil {
				log.Error("VoteMain Fail %s", err.Error())
				continue
			} else if voted {
				// 本地为主备，打开Pool
				if !onlineser.Pool().IsOpen() {
					onlineser.Pool().Open()
				}
			} else {
				// 本地不是主备，关闭Pool
				if onlineser.Pool().IsOpen() {
					onlineser.Pool().Close()
				}
			}
			// 更新活跃RoleApp列表
			if err := onlineser.Pool().LoadAppCache(); err != nil {
				log.Error("%s", err)
			}
			continue
		}
	}
}

// 选主，返回本地是否为主备
func (a *OnlineApp) VoteMain(ctx context.Context) (bool, error) {
	// 临时参数
	appid := config.GetAppID().String()
	svcid := config.GetSvcID().String()
	// 获取实际主备
	leaseAppId := ""
	if lv, err := lease.GetLease(ctx, a.cfg.VoteLease); err != nil {
		return false, errors.Wrapf(err, "Vote GetLease %s", a.cfg.VoteLease)
	} else {
		leaseAppId = string(lv)
	}
	// 计算逻辑主备
	logicAppId := ""
	if routeinss, err := router.Instance().GetAllInstances(svcid); err != nil || len(routeinss) == 0 {
		return false, errors.Wrapf(err, "Vote GetInss %s", svcid)
	} else {
		routehelper := mesh.NewRouteHelper()
		for _, ins := range routeinss {
			routehelper.Add(ins.GetIDStr())
		}
		routehelper.Compress()
		if dstappid, err := routehelper.RouteSingle(); err != nil {
			return false, errors.Wrapf(err, "RouteSingle Fail")
		} else {
			logicAppId = dstappid
		}
	}
	// 逻辑与实际一致
	if logicAppId == leaseAppId {
		if logicAppId == appid {
			return true, nil
		}
		return false, nil
	}
	// 不一致 开始调整
	log.Info("LeaseAppId %s , LogicAppId %s, local %s", leaseAppId, logicAppId, appid)
	if leaseAppId == appid {
		log.Info("Local is LeaseAppId %s, Free Lease %s...", leaseAppId, a.cfg.VoteLease)
		if err := lease.FreeLease(ctx, a.cfg.VoteLease); err != nil {
			log.Error("Local %s Free Lease %s Fail %s", appid, a.cfg.VoteLease, err.Error())
		}
		return false, nil
	} else if logicAppId == appid {
		log.Info("Local is LogicAppId %s, Alloc Lease %s...", leaseAppId, a.cfg.VoteLease)
		if onlineser.Pool().IsOpen() {
			onlineser.Pool().Close()
		}
		if err := lease.AllocLease(ctx, a.cfg.VoteLease, lease.DefaultLeaseMoveOutFunc); err != nil {
			log.Error("Local %s Alloc Lease %s Fail %s", appid, a.cfg.VoteLease, err.Error())
			return false, nil
		}
		return true, nil
	}
	return false, nil
}

func (d *OnlineApp) Alive(app string, svc string) bool {
	return true
}

func Run(ctx context.Context) error {
	return runtime.New().Run(ctx, newApp(), []server.Server{onlineser.New()})
}
