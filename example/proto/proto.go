package proto

const (
	MutilAppFuncID  = "mutil"  // 并发
	AsyncAppFuncID  = "async"  // 单线程异步
	ActorAppFuncID  = "actor"  // Actor模型
	GateAppFuncID   = "gate"   // 服务网关
	SimpleAppFuncID = "simple" //
)

const (
	RHMeta_Client  = "_rhmeta_client"
	RHMeta_ActorID = "_rhmeta_actorid"
)
