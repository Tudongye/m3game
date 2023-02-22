package balance

import "google.golang.org/grpc/balancer"

const (
	BalanceAttrKey = "m3gblattrkey"
)

type BalanceAttrValue struct {
	IDStr string
}

func init() {
	balancer.Register(newM3GPikerBuilder())
}

/*
	RT_P2P 点对点
	RT_RAND 随机
	RT_HASH 哈希
	RT_BROAD 广播？
	RT_MUTIL 多播？
	RT_SINGLE 单例
	RT_DIRECT 直连？
*/
