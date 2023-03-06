package runtime

import (
	"context"
	"fmt"
	"m3game/app"
	"m3game/config"
	_ "m3game/mesh/balance"
	"m3game/resource"
	"m3game/runtime/plugin"
	"m3game/runtime/transport"
	"m3game/server"
	"m3game/util"
	"m3game/util/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	_runtime *Runtime
)

func init() {
	_runtime = &Runtime{
		servers: make(map[string]server.Server),
	}
}

type RuntimeCfg struct {
	Resource  map[string]interface{}            `toml:"Resource"`
	App       map[string]interface{}            `toml:"App"`
	Server    map[string]map[string]interface{} `toml:"Server"`
	Transport map[string]interface{}            `toml:"Transport"`
}

func SendInterFunc(sctx *transport.Sender) error {
	return transport.SendInterFunc(sctx)
}

func ShutDown() error {
	_runtime.cancel()
	return nil
}

func Reload() error {
	err := reload()
	if err != nil {
		log.Error("Runtime.Reload Fail:%s", err.Error())
	}
	return nil
}

func PreExit() error {
	return nil
}

func Run(app app.App, servers []server.Server) error {
	ctx, cancel := context.WithCancel(context.Background())

	log.Fatal("Runtime.Init...")
	_runtime.cancel = cancel
	_runtime.app = app
	for _, server := range servers {
		if err := _runtime.registerServer(server); err != nil {
			log.Error("registerServer %s err %s", server.Name(), err.Error())
			return err
		}
	}

	log.Fatal("Resource.Load...")
	v := *config.GetConfig()
	var cfg RuntimeCfg
	if err := v.Unmarshal(&cfg); err != nil {
		log.Error("UnMarshal RuntimeCfg %s", err.Error())
		return err
	}

	log.Fatal("Resource.Load...")
	if err := resource.Init(cfg.Resource); err != nil {
		log.Error("Runtime.Resource.Init %s err %s", cfg.Resource, err.Error())
		return err
	}

	log.Fatal("Transport.Init...")
	if err := transport.Init(cfg.Transport, _runtime); err != nil {
		log.Error("Transport.Init err %s", err.Error())
		return err
	}

	log.Fatal("Plugin.Init...")
	if err := plugin.InitPlugins(v); err != nil {
		log.Error("InitPlugins err %s", err.Error())
		return err
	}

	log.Fatal("Server.Init...")
	for _, server := range servers {
		log.Info("Server.Init.%s...", server.Name())
		if err := server.Init(cfg.Server[string(server.Type())], app); err != nil {
			log.Error("Server.Init %s err %s", server.Name(), err.Error())
			return err
		}
		if err := transport.RegisterServer(server.TransportRegister()); err != nil {
			log.Error("Transport.RegisterServer %s err %s", server.Name(), err.Error())
			return err
		}
	}

	log.Fatal("App.Init...")
	if err := app.Init(cfg.App); err != nil {
		log.Error("App.Init err %s", err.Error())
		return err
	}
	var wg sync.WaitGroup

	log.Fatal("Transport.Start...")
	if err := transport.Start(&wg); err != nil {
		log.Error("Transport.Start err %s", err.Error())
		return err
	}

	log.Fatal("Server.Start...")
	for _, server := range servers {
		log.Info("Server.Start.%s...", server.Name())
		if err := server.Start(&wg); err != nil {
			log.Error("Server.Start %s err %s", err.Error())
			return err
		}
	}

	log.Fatal("App.Start.%s...", app.IDStr())
	if err := app.Start(&wg); err != nil {
		log.Error("App.Start err %s", err.Error())
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			_runtime.app.Stop()
			for _, server := range _runtime.servers {
				server.Stop()
			}
		}
	}()
	go signalProc()

	log.Fatal("Wait...")
	wg.Wait()
	return nil
}

type Runtime struct {
	app     app.App
	servers map[string]server.Server
	cancel  context.CancelFunc
}

func (r *Runtime) HealthCheck(idstr string) bool {
	if r.app.IDStr() != idstr {
		return false
	}
	return r.app.HealthCheck()
}

func (r *Runtime) RecvInterFunc(recv *transport.Reciver) (resp interface{}, err error) {
	for _, server := range r.servers {
		if server == recv.Info().Server {
			return r.app.RecvInterFunc(recv, server.RecvInterFunc)
		}
	}
	return nil, fmt.Errorf("Can't find Server")
}

func (r *Runtime) registerServer(s server.Server) error {
	if _, ok := r.servers[s.Name()]; ok {
		return fmt.Errorf("Register repeated ServerName %s", string(s.Name()))
	}
	r.servers[s.Name()] = s
	return nil
}

func signalProc() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, util.SysCallSIGUSR1(), util.SysCallSIGUSR2())

	for s := range sigs {
		log.Fatal("Recv sig %s", s.String())
		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT:
			ShutDown()
		case util.SysCallSIGUSR1():
			Reload()
		case util.SysCallSIGUSR2():
			PreExit()
		}
	}
}

func reload() error {
	log.Fatal("Runtime.Reload...")
	log.Fatal("Config.Reload...")
	if err := config.Reload(); err != nil {
		return err
	}
	v := *config.GetConfig()
	var cfg RuntimeCfg
	if err := v.Unmarshal(&cfg); err != nil {
		log.Error("UnMarshal RuntimeCfg %s", err.Error())
		return err
	}

	log.Fatal("Resource.Reload...")
	if err := resource.ReLoad(cfg.Resource); err != nil {
		return err
	}
	log.Fatal("Transport.Reload...")
	if err := transport.Reload(cfg.Transport); err != nil {
		return err
	}
	log.Fatal("Plugin.Reload...")
	if err := plugin.Reload(v); err != nil {
		return err
	}
	log.Fatal("Server.Reload...")
	for _, server := range _runtime.servers {
		if err := server.Reload(cfg.Server[string(server.Type())]); err != nil {
			return err
		}
	}
	log.Fatal("App.Reload...")
	if err := _runtime.app.Reload(cfg.App); err != nil {
		return err
	}
	log.Fatal("Runtime.Reload Succ...")
	return nil
}
