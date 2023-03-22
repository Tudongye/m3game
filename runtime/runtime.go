package runtime

import (
	"context"
	"fmt"
	"m3game/config"
	"m3game/meta"
	"m3game/plugins/log"
	"m3game/plugins/router"
	"m3game/plugins/shape"
	"m3game/plugins/trace"
	"m3game/runtime/app"
	"m3game/runtime/client"
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

	"google.golang.org/grpc"
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

func ShutDown(s string) error {
	log.Info("ShutDown %s...", s)
	_runtime.cancel()
	return nil
}

func Reload() error {
	log.Info("Reload...")
	err := _runtime.reload()
	if err != nil {
		log.Error("Runtime.Reload Fail:%s", err.Error())
	}
	return nil
}

func PreExit() error {
	return nil
}

func Addr() string {
	return _runtime.transport.Addr()
}

func Run(c context.Context, app app.App, servers []server.Server) error {
	ctx, cancel := context.WithCancel(c)

	log.Info("config.Init...")
	config.Init()

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
	if resourcemgr, err := resource.New(cfg.Options.Resource); err != nil {
		log.Error("Runtime.Resource.Init %s err %s", cfg.Options.Resource, err.Error())
		return err
	} else {
		_runtime.resourcemgr = resourcemgr
		if err := _runtime.resourcemgr.ReLoad(cfg.Options.Resource); err != nil {
			log.Error("Runtime.Resource.Reload %s err %s", cfg.Options.Resource, err.Error())
			return err
		}
	}

	log.Info("Transport.New...")
	if trans, err := transport.New(cfg.Transport, _runtime); err != nil {
		log.Error("Transport.New err %s", err.Error())
		return err
	} else {
		_runtime.transport = *trans
	}

	log.Info("Plugins.Init...")
	if err := plugin.InitPlugins(v); err != nil {
		log.Error("InitPlugins err %s", err.Error())
		return err
	}

	log.Info("SetupPluginInterceptor...")
	if err := _runtime.setupPluginInterceptor(cfg); err != nil {
		log.Error("setupPluginInterceptor err %s", err.Error())
		return err
	}

	log.Info("Server.Init...")
	for _, server := range servers {
		log.Info("Server.Init.%s...", server.Name())
		if err := server.Init(cfg.Server[string(server.Type())], app); err != nil {
			log.Error("Server.Init %s err %s", server.Name(), err.Error())
			return err
		}
	}

	log.Info("App.Init...")
	if err := app.Init(cfg.App); err != nil {
		log.Error("App.Init err %s", err.Error())
		return err
	}
	var wg sync.WaitGroup

	log.Info("Transport.Prepare...")
	if err := _runtime.transport.Prepare(ctx); err != nil {
		log.Error("Transport.Prepare err %s", err.Error())
		return err
	}

	log.Info("Server.Register...")
	for _, server := range servers {
		if err := _runtime.transport.RegisterServer(server.TransportRegister()); err != nil {
			log.Error("Transport.RegisterServer %s err %s", server.Name(), err.Error())
			return err
		}
	}

	log.Info("Transport.Start...")
	wg.Add(1)
	go func() {
		defer wg.Done()
		_runtime.transport.Start(ctx)
	}()

	log.Info("Server.Prepare...")
	for _, ser := range servers {
		log.Info("Server.Prepare.%s...", ser.Name())
		if err := ser.Prepare(ctx); err != nil {
			log.Error("Server.Prepare %s err %s", err.Error())
			return err
		}
	}
	log.Info("Server.Start...")
	for _, ser := range servers {
		log.Info("Server.Start.%s...", ser.Name())
		wg.Add(1)
		go func(s server.Server) {
			defer wg.Done()
			s.Start(ctx)
		}(ser)
	}

	log.Info("App.Prepare...")
	if err := app.Prepare(ctx); err != nil {
		log.Error("App.Prepare err %s", err.Error())
		return err
	}

	log.Info("App.Start...")
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Start(ctx)
	}()

	log.Info("Router.Register...")
	if err := router.Register(config.GetAppID().String(), config.GetSvcID().String(), _runtime.transport.Addr(), map[string]string{meta.M3AppVer.String(): config.GetVer()}); err != nil {
		log.Error("Router.Register err %s", err.Error())
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			log.Info("Recv Done...")
			log.Info("Plugin.Stop...")
			plugin.Destroy()
			log.Info("Doned")
		}
	}()
	go signalProc()

	log.Info("Wait...")
	wg.Wait()
	return nil
}

type Runtime struct {
	app         app.App
	servers     map[string]server.Server
	cancel      context.CancelFunc
	transport   transport.Transport
	resourcemgr *resource.ResourceMgr
}

func (r *Runtime) HealthCheck(idstr string) bool {
	if string(config.GetAppID()) != idstr {
		return false
	}
	return r.app.HealthCheck()
}

func (r *Runtime) registerServer(s server.Server) error {
	if _, ok := r.servers[s.Name()]; ok {
		return fmt.Errorf("Register repeated ServerName %s", string(s.Name()))
	}
	r.servers[s.Name()] = s
	return nil
}

func (r *Runtime) setupPluginInterceptor(cfg RuntimeCfg) error {
	client.RegisterClientInterceptor(r.ClientInterceptor)

	log.Info("SetupPluginInterceptor.Shape...")
	if err := shape.Setup(cfg.Options.Shape); err != nil {
		return err
	}
	_runtime.transport.RegisterServerInterceptor(shape.ServerInterceptor())
	client.RegisterClientInterceptor(shape.ClientInterceptor())

	log.Info("SetupPluginInterceptor.Trace...")
	_runtime.transport.RegisterServerInterceptor(trace.ServerInterceptor())
	client.RegisterClientInterceptor(trace.ClientInterceptor())

	return nil
}

func (r *Runtime) ServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	ser := info.Server.(server.Server)
	ctx = server.WithServer(ctx, ser)
	return ser.ServerInterceptor(ctx, req, info, handler)
}

func (r *Runtime) ClientInterceptor(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ser := server.ParseServer(ctx)
	if ser == nil {
		// if not server call, direct call transport
		return _runtime.transport.ClientInterceptor(ctx, method, req, resp, cc, invoker, opts...)
	}
	f := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return _runtime.transport.ClientInterceptor(ctx, method, req, resp, cc, invoker, opts...)
	}
	return ser.ClientInterceptor(ctx, method, req, resp, cc, f, opts...)
}

func (r *Runtime) reload() error {
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
	if err := r.resourcemgr.ReLoad(cfg.Options.Resource); err != nil {
		return err
	}
	log.Info("Transport.Reload...")
	if err := r.transport.Reload(cfg.Transport); err != nil {
		return err
	}
	log.Info("Plugin.Reload...")
	if err := plugin.Reload(v); err != nil {
		return err
	}
	log.Info("Server.Reload...")
	for _, server := range r.servers {
		if err := server.Reload(cfg.Server[string(server.Type())]); err != nil {
			return err
		}
	}
	log.Info("App.Reload...")
	if err := r.app.Reload(cfg.App); err != nil {
		return err
	}
	log.Info("Runtime.Reload Succ...")
	return nil
}

func signalProc() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, util.SysCallSIGUSR1(), util.SysCallSIGUSR2())

	for s := range sigs {
		log.Info("Recv sig %s", s.String())
		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT:
			ShutDown("Recv Sig")
		case util.SysCallSIGUSR1():
			Reload()
		case util.SysCallSIGUSR2():
			PreExit()
		}
	}
}
