package http2trans

import (
	"m3game/meta"
	"m3game/plugins/log"
	"m3game/runtime/mesh"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/metadata"
)

const (
	BalanceAttrKey = "m3gblattrkey"
)

type BalanceAttrValue struct {
	AppID string
}

func init() {
	balancer.Register(newM3GPikerBuilder())
}

/*
	RT_P2P 点对点
	RT_RAND 随机
	RT_HASH 哈希
	RT_BROAD 广播？
	RT_MUTIL 多播？
	RT_SINGLE 单例
*/

const (
	Balance_m3g = "balance_m3g"
)

func newM3GPikerBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Balance_m3g, &M3GPickerBuilder{}, base.Config{HealthCheck: true})
}

type M3GPickerBuilder struct {
}

func (p *M3GPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	subConns := make(map[string]balancer.SubConn)
	routeHelper := mesh.NewRouteHelper()
	for subConn, conInfo := range info.ReadySCs {
		v := conInfo.Address.BalancerAttributes.Value(BalanceAttrKey)
		if v == nil {
			continue
		}
		idstr := v.(*BalanceAttrValue).AppID
		subConns[idstr] = subConn
		routeHelper.Add(idstr)
	}
	routeHelper.Compress()
	return &M3GPicker{
		subConns:    subConns,
		routeHelper: routeHelper,
	}
}

type M3GPicker struct {
	subConns    map[string]balancer.SubConn
	routeHelper *mesh.RouteHelper
}

func (p *M3GPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	var pickret balancer.PickResult
	if metad, ok := metadata.FromOutgoingContext(info.Ctx); !ok {
		log.Error("Can't get outgoing metadata from ctx")
		return pickret, balancer.ErrNoSubConnAvailable
	} else {
		vlist := metad[string(meta.M3RouteType)]
		if len(vlist) != 1 {
			log.Error("M3RouteType is invalid")
			return pickret, balancer.ErrNoSubConnAvailable
		}
		switch vlist[0] {
		case mesh.RouteTypeP2P.String():
			return p.pickP2P(metad)
		case mesh.RouteTypeRandom.String():
			return p.pickRandom(metad)
		case mesh.RouteTypeHash.String():
			return p.pickHash(metad)
		case mesh.RouteTypeSingle.String():
			return p.pickSingle(metad)
		default:
			log.Error("Unknow RouteType %s", vlist[0])
			return pickret, balancer.ErrNoSubConnAvailable
		}
	}
}

func (p *M3GPicker) pickP2P(metad metadata.MD) (balancer.PickResult, error) {
	vlist := metad[string(meta.M3RouteDstApp)]
	if len(vlist) != 1 {
		log.Error("Invaild Para")
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	if dstappid, err := p.routeHelper.RouteP2P(vlist[0]); err != nil {
		log.Error("%s", err.Error())
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	} else {
		return balancer.PickResult{SubConn: p.subConns[dstappid]}, nil
	}
}

func (p *M3GPicker) pickRandom(metad metadata.MD) (balancer.PickResult, error) {
	if dstappid, err := p.routeHelper.RouteRandom(); err != nil {
		log.Error("%s", err.Error())
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	} else {
		return balancer.PickResult{SubConn: p.subConns[dstappid]}, nil
	}
}

func (p *M3GPicker) pickHash(metad metadata.MD) (balancer.PickResult, error) {
	vlist := metad[string(meta.M3RouteHashKey)]
	if len(vlist) != 1 {
		log.Error("Invaild Para")
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	if dstappid, err := p.routeHelper.RouteHash(vlist[0]); err != nil {
		log.Error("%s", err.Error())
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	} else {
		return balancer.PickResult{SubConn: p.subConns[dstappid]}, nil
	}
}

func (p *M3GPicker) pickSingle(metad metadata.MD) (balancer.PickResult, error) {
	if dstappid, err := p.routeHelper.RouteSingle(); err != nil {
		log.Error("%s", err.Error())
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	} else {
		return balancer.PickResult{SubConn: p.subConns[dstappid]}, nil
	}
}
