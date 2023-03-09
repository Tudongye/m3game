package dirapp

import (
	_ "m3game/broker/nats"
	"m3game/demo/dirapp/dirserver"
	dproto "m3game/demo/proto"
	_ "m3game/log/zap"
	"m3game/mesh/router"
	_ "m3game/mesh/router/consul"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/server"
	_ "m3game/shape/sentinel"
	"sync"
)

func newApp() *DirApp {
	return &DirApp{
		App: app.New(dproto.DirAppFuncID),
	}
}

type DirApp struct {
	app.App
}

func (d *DirApp) Start(wg *sync.WaitGroup) error {
	if err := router.Register(d); err != nil {
		return err
	}
	return nil
}
func (d *DirApp) HealthCheck() bool {
	return true
}
func Run() error {
	runtime.Run(newApp(), []server.Server{dirserver.New()})
	return nil
}
