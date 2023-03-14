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
	_ "m3game/plugins/log/zap"
	_ "m3game/plugins/metric/prometheus"
	_ "m3game/plugins/router/consul"
	_ "m3game/plugins/shape/sentinel"
	mruntime "m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/rpc"
	"m3game/runtime/server"
	"strings"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var (
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
}

func (d *GateApp) Init(c map[string]interface{}) error {

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
	return nil
}

func (d *GateApp) LogicCall(in *metapb.CSMsg) (*metapb.CSMsg, error) {
	if !rpc.IsCSFullMethod(in.Method) {
		return nil, fmt.Errorf("Method %s invaild", in.Method)
	}
	if strings.HasPrefix(in.Method, "/proto.ActorRegSer") {
		return gate.CallGrpcCli(context.Background(), actorregcli.Conn(), in)
	}
	if strings.HasPrefix(in.Method, "/proto.ActorSer") {
		return gate.CallGrpcCli(context.Background(), actorcli.Conn(), in)
	}
	if strings.HasPrefix(in.Method, "/proto.MultiSer") {
		return gate.CallGrpcCli(context.Background(), multicli.Conn(), in)
	}
	return nil, fmt.Errorf("Unknow Method %s", in.Method)
}

func (d *GateApp) AuthCall(req *metapb.AuthReq) (*metapb.AuthRsp, error) {
	rsp := &metapb.AuthRsp{
		PlayerID: fmt.Sprintf("PlayerID-%s", req.Token),
	}
	return rsp, nil
}

func (d *GateApp) Stop() error {
	return nil
}

func (d *GateApp) HealthCheck() bool {
	return true
}

func Run() error {
	mruntime.Run(newApp(), []server.Server{gateser.New()})
	return nil
}

func CustomMatcher(key string) (string, bool) {
	if len(key) > 2 && key[:2] == "M3" {
		return key, true
	}
	return runtime.DefaultHeaderMatcher(key)
}
