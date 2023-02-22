package client

import (
	"m3game/proto/pb"
	"m3game/runtime"
	"m3game/runtime/transport"
	"m3game/server"
)

func SendInterFunc(sender *transport.Sender) error {
	sctx := server.ParseContext(sender.Ctx())
	if sctx == nil {
		return runtime.SendInterFunc(sender)
	}
	s := sctx.Server()
	return s.SendInterFunc(sender)
}

func CreateRouteHead_Random(srcins *pb.RouteIns, dstsvc *pb.RouteSvc) *pb.RouteHead {
	var routehead pb.RouteHead
	routehead.SrcIns = srcins
	routehead.DstSvc = dstsvc
	routehead.RouteType = pb.RouteType_RT_RAND
	routehead.RoutePara = &pb.RoutePara{
		RouteRandHead: []*pb.RouteRandHead{
			{
				Pass: "",
			},
		},
	}
	return &routehead
}

func CreateRouteHead_P2P(srcins *pb.RouteIns, dstsvc *pb.RouteSvc, dstins *pb.RouteIns) *pb.RouteHead {
	var routehead pb.RouteHead
	routehead.SrcIns = srcins
	routehead.DstSvc = dstsvc
	routehead.RouteType = pb.RouteType_RT_P2P
	routehead.RoutePara = &pb.RoutePara{
		RouteP2PHead: []*pb.RouteP2PHead{
			{
				DstIns: dstins,
			},
		},
	}
	return &routehead
}

func CreateRouteHead_Hash(srcins *pb.RouteIns, dstsvc *pb.RouteSvc, hashkey string) *pb.RouteHead {
	var routehead pb.RouteHead
	routehead.SrcIns = srcins
	routehead.DstSvc = dstsvc
	routehead.RouteType = pb.RouteType_RT_HASH
	routehead.RoutePara = &pb.RoutePara{
		RouteHashHead: []*pb.RouteHashHead{
			{
				HashKey: hashkey,
			},
		},
	}
	return &routehead
}

func CreateRouteHead_Single(srcins *pb.RouteIns, dstsvc *pb.RouteSvc) *pb.RouteHead {
	var routehead pb.RouteHead
	routehead.SrcIns = srcins
	routehead.DstSvc = dstsvc
	routehead.RouteType = pb.RouteType_RT_SINGLE
	routehead.RoutePara = &pb.RoutePara{
		RouteSingleHead: []*pb.RouteSingleHead{
			{
				Pass: "",
			},
		},
	}
	return &routehead
}
