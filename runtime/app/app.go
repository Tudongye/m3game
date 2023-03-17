package app

import (
	"context"
)

type App interface {
	Init(cfg map[string]interface{}) error // 初始化
	Prepare(ctx context.Context) error     // 启动
	Start(ctx context.Context)             // 启动
	Stop() error                           // 停止
	Reload(map[string]interface{}) error   // 重载                                                   // RPC主调拦截器
	HealthCheck() bool                     // 健康检查
}

func New(pfuncid string) *appBase {
	return &appBase{}
}

type appBase struct {
}

var (
	_ App = (*appBase)(nil)
)

func (a *appBase) Init(cfg map[string]interface{}) error {
	return nil
}

func (a *appBase) Name() string {
	return ""
}

func (a *appBase) Prepare(context.Context) error {
	return nil
}

func (a *appBase) Start(context.Context) {
}

func (a *appBase) Stop() error {
	return nil
}

func (a *appBase) Reload(map[string]interface{}) error {
	return nil
}

func (a *appBase) HealthCheck() bool {
	return false
}
