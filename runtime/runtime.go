package runtime

import (
	"context"
	"fmt"
	"log"
	"m3game/app"
	"m3game/config"
	_ "m3game/mesh/balance"
	"m3game/resource"
	"m3game/runtime/plugin"
	"m3game/runtime/transport"
	"m3game/server"
)

func init() {
	_instance = &Runtime{
		servers: make(map[string]server.Server),
	}
}

var (
	_instance *Runtime
)

type Runtime struct {
	app     app.App
	servers map[string]server.Server
	cancel  context.CancelFunc
}

type RuntimeCfg struct {
	Resource  map[string]interface{}            `toml:"Resource"`
	App       map[string]interface{}            `toml:"App"`
	Server    map[string]map[string]interface{} `toml:"Server"`
	Transport map[string]interface{}            `toml:"Transport"`
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

func (r *Runtime) RegisterServer(s server.Server) error {
	if _, ok := r.servers[s.Name()]; ok {
		return fmt.Errorf("Register repeated ServerName %s", string(s.Name()))
	}
	r.servers[s.Name()] = s
	return nil
}

func RecvInterFunc(trecv *transport.Reciver) (resp interface{}, err error) {
	return _instance.RecvInterFunc(trecv)
}

func SendInterFunc(sctx *transport.Sender) error {
	return transport.SendInterFunc(sctx)
}

func ShutDown() error {
	_instance.cancel()
	return nil
}

func Stop() error {
	_instance.app.Stop()
	for _, server := range _instance.servers {
		server.Stop()
	}
	return nil
}

func Run(app app.App, servers []server.Server) error {
	ctx, cancel := context.WithCancel(context.Background())
	_instance.cancel = cancel
	_instance.app = app
	for _, server := range servers {
		if err := _instance.RegisterServer(server); err != nil {
			log.Println(err.Error())
			return err
		}
	}

	v := *config.GetConfig()
	var cfg RuntimeCfg
	if err := v.Unmarshal(&cfg); err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println("Resource.Load...")
	if err := resource.Init(cfg.Resource); err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println("Transport.Init...")
	if err := transport.Init(cfg.Transport, _instance); err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println("Plugin.Init...")
	if err := plugin.InitPlugins(v); err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println("App.Init...")
	if err := app.Init(cfg.App); err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println("Server.Init...")
	for _, server := range servers {
		log.Printf("Server.Init.%s...\n", server.Name())
		if err := server.Init(cfg.Server[string(server.Type())], app); err != nil {
			log.Println(err.Error())
			return err
		}
		if err := transport.RegistServer(server.TransportRegister()); err != nil {
			log.Println(err.Error())
			return err
		}
	}

	log.Println("Transport.Start...")
	if err := transport.Start(); err != nil {
		log.Println(err.Error())
		return err
	}

	log.Printf("App.Start.%s...\n", app.IDStr())
	if err := app.Start(); err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println("Server.Start...")
	for _, server := range servers {
		log.Printf("Server.Start.%s...\n", server.Name())
		if err := server.Start(); err != nil {
			log.Println(err.Error())
			return err
		}
	}

	log.Println("Wait...")
	select {
	case <-ctx.Done():
		Stop()
	}
	return nil
}
