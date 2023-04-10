package errs

//go:generate stringer -type=Code
import (
	"fmt"

	"github.com/pkg/errors"
)

type Code uint64

const (
	Unknow                           Code = iota // 未知错误
	MsgUnmarshFail                   Code = 1    // 消息解码失败
	MsgMarshFail                     Code = 2    // 消息编码失败
	MeshNoAvalibleDstApp             Code = 101  // Mesh寻找路由未找到目标App
	MeshInitFail                     Code = 102  // Mesh初始化失败
	ActorCantAutoCreate              Code = 201  // Actor模式 禁止创建Actor
	TransportInitFail                Code = 301  // Transport初始化失败
	TransportCliCantFindTopic        Code = 302  // Transport ClientInterceptor 无法找到Topic
	TransportRegisterSerFail         Code = 303  // Transport 向Grpc注册服务失败
	TransportInsHasNewed             Code = 304  // Transport实例已创建
	TransportSetupFail               Code = 305  // Transport创建实例失败
	BrokerSerRegisterSerFail         Code = 401  // BrokerSer 向BrokerSer注册服务失败
	BrokerSerSetBrokerFail           Code = 402  // BrokerSer 设置Broker
	BrokerSerHandlerNotFind          Code = 403  // BrokerSer 未找到Handler
	BrokerSerClose                   Code = 405  // BrokerSer 已关闭
	NatsTransportReqCtxDone          Code = 406  // NatsTrans Ctx已取消
	NatsTransportReqTimeOut          Code = 407  // NatsTrans 请求已超时
	NatsTransportAckNotfindCB        Code = 408  // NatsTrans Ack回包找不到对应的CallBack
	NatsTransportAckCBTimeOut        Code = 409  // NatsTrans Ack回包CallBack已超时
	NatsTransportBalanceErr          Code = 410  // NatsTrans 路由寻路异常
	NatsTransportBalanceNoAva        Code = 411  // NatsTrans 路由寻路未找到可用路由
	NatsTransportBalanceNotFind      Code = 412  // NatsTrans 未找到对应的Balancer
	RPCCantFindHashKey               Code = 501  // RPC未找到HashKey
	RPCMethodNotRegister             Code = 502  // Method没有注册到RPC
	RPCCallFuncFail                  Code = 503  // RPC调起用户Func失败
	RPCInJectFail                    Code = 504  // RPC注入异常
	PluginsInitFail                  Code = 601  // Plugins初始化失败
	PluginsFactoryNotRegister        Code = 602  // Plugin的Factory没有注册
	PluginsSetupFail                 Code = 603  // Plugin调起Setup创建实例失败
	PluginsAddPluginFail             Code = 604  // Plugin添加Plugin失败
	PluginsReloadFail                Code = 605  // Plugins重新加载配置失败
	ActorRuntimeAllocLeaseFail       Code = 701  // ActorRuntime申请Lease失败
	ActorRuntimeFreeLeaseFail        Code = 702  // ActorRuntime释放Lease失败
	ActorRuntimeCallHandleActorDone  Code = 703  // ActorRuntime调用Handle失败Actor已退出
	ActorRuntimeCallHandleRPCDone    Code = 704  // ActorRuntime调用Handle失败RPC已退出
	ActorRuntimePushReqFailActorDone Code = 705  // ActorRuntime发送请求失败Actor已退出
	ActorRuntimePushReqFailChanFull  Code = 706  // ActorRuntime发送请求失败Chan已满
	ActorRuntimeKickFailActorDone    Code = 707  // ActorRuntime踢下线失败Actor已退出
	ActorRuntimeKickFailRPCDone      Code = 708  // ActorRuntime踢下线失败RPC已退出
	ActorMoveOutFail                 Code = 801  // Actor执行MoveOut数据迁出
	ActorOnInitFail                  Code = 802  // Actor执行OnInit失败
	ActorServerInitFail              Code = 803  // ActorServer初始化失败
	ActorKickNoFindActor             Code = 804  // Actor踢下线，未找到对应Actor
	RuntimeRegisterRepatedServer     Code = 901  // Runtime注册了重复的Server
	ResourceLoadFail                 Code = 1001 // Resource加载资源错误
	MetaRouteAppParseFail            Code = 1101 // Meta RouteApp解析失败
	MetaRouteSvcParseFail            Code = 1102 // Meta RouteSvc解析失败
	MetaRouteWorldParseFail          Code = 1103 // Meta RouteWorld解析失败
	BrokerInsIsNill                  Code = 5001 // Broker实例不存在
	BrokerInsHasNewed                Code = 5002 // Broker实例已创建
	NatsSetupFail                    Code = 5003 // Nats创建实例失败
	DBInsHasNewed                    Code = 5101 // DB实例已创建
	DBKeyNotFound                    Code = 5102 // 未找到对应数据
	DBDuplicateEntry                 Code = 5103 // 重复插入主键
	MongoSetupFail                   Code = 5104 // Mongo创建实例失败
	RedisSetupFail                   Code = 5105 // Redis创建实例失败
	RedisDelFail                     Code = 5106 // Redis Del失败
	GateInsHasNewed                  Code = 5201 // Gate实例已创建
	GrpcGateSetUpFail                Code = 5202 // GrpcGate创建实例失败
	GrpcGateConnClosed               Code = 5203 // GrpcGate Conn已关闭
	GrpcGateSendFailRPCDone          Code = 5204 // GrpcGate 发送消息失败RPC已退出
	GrpcGateSendFailChanFull         Code = 5205 // GrpcGate 发送消息失败队列已满
	LeaseInsIsNill                   Code = 5301 // Lease实例不存在
	LeaseInsHasNewed                 Code = 5302 // Lease实例已创建
	EtcdSetupFail                    Code = 5303 // Etcd创建实例失败
	EtcdIsClosed                     Code = 5304 // Etcd已关闭
	EtcdAllocLeaseFail               Code = 5305 // Etcd申请Lease失败
	LogInsHasNewed                   Code = 5401 // Log实例已创建
	ZlogSetupFail                    Code = 5402 // Zlog创建实例失败
	MetricInsHasNewed                Code = 5501 // Metric实例已创建
	PromSetupFaul                    Code = 5502 // Prometheus实例创建失败
	PromRegisterConsulFail           Code = 5503 // Prometheus注册Consul失败
	RouterInsHasNewed                Code = 5601 // Router实例已创建
	ConsulSetupFail                  Code = 5602 // Consul实例创建失败
	ConsulRegisterAppFail            Code = 5603 // Consul注册App失败
	ConsulGetAllInstanceFail         Code = 5604 // Consul获取Svc的App列表失败
	ShapeInsHasNewed                 Code = 5701 // Shape实例已创建
	ShapeRuleInitFail                Code = 5702 // Shape规则读取失败
	SentinelSteupFail                Code = 5703 // Sentinel实例创建失败
	SentinelRegisterRuleFail         Code = 5704 // Sentinel读规则失败
	TraceInsHasNewed                 Code = 5801 // Trace实例已创建
)

func (m Code) New(format string, params ...interface{}) *M3Err {
	return &M3Err{
		error: errors.Errorf("%d-%s: %s", m, m.String(), fmt.Sprintf(format, params...)),
		Code:  m,
	}
}

func (m Code) Wrap(err error, format string, params ...interface{}) *M3Err {
	return &M3Err{
		error: errors.Wrapf(err, "%d-%s: %s", m, m.String(), fmt.Sprintf(format, params...)),
		Code:  m,
	}
}

func (m Code) Is(err error) bool {
	if e, ok := err.(*M3Err); !ok {
		return false
	} else if e.Code == m {
		return true
	}
	return false
}

type M3Err struct {
	error
	Code Code
}

func New(code int, format string, params ...interface{}) *M3Err {
	return &M3Err{
		error: errors.Errorf(format, params...),
		Code:  Code(code),
	}
}
