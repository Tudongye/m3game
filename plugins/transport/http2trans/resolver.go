package http2trans

import (
	"fmt"
	"m3game/plugins/log"
	"m3game/plugins/router"
	"sync"
	"time"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

func init() {
	resolver.Register(NewBuilder())
}

type Builder struct {
}

func NewBuilder() resolver.Builder {
	return &Builder{}
}

func (cb *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	cr := &Resolver{
		svc:                  target.Authority,
		cc:                   cc,
		disableServiceConfig: opts.DisableServiceConfig,
	}
	cr.wg.Add(1)
	go cr.watcher()
	return cr, nil
}

func (cb *Builder) Scheme() string {
	return "router"
}

type Resolver struct {
	wg                   sync.WaitGroup
	cc                   resolver.ClientConn
	svc                  string
	disableServiceConfig bool
}

func (cr *Resolver) watcher() {

	for {
		<-time.After(time.Second * 1)
		var newAddrs []resolver.Address
		instances, err := router.Instance().GetAllInstances(cr.svc)
		if err != nil {
			log.Error("GetAllInstance %s fail err : %s", cr.svc, err.Error())
			continue
		}
		for _, instance := range instances {
			addr := fmt.Sprintf("%v:%v", instance.GetHost(), instance.GetPort())
			newAddrs = append(newAddrs, resolver.Address{
				BalancerAttributes: attributes.New(BalanceAttrKey,
					&BalanceAttrValue{
						AppID: instance.GetAppID(),
					}),
				Addr: addr,
			})

		}
		cr.cc.NewAddress(newAddrs)
	}
}

func (cr *Resolver) ResolveNow(opt resolver.ResolveNowOptions) {
}

func (cr *Resolver) Close() {
}
