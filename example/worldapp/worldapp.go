package worldapp

import (
	"context"
	"m3game/example/proto"
	"m3game/example/worldapp/worldser"
	_ "m3game/plugins/broker/nats"
	_ "m3game/plugins/gate/grpcgate"
	_ "m3game/plugins/log/zap"
	_ "m3game/plugins/metric/prometheus"
	_ "m3game/plugins/router/consul"
	_ "m3game/plugins/shape/sentinel"
	_ "m3game/plugins/transport/http2trans"
	_ "m3game/plugins/transport/natstrans"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/server"
)

// 创建App实体
func newApp() *WorldApp {
	return &WorldApp{
		App: app.New(proto.WorldAppFuncID), // 指定App的FuncID
	}
}

type WorldApp struct {
	app.App
}

// 健康检测
func (d *WorldApp) Alive(app string, svc string) bool {
	return true
}

func Run(ctx context.Context) error {
	wser := worldser.New()
	wser.RunWorld(ctx)
	runtime.New().Run(ctx, newApp(), []server.Server{wser})
	return nil
}
