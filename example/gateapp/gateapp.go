package gateapp

import (
	"context"
	"fmt"
	"io"
	"m3game/plugins/agent"
	_ "m3game/plugins/agent/httpagent"
	_ "m3game/plugins/broker/nats"
	"m3game/example/actorapp/actorcli"
	"m3game/example/actorapp/actorregcli"
	"m3game/example/asyncapp/asynccli"
	"m3game/example/mutilapp/mutilcli"
	"m3game/example/proto"
	"m3game/example/proto/pb"
	"m3game/plugins/log"
	_ "m3game/plugins/log/zap"
	_ "m3game/plugins/router/consul"
	_ "m3game/plugins/metric/prometheus"
	_ "m3game/plugins/shape/sentinel"
	mruntime "m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/client"
	"m3game/runtime/server"
	_ "m3game/plugins/trace/stdout"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

var (
	_cfg        GateAppCfg
	_listenaddr string
	_onebuf     = make([]byte, 1)
)

func newApp() *GateApp {
	return &GateApp{
		App: app.New(proto.GateAppFuncID),
	}
}

type GateApp struct {
	app.App
	gwServer *http.Server
}

type GateAppCfg struct {
	GatePort int `mapstructure:"GatePort"`
}

func (c *GateAppCfg) CheckVaild() error {
	if c.GatePort == 0 {
		return fmt.Errorf("GatePort cant be 0")
	}
	return nil
}

func (d *GateApp) Init(c map[string]interface{}) error {
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return errors.Wrap(err, "GateApp Decode Cfg")
	}
	if err := _cfg.CheckVaild(); err != nil {
		return err
	}
	return nil
}

func (d *GateApp) Start(wg *sync.WaitGroup) error {
	if err := mutilcli.Init(d.RouteIns(), client.GenMetaOption(proto.RHMeta_Client, "1")); err != nil {
		return err
	}
	if err := asynccli.Init(d.RouteIns(), client.GenMetaOption(proto.RHMeta_Client, "1")); err != nil {
		return err
	}
	if err := actorcli.Init(d.RouteIns(), client.GenMetaOption(proto.RHMeta_Client, "1")); err != nil {
		return err
	}
	if err := actorregcli.Init(d.RouteIns(), client.GenMetaOption(proto.RHMeta_Client, "1")); err != nil {
		return err
	}
	gwmux := runtime.NewServeMux()
	if err := pb.RegisterMutilSerHandler(context.Background(), gwmux, mutilcli.Conn()); err != nil {
		return errors.Wrap(err, "RegisterMutilSerHandler fail")
	}
	if err := pb.RegisterActorSerHandler(context.Background(), gwmux, actorcli.Conn()); err != nil {
		return errors.Wrap(err, "RegisterActorSerHandler fail")
	}
	if err := pb.RegisterActorRegSerHandler(context.Background(), gwmux, actorregcli.Conn()); err != nil {
		return errors.Wrap(err, "RegisterActorRegSerHandler fail")
	}
	_listenaddr = fmt.Sprintf("127.0.0.1:%d", _cfg.GatePort)
	d.gwServer = &http.Server{
		Addr:    _listenaddr,
		Handler: gwmux,
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("GateSer Listen %s", _listenaddr)
		if err := d.gwServer.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Info("GateWay Server Close")
			} else {
				panic(fmt.Sprintf("GateWay Server Err %s", err.Error()))
			}
		}
	}()
	agent.SetAuther(agentAuther)
	agent.SetCaller(agentCaller)
	return nil
}

func agentAuther(para agent.AuthPara) (string, error) {
	return para.Token, nil
}

func agentCaller(method string, uid string, reqbuf io.ReadCloser) (io.ReadCloser, error) {
	client := &http.Client{}
	if req, err := http.NewRequest("POST", fmt.Sprintf("http://%s%s", _listenaddr, method), io.MultiReader(reqbuf)); err != nil {
		return nil, err
	} else if rsp, err := client.Do(req); err != nil {
		return nil, err
	} else {
		return rsp.Body, nil
	}
}

func (d *GateApp) Stop() error {
	return d.gwServer.Shutdown(context.Background())
}

func (d *GateApp) HealthCheck() bool {
	return true
}

func Run() error {
	mruntime.Run(newApp(), []server.Server{})
	return nil
}
