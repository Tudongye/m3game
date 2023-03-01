package client

import (
	"context"
	"m3game/proto/pb"
	"m3game/runtime"
	"m3game/runtime/transport"
	"m3game/server"
	"m3game/util/log"

	"google.golang.org/grpc"
)

func SendInteror(f func(*transport.Sender, func(*transport.Sender) error) error) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		s, err := transport.NewSender(ctx, method, req, resp, cc, invoker, opts)
		if err != nil {
			log.Error("NewSender err %s", err.Error())
			return err
		}
		return f(s, sendInterFunc)
	}
}

func sendInterFunc(sender *transport.Sender) error {
	sctx := server.ParseContext(sender.Ctx())
	if sctx == nil {
		return runtime.SendInterFunc(sender)
	}
	s := sctx.Server()
	return s.SendInterFunc(sender)
}

func NewRouteHeadRandom(srcins *pb.RouteIns, dstsvc *pb.RouteSvc) *pb.RouteHead {
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

func NewRouteHeadP2P(srcins *pb.RouteIns, dstsvc *pb.RouteSvc, dstins *pb.RouteIns) *pb.RouteHead {
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

func NewRouteHeadHash(srcins *pb.RouteIns, dstsvc *pb.RouteSvc, hashkey string) *pb.RouteHead {
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

func NewRouteHeadSingle(srcins *pb.RouteIns, dstsvc *pb.RouteSvc) *pb.RouteHead {
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

func NewRouteHeadMutil(srcins *pb.RouteIns, topic string) *pb.RouteHead {
	var routehead pb.RouteHead
	routehead.SrcIns = srcins
	routehead.RouteType = pb.RouteType_RT_MUTIL
	routehead.RoutePara = &pb.RoutePara{
		RouteMutilHead: []*pb.RouteMutilHead{
			{
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
		RouteBroadHead: []*pb.RouteBroadHead{
			{
				Pass: "",
			},
		},
	}
	return &routehead
}
