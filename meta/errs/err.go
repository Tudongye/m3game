package errs

//go:generate stringer -type=Code
import (
	"fmt"

	"github.com/pkg/errors"
)

type Code uint64

const (
	Unknow                           Code = iota // 未知错误
	MsgUnmarshFail                               // 消息解码失败
	MsgMarshFail                                 // 消息编码失败
	MeshNoAvalibleDstApp             Code = 101  // Mesh寻找路由未找到目标App
	MeshInitFail                                 // Mesh初始化失败
	ActorCantAutoCreate              Code = 201  // Actor模式 禁止创建Actor
	TransportInitFail                Code = 301  // Transport初始化失败
	TransportCliCantFindTopic                    // Transport ClientInterceptor 无法找到Topic
	TransportRegisterSerFail                     // Transport 向Grpc注册服务失败
	BrokerSerRegisterSerFail         Code = 401  // BrokerSer 向BrokerSer注册服务失败
	BrokerSerSetBrokerFail                       // BrokerSer 设置Broker
	BrokerSerHandlerNotFind                      // BrokerSer 未找到Handler
	RPCCantFindHashKey               Code = 501  // RPC未找到HashKey
	RPCMethodNotRegister                         // Method没有注册到RPC
	RPCCallFuncFail                              // RPC调起用户Func失败
	RPCInJectFail                                // RPC注入异常
	PluginsInitFail                  Code = 601  // Plugins初始化失败
	PluginsFactoryNotRegister                    // Plugin的Factory没有注册
	PluginsSetupFail                             // Plugin调起Setup创建实例失败
	PluginsAddPluginFail                         // Plugin添加Plugin失败
	PluginsReloadFail                            // Plugins重新加载配置失败
	ActorRuntimeAllocLeaseFail       Code = 701  // ActorRuntime申请Lease失败
	ActorRuntimeFreeLeaseFail                    // ActorRuntime释放Lease失败
	ActorRuntimeCallHandleActorDone              // ActorRuntime调用Handle失败Actor已退出
	ActorRuntimeCallHandleRPCDone                // ActorRuntime调用Handle失败RPC已退出
	ActorRuntimePushReqFailActorDone             // ActorRuntime发送请求失败Actor已退出
	ActorRuntimePushReqFailChanFull              // ActorRuntime发送请求失败Chan已满
	ActorRuntimeKickFailActorDone                // ActorRuntime踢下线失败Actor已退出
	ActorRuntimeKickFailRPCDone                  // ActorRuntime踢下线失败RPC已退出
	ActorMoveOutFail                 Code = 801  // Actor执行MoveOut数据迁出
	ActorOnInitFail                              // Actor执行OnInit失败
	ActorServerInitFail                          // ActorServer初始化失败
	ActorKickNoFindActor                         // Actor踢下线，未找到对应Actor
	RuntimeRegisterRepatedServer     Code = 901  // Runtime注册了重复的Server
	ResourceLoadFail                 Code = 1001 // Resource加载资源错误
	MetaRouteAppParseFail            Code = 1101 // Meta RouteApp解析失败
	MetaRouteSvcParseFail                        // Meta RouteSvc解析失败
	MetaRouteWorldParseFail                      // Meta RouteWorld解析失败
	BrokerInsIsNill                  Code = 5001 // Broker实例不存在
	BrokerInsHasNewed                            // Broker实例已创建
	NatsSetupFail                                // Nats创建实例失败
	DBInsHasNewed                    Code = 5101 // DB实例已创建
	DBKeyNotFound                                // 未找到对应数据
	DBDuplicateEntry                             // 重复插入主键
	MongoSetupFail                               // Mongo创建实例失败
	RedisSetupFail                               // Redis创建实例失败
	RedisDelFail                                 // Redis Del失败
	GateInsHasNewed                  Code = 5201 // Gate实例已创建
	GrpcGateSetUpFail                            // GrpcGate创建实例失败
	GrpcGateConnClosed                           // GrpcGate Conn已关闭
	GrpcGateSendFailRPCDone                      // GrpcGate 发送消息失败RPC已退出
	GrpcGateSendFailChanFull                     // GrpcGate 发送消息失败队列已满
	LeaseInsIsNill                   Code = 5301 // Lease实例不存在
	LeaseInsHasNewed                             // Lease实例已创建
	EtcdSetupFail                                // Etcd创建实例失败
	EtcdIsClosed                                 // Etcd已关闭
	EtcdAllocLeaseFail                           // Etcd申请Lease失败
	LogInsHasNewed                   Code = 5401 // Log实例已创建
	ZlogSetupFail                                // Zlog创建实例失败
	MetricInsHasNewed                Code = 5501 // Metric实例已创建
	PromSetupFaul                                // Prometheus实例创建失败
	PromRegisterConsulFail                       // Prometheus注册Consul失败
	RouterInsHasNewed                Code = 5601 // Router实例已创建
	ConsulSetupFail                              // Consul实例创建失败
	ConsulRegisterAppFail                        // Consul注册App失败
	ConsulGetAllInstanceFail                     // Consul获取Svc的App列表失败
	ShapeInsHasNewed                 Code = 5701 // Shape实例已创建
	ShapeRuleInitFail                            // Shape规则读取失败
	SentinelSteupFail                            // Sentinel实例创建失败
	SentinelRegisterRuleFail                     // Sentinel读规则失败
	TraceInsHasNewed                 Code = 5801 // Trace实例已创建
)

type M3Err struct {
	error
	Code Code
}

func (m Code) New(format string, params ...interface{}) *M3Err {
	return &M3Err{
		error: errors.Errorf("%s: %s", m.String(), fmt.Sprintf(format, params...)),
		Code:  m,
	}
}

func (m Code) Wrap(err error, format string, params ...interface{}) *M3Err {
	return &M3Err{
		error: errors.Wrapf(err, "%s: %s", m.String(), fmt.Sprintf(format, params...)),
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
