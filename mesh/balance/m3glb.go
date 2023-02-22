package balance

import (
	"m3game/proto/pb"
	"m3game/runtime/transport"
	"math/rand"

	"github.com/serialx/hashring"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

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
	singleIDstr := ""
	var idstrlist []string
	subConns := make(map[string]balancer.SubConn)
	for subConn, conInfo := range info.ReadySCs {
		v := conInfo.Address.BalancerAttributes.Value(BalanceAttrKey)
		if v == nil {
			continue
		}
		idstr := v.(*BalanceAttrValue).IDStr
		subConns[idstr] = subConn
		idstrlist = append(idstrlist, idstr)
		if idstr < singleIDstr || singleIDstr == "" {
			singleIDstr = idstr
		}
	}
	return &M3GPicker{
		subConns:    subConns,
		hashRing:    hashring.New(idstrlist),
		idstrlist:   idstrlist,
		singleIDstr: singleIDstr,
	}
}

type M3GPicker struct {
	subConns    map[string]balancer.SubConn
	hashRing    *hashring.HashRing
	idstrlist   []string
	singleIDstr string
}

func (p *M3GPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	var ret balancer.PickResult
	sender := transport.ParseSender(info.Ctx)
	if sender == nil {
		return ret, nil
	}
	if sender.RouteHead().RouteType == pb.RouteType_RT_P2P {
		return M3G_P2P_Pick(p, sender.RouteHead().RoutePara.RouteP2PHead)
	} else if sender.RouteHead().RouteType == pb.RouteType_RT_RAND {
		return M3G_RAND_Pick(p, sender.RouteHead().RoutePara.RouteRandHead)
	} else if sender.RouteHead().RouteType == pb.RouteType_RT_HASH {
		return M3G_HASH_Pick(p, sender.RouteHead().RoutePara.RouteHashHead)
	} else if sender.RouteHead().RouteType == pb.RouteType_RT_SINGLE {
		return M3G_SINGLE_Pick(p, sender.RouteHead().RoutePara.RouteSingleHead)
	}
	return ret, nil
}

func M3G_P2P_Pick(p *M3GPicker, rps []*pb.RouteP2PHead) (balancer.PickResult, error) {
	var ret balancer.PickResult
	if len(rps) != 1 {
		return ret, nil
	}
	for IDStr, subConn := range p.subConns {
		if IDStr == rps[0].DstIns.IDStr {
			return balancer.PickResult{SubConn: subConn}, nil
		}
	}
	return ret, nil
}

func M3G_RAND_Pick(p *M3GPicker, rps []*pb.RouteRandHead) (balancer.PickResult, error) {
	return balancer.PickResult{SubConn: p.subConns[p.idstrlist[rand.Int()%len(p.idstrlist)]]}, nil
}

func M3G_HASH_Pick(p *M3GPicker, rps []*pb.RouteHashHead) (balancer.PickResult, error) {
	var ret balancer.PickResult
	if len(rps) != 1 {
		return ret, nil
	}
	if idstr, ok := p.hashRing.GetNode(rps[0].HashKey); ok {
		return balancer.PickResult{SubConn: p.subConns[idstr]}, nil
	}
	return ret, nil
}

func M3G_SINGLE_Pick(p *M3GPicker, rps []*pb.RouteSingleHead) (balancer.PickResult, error) {
	return balancer.PickResult{SubConn: p.subConns[p.singleIDstr]}, nil
}
