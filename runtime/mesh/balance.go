package mesh

import (
	"m3game/meta"
	"m3game/plugins/log"

	"github.com/pkg/errors"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/metadata"
)

const (
	BalanceAttrKey = "m3gblattrkey"
)

type BalanceAttrValue struct {
	IDStr string
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

var (
	_err_parsesenderfail   = errors.New("_err_parsesenderfail")
	_err_routerheadinvaild = errors.New("_err_routerheadinvaild")
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
	routeHelper := NewRouteHelper()
	for subConn, conInfo := range info.ReadySCs {
		v := conInfo.Address.BalancerAttributes.Value(BalanceAttrKey)
		if v == nil {
			continue
		}
		idstr := v.(*BalanceAttrValue).IDStr
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
	routeHelper *RouteHelper
}

func (p *M3GPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	var ret balancer.PickResult
	if md, ok := metadata.FromOutgoingContext(info.Ctx); !ok {
		log.Error("Can't get outgoing ctx")
		return ret, balancer.ErrNoSubConnAvailable
	} else {
		vlist := md[string(meta.M3RouteType)]
		if len(vlist) != 1 {
			log.Error("M3RouteType is invalid")
			return ret, balancer.ErrNoSubConnAvailable
		}
		switch vlist[0] {
		case meta.RouteTypeP2P.String():
			return p.pickP2P(md)
		case meta.RouteTypeRandom.String():
			return p.pickRandom(md)
		case meta.RouteTypeHash.String():
			return p.pickHash(md)
		case meta.RouteTypeSingle.String():
			return p.pickSingle(md)
		default:
			log.Error("Unknow RouteType %s", vlist[0])
			return ret, balancer.ErrNoSubConnAvailable
		}
	}
}

func (p *M3GPicker) pickP2P(md metadata.MD) (balancer.PickResult, error) {
	vlist := md[string(meta.M3RouteDstApp)]
	if len(vlist) != 1 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	if dstappid, err := p.routeHelper.RouteP2P(vlist[0]); err != nil {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	} else {
		return balancer.PickResult{SubConn: p.subConns[dstappid]}, nil
	}
}

func (p *M3GPicker) pickRandom(md metadata.MD) (balancer.PickResult, error) {
	if dstappid, err := p.routeHelper.RouteRandom(); err != nil {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	} else {
		return balancer.PickResult{SubConn: p.subConns[dstappid]}, nil
	}
}

func (p *M3GPicker) pickHash(md metadata.MD) (balancer.PickResult, error) {
	vlist := md[string(meta.M3RouteHashKey)]
	if len(vlist) != 1 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	if dstappid, err := p.routeHelper.RouteHash(vlist[0]); err != nil {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	} else {
		return balancer.PickResult{SubConn: p.subConns[dstappid]}, nil
	}
}

func (p *M3GPicker) pickSingle(md metadata.MD) (balancer.PickResult, error) {
	if dstappid, err := p.routeHelper.RouteSingle(); err != nil {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	} else {
		return balancer.PickResult{SubConn: p.subConns[dstappid]}, nil
	}
}
