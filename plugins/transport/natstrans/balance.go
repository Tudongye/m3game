package natstrans

import (
	"context"
	"m3game/meta"
	"m3game/meta/errs"
	"m3game/plugins/log"
	"m3game/plugins/router"
	"m3game/runtime/mesh"
	"sync"
	"time"

	"google.golang.org/grpc/metadata"
)

func NewBalancer(target string) *NatsBalancer {
	b := &NatsBalancer{
		target: target,
	}
	if err := b.Watch(); err != nil {
		log.Fatal(err.Error())
	}
	return b
}

type NatsBalancer struct {
	target      string
	mu          sync.RWMutex
	routeHelper *mesh.RouteHelper
}

// 路由
func (b *NatsBalancer) Watch() error {
	go func() {
		for {
			<-time.After(time.Second * 1)
			instances, err := router.Instance().GetAllInstances(b.target)
			if err != nil {
				log.Error("GetAllInstance %s fail err : %s", b.target, err.Error())
				continue
			}
			rh := mesh.NewRouteHelper()
			for _, instance := range instances {
				rh.Add(instance.GetAppID())
			}
			rh.Compress()
			b.mu.Lock()
			b.routeHelper = rh
			b.mu.Unlock()
		}
	}()
	return nil
}

// 选路
func (b *NatsBalancer) Pick(ctx context.Context) (string, bool, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if metad, ok := metadata.FromOutgoingContext(ctx); !ok {
		log.Error("Can't get outgoing metadata from ctx")
		return "", false, errs.NatsTransportBalanceErr.New("")
	} else {
		vlist := metad[string(meta.M3RouteType)]
		if len(vlist) != 1 {
			log.Error("M3RouteType is invalid")
			return "", false, errs.NatsTransportBalanceErr.New("")
		}
		switch vlist[0] {
		case mesh.RouteTypeP2P.String():
			return b.pickP2P(metad)
		case mesh.RouteTypeRandom.String():
			return b.pickRandom(metad)
		case mesh.RouteTypeHash.String():
			return b.pickHash(metad)
		case mesh.RouteTypeSingle.String():
			return b.pickSingle(metad)
		case mesh.RouteTypeBroad.String():
			return b.pickBroad(metad)
		default:
			log.Error("Unknow RouteType %s", vlist[0])
			return "", false, errs.NatsTransportBalanceErr.New("")
		}
	}
}

func (b *NatsBalancer) pickP2P(metad metadata.MD) (string, bool, error) {
	vlist := metad[string(meta.M3RouteDstApp)]
	if len(vlist) != 1 {
		log.Error("Invaild Para")
		return "", false, errs.NatsTransportBalanceNoAva.New("")
	}
	if dstappid, err := b.routeHelper.RouteP2P(vlist[0]); err != nil {
		log.Error("%s", err.Error())
		return "", false, errs.NatsTransportBalanceNoAva.New("")
	} else {
		return dstappid, false, nil
	}
}

func (b *NatsBalancer) pickRandom(metad metadata.MD) (string, bool, error) {
	if dstappid, err := b.routeHelper.RouteRandom(); err != nil {
		log.Error("%s", err.Error())
		return "", false, errs.NatsTransportBalanceNoAva.New("")
	} else {
		return dstappid, false, nil
	}
}

func (b *NatsBalancer) pickHash(metad metadata.MD) (string, bool, error) {
	vlist := metad[string(meta.M3RouteHashKey)]
	if len(vlist) != 1 {
		log.Error("Invaild Para")
		return "", false, errs.NatsTransportBalanceNoAva.New("")
	}
	if dstappid, err := b.routeHelper.RouteHash(vlist[0]); err != nil {
		log.Error("%s", err.Error())
		return "", false, errs.NatsTransportBalanceNoAva.New("")
	} else {
		return dstappid, false, nil
	}
}

func (b *NatsBalancer) pickSingle(metad metadata.MD) (string, bool, error) {
	if dstappid, err := b.routeHelper.RouteSingle(); err != nil {
		log.Error("%s", err.Error())
		return "", false, errs.NatsTransportBalanceNoAva.New("")
	} else {
		return dstappid, false, nil
	}
}

func (b *NatsBalancer) pickBroad(metad metadata.MD) (string, bool, error) {
	return "", true, nil
}
