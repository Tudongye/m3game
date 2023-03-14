package client

import (
	"context"
	"m3game/meta"
	"m3game/plugins/broker"
	"m3game/runtime/transport"

	"google.golang.org/grpc/metadata"
)

func FillRouteHeadRandom(ctx context.Context, srcapp meta.RouteApp, dstsvc meta.RouteSvc, isnty meta.IsNty) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		meta.M3RouteType.String(), meta.RouteTypeRandom.String(),
		meta.M3RouteSrcApp.String(), srcapp.String(),
		meta.M3RouteDstSvc.String(), dstsvc.String(),
		meta.M3RouteIsNty.String(), isnty.String(),
	)
}

func FillRouteHeadP2P(ctx context.Context, srcapp meta.RouteApp, dstsvc meta.RouteSvc, dstapp meta.RouteApp, isnty meta.IsNty) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		meta.M3RouteType.String(), meta.RouteTypeP2P.String(),
		meta.M3RouteSrcApp.String(), srcapp.String(),
		meta.M3RouteDstSvc.String(), dstsvc.String(),
		meta.M3RouteDstApp.String(), dstapp.String(),
		meta.M3RouteIsNty.String(), isnty.String(),
	)
}

func FillRouteHeadHash(ctx context.Context, srcapp meta.RouteApp, dstsvc meta.RouteSvc, hashkey string, isnty meta.IsNty) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		meta.M3RouteType.String(), meta.RouteTypeHash.String(),
		meta.M3RouteSrcApp.String(), srcapp.String(),
		meta.M3RouteDstSvc.String(), dstsvc.String(),
		meta.M3RouteHashKey.String(), hashkey,
		meta.M3RouteIsNty.String(), isnty.String(),
	)
}

func FillRouteHeadSingle(ctx context.Context, srcapp meta.RouteApp, dstsvc meta.RouteSvc, isnty meta.IsNty) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		meta.M3RouteType.String(), meta.RouteTypeHash.String(),
		meta.M3RouteSrcApp.String(), srcapp.String(),
		meta.M3RouteDstSvc.String(), dstsvc.String(),
		meta.M3RouteIsNty.String(), isnty.String(),
	)
}

func FillRouteHeadMulti(ctx context.Context, srcapp meta.RouteApp, topic string) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		meta.M3RouteType.String(), meta.RouteTypeMulti.String(),
		meta.M3RouteSrcApp.String(), srcapp.String(),
		meta.M3RouteTopic.String(), topic,
		meta.M3RouteIsNty.String(), meta.IsNtyTrue.String(),
	)
}

func FillRouteHeadBroad(ctx context.Context, srcapp meta.RouteApp, dstsvc meta.RouteSvc) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		meta.M3RouteType.String(), meta.RouteTypeBroad.String(),
		meta.M3RouteSrcApp.String(), srcapp.String(),
		meta.M3RouteDstSvc.String(), dstsvc.String(),
		meta.M3RouteTopic.String(), broker.GenTopic(transport.BrokerSerTopic(dstsvc.String())),
		meta.M3RouteIsNty.String(), meta.IsNtyTrue.String(),
	)
}
