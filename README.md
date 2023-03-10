# m3game

一个基于Golang和Grpc的游戏后端框架。

A game framework using Golang and Grpc

M3Game是一个采用Golang重构游戏后端框架的尝试，其旨在探索基于Golang的游戏后台开发方案。

框架分为GameLogic，Frame-Runtime，Custom-Plugin三层。Frame-Runtime为框架驱动层，负责消息驱动，服务网格，插件管理等核心驱动工作。Custom-Plugin为自定义插件层，框架层将第三方服务抽象为多种自定义插件接口，插件层根据实际的基础设施来进行实现。GameLogic为游戏逻辑层，用于承载实际的业务逻辑。框架使用protobuf来生成脚手架，通过引入pb.Option等方式将业务逻辑自动注入到框架层中。

优势：

1，简单但不简陋。框架包含了一个重度游戏后端的完备实现。

2、自动化的逻辑注入。借助pb的自定义选项，业务逻辑只需要很少的代码，就可以自动的注入到框架层

3、没有自定义代码生成器。框架的代码生成和逻辑注入只依赖原生的protobuf和grpc，不需要额外安装定制化工具

![未命名文件 (2)](https://user-images.githubusercontent.com/16680818/222721483-8f14f7f2-7bb9-4eb2-8688-1367a67ed2ac.png)

Mutil，Async，Actor-Server: 游戏后台常见的业务模式，分别对应并发，单线程异步，Actor模式

App: 用于承载业务逻辑的服务实体，是服务网格中的独立个体，由“环境ID.区服ID.功能ID.实例ID”唯一标识。一个App可以承载一个或多个Server

Client：RPC客户端，由服务提供方编写，包含一些参数校验，和路由规则

ResourceLoader: 可线上热更新的资源加载器，一般用于GameLogic Config的管理

Runtime: 框架驱动器

Transport: 提供服务之间Req-Rsp式RPC调用能力，采用tcp/GrpcSer实现一对一传输

BroekerSer：提供服务之间单向Ntify式RPC调用能力，采用Broker-plugin实现一对多传输

Mesh：服务网格，内含一组路由规则，以及规则对应的选路逻辑。采用Router-Plugin实现服务发现和服务注册

ResourceMgr: 资源管理器

PluginMgr：插件管理器

Router-Plugin： 路由组件，提供服务注册和服务发现的能力。当前有一个Consul实现

DB-Plugin: 存储组件，提供数据存储能力，当前有一个内存数据库实现

Broker-Plugin：消息队列组件，提供针对主题的发布和订阅功能，当前有一个Nats实现

Log-Plugin: 日志组件。

Trace-Plugin: 链路追踪组件

Metric-Plugin: 监控组件

## M3包依赖

![image](https://user-images.githubusercontent.com/16680818/224403501-995bb24a-8199-416b-98a6-082ea095376b.png)

## Example

example 是一组简单服务的样例，用来展示M3框架提供的基础能力。

example/simpleapp 是一个HelloWorld服务。

example/mutilapp 是一个并发服务，提供Hello，TraceHello(链路追踪)，BreakHello(熔断限流) 接口

example/asyncapp 是一个单线程异步服务，提供PostChannel(广播处理)，SSPullChannel(单线程阻塞) 接口

example/actorapp 是一个Actor模型服务，提供 Register(一个App部署多个Server)，Login(DB数据加载)，ModifyName(自动置脏标记)，LvUp(自动置脏标记)，GetInfo(资源配置)，PostChannel(广播发送)，PullChannel(服务间RPC调用)。

example/gateapp 是一个网关服务，对外提供Http接口(服务网关)访问内部服务。

example/test 是一个模拟客户端发包程序，内置多个测试用例。

## HelloWorld

以example/simpleapp为例

Step1、定义服务 proto，生成pb文件

```
// example/proto/simpleapp.proto
syntax = "proto3";

package proto;

import "pkg.proto";		// 框架文件
import "options.proto";		// 框架文件

option go_package = "proto/pb";

// 定义SimpleSer服务
service SimpleSer {
    rpc HelloWorld(HelloWorld.Req) returns (HelloWorld.Rsp);	 // 定义接口
}

// 定义RPC
message HelloWorld {
    option (rpc_option).route_key = "";
    message Req {
        RouteHead RouteHead = 1;
        string Req = 2;
    }
    message Rsp {
        RouteHead RouteHead = 1;
        string Rsp = 2;
    }
}

```

Step2、编写App代码

```
// example/simpleapp/simpleapp.go
package simpleapp

import (
	"m3game/example/proto"
	"m3game/example/simpleapp/simpleser"
	_ "m3game/plugins/broker/nats"
	_ "m3game/plugins/router/consul"
	"m3game/runtime"
	"m3game/runtime/app"
	"m3game/runtime/server"
)

// 创建App实体
func newApp() *SimpleApp {
	return &SimpleApp{
		App: app.New(proto.SimpleAppFuncID), // 指定App的FuncID
	}
}

type SimpleApp struct {
	app.App
}

// 健康检测
func (d *SimpleApp) HealthCheck() bool {
	return true
}

func Run() error {
	// 启动一个 包含了simpleser的SimpleApp
	runtime.Run(newApp(), []server.Server{simpleser.New()})
	return nil
}

```

Step3、定义服务实体simpleser

```
// example/simpleapp/simpleser
package simpleser

import (
	"context"
	"fmt"
	"m3game/example/proto/pb"
	"m3game/runtime/rpc"
	"m3game/runtime/server/mutil"

	"google.golang.org/grpc"
)

func init() {
	// 注册RPC信息到框架层
	if err := rpc.RegisterRPCSvc(pb.File_simple_proto.Services().Get(0)); err != nil {
		panic(fmt.Sprintf("RegisterRPCSvc SimpleSer %s", err.Error()))
	}
}

func New() *SimpleSer {
	return &SimpleSer{
		Server: mutil.New("SimpleSer"), // 以MutilSer为基础构建SimpleSer
	}
}

type SimpleSer struct {
	*mutil.Server
	pb.UnimplementedSimpleSerServer
}

// 实现HelloWorld接口
func (d *SimpleSer) HelloWorld(ctx context.Context, in *pb.HelloWorld_Req) (*pb.HelloWorld_Rsp, error) {
	out := new(pb.HelloWorld_Rsp)
	out.Rsp = fmt.Sprintf("HelloWorld , %s", in.Req)
	return out, nil
}

// 将SimpleSer注册到grpcser
func (s *SimpleSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		pb.RegisterSimpleSerServer(t, s)
		return nil
	}
}
```

step4 制作配置文件

```
[Transport]
Addr = "127.0.0.1:22105"	// 内部监听端口
BroadTimeOut = 5		// 广播处理超时

[Options]
[[Options.Mesh]]
WatcherInterSecond = 1		// 服务发现间隔

[Plugin]
[[Plugin.Router.router_consul]]
ConsulHost = "127.0.0.1:8500"	

[[Plugin.Broker.broker_nats]]
NatsURL = "127.0.0.1:4222"
```

Step5 编译运行

```
go build .

./main -idstr example.world1.simple.1 -conf ../../config/simpleapp.toml
```

![image](https://user-images.githubusercontent.com/16680818/224407634-3c464a0d-17bb-4f1b-8668-92a54a50d612.png)


## RPC驱动

在M3中所有的跨服务功能调用都依托RPC进行，RPC接口通过pb-grpc生成

如下是一个RPC定义的proto。
```
// 定义SimpleSer服务
service SimpleSer {
    rpc HelloWorld(HelloWorld.Req) returns (HelloWorld.Rsp);	 // 定义接口
}

// 定义RPC
message HelloWorld {
    option (rpc_option).route_key = "";
    message Req {
        RouteHead RouteHead = 1;
        string Req = 2;
    }
    message Rsp {
        RouteHead RouteHead = 1;
        string Rsp = 2;
    }
}
```
游戏RPC驱动通过ServerInterceptor注入Grpc，

游戏服务端驱动参看 transport/transport.go: RecvInterceptor。

![未命名文件 (11)](https://user-images.githubusercontent.com/16680818/222907647-cd2cf32e-c633-4cc8-95f5-187a10251e1f.png)

游戏客户端端驱动参看 client/client.go: SendInterceptor。

![未命名文件 (9)](https://user-images.githubusercontent.com/16680818/222907580-1d82955a-ef8f-45da-a897-e99a2f13b55c.png)

其中rpc_option是M3为了减少重复编码而添加的自定义选项，使用反射注入到框架层。自定义选项相关定义参看 options.proto，相关逻辑参看runtime/rpc.

当客户端调用RPCCall时，M3框架会自动根据协议文件内容填充路由参数。


## 三种业务模型

游戏后台服务常见的业务模型有 Mutil 多线程，Async 单线程异步，Actor 模式 三种（暂时没见过更复杂的模型）

### Mutil

Mutil 多线程模型，主要用于无状态服务，M3采用原生Grpc服务实现。参考实现 example/mutilapp/mutilser

### Async

Async 单线程异步，使用这类模型的服务不允许并发的执行RPC调用。参考实现 example/asyncapp/asyncser

M3在Async服务的RPC驱动链中加入了资源锁。通过资源锁确保同一时间只有一个RPC调用再执行

![未命名文件 (12)](https://user-images.githubusercontent.com/16680818/222913602-eca183aa-c449-4d30-af10-c2579fdc4346.png)

### Actor

Actor模型。使用这类模型的服务将RPC调用和游戏实体绑定，实体内部串行，实体之间并发。参考实现 example/actorapp/actorser

M3为每个Actor分配一个执行Goroutine，并引入ActorRuntime和ActorMgr对Actor进行管理，前者用于管理单个Actor的执行Goroutine，后者用于管理整个Actor池。

M3在Actor服务的RPC调用链中加入了Actor管理逻辑，业务层逻辑都在Actor自己的Goroutine中执行。

![未命名文件 (13)](https://user-images.githubusercontent.com/16680818/222914612-a50f88b5-ad3f-4dc9-9b65-35078f83605d.png)


## 服务发现与路由

### Mesh

Mesh使用Router插件进行服务注册和服务发现，Router插件是M3的必要插件，plugins/router/consul是一个基于Consul的Rotuer实现。

M3使用Grpc的Resolver & Picker方式将Mesh与RPC路由相关联，相关逻辑参看runtime/mesh/resolver.go，balance.go

当前支持 P2P，Random，Hash，BroadCast，MutilCast，Single路由模式

|  路由模式   | 选路参数  | 选路规则  |
|  ----  | ----  | ----  |
| P2P  | 目标实例ID | 直接寻路 |
| Random  | 目标服务ID | 在目标服务中随机 |
| Hash  | 目标服务ID & 哈希Key | 在目标服务中按哈希key，一致性哈希映射寻路 |
| BroadCast  | 目标服务ID | 对目标服务所有实例广播 |
| MutilCast  | 目标TopicID | 对订阅目标TopciID的所有实例广播 |
| Single  | 目标服务ID | 对目标服务中ID最小的实例寻路 |

### 广播

M3基于Broker插件，实现了GrpcSer兼容的BrokerSer，用于处理BroadCast和MutilCast等单向Notify式RPC调用。plugins/broker/nats 是一个基于Nats的Broker实现。

M3使用BrokerSer来处理广播，BrokerSer的相关实现参看 runtime/transport/brokerser.go。 

![未命名文件 (6)](https://user-images.githubusercontent.com/16680818/224411628-ce6afe7c-67b5-425e-bf32-003c600b08b5.png)


## 资源管理

M3中的资源指由GameLogic定义，在服务运行过程中需要实时热更新的资源文件。一般用于GameLogic的配置管理。

ResourceMgr使用双缓冲区模型，一主一备，主缓冲区用于资源访问，备缓冲区用于资源更新，每次热更新后主备缓冲区交换。相关逻辑参看resource/resourcemgr.go

M3对于资源的访问需要附带上下文context用于确认是资源访问还是资源更新

M3对于资源文件格式没有要求，只要求资源管理器提供Load接口，example/loader/titlecfgloader.go是一个对于json配置文件的资源加载器样例。

![未命名文件 (7)](https://user-images.githubusercontent.com/16680818/224412683-4511817c-55b9-4657-915d-d1d6d55cadec.png)


## 存储定义

M3采用pb管理游戏实体的DB存储结构。如下是一个简单实体的结构定义。相关实现参看example/actorapp/actor

当前M3要求DB结构所有一级字段必须是string（必须是主键） 或 proto.Message（pb类型不可重复），且DB结构必须设置一个string类型的主键。

```

message ActorDB {
    option (db_primary_key) = "ActorID";
    string ActorID = 1;
    ActorName ActorName = 2;
    ActorInfo ActorInfo = 3;
}

message ActorName {
    string Name = 1;
}

message ActorInfo {
    int32 Level = 1;
}


```

### DB结构注入

M3使用DB插件来进行实体数据的落地，M3通过PbReflect自动感知实体数据的DB结构Meta，DB插件根据Meta来对实体数据进行CRUD操作。DBMeta的生成逻辑参看 db/dbmeta.go

```
type DBMeta[T proto.Message] struct {
	objName   string
	table     string                                  // DB表名
	keyField  string                                  // 主键，强制为string
	allFields []string                                // 所有数据键
	allPBName map[string]string                       // 类型名到字段名映射
	creater   func() T                                // 游戏实体工场
	fieldds   map[string]protoreflect.FieldDescriptor // 游戏实体字段反射信息
}

type DB interface {
	Read(meta DBMetaInter, key string, filters ...string) (proto.Message, error)
	Update(meta DBMetaInter, key string, obj proto.Message, filters ...string) error
	Create(meta DBMetaInter, key string, obj proto.Message, filters ...string) error
	Delete(meta DBMetaInter, key string) error
}
```

### Wraper

Wraper，对数据的ORM级封装，采用pb反射&泛型极大的简化了DB相关操作，同时封装了自动化的置脏管理。example/actorapp/actor是一个基于Wraper的实体样例

如下是Wraper定义

```
type Wraper[T proto.Message] struct {
	key    string          // 实体Key
	obj    T               // 实体pb.Message
	meta   *db.DBMeta[T]   // Meta
	dirtys map[string]bool // 置脏标记
}


func (w *Wraper[T]) Update(db db.DB) error	 // CRUD操作
func (w *Wraper[T]) Create(db db.DB) error
func (w *Wraper[T]) Delete(db db.DB) error
func (w *Wraper[T]) Read(db db.DB) error

func KeySetter[T proto.Message](wraper *Wraper[T], value string) error	 // key字段操作
func KeyGetter[T proto.Message](wraper *Wraper[T]) (string, error)
func Setter[P, T proto.Message](wraper *Wraper[T], value P) error	// 普通字段操作
func Getter[P, T proto.Message](wraper *Wraper[T]) (P, error)		 
```
使用方式如下，以前述ActorDB为例

```
func actorDBCreater() *pb.ActorDB {
	return &pb.RoleDB{		// 所有的一级结构体都要初始化
		ActorID:       "",
		ActorName:     &pb.ActorName{},
		ActorInfo:     &pb.ActorInfo{},
	}
}
actormeta := db.NewMeta("actor_table", actorDBCreater)
wp := wraper.New(actormeta, "ActorID123")	// 构建Wraper

// 读数据
dbplugin := plugin.GetDBPlugin()
wp.Read(dbplugin)

// 修改用户名
actorname, _ := wraper.Getter[*pb.ActorName](wp)	 // 参看前述 要求DB的一级pb字段类型不能重复
actorname.Name = "王小明"
wp.Setter(a.wraper, actorname)

// 脏字段写回
if wp.HasDirty() {
	wp.Update(dbplugin)
}
```


## 熔断限流

M3采用Shape组件进行流量管理，Shape组件采用Interceptor方式注入RPC调用链。shape/Sentinel是一个基于Sentinel实现的shape组件。

流量管理针对RPC进行，规则分为限流规则 FlowRule 和 熔断规则 BreakRule。

如下是对example/mutilapp的BreakHello的流量管理配置

```
[Rules]
Method = "/proto.MutilSer/BreakHello"	// RPC方法
[[Rules.FlowRules]]			// 限流规则
Threshold = 2				// 限流阈值
StatIntervalMs = 1000			// 统计周期
MaxQueueWaitMs = 0			// 限流最大等待时长
[[Rules.BreakRules]]			// 熔断规则
Threshold = 1				// 熔断阈值
Strategy = "ErrorCount"			// 熔断规则，错误请求数
StatIntervalMs = 1000			// 统计周期
RetryTimeOutMs = 2000			// 熔断恢复市场
MinRequestNum = 2			// 熔断生效最小请求次数
```

## 监控统计

M3采用Metric组件来进行监控统计，对于统计项分为Counter，Guage，Histogram，Summary四类。

metric/prometheus 是一个基于 prometheus 实现的Metric

```
type StatCounter interface {	// 计数器
	Add(float64)
	Inc()
}

type StatGauge interface {	// 测量器
	Set(float64)
	Sub(float64)
	Dec()
	Add(float64)
	Inc()
}

type StatHistogram interface {	// 直方图
	Observe(float64)
}

type StatSummary interface {	// 点分数
	Observe(float64)
}
```

## 链路追踪

M3的链路追踪采用Open-Telemetry方案，trace/stdout是一个直接向控制台打印的trace。

对于RPC是否启用Trace，采用pb.Option的方式进行定义。如下是一个开启Trace的TraceHello RPC的定义

```
message TraceHello {
    option (rpc_option).route_key = "";
    option (rpc_option).trace = true;	// 打开Trace
    message Req {
        RouteHead RouteHead = 1;
        string Req = 2;
    }
    message Rsp {
        RouteHead RouteHead = 1;
        string Rsp = 2;
    }
}
```

## 本地日志

M3采用Log组件进行本地日志管理，日志分为DEBUG，INFO，WARN，ERROR，FATAL 五个级别，log/zap 是一个基于zap实现的Log组件样例。

```
type Logger interface {
	Output(depth Depth, lv LogLv, plus LogPlus, format string, v ...interface{})
	SetLevel(level LogLv)
	GetLevel() LogLv
}

func Debug(format string, v ...interface{})	// 调试日志，只在开发环境开启
func Info(format string, v ...interface{}) 	// 重要行为日志，生产环境开启
func Warn(format string, v ...interface{}) 	// 警告日志，如果遇到问题，用于辅助检查
func Error(format string, v ...interface{})	// 错误日志，明确的逻辑异常，高度关注
func Fatal(format string, v ...interface{})	// 致命错误，必须立机告警处理

func DebugP(plus LogPlus, format string, v ...interface{})
func InfoP(plus LogPlus, format string, v ...interface{})
func WarnP(plus LogPlus, format string, v ...interface{})
func ErrorP(plus LogPlus, format string, v ...interface{})
func FatalP(plus LogPlus, format string, v ...interface{})
```

## 异步任务

## 服务网关

## 压测

## Demo(TODO)

M3Game是为了解决游戏开发时遇到的具体问题而构建的框架。为了更好的暴露问题,并验证解决方案，M3构建了一个重度游戏后端Demo作为集群化解决方案的载体。

Demo是一个分区式游戏，玩家(Role)数据按小区(World)隔离，玩家可以自由组建社团(Club)，核心玩法采用跨区匹配(Match)开单局(Fight)方式进行

游戏实体可以分为 玩家(Role)，小区(World)，社团(Club)，小区玩家关系(WorldRole)，社团玩家关系(ClubRole)，单局(Fight)

![未命名文件 (3)](https://user-images.githubusercontent.com/16680818/223912107-3d6c8c5c-7eb8-45a1-a820-75c49652257e.png)

服务实例包括DirApp(导航服务)，RoleApp(玩家服务)，ClubApp(社团服务)， ClubRoleApp(社团玩家服务)，WorldApp(小区服务)，WorlRoledApp(小区玩家服务)，MatchApp(匹配服务)，FightApp(战斗服务)，ZoneApp(战斗集群服务)

其中DirApp，ClubRoleApp，WorldRoleApp，MatchApp，ZoneApp为无状态服务，RoleApp，FightApp为激发式有状态服务(负载受玩家行为影响)，ClubApp 为常驻式动态负载有状态服务(负载不受玩家行为影响，且负载动态可变)，WorldApp为常驻式固定负载有状态服务(负载不受玩家行为影响，且负载固定)

![未命名文件 (4)](https://user-images.githubusercontent.com/16680818/223912598-982bc454-409e-46ec-b54b-84238194d582.png)

## 灰度发布

## 容灾

## 动态伸缩

## 热更新

## 集群部署
