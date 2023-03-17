package meta

const (
	META_CLIENT       = "_client"      // 是否为客户端请求
	META_ACTORID      = "_actorid"     // Actor实例ID
	META_CREATE_ACTOR = "_createactor" // 创建Actor
)

const (
	META_FLAG_TRUE  = "true"
	META_FLAG_FALSE = "false"
)

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
	M3RouteTopic   M3Meta = "m3routetopic"   // 目标主题，Multi
	M3RouteIsNty   M3Meta = "m3routeisnty"

	M3ActorActorID M3Meta = "m3actoractorid"
	M3PlayerID     M3Meta = "m3playerid"
	M3ClientSerial M3Meta = "m3clientserial" // 客户端序列号

	// Mesh相关
	M3AppVer M3Meta = "m3appver"
)
