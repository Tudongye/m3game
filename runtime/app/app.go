package app

import (
	"fmt"
	"m3game/config"
	"m3game/proto/pb"
	"m3game/runtime/transport"
	"m3game/util"
	"sync"
)

type App interface {
	Init(cfg map[string]interface{}) error                                                                // 初始化
	Start(wg *sync.WaitGroup) error                                                                       // 启动
	Stop() error                                                                                          // 停止
	Reload(map[string]interface{}) error                                                                  // 重载
	RecvInterFunc(*transport.Reciver, func(*transport.Reciver) (interface{}, error)) (interface{}, error) // RPC被调拦截器
	SendInterFunc(*transport.Sender, func(sender *transport.Sender) error) error                          // RPC主调拦截器
	HealthCheck() bool                                                                                    // 健康检查

	IDStr() string              // 实例ID
	RouteIns() *pb.RouteIns     // 实例信息
	RouteSvc() *pb.RouteSvc     // 服务信息
	RouteWorld() *pb.RouteWorld // 服务区服
}

func New(pfuncid string) *appBase {
	idstr := config.GetIDStr()
	if idstr == "" {
		panic("Flag IDStr not find")
	}
	if envid, worldid, funcid, insid, err := util.AppStr2ID(idstr); err != nil {
		panic(fmt.Errorf("IDStr %s not vaild %w", idstr, err))
	} else if funcid != pfuncid {
		panic(fmt.Errorf("IDStr %s not vaild %w", idstr, err))
	} else {
		return &appBase{
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
}

type appBase struct {
	routeins   *pb.RouteIns
	routesvc   *pb.RouteSvc
	routeworld *pb.RouteWorld
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

func (a *appBase) Start(wg *sync.WaitGroup) error {
	return nil
}

func (a *appBase) Stop() error {
	return nil
}

func (a *appBase) Reload(map[string]interface{}) error {
	return nil
}

func (a *appBase) RecvInterFunc(rec *transport.Reciver, f func(*transport.Reciver) (resp interface{}, err error)) (resp interface{}, err error) {
	return f(rec)
}

func (a *appBase) HealthCheck() bool {
	return false
}

func (a *appBase) SendInterFunc(sender *transport.Sender, f func(sender *transport.Sender) error) error {
	return f(sender)
}

func (a *appBase) RouteIns() *pb.RouteIns {
	return a.routeins
}

func (a *appBase) RouteSvc() *pb.RouteSvc {
	return a.routesvc
}
func (a *appBase) RouteWorld() *pb.RouteWorld {
	return a.routeworld
}

func (a *appBase) IDStr() string {
	return a.routeins.IDStr
}
