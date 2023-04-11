package meta

type M3Meta string

func (r M3Meta) String() string {
	return string(r)
}

const (
	// 路由相关
	M3RouteType    M3Meta = "m3routetype"    // 路由方式
	M3RouteSrcApp  M3Meta = "m3routesrcapp"  // 源App
	M3RouteDstSvc  M3Meta = "m3routedstsvc"  // 目标Svc
	M3RouteDstApp  M3Meta = "m3routedstapp"  // 目标App，P2P
	M3RouteHashKey M3Meta = "m3routehashkey" // 哈希Key, Hash
	M3RouteIsNty   M3Meta = "m3routeisnty"   // 是Nty, 未实装 TODO
	M3GateMsg      M3Meta = "m3gatemsg"      // 是Gate消息
	// Mesh相关
	M3AppVer M3Meta = "m3appver"
)
