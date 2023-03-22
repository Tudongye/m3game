package uidapp

import (
	"context"
	"m3game/config"
	"m3game/demo/proto"
	"m3game/demo/uidapp/uidser"
	_ "m3game/plugins/broker/nats"
	_ "m3game/plugins/db/redis"
	"m3game/plugins/lease"
	_ "m3game/plugins/lease/etcd"
	"m3game/plugins/log"
	_ "m3game/plugins/log/zap"
	"m3game/plugins/router"
	_ "m3game/plugins/router/consul"
	_ "m3game/plugins/shape/sentinel"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/mesh"
	"m3game/runtime/server"
	"m3game/util"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

var (
	_cfg AppCfg
)

func newApp() *UidApp {
	return &UidApp{
		App:  app.New(proto.UidFuncID),
		exit: make(chan struct{}, 1),
	}
}

type UidApp struct {
	app.App
	exit chan struct{}
}

type AppCfg struct {
	PrePareTime int    `mapstructure:"PrePareTime"`
	VoteLease   string `mapstructure:"VoteLease"`
}

func (c *AppCfg) checkValid() error {
	if err := util.InEqualInt(c.PrePareTime, 0, "PrePareTime"); err != nil {
		return err
	}
	if err := util.InEqualStr(c.VoteLease, "", "VoteLease"); err != nil {
		return err
	}
	return nil
}

func (a *UidApp) Init(cfg map[string]interface{}) error {
	if err := mapstructure.Decode(cfg, &_cfg); err != nil {
		return errors.Wrap(err, "App Decode Cfg")
	}
	if err := _cfg.checkValid(); err != nil {
		return err
	}
	if err := uidser.Init(cfg); err != nil {
		return err
	}
	return nil
}

func (d *UidApp) Start(ctx context.Context) {
	log.Info("UidApp PrepareTime %d", _cfg.PrePareTime)
	time.Sleep(time.Duration(_cfg.PrePareTime) * time.Second)
	log.Info("UidApp Ready")
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			// 插件检查
			if router.Get().Factory().CanDelete(router.Get()) {
				runtime.ShutDown("Router Delete")
				return
			}
			// 选主逻辑
			if d.VoteMain(ctx) {
				// 本地为主备，打开Pool
				if !uidser.Pool().IsOpen() {
					uidser.Pool().Open()
				}
			} else {
				// 本地不是主备，关闭Pool
				if uidser.Pool().IsOpen() {
					uidser.Pool().Close()
				}
			}
			continue
		}
	}
}

// 选主，返回本地是否为主备
func (d *UidApp) VoteMain(ctx context.Context) bool {
	// 获取实际主备
	leaseAppId := ""
	if lv, err := lease.GetLease(ctx, _cfg.VoteLease); err != nil {
		log.Error("Vote GetLease %s Fail %s", _cfg.VoteLease, err.Error())
		return false
	} else {
		leaseAppId = string(lv)
	}
	// 计算逻辑主备
	logicAppId := ""
	if routeinss, err := router.GetAllInstances(config.GetSvcID().String()); err != nil || len(routeinss) == 0 {
		log.Error("Vote GetInss %s Fail %s", config.GetSvcID().String(), err.Error())
		return false
	} else {
		routehelper := mesh.NewRouteHelper()
		for _, ins := range routeinss {
			routehelper.Add(ins.GetIDStr())
		}
		routehelper.Compress()
		if dstappid, err := routehelper.RouteSingle(); err != nil {
			log.Error("RouteSingle Fail %s", err.Error())
			return false
		} else {
			logicAppId = dstappid
		}
	}
	// 逻辑与实际一致
	if logicAppId == leaseAppId {
		if logicAppId == config.GetAppID().String() {
			return true
		}
		return false
	}
	// 不一致 开始调整
	log.Info("LeaseAppId %s , LogicAppId %s, local %s", leaseAppId, logicAppId, config.GetAppID().String())
	if leaseAppId == config.GetAppID().String() {
		log.Info("Local is LeaseAppId %s, Free Lease %s...", leaseAppId, _cfg.VoteLease)
		if err := lease.FreeLease(ctx, _cfg.VoteLease); err != nil {
			log.Error("Local %s Free Lease %s Fail %s", config.GetAppID().String(), _cfg.VoteLease, err.Error())
		}
		return false
	} else if logicAppId == config.GetAppID().String() {
		log.Info("Local is LogicAppId %s, Alloc Lease %s...", leaseAppId, _cfg.VoteLease)
		if uidser.Pool().IsOpen() {
			uidser.Pool().Close()
		}
		if err := lease.AllocLease(ctx, _cfg.VoteLease, lease.DefaultLeaseMoveOutFunc); err != nil {
			log.Error("Local %s Alloc Lease %s Fail %s", config.GetAppID().String(), _cfg.VoteLease, err.Error())
			return false
		}
		return true
	}
	return false
}

func (d *UidApp) HealthCheck() bool {
	return true
}

func Run(ctx context.Context) error {
	runtime.Run(ctx, newApp(), []server.Server{uidser.New()})
	return nil
}
