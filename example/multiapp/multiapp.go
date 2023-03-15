package multiapp

import (
	"m3game/example/multiapp/multiser"
	"m3game/example/proto"
	_ "m3game/plugins/broker/nats"
	"m3game/plugins/log"
	_ "m3game/plugins/log/zap"
	"m3game/plugins/router"
	_ "m3game/plugins/router/consul"
	_ "m3game/plugins/shape/sentinel"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/server"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

var (
	_cfg AppCfg
)

func newApp() *MultiApp {
	return &MultiApp{
		App:  app.New(proto.MultiAppFuncID),
		exit: make(chan struct{}, 1),
	}
}

type MultiApp struct {
	app.App
	exit chan struct{}
}
type AppCfg struct {
	PrePareTime int `mapstructure:"PrePareTime"`
}

func (c *AppCfg) CheckVaild() error {
	if c.PrePareTime == 0 {
		return errors.New("PrePareTime cant be 0")
	}
	return nil
}

func (a *MultiApp) Init(cfg map[string]interface{}) error {
	if err := mapstructure.Decode(cfg, &_cfg); err != nil {
		return errors.Wrap(err, "App Decode Cfg")
	}
	if err := _cfg.CheckVaild(); err != nil {
		return err
	}
	return nil
}
func (d *MultiApp) Start(wg *sync.WaitGroup) error {
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("MultiApp PrepareTime %d", _cfg.PrePareTime)
		time.Sleep(time.Duration(_cfg.PrePareTime) * time.Second)
		log.Info("MultiApp Ready")
		t := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-d.exit:
				return
			case <-t.C:
				// 插件检查
				if router.Get().Factory().CanDelete(router.Get()) {
					t.Stop()
					runtime.ShutDown()
				}
				continue
			}
		}
	}()
	return nil
}

func (d *MultiApp) Stop() error {
	select {
	case d.exit <- struct{}{}:
		return nil
	default:
		return nil
	}
}

func (d *MultiApp) HealthCheck() bool {
	return true
}

func Run() error {
	runtime.Run(newApp(), []server.Server{multiser.New()})
	return nil
}
