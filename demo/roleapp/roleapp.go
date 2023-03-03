package roleapp

import (
	"m3game/app"
	"m3game/broker/nats"
	"m3game/client"
	"m3game/db/cache"
	"m3game/demo/dirapp/dirclient"
	"m3game/demo/loader"
	"m3game/demo/mapapp/mapclient"
	dproto "m3game/demo/proto"
	"m3game/demo/roleapp/rolechclient"
	"m3game/demo/roleapp/rolechserver"
	"m3game/demo/roleapp/roleserver"
	"m3game/mesh/router/consul"
	"m3game/proto"
	"m3game/runtime"
	"m3game/runtime/plugin"
	"m3game/server"
	"sync"
)

func newApp() *RoleApp {
	return &RoleApp{
		App: app.New(dproto.RoleAppFuncID),
	}
}

type RoleApp struct {
	app.App
}

func (d *RoleApp) Start(wg *sync.WaitGroup) error {
	router := plugin.GetRouterPlugin()
	if router != nil {
		if err := router.Register(d); err != nil {
			return err
		}
	}
	if err := dirclient.Init(d.RouteIns(), client.GenMetaClientOption(proto.META_FLAG_FALSE)); err != nil {
		return err
	}
	if err := mapclient.Init(d.RouteIns()); err != nil {
		return err
	}
	if err := rolechclient.Init(d.RouteIns(), client.GenMetaClientOption(proto.META_FLAG_FALSE)); err != nil {
		return err
	}
	return nil
}
func (d *RoleApp) HealthCheck() bool {
	return true
}
func Run() error {
	plugin.RegisterFactory(&consul.Factory{})
	plugin.RegisterFactory(&cache.Factory{})
	plugin.RegisterFactory(&nats.Factory{})

	loader.RegisterLocationCfg()

	runtime.Run(newApp(), []server.Server{roleserver.New(), rolechserver.New()})
	return nil
}
