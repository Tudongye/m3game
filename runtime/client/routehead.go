package client

import (
	"m3game/proto/pb"
)

func NewRouteHeadRandom(srcins *pb.RouteIns, dstsvc *pb.RouteSvc) *pb.RouteHead {
	var routehead pb.RouteHead
	routehead.SrcIns = srcins
	routehead.DstSvc = dstsvc
	routehead.RouteType = pb.RouteType_RT_RAND
	routehead.RoutePara = &pb.RoutePara{
		Para: &pb.RoutePara_RouteRandHead{
			RouteRandHead: &pb.RouteRandHead{
				Pass: "",
			},
		},
	}
	return &routehead
}

func NewRouteHeadP2P(srcins *pb.RouteIns, dstsvc *pb.RouteSvc, dstins *pb.RouteIns) *pb.RouteHead {
	var routehead pb.RouteHead
	routehead.SrcIns = srcins
	routehead.DstSvc = dstsvc
	routehead.RouteType = pb.RouteType_RT_P2P
	routehead.RoutePara = &pb.RoutePara{
		Para: &pb.RoutePara_RouteP2PHead{
			RouteP2PHead: &pb.RouteP2PHead{
				DstIns: dstins,
			},
		},
	}
	return &routehead
}

func NewRouteHeadHash(srcins *pb.RouteIns, dstsvc *pb.RouteSvc, hashkey string) *pb.RouteHead {
	var routehead pb.RouteHead
	routehead.SrcIns = srcins
	routehead.DstSvc = dstsvc
	routehead.RouteType = pb.RouteType_RT_HASH
	routehead.RoutePara = &pb.RoutePara{
		Para: &pb.RoutePara_RouteHashHead{
			RouteHashHead: &pb.RouteHashHead{
				HashKey: hashkey,
			},
		},
	}
	return &routehead
}

func NewRouteHeadSingle(srcins *pb.RouteIns, dstsvc *pb.RouteSvc) *pb.RouteHead {
	var routehead pb.RouteHead
	routehead.SrcIns = srcins
	routehead.DstSvc = dstsvc
	routehead.RouteType = pb.RouteType_RT_SINGLE
	routehead.RoutePara = &pb.RoutePara{
		Para: &pb.RoutePara_RouteSingleHead{
			RouteSingleHead: &pb.RouteSingleHead{
				Pass: "",
			},
		},
	}
	return &routehead
}

func NewRouteHeadMutil(srcins *pb.RouteIns, topic string) *pb.RouteHead {
	var routehead pb.RouteHead
	routehead.SrcIns = srcins
	routehead.RouteType = pb.RouteType_RT_MUTIL
	routehead.RoutePara = &pb.RoutePara{
		Para: &pb.RoutePara_RouteMutilHead{
			RouteMutilHead: &pb.RouteMutilHead{
				Topic: topic,
			},
		},
	}
	return &routehead
}

func NewRouteHeadBroad(srcins *pb.RouteIns, dstsvc *pb.RouteSvc) *pb.RouteHead {
	var routehead pb.RouteHead
	routehead.SrcIns = srcins
	routehead.DstSvc = dstsvc
	routehead.RouteType = pb.RouteType_RT_BROAD
	routehead.RoutePara = &pb.RoutePara{
		Para: &pb.RoutePara_RouteBroadHead{
			RouteBroadHead: &pb.RouteBroadHead{
				Pass: "",
			},
		},
	}
	return &routehead
}
