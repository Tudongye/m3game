# m3game

一个基于Golang和Grpc的游戏后端框架。

A GameServer framework built using Golang and GRPC

M3Game是一个采用Golang构建游戏后端的尝试，期望能探索出一条Golang游戏后台的开发方案。

框架分为GameLogic，Frame-Runtime，Custom-Plugin三层。Frame-Runtime为框架驱动层，负责消息驱动，服务网格，插件管理等核心驱动工作。Custom-Plugin为自定义插件层，框架层将第三方服务抽象为多种插件接口，插件层根据实际的基础设施来进行实现。GameLogic为游戏逻辑层，用于承载实际的业务逻辑。框架使用protobuf来生成脚手架，可以通过在pb中添加Option的方式将业务层接口自动注入到框架层。

优势：

1，更加贴近实际业务。

2、自动化的逻辑注入。借助pb的自定义选项，业务逻辑只需要很少的代码，就可以自动的注入到框架层

3、更通用的技术和更低的门槛。M3基于golang主流的protobuf和grpc进行构建，没有繁琐的代码生成工具，上手门槛低。

4、这里有一个很有意思的置脏管理模块，只需要在pb中定义好数据和脏标记，就可以轻松实现置脏&批量写回功能。

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

DB-Plugin: 存储组件，提供数据存储能力，当前有内存数据库，redis，mongo实现

Broker-Plugin：消息队列组件，提供针对主题的发布和订阅功能，当前有一个Nats实现

Log-Plugin: 日志组件，当前有一个Zap实现。

Trace-Plugin: 链路追踪组件，当前接入opentelemetry标准。

Metric-Plugin: 监控组件，当前有一个prometheus实现

Shape-Plugin：流量治理组件，当前有一个sentinel实现

Gate-Plugin：服务网关组件，当前有一个grpc-stream实现

Lease-Plugin：租约管理组件，当前有一个etcd实现

## 集群化部署架构

