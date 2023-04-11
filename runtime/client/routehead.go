package client

import (
	"context"
	"m3game/meta"
	"m3game/runtime/mesh"

	"google.golang.org/grpc/metadata"
)

func FillRouteHeadRandom(ctx context.Context, srcapp mesh.RouteApp, dstsvc mesh.RouteSvc, isnty mesh.IsNty) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		meta.M3RouteType.String(), mesh.RouteTypeRandom.String(),
		meta.M3RouteSrcApp.String(), srcapp.String(),
		meta.M3RouteDstSvc.String(), dstsvc.String(),
		meta.M3RouteIsNty.String(), isnty.String(),
	)
}

func FillRouteHeadP2P(ctx context.Context, srcapp mesh.RouteApp, dstsvc mesh.RouteSvc, dstapp mesh.RouteApp, isnty mesh.IsNty) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		meta.M3RouteType.String(), mesh.RouteTypeP2P.String(),
		meta.M3RouteSrcApp.String(), srcapp.String(),
		meta.M3RouteDstSvc.String(), dstsvc.String(),
		meta.M3RouteDstApp.String(), dstapp.String(),
		meta.M3RouteIsNty.String(), isnty.String(),
	)
}

func FillRouteHeadHash(ctx context.Context, srcapp mesh.RouteApp, dstsvc mesh.RouteSvc, hashkey string, isnty mesh.IsNty) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		meta.M3RouteType.String(), mesh.RouteTypeHash.String(),
		meta.M3RouteSrcApp.String(), srcapp.String(),
		meta.M3RouteDstSvc.String(), dstsvc.String(),
		meta.M3RouteHashKey.String(), hashkey,
		meta.M3RouteIsNty.String(), isnty.String(),
	)
}

func FillRouteHeadSingle(ctx context.Context, srcapp mesh.RouteApp, dstsvc mesh.RouteSvc, isnty mesh.IsNty) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		meta.M3RouteType.String(), mesh.RouteTypeSingle.String(),
		meta.M3RouteSrcApp.String(), srcapp.String(),
		meta.M3RouteDstSvc.String(), dstsvc.String(),
		meta.M3RouteIsNty.String(), isnty.String(),
	)
}

func FillRouteHeadBroad(ctx context.Context, srcapp mesh.RouteApp, dstsvc mesh.RouteSvc) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		meta.M3RouteType.String(), mesh.RouteTypeBroad.String(),
		meta.M3RouteSrcApp.String(), srcapp.String(),
		meta.M3RouteDstSvc.String(), dstsvc.String(),
		meta.M3RouteIsNty.String(), "1",
	)
}
