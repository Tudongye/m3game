package monitor

const (
	CallPRCTotal         = "CallPRCTotal"         // 接受RPC数量
	HandleRPCTotal       = "HandleRPCTotal"       // 发送RPC数量
	CallRPCFailTotal     = "CallRPCFailTotal"     // 调起RPC失败数量
	HandleRPCFailTotal   = "HandleRPCFailTotal"   // 处理RPC失败数量
	HandleHealthRPCTotal = "HandleHealthRPCTotal" // 处理心跳检测数量
	ActorRuntimeTotal    = "ActorRuntimeTotal"    // Actor模式Actor数量
)
