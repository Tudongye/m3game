package simpleapp

import (
	"m3game/example/proto"
	"m3game/example/simpleapp/simpleser"
	_ "m3game/plugins/broker/nats"
	_ "m3game/plugins/router/consul"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/server"
)

// 创建App实体
func newApp() *SimpleApp {
	return &SimpleApp{
		App: app.New(proto.SimpleAppFuncID), // 指定App的FuncID
	}
}

type SimpleApp struct {
	app.App
}

// 健康检测
func (d *SimpleApp) HealthCheck() bool {
	return true
}

func Run() error {
	// 启动一个 包含了simpleser的SimpleApp
	runtime.Run(newApp(), []server.Server{simpleser.New()})
	return nil
}
