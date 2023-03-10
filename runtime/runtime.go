package runtime

import (
	"context"
	"fmt"
	"m3game/config"
	"m3game/plugins/log"
	"m3game/plugins/router"
	"m3game/plugins/shape"
	"m3game/plugins/trace"
	"m3game/runtime/app"
	"m3game/runtime/mesh"
	"m3game/runtime/plugin"
	"m3game/runtime/resource"
	"m3game/runtime/server"
	"m3game/runtime/transport"
	"m3game/util"
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

type RuntimeOptionCfg struct {
	Mesh     map[string]interface{} `toml:"Mesh"`
	Shape    map[string]interface{} `toml:"Shape"`
	Resource map[string]interface{} `toml:"Resource"`
}

type RuntimeCfg struct {
	Transport map[string]interface{}            `toml:"Transport"`
	Options   RuntimeOptionCfg                  `toml:"Options"`
	App       map[string]interface{}            `toml:"App"`
	Server    map[string]map[string]interface{} `toml:"Server"`
}

func SendInterFunc(sctx *transport.Sender) error {
	return transport.SendInterFunc(sctx)
}

func ShutDown() error {
	log.Info("ShutDown...")
	_runtime.cancel()
	return nil
}

func Reload() error {
	log.Info("Reload...")
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

	log.Info("Runtime.Init...")
	_runtime.cancel = cancel
	_runtime.app = app
	for _, server := range servers {
		if err := _runtime.registerServer(server); err != nil {
			log.Error("registerServer %s err %s", server.Name(), err.Error())
			return err
		}
	}

	log.Info("RuntimeCfg.Load...")
	v := *config.GetConfig()
	var cfg RuntimeCfg
	if err := v.Unmarshal(&cfg); err != nil {
		log.Error("UnMarshal RuntimeCfg err %s", err.Error())
		return err
	}

	log.Info("Mesh.Init...")
	if err := mesh.Init(cfg.Options.Mesh); err != nil {
		log.Error("Mesh.Init err %s", err.Error())
		return err

	}

	log.Info("Resource.Init...")
	if err := resource.Init(cfg.Options.Resource); err != nil {
		log.Error("Runtime.Resource.Init %s err %s", cfg.Options.Resource, err.Error())
		return err
	}

	log.Info("Transport.Init...")
	if err := transport.Init(cfg.Transport, _runtime); err != nil {
		log.Error("Transport.Init err %s", err.Error())
		return err
	}

	log.Info("Plugins.Init...")
	if err := plugin.InitPlugins(v); err != nil {
		log.Error("InitPlugins err %s", err.Error())
		return err
	}

	log.Info("SetupPluginInterceptor...")
	if err := setupPluginInterceptor(cfg); err != nil {
		log.Error("setupPluginInterceptor err %s", err.Error())
		return err
	}

	log.Info("Transport.CreateSer...")
	if err := transport.CreateSer(); err != nil {
		log.Error("Transport.CreateSer err %s", err.Error())
		return err
	}

	log.Info("Server.Init...")
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

	log.Info("App.Init...")
	if err := app.Init(cfg.App); err != nil {
		log.Error("App.Init err %s", err.Error())
		return err
	}
	var wg sync.WaitGroup

	log.Info("Transport.Start...")
	if err := transport.Start(&wg); err != nil {
		log.Error("Transport.Start err %s", err.Error())
		return err
	}

	log.Info("Server.Start...")
	for _, server := range servers {
		log.Info("Server.Start.%s...", server.Name())
		if err := server.Start(&wg); err != nil {
			log.Error("Server.Start %s err %s", err.Error())
			return err
		}
	}

	log.Info("App.Start.%s...", app.IDStr())
	if err := app.Start(&wg); err != nil {
		log.Error("App.Start err %s", err.Error())
		return err
	}

	log.Info("Router.Register...")
	if err := router.Register(app); err != nil {
		log.Error("Router.Register err %s", err.Error())
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			log.Info("Recv Done...")
			log.Info("App.Stop...")
			_runtime.app.Stop()
			for _, server := range _runtime.servers {
				log.Info("Server.Stop %s...", server.Name())
				server.Stop()
			}
			log.Info("Transport.Stop...")
			transport.ShutDown()
			log.Info("Doned")
		}
	}()
	go signalProc()

	log.Info("Wait...")
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
		log.Info("Recv sig %s", s.String())
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
	log.Info("Runtime.Reload...")
	log.Info("Config.Reload...")
	if err := config.Reload(); err != nil {
		return err
	}
	v := *config.GetConfig()
	var cfg RuntimeCfg
	if err := v.Unmarshal(&cfg); err != nil {
		log.Error("UnMarshal RuntimeCfg %s", err.Error())
		return err
	}

	log.Info("Resource.Reload...")
	if err := resource.ReLoad(cfg.Options.Resource); err != nil {
		return err
	}
	log.Info("Transport.Reload...")
	if err := transport.Reload(cfg.Transport); err != nil {
		return err
	}
	log.Info("Plugin.Reload...")
	if err := plugin.Reload(v); err != nil {
		return err
	}
	log.Info("Server.Reload...")
	for _, server := range _runtime.servers {
		if err := server.Reload(cfg.Server[string(server.Type())]); err != nil {
			return err
		}
	}
	log.Info("App.Reload...")
	if err := _runtime.app.Reload(cfg.App); err != nil {
		return err
	}
	log.Info("Runtime.Reload Succ...")
	return nil
}

func setupPluginInterceptor(cfg RuntimeCfg) error {
	log.Info("SetupPluginInterceptor.Shape...")
	if err := shape.Setup(cfg.Options.Shape); err != nil {
		return err
	}
	transport.RegisterServerInterceptor(shape.ServerInterceptor())
	transport.RegisterClientInterceptor(shape.ClientInterceptor())

	log.Info("SetupPluginInterceptor.Trace...")
	transport.RegisterServerInterceptor(trace.ServerInterceptor())
	transport.RegisterClientInterceptor(trace.ClientInterceptor())

	return nil
}
