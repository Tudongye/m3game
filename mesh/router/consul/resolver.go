package consul

import (
	"fmt"
	"log"
	"m3game/mesh/balance"
	"sync"

	"github.com/hashicorp/consul/api"
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
		lastIndex:            0,
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
	lastIndex            uint64
}

func (cr *Resolver) watcher() {
	for {
		services, metainfo, err := _instance.client.Health().Service(cr.svc, "", true, &api.QueryOptions{WaitIndex: cr.lastIndex})
		if err != nil {
			log.Println(err)
			return
		}
		cr.lastIndex = metainfo.LastIndex
		var newAddrs []resolver.Address
		for _, service := range services {
			addr := fmt.Sprintf("%v:%v", service.Service.Address, service.Service.Port)
			newAddrs = append(newAddrs, resolver.Address{
				BalancerAttributes: attributes.New(balance.BalanceAttrKey,
					balance.BalanceAttrValue{
						IDStr: service.Service.ID,
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