![未命名文件 (12)](https://user-images.githubusercontent.com/16680818/227127493-b41dbfa7-3d65-4b1a-9331-0baf45b4b8fd.png)

## M3内部依赖

![graphviz](https://user-images.githubusercontent.com/16680818/226848482-d1facfba-8e86-4206-96f8-c786e385e862.png)

## 感谢GPT的CR

在GPT的帮助下，对runtime做了一轮优化

![企业微信截图_167946848074](https://user-images.githubusercontent.com/16680818/226848753-bdae1991-0396-426f-86b9-3bc2cde751e2.png)

## HelloWorld

以example/simpleapp为例

Step1、定义服务 proto，生成pb文件

```
// example/proto/simpleapp.proto
syntax = "proto3";
package proto;
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
        string Req = 1;
    }
    message Rsp {
        string Rsp = 1;
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

# 单实例开发方案(已完成)

## RPC驱动

在M3中所有的跨服务调用都依托RPC进行，RPC接口通过pb-grpc生成。M3框架的附加信息都存储在RPC的metadata中。

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
        string Req = 1;
    }
    message Rsp {
        string Rsp = 1;
    }
}
```

业务层通过编写rpc_option将RPC接口注入框架层，解析相关逻辑参看runtime/rpc。rpc_option定义如下

```
message M3GRPCOption {
    string route_key = 1;	// Hash路由时的key字段名
    bool ntf = 2;		// 是否是单向Nty
    bool trace = 3;		// 是否开启链路追踪
    bool cs = 4;		// 是否支持客户端访问
}
```

M3框架通过rpc注入和泛型编程，大大简化了业务层进行RPC调用时的操作，如下是对hello接口进行"随机选址"的RPCCall调用

```
func Hello(ctx context.Context, hellostr string, opts ...grpc.CallOption) (string, error) {
	var in pb.Hello_Req
	in.Req = hellostr
	// RPCCallRandom 接受泛型参数in,返回泛型参数out。
	// RPCCall通过入参in获取到对应的rpc_option，自动填充选址参数，并对常见RPC异常进行前置处理。
	out, err := client.RPCCallRandom(_client, _client.Hello, ctx, &in, opts...)
	if err != nil {
		return "", err
	} else {
		return out.Rsp, nil
	}
}
```

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

M3在Actor服务的RPC调用链中加入了Actor管理逻辑。对于Actor的RPC调用都在Actor自己的Goroutine中执行。

引入Lease-plugin可以保证一个Actor在分布式环境下至多只会在一个App上运行。参看rumtime/server/actor

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


## 实体存储

M3采用pb来定义游戏实体的DB存储结构。如下是一个简单实体的结构定义。相关实现参看example/actorapp/actor

```
message ActorDB {
    string ActorID = 1 [(dbfield_option) = { flag: "FActorID", primary: true }];	// 主键
    string Name    = 2 [(dbfield_option) = { flag: "FActorName" }];
    int32 Level    = 3 [(dbfield_option) = { flag: "FActorLevel" }];
}

enum AcFlag {
    FActorMin   = 0;
    FActorID    = 1;
    FActorName  = 2;
    FActorLevel = 3;
}
```

### DB结构注入

M3使用DB插件来对实体数据进行落地，M3根据实体的PB结构生成对应的dbmeta，DB插件不用感知业务数据的具体类型，直接根据Meta就可以对实体数据进行CRUD操作。

DBMeta的生成逻辑参看 db/dbmeta.go

```
type DBMetaInter interface {
	Setter(msg proto.Message, flag int32, data interface{}) // 赋值
	Getter(msg proto.Message, flag int32) interface{}       // 读取
	FlagKind(flag int32) protoreflect.Kind                  // 获取字段类型
	FlagName(flag int32) string                             // 获取字段类型
	KeyFlag() int32                                         // 主键字段
	AllFlags() []int32                                      // 所有字段名
	New() proto.Message
	Table() string
}
type DB interface {
	plugin.PluginIns
	Read(ctx context.Context, meta DBMetaInter, key interface{}, flags ...int32) (proto.Message, error)
	Update(ctx context.Context, meta DBMetaInter, key interface{}, obj proto.Message, flags ...int32) error
	Create(ctx context.Context, meta DBMetaInter, key interface{}, obj proto.Message) error
	Delete(ctx context.Context, meta DBMetaInter, key interface{}) error

	ReadMany(ctx context.Context, meta DBMetaInter, filters interface{}, flags ...int32) ([]proto.Message, error)
}
```

### Wraper

Wraper，对数据的ORM级封装，采用反射&泛型极大的简化了DB操作，同时封装了一套置脏管理。example/actorapp/actor是一个基于Wraper的实体样例

如下是Wraper定义

```
type Wraper[TM proto.Message, TF Flag] struct {
	meta   *WraperMeta[TM, TF] // Meta
	key    interface{}         // 主键值
	obj    TM                  // 原始数据
	dirtys map[TF]bool         // 脏标记
}
func (w *Wraper[TM, TF]) Set(flag TF, value interface{}) 	
func (w *Wraper[TM, TF]) Get(flag TF) interface{}
func (w *Wraper[TM, TF]) Update(db db.DB) error	 // CRUD操作
func (w *Wraper[TM, TF]) Create(db db.DB) error
func (w *Wraper[TM, TF]) Delete(db db.DB) error
func (w *Wraper[TM, TF]) Read(db db.DB) error
```

使用方式如下，以前述ActorDB为例

pb定义
```
message ActorDB {
    string ActorID = 1 [(dbfield_option) = { flag: "FActorID", primary: true }];	// 主键
    string Name    = 2 [(dbfield_option) = { flag: "FActorName" }];
    int32 Level    = 3 [(dbfield_option) = { flag: "FActorLevel" }];
}

enum AcFlag {
    FActorMin   = 0;
    FActorID    = 1;
    FActorName  = 2;
    FActorLevel = 3;
}
```

```
dbmeta := db.NewMeta[*pb.ActorDB]("actor_table")
wrapermeata := db.NewWraperMeta[*pb.ActorDB, pb.AcFlag](db,eta)
wp := wrapermeata.New("ActorID123")	// 构建Wraper
// 读数据
dbplugin := plugin.GetDBPlugin()
wp.Read(ctx, dbplugin)
// 修改用户名
wp.Setter(pb.AcFlag_FActorName, "小明")
// 脏字段写回
if wp.IsDirty() {
	wp.Update(ctx, dbplugin)
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

metric/prometheus 是一个基于prometheu 实现的Metric

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

## 服务网关

服务网关用于管理与客户端的连接，并将客户端请求转化为Grpc-Reply请求。plugins/gate/grpcgate是一个 基于Grpc-Stream建立的Gate组件。

example/gateapp实现了一套将客户端请求转化为Grpc-Reply请求的通用方案。

```
type Gate interface {
	GetConn(playerid string) CSConn
}
type CSConn interface {
	Send(ctx context.Context, msg *metapb.CSMsg) error
	Kick()
}
type GateReciver interface {
	AuthCall(*metapb.AuthReq) (*metapb.AuthRsp, error)	// 建立连接时的鉴权接口
	LogicCall(*metapb.CSMsg) (*metapb.CSMsg, error)		// 将客户端请求转化为Grpc-Reply请求
}
```

## 租约管理

为了解决分布式系统下的数据一致性问题，M3引入了租约(悲观锁)，Lease-plugin，plugins/lease/etcd 是一个基于etcd的实现。

使用租约来保护数据的所有权，可以保证在同一时间，整个分布式系统中最多只会有一个App可以操作该数据。

```
type LeaseMoveOutFunc func(context.Context) ([]byte, error) // 租约退出回调
type Lease interface {
	plugin.PluginIns
	AllocLease(ctx context.Context, id string, f LeaseMoveOutFunc) error // 获取租约
	FreeLease(ctx context.Context, id string) error                      // 释放租约
	KickLease(ctx context.Context, id string) ([]byte, error)            // 要求释放租约
	RecvKickLease(ctx context.Context, id string) ([]byte, error)        // 接受释放租约消息
	GetLease(ctx context.Context, id string) ([]byte, error)	     // 获取租约内容
}
type LeaseReciver interface {
	SendKickLease(ctx context.Context, id string, app string) ([]byte, error) // 发送释放租约消息
}
```

![未命名文件 (10)](https://user-images.githubusercontent.com/16680818/225224998-70fecb14-d28c-47d7-a49f-16516e3a53ae.png)

## 热更新

这里说的热更新指的是在不影响服务能力前提下，对线上程序的逻辑代码进行更新。对于使用编译型语言进行业务逻辑开发的程序，常见的热更新方式有

1、脚本虚拟机。在编译型语言中引入脚本语言，比如Lua，Python。在合理的位置“打桩”，当程序执行到“桩”点时，就会进入脚本虚拟机中，这时可以利用脚本语言的特性来实现程序热更新。

2、动态链接库。大部分的编译型语言都支持生成、链接动态库文件。可以通过线上更新挂载的动态库文件，实现程序热更新。

3、直接莽。线上更新程序文件，并重启程序，以实现业务逻辑更新的效果。这种方式必然会带来服务能力的损失，但是可以通过一些方式降低损失，比如滚动更新，共享内存等。

热更新方式的选择更新基于技术选型（比如你用skynet，那就没啥好说的了）。M3本身并不提供热更新能力，但是在demo/roleapp/roleser中，给出一个使用Lua & “打桩”的方式，实现拒绝RoleId最后一位数为“1”的玩家登陆的热更新实现。



## 自动化测试

## Example

example 是一组简单服务的样例，用来展示M3框架的单实例开发方案。

example/simpleapp 是一个HelloWorld服务。

example/mutilapp 是一个并发服务，提供Hello，TraceHello(链路追踪)，BreakHello(熔断限流) 接口

example/asyncapp 是一个单线程异步服务，提供PostChannel(广播处理)，SSPullChannel(单线程阻塞) 接口

example/actorapp 是一个Actor模型服务，提供 Register(一个App部署多个Server)，Login(DB数据加载)，ModifyName(自动置脏标记)，LvUp(自动置脏标记)，GetInfo(资源配置)，PostChannel(广播发送)，PullChannel(服务间RPC调用)。ActorApp加入了租约插件，确保每个Actor最多只会在一个ActorApp上存在

example/gateapp 是一个网关服务，客户端可以通过gprc-stream方式与网关建立长连接。

example/test 是一个模拟客户端发包程序，内置多个测试用例。

![未命名文件 (9)](https://user-images.githubusercontent.com/16680818/224889189-950ed58b-2b9f-470d-a096-282cd849767e.png)

example使用方式

```
1、修改example/config中 nats和router接口地址
2、依次启动mutilapp,asyncapp,actorapp,gateapp的main/start.sh
3、到test/main目录下执行测试用例命令
./main -testmode Hello -agenturl 127.0.0.1:22000 // helloworld用例
./main -testmode Trace -agenturl 127.0.0.1:22000 // helloworld链路追踪
./main -testmode Break -agenturl 127.0.0.1:22000 // helloworld流量治理
./main -testmode ActorCommon -agenturl 127.0.0.1:22000 // 注册，登陆，改名，升级，服务端到客户端主动通知
./main -testmode ActorBroadCast -agenturl 127.0.0.1:22000 // 注册，登陆，广播
./main -testmode ActorMove -agenturl 127.0.0.1:22000 // 数据一致性，一个Actor两个ActorApp之间进行服务迁移（需要启动ActorApp1 和 ActorApp2）
```

# 集群化部署方案(进行中)

## 集群部署

游戏后端服务的核心功能就是对业务数据进行增删改查。当服务采用多机部署时，就会引入对同一份数据的并发操作问题。而集群化部署所要处理的问题就是在多机环境下，如何分配数据管理权。

无状态服务，其每次请求的处理结果不依赖上下文，不需要长期占有数据的管理权，一般采用加锁 或者 CompareAndSwap来处理并发问题

有状态服务，其每次请求的处理结果依赖上下文，需要长期占有数据的管理权。根据其管理的数据不同，分为三种类型：

1、元数据。这类数据的总量小，数据管理权可以集中在一台机器上处理。集群部署时，一般采用主从模式，当主备宕机时，数据管理权整体迁移到备机。

2、轻量级数据。这类数据的总量大，且数据管理权的跨机迁移成本低（时间成本，资源成本），数据管理权分散在多台机器。集群部署时，一般采用对等部署，基于一致性哈希进行寻址，当机器增减时，数据管理权会动态调整。

3、重量级数据。这类数据的总量大，且数据管理权的跨机迁移成本高，数据管理权分散在多台机器。集群部署时，一般会专门指定一个管理进程（管理进程采用元数据方式部署），用于处理数据管理权的调度，尽量减少跨机的管理权迁移

对于集群部署方式，M3采用K8s的方式部署。本节的demo已经支持这种部署方式，详情参看"部署方式"


## 灰度发布

## 容灾

## 动态伸缩

# Demo

为了更好的暴露问题,并验证解决方案，M3构建了一个重度游戏后端Demo作为集群化解决方案的载体。

Demo是一个全服互通游戏，玩家(Role)可以自由组建社团(Club)，核心玩法采用匹配(Match)开单局(Fight)方式进行。

在部署上，所有的服务都采用集群化部署以支持容灾恢复和弹性伸缩能力。

游戏实体分为 玩家(Role)，社团(Club)，单局(Fight)。玩家(Role)实体只有对应玩家在线时才会激活，社团(Club)实体一经创建常驻激活，直到被解散，单局(Fight)实体在玩家开启单局期间才会激活且激活期间不易发生服务迁移。

服务实例包括GateApp(网关服务), UidApp(id管理服务)，OnlineApp(在线管理服务)，RoleApp(玩家服务)，ClubApp(社团服务)，MatchApp(匹配服务)，ZoneApp(战斗集群服务)，FightApp(战斗服务)

## 简单介绍一下

游戏后台服务一般分为玩法服务和外围服务。

玩法服务指与客户端表现直接关联的状态同步类服务，比如MMO的地图服务，Moba的单局服务等。这类服务同质化高，易于抽象，经常与客户端共用逻辑代码，甚至可以由unity，ue等客户端游戏引擎直接生成。

外围服务指的是玩法服务以外用于承载游戏逻辑的服务，比如吃鸡的大厅服务，好友服务，战队服务等。这类服务以数据管理为核心，与具体的业务逻辑相关比如数值成长，运营活动，很难抽象为通用架构，是游戏后端开发的主要工作。

## 外围服务

在本demo中外围服务包括GateApp, UidApp, OnlineApp, RoleApp, ClubApp，ClubMgrApp组成的部分，管理玩家(Role) 和 社团(Club)数据。

![未命名文件 (12)](https://user-images.githubusercontent.com/16680818/227127493-b41dbfa7-3d65-4b1a-9331-0baf45b4b8fd.png)

GateApp: 网关服务，无状态服务，客户端任意链接

UidApp: Id管理服务，包括玩家Openid到RoleId的映射，ClubId的分配

RoleApp：玩家服务，以Role为单位的Actor服务

OnlineApp：在线管理，维护Role在线状态，提供大量级租约服务。

ClubApp：社团服务，将Club划分为有限个Slot，以Slot为单位的Actor服务。

ClubRoleApp：社团玩家服务，管理社团和玩家的关联关系。

### 服务接口协议

首先编写服务接口协议

```
# GateApp
service GateSer {
    rpc SendToCli(SendToCli.Req) returns (SendToCli.Rsp);	// 向客户端主动推送
}
# UidApp
service UidSer {
    rpc AllocRoleId(AllocRoleId.Req) returns (AllocRoleId.Rsp); // 分配RoleID
    rpc AllocClubId(AllocClubId.Req) returns (AllocClubId.Rsp); // 分配ClubID
}
# RoleApp
service RoleSer {
    rpc RoleLogin(RoleLogin.Req) returns (RoleLogin.Rsp);   // 登陆注册
    rpc RoleGetInfo(RoleGetInfo.Req) returns (RoleGetInfo.Rsp); // 获取详情
    rpc RoleModifyName(RoleModifyName.Req) returns (RoleModifyName.Rsp);    // 改名
    rpc RolePowerUp(RolePowerUp.Req) returns (RolePowerUp.Rsp);    // 战力提升
    rpc RoleKick(RoleKick.Req) returns (RoleKick.Rsp);    // 服务迁移
    rpc RoleGetClubInfo(RoleGetClubInfo.Req) returns (RoleGetClubInfo.Rsp); // 获取社团信息
    rpc RoleGetClubList(RoleGetClubList.Req) returns (RoleGetClubList.Rsp); // 获取社团列表
    rpc RoleGetClubRoleInfo(RoleGetClubRoleInfo.Req) returns (RoleGetClubRoleInfo.Rsp); // 获取玩家社团信息
    rpc RoleCreateClub(RoleCreateClub.Req) returns (RoleCreateClub.Rsp); // 创建社团
    rpc RoleJoinClub(RoleJoinClub.Req) returns (RoleJoinClub.Rsp); // 加入社团
    rpc RoleExitClub(RoleExitClub.Req) returns (RoleExitClub.Rsp); // 退出社团
    rpc RoleCancelClub(RoleCancelClub.Req) returns (RoleCancelClub.Rsp); // 解散社团
}
# OnlineApp
service OnlineSer {
    rpc OnlineCreate(OnlineCreate.Req) returns (OnlineCreate.Rsp);   // 创建在线状态
    rpc OnlineRead(OnlineRead.Req) returns (OnlineRead.Rsp);   // 获取在线情况
    rpc OnlineDelete(OnlineDelete.Req) returns (OnlineDelete.Rsp);   // 删除在线状态
}
# ClubApp
service ClubSer {
    rpc ClubCreate(ClubCreate.Req) returns (ClubCreate.Rsp);   // 创建社团
    rpc ClubGetInfo(ClubGetInfo.Req) returns (ClubGetInfo.Rsp);   // 创建社团
    rpc ClubJoin(ClubJoin.Req) returns (ClubJoin.Rsp);   // 加入社团
    rpc ClubExit(ClubExit.Req) returns (ClubExit.Rsp);   // 退出社团
    rpc ClubCancel(ClubCancel.Req) returns (ClubCancel.Rsp);   // 解散社团
}
service ClubDaemonSer {
    rpc ClubKick(ClubKick.Req) returns (ClubKick.Rsp);    // 服务迁移
}
# ClubRoleApp
service ClubRoleSer {
    rpc ClubRoleRead(ClubRoleRead.Req) returns (ClubRoleRead.Rsp);   // 查询Role归属Club
    rpc ClubRoleCreate(ClubRoleCreate.Req) returns (ClubRoleCreate.Rsp);   // 创建Role-Club关系
    rpc ClubRoleDelete(ClubRoleDelete.Req) returns (ClubRoleDelete.Rsp);   // 删除Role-Club关系
}
```

### UidApp

UidApp 管理玩家OpenId（社交账户Id）到RoleId（游戏角色Id）的映射关系，以及ClubId(社团Id)的生成。其负载与单位时间内玩家登陆次数和创建社团次数相关，署于小负载服务。类比元数据类服务，这里采用主从部署，使用Lease（租约）保证数据一致性，使用Single（最小ID）方式选主。参看 demo/uidapp

### OnlineApp

OnlineApp 管理玩家的在线状态（玩家所处的RoleApp信息），其提供了一种大规模Lease服务，用来保证玩家数据在RoleApp上的一致性，当玩家登陆时RoleApp会先在OnlineApp申请Lease后再提供服务(如果Lease冲突，则把原Lease踢下线)。其负载与单位时间内玩家登陆次数相关，署于小负载服务。类比元数据类服务，这里采用主从部署，使用Lease（租约）保证数据一致性，使用Single（最小ID）方式选主。参看 demo/onlineapp

### RoleApp

RoleApp 持有玩家RoleDB的数据管理权，负责RoleDB层面的游戏逻辑，其负载与单位时间内所有在线玩家的总操作次数相关，署于大负载服务，且数据迁移成本低。类比轻量数据服务，这里采用对等部署，使用OnlineApp提供的Lease服务保证数据一致性，使用Hash（一致性哈希）方式寻址。参看demo/roleapp

### GateApp

GateApp 管理玩家链接，本质上是个代理服务，没有数据管理权需求，署于无状态服务。这里采用对等部署。参看demo/gateapp

### ClubRoleApp

ClubRoleApp 持有玩家与社团的映射关系，没有数据管理权，实时查询ClubApp的相关数据表。这里采用对等部署，随机路由。参看demo/clubroleapp

### ClubApp

ClubApp 管理社团数据，社团数量动态变化，ClubApp将社团按照Slot划分，以Slot为服务单位，采用Actor+Lease模式构建，对等部署。参看demo/clubapp

### Test

Test 测试客户端

```
sh test.sh Test1       // 单次 关键路径（登陆，修改Role数据，拉取Role数据）测试
sh test.sh MutilTest1  // 100TPS 10000次 关键路径（登陆，修改Role数据，拉取Role数据）测试
sh test.sh Test2       // 单次 社团路径（登陆，创建社团，退出社团）测试
```

### Demo部署

Demo采用Docker方式进行交付。采用Helm编排服务架构。采用K8s进行程序部署。

#### DockerBuild

使用demo/dockerbuild.sh脚本即可Docker容器构建。如下是使用的dockerfile

'''
# 基础镜像，包含一套dev环境
FROM golang:1.20-rc 
# Transport 40001 Metric 40002 Gate 40003
EXPOSE 40001/tcp 40002/tcp 40003/tcp
# 
COPY uidapp/main/main /go/bin/demo/uidapp/main/main
COPY roleapp/main/main /go/bin/demo/roleapp/main/main
COPY onlineapp/main/main /go/bin/demo/onlineapp/main/main
COPY gateapp/main/main /go/bin/demo/gateapp/main/main
COPY clubapp/main/main /go/bin/demo/clubapp/main/main
COPY clubroleapp/main/main /go/bin/demo/clubroleapp/main/main
COPY test/main/main /go/bin/demo/test/main/main
COPY resource /go/bin/demo/resource
COPY deploy /go/bin/demo/deploy
'''

#### Helm

Helm相关配置在demo/helm/m3demo。部署前需要修改demo/helm/m3demo/values.yaml中的image字段
```
m3demoimage: m3demo:latest
```

#### 部署

![image](https://user-images.githubusercontent.com/16680818/228788248-4fbc57fd-d4d1-47c1-b3a3-4f773b846390.png)




