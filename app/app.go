package app

import (
	"fmt"
	"log"
	"m3game/config"
	"m3game/proto/pb"
	"m3game/runtime/transport"
	"m3game/util"
)

type App interface {
	Init(cfg map[string]interface{}) error                                                                // 初始化
	Start() error                                                                                         // 启动
	Stop() error                                                                                          // 停止
	Reload() error                                                                                        // 重载
	RecvInterFunc(*transport.Reciver, func(*transport.Reciver) (interface{}, error)) (interface{}, error) // RPC被调拦截器
	SendInterFunc(*transport.Sender, func(sender *transport.Sender) error) error                          // RPC主调拦截器
	HealthCheck() bool                                                                                    // 健康检查

	IDStr() string              // 实例ID
	RouteIns() *pb.RouteIns     // 实例信息
	RouteSvc() *pb.RouteSvc     // 服务信息
	RouteWorld() *pb.RouteWorld // 服务区服
}

func CreateDefaultApp(pfuncid string) *DefaultApp {
	idstr := config.GetIDStr()
	if idstr == "" {
		log.Panic("Flag IDStr not find")
	}
	if envid, worldid, funcid, insid, err := util.AppStr2ID(idstr); err != nil {
		log.Panic(fmt.Errorf("IDStr %s not vaild %w", idstr, err))
	} else if funcid != pfuncid {
		log.Panic(fmt.Errorf("IDStr %s not vaild %w", idstr, err))
	} else {
		return &DefaultApp{
			routeins: &pb.RouteIns{
				EnvID:   envid,
				WorldID: worldid,
				FuncID:  funcid,
				InsID:   insid,
				IDStr:   idstr,
			},
			routesvc: &pb.RouteSvc{
				EnvID:   envid,
				WorldID: worldid,
				FuncID:  funcid,
			},
			routeworld: &pb.RouteWorld{
				EnvID:   envid,
				WorldID: worldid,
			},
		}
	}
	return nil
}

type DefaultApp struct {
	routeins   *pb.RouteIns
	routesvc   *pb.RouteSvc
	routeworld *pb.RouteWorld
}

var (
	_ App = (*DefaultApp)(nil)
)

func (a *DefaultApp) Init(cfg map[string]interface{}) error {
	return nil
}

func (a *DefaultApp) Name() string {
	return ""
}

func (a *DefaultApp) Start() error {
	return nil
}

func (a *DefaultApp) Stop() error {
	return nil
}

func (a *DefaultApp) Reload() error {
	return nil
}

func (a *DefaultApp) RecvInterFunc(rec *transport.Reciver, f func(*transport.Reciver) (resp interface{}, err error)) (resp interface{}, err error) {
	return f(rec)
}

func (a *DefaultApp) HealthCheck() bool {
	return false
}

func (a *DefaultApp) SendInterFunc(sender *transport.Sender, f func(sender *transport.Sender) error) error {
	return f(sender)
}

func (a *DefaultApp) RouteIns() *pb.RouteIns {
	return a.routeins
}

func (a *DefaultApp) RouteSvc() *pb.RouteSvc {
	return a.routesvc
}
func (a *DefaultApp) RouteWorld() *pb.RouteWorld {
	return a.routeworld
}

func (a *DefaultApp) IDStr() string {
	return a.routeins.IDStr
}
