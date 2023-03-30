package runtime

import (
	"context"
	"m3game/config"
	"m3game/meta"
	"m3game/meta/errs"
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
		return err
	}
	return nil
}

func PreExit() error {
	return nil
}

func Addr() string {
	return _runtime.transport.Addr()
}

func Host() string {
	return _runtime.transport.Host()
}

func Port() int {
	return _runtime.transport.Port()
}

func New() *Runtime {
	if _runtime != nil {
		return _runtime
	}
	_runtime = &Runtime{
		servers: make(map[string]server.Server),
	}
	return _runtime
}

type Runtime struct {
	app         app.App
	servers     map[string]server.Server
	cancel      context.CancelFunc
	transport   *transport.Transport
	resourcemgr *resource.ResourceMgr
}

func (r *Runtime) Run(c context.Context, app app.App, servers []server.Server) error {
	ctx, cancel := context.WithCancel(c)

	// 加载命令行配置
	log.Info("config.Init...")
	config.Init()

	// 初始化Runtime
	log.Info("Runtime.Init...")
	r.cancel = cancel
	r.app = app
	for _, server := range servers {
		if err := r.registerServer(server); err != nil {
			log.Error("RegisterServer %s err %s", server.Name(), err.Error())
			return err
		}
	}

	// 加载Runtime配置
	log.Info("RuntimeCfg.Load...")
	v := *config.GetConfig()
	var cfg RuntimeCfg
	if err := v.Unmarshal(&cfg); err != nil {
		log.Error("UnMarshal RuntimeCfg err %s", err.Error())
		return err
	}

	// 初始化服务网格
	log.Info("Mesh.Init...")
	if err := mesh.Init(cfg.Options.Mesh); err != nil {
		log.Error("Mesh.Init err %s", err.Error())
		return err
	}

	// 初始化资源加载器
	log.Info("Resource.Init...")
	if resourcemgr, err := resource.New(cfg.Options.Resource); err != nil {
		log.Error("Runtime.Resource.Init %s err %s", cfg.Options.Resource, err.Error())
		return err
	} else {
		r.resourcemgr = resourcemgr
		if err := r.resourcemgr.ReLoad(cfg.Options.Resource); err != nil {
			log.Error("Runtime.Resource.Reload %s err %s", cfg.Options.Resource, err.Error())
			return err
		}
	}

	// 创建传输层对象
	log.Info("Transport.New...")
	if trans, err := transport.New(cfg.Transport, r); err != nil {
		log.Error("Transport.New err %s", err.Error())
		return err
	} else {
		r.transport = trans
	}

	// 初始化插件
	log.Info("Plugins.Init...")
	if err := plugin.InitPlugins(ctx, v); err != nil {
		log.Error("InitPlugins err %s", err.Error())
		return err
	}

	// 挂载Grpc拦截器
	log.Info("SetupPluginInterceptor...")
	if err := r.setupPluginInterceptor(cfg); err != nil {
		log.Error("setupPluginInterceptor err %s", err.Error())
		return err
	}

	// 初始化业务Server
	log.Info("Server.Init...")
	for _, server := range servers {
		log.Info("Server.Init.%s...", server.Name())
		if err := server.Init(cfg.Server[string(server.Type())], app); err != nil {
			log.Error("Server.Init %s err %s", server.Name(), err.Error())
			return err
		}
	}

	// 初始化业务App
	log.Info("App.Init...")
	if err := app.Init(cfg.App); err != nil {
		log.Error("App.Init err %s", err.Error())
		return err
	}

	// 传输层预启动
	log.Info("Transport.Prepare...")
	if err := r.transport.Prepare(ctx); err != nil {
		log.Error("Transport.Prepare err %s", err.Error())
		return err
	}

	// 将业务Server注册到传输层
	log.Info("Server.Register...")
	for _, server := range servers {
		if err := r.transport.RegisterServer(server.TransportRegister()); err != nil {
			log.Error("Transport.RegisterServer %s err %s", server.Name(), err.Error())
			return err
		}
	}

	var wg sync.WaitGroup
	// 启动传输层
	log.Info("Transport.Start...")
	wg.Add(1)
	go func() {
		defer wg.Done()
		r.transport.Start(ctx)
	}()

	// 业务Server预启动
	log.Info("Server.Prepare...")
	for _, ser := range servers {
		log.Info("Server.Prepare.%s...", ser.Name())
		if err := ser.Prepare(ctx); err != nil {
			log.Error("Server.Prepare %s err %s", err.Error())
			return err
		}
	}

	// 业务Server启动
	log.Info("Server.Start...")
	for _, ser := range servers {
		log.Info("Server.Start.%s...", ser.Name())
		wg.Add(1)
		go func(s server.Server) {
			defer wg.Done()
			s.Start(ctx)
		}(ser)
	}

	// 业务App预启动
	log.Info("App.Prepare...")
	if err := app.Prepare(ctx); err != nil {
		log.Error("App.Prepare err %s", err.Error())
		return err
	}

	// 业务App启动
	log.Info("App.Start...")
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Start(ctx)
	}()

	// 向服务网格注册本地实例
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

func (r *Runtime) HealthCheck(idstr string) bool {
	if string(config.GetAppID()) != idstr {
		return false
	}
	return r.app.HealthCheck()
}

func (r *Runtime) registerServer(s server.Server) error {
	if _, ok := r.servers[s.Name()]; ok {
		return errs.RuntimeRegisterRepatedServer.New("Register repeated ServerName %s", string(s.Name()))
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
	r.transport.RegisterServerInterceptor(shape.ServerInterceptor())
	client.RegisterClientInterceptor(shape.ClientInterceptor())

	log.Info("SetupPluginInterceptor.Trace...")
	r.transport.RegisterServerInterceptor(trace.ServerInterceptor())
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
		return r.transport.ClientInterceptor(ctx, method, req, resp, cc, invoker, opts...)
	}
	f := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return r.transport.ClientInterceptor(ctx, method, req, resp, cc, invoker, opts...)
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
