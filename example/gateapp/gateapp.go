package gateapp

import (
	"context"
	"fmt"
	"m3game/config"
	"m3game/example/actorapp/actorcli"
	"m3game/example/actorapp/actorregcli"
	"m3game/example/gateapp/gateser"
	"m3game/example/multiapp/multicli"
	"m3game/example/proto"
	"m3game/meta/metapb"
	_ "m3game/plugins/broker/nats"
	"m3game/plugins/gate"
	_ "m3game/plugins/gate/grpcgate"
	"m3game/plugins/log"
	_ "m3game/plugins/log/zap"
	_ "m3game/plugins/metric/prometheus"
	"m3game/plugins/router"
	_ "m3game/plugins/router/consul"
	_ "m3game/plugins/shape/sentinel"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/rpc"
	"m3game/runtime/server"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	_cfg AppCfg
)

func newApp() *GateApp {
	return &GateApp{
		App:  app.New(proto.GateAppFuncID),
		exit: make(chan struct{}, 1),
	}
}

type GateApp struct {
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

func (a *GateApp) Init(cfg map[string]interface{}) error {
	if err := mapstructure.Decode(cfg, &_cfg); err != nil {
		return errors.Wrap(err, "App Decode Cfg")
	}
	if err := _cfg.CheckVaild(); err != nil {
		return err
	}
	return nil
}

func (d *GateApp) Start(wg *sync.WaitGroup) error {
	if err := multicli.Init(config.GetAppID(), grpc.WithCodec(&gate.GateCodec{})); err != nil {
		return err
	}
	if err := actorcli.Init(config.GetAppID(), grpc.WithCodec(&gate.GateCodec{})); err != nil {
		return err
	}
	if err := actorregcli.Init(config.GetAppID(), grpc.WithCodec(&gate.GateCodec{})); err != nil {
		return err
	}
	gate.SetReciver(d)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("GateApp PrepareTime %d", _cfg.PrePareTime)
		time.Sleep(time.Duration(_cfg.PrePareTime) * time.Second)
		log.Info("GateApp Ready")
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
				if gate.Get().Factory().CanDelete(gate.Get()) {
					t.Stop()
					runtime.ShutDown()
				}
				continue
			}
		}
	}()
	return nil
}

func (d *GateApp) Stop() error {
	select {
	case d.exit <- struct{}{}:
		return nil
	default:
		return nil
	}
}
func (d *GateApp) LogicCall(in *metapb.CSMsg) (*metapb.CSMsg, error) {
	if !rpc.IsCSFullMethod(in.Method) {
		return nil, fmt.Errorf("Method %s invaild", in.Method)
	}
	var out *metapb.CSMsg
	var err error
	if strings.HasPrefix(in.Method, "/proto.ActorRegSer") {
		out, err = gate.CallGrpcCli(context.Background(), actorregcli.Conn(), in)
	}
	if strings.HasPrefix(in.Method, "/proto.ActorSer") {
		out, err = gate.CallGrpcCli(context.Background(), actorcli.Conn(), in)
	}
	if strings.HasPrefix(in.Method, "/proto.MultiSer") {
		out, err = gate.CallGrpcCli(context.Background(), multicli.Conn(), in)
	}
	if err != nil {
		return out, err
	} else if out != nil {
		out.Metas = in.Metas
		return out, nil
	}
	return nil, fmt.Errorf("Unknow Method %s", in.Method)
}

func (d *GateApp) AuthCall(req *metapb.AuthReq) (*metapb.AuthRsp, error) {
	rsp := &metapb.AuthRsp{
		PlayerID: fmt.Sprintf("PlayerID-%s", req.Token),
	}
	return rsp, nil
}

func (d *GateApp) HealthCheck() bool {
	return true
}

func Run() error {
	runtime.Run(newApp(), []server.Server{gateser.New()})
	return nil
}
