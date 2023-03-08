# m3game

A game framework using Golang and Grpc

# m3game

一个基于Golang和Grpc的游戏后端框架。

框架分为GameLogic, Frame-Runtime, Custom-Plugin三层。Frame-Runtime为框架驱动层，负责消息驱动，服务网格，插件管理等核心驱动工作。Custom-Plugin为自定义插件层，框架层将第三方服务抽象为多种自定义插件接口，插件层根据实际的基础设施来进行实现。GameLogic为游戏逻辑层，用于承载实际的业务逻辑。框架使用protobuf来生成脚手架，通过引入自定义Option等方式将业务逻辑自动注入到框架层中。

优势：

1，简单但不简陋。框架包含了一个重度游戏后端的完备实现。

2、自动化的逻辑注入。借助pb的自定义选项，业务逻辑只需要很少的代码，就可以自动的注入到框架层

3、拒绝定制化工具。框架的代码生成和逻辑注入只依赖原生的protobuf和grpc，不需要额外安装定制化工具

![未命名文件 (2)](https://user-images.githubusercontent.com/16680818/222721483-8f14f7f2-7bb9-4eb2-8688-1367a67ed2ac.png)

Mutil,Async,Actor-Server: 游戏后台常见的业务模式，分别对应并发，单线程异步，Actor模式

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

![image](https://user-images.githubusercontent.com/16680818/223596615-35c4111e-0d8f-489a-85a5-78f6932145d3.png)


## 一个简单的样例

demo/dirapp 是一个无状态的并发服务，该服务提供Hello的RPC接口，其在业务层包含一个App 和 一个MutilServer

Step1、定义服务 proto

```
// demo/proto/dirapp.proto
syntax = "proto3";

package proto;

option go_package = "proto/pb";


import "pkg.proto";     // 框架层基础定义
import "options.proto"; // 自定义选线

// 定义服务与RPC路由
service DirSer {
    rpc Hello(Hello.Req) returns (Hello.Rsp) ;
}

// 定义RPC参数
message Hello {
    option (rpc_option).route_key = ""; // 路由Key字段名
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
Step2、定义逻辑实体DirSer
```
package dirserver

// 创建服务实体
func New() *DirSer {
	return &DirSer{
		Server: mutil.New("DirSer"),  // Mutil，并发服务
	}
}

type DirSer struct {
	*mutil.Server
	dpb.UnimplementedDirSerServer
}

// 接口文件由pb自动生成，业务层自行实现
func (d *DirSer) Hello(ctx context.Context, in *dpb.Hello_Req) (*dpb.Hello_Rsp, error) {
	out := new(dpb.Hello_Rsp)
	sctx := server.ParseContext(ctx)
	out.Rsp = fmt.Sprintf("Hello , %s", in.Req)
	if sctx != nil {
		if v, ok := sctx.Reciver().Metas().Get(proto.META_CLIENT); ok && v == proto.META_FLAG_TRUE {
			out.Rsp = fmt.Sprintf("Hello Client , %s", in.Req)
		}
	}
	return out, nil
}
// 逻辑注入接口
func (s *DirSer) TransportRegister() func(grpc.ServiceRegistrar) error {
	return func(t grpc.ServiceRegistrar) error {
		dpb.RegisterDirSerServer(t, s)
		return nil
	}
}
```
Step3、定义服务实体DirApp
```
package dirapp

// 创建DirApp实体
func newApp() *DirApp {
	return &DirApp{
		App: app.New(dproto.DirAppFuncID),  
	}
}

type DirApp struct {
	app.App
}

func (d *DirApp) Start(wg *sync.WaitGroup) error {
	router := plugin.GetRouterPlugin()
	if router != nil {
  // 服务注册
		if err := router.Register(d); err != nil {
			return err
		}
	}
	return nil
}
func (d *DirApp) HealthCheck() bool {
	return true
}

// 启动服务
func Run() error {
	plugin.RegisterFactory(&consul.Factory{})  // 注册服务发现组件
	plugin.RegisterFactory(&nats.Factory{})    // 注册broker组件
	runtime.Run(newApp(), []server.Server{dirserver.New()})   // 逻辑注入到框架层运行
	return nil
}
```
如下是从实例2,ClientApp 向 实例1 DirApp 发起RPC调用的调用链

DirSer 和 DirClient 是由dir.proto生成RPC调用服务端和客户端，protobuf保证双端协议一致。

DirApp 和 ClientApp 都是服务网格中的服务实体，Router会根据服务状态和路由策略最终选取一个服务实体发送请求。

Rumtime为框架驱动，根据RPC请求的性质，选择不同的传输路径（比如单向Ntify，广播，多播等）

Transport 内建了一个绑定在TcpConn的GrpcSer，用于服务实体间通讯

![未命名文件 (6)](https://user-images.githubusercontent.com/16680818/222782344-279fe08d-73f9-40f6-8bf2-5e3d4d56510e.png)


## RPC驱动

在M3中所有的跨服务功能调用都依托RPC进行，RPC接口通过pb-grpc生成

如下是一个RPC定义的proto。
```
// 定义服务与RPC路由
service DirSer {
    rpc Hello(Hello.Req) returns (Hello.Rsp) ;
}

// 定义RPC参数
message Hello {
    option (rpc_option).route_key = ""; // 当使用Hash路由时，路由Key字段名
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

其中rpc_option是M3为了减少重复编码而添加的自定义选项（大部分框架都使用定制化的代码生成工具，这使得那些框架很难被集成到原先的代码中）。自定义选项相关定义参看 options.proto，相关逻辑参看client/meta.go.调用RPCCall，M3框架会自动根据协议文件内容填充路由参数。


## 三种业务模型

游戏后台服务常见的业务模型有 Mutil 多线程，Async 单线程异步，Actor 模式 三种（暂时没见过更复杂的模型）

### Mutil

Mutil 多线程模型，主要用于无状态服务，M3采用原生Grpc服务实现。

### Async

Async 单线程异步，使用这类模型的服务不允许并发的执行RPC调用

M3在Async服务的RPC驱动链中加入了资源锁。通过资源锁确保同一时间只有一个RPC调用再执行

![未命名文件 (12)](https://user-images.githubusercontent.com/16680818/222913602-eca183aa-c449-4d30-af10-c2579fdc4346.png)

### Actor

Actor模型。使用这类模型的服务将RPC调用和游戏实体绑定，实体内部串行，实体之间并发。

M3为每个Actor分配一个执行Goroutine，并引入ActorRuntime和ActorMgr对Actor进行管理，前者用于管理单个Actor的执行Goroutine，后者用于管理整个Actor池。

M3在Actor服务的RPC调用链中加入了Actor管理逻辑，业务层逻辑都在Actor自己的Goroutine中执行。

![未命名文件 (13)](https://user-images.githubusercontent.com/16680818/222914612-a50f88b5-ad3f-4dc9-9b65-35078f83605d.png)


## 服务发现与路由

### Mesh

Mesh使用Router插件进行服务注册和服务发现，Router插件是必要插件，mesh/router/consul是一个基于Consul的Rotuer实现。

M3使用Grpc的Resolver- Balancer.Picker方式将服务网格与RPC路由相关联，相关逻辑参看mesh/resolver.go,balance.go

当前支持 P2P,Random,Hash,BroadCast,MutilCast,Single路由模式

|  路由模式   | 选路参数  | 选路规则  |
|  ----  | ----  | ----  |
| P2P  | 目标实例ID | 直接寻路 |
| Random  | 目标服务ID | 在目标服务中随机 |
| Hash  | 目标服务ID & 哈希Key | 在目标服务中按哈希key，一致性哈希映射寻路 |
| BroadCast  | 目标服务ID | 对目标服务所有实例广播 |
| MutilCast  | 目标TopicID | 对订阅目标TopciID的所有实例广播 |
| Single  | 目标服务ID | 对目标服务中ID最小的实例寻路 |

### 广播

M3基于Broker插件，实现了GrpcSer兼容的BrokerSer，用于处理BroadCast和MutilCast等单向Notify式RPC调用。

M3采用Interceptor的方式将BrokerSer注入RPC调用链，BrokerSer的相关实现参看 runtime/transport/brokerser.go。 broker/nats 是一个基于Nats的Broker实现。

## 资源管理

M3使用ResourceMgr进行资源管理，在M3中的资源指由GameLogic定义，在服务运行过程中需要实时热更新的资源文件。一般用于GameLogic的配置管理。

ResourceMgr使用双缓冲区模型，一主一备，主缓冲区用于资源访问，备缓冲区用于资源更新，每次热更新后主备缓冲区交换。相关逻辑参看resource/resourcemgr.go

M3对于资源的访问需要附带上下文context用于确认是资源访问还是资源更新

M3对于资源文件格式没有要求，只要求资源管理器提供Load接口，demo/loader/locationcfg.go是一个对于json配置文件的资源管理器样例。

```
type ResLoader interface {
	Load(ctx context.Context, cfgpath string) error // 资源更新
	Name() string
}
```


## 数据存储

M3采用pb管理游戏实体的DB存储结构。如下是一个简单实体的结构定义。相关实现参看demo/roleapp/roleser/roleactor.go

当前M3要求DB结构所有一级字段必须是string（必须是主键） 或 proto.Message（pb类型不可重复）,且DB结构必须设置一个string类型的主键。

```
message RoleDB {
    option (db_primary_key) = "RoleID";		 // DB主键
    string RoleID = 1;
    RoleName RoleName = 2;
    LocationInfo LocationInfo = 3;
}

message RoleName {
    string Name = 1;
}

message LocationInfo {
    int32 Location = 1;
    string LocateName = 2;
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
```

### Wraper

Wraper，对数据的ORM级封装，采用pb反射&泛型极大的简化了DB相关操作，同时封装了置脏管理。如下是Wraper定义

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
使用方式如下，以前述RoleDB为例

```
func roleDBCreater() *pb.RoleDB {
	return &pb.RoleDB{		// 所有的一级结构体都要初始化
		RoleID:       "",
		RoleName:     &pb.RoleName{},
		LocationInfo: &pb.LocationInfo{},
	}
}
rolemeta := db.NewMeta("role_table", roleDBCreater)
wp := wraper.New(rolemeta, "RoleID")	// 构建Wraper

// 读数据
dbplugin := plugin.GetDBPlugin()
wp.Read(dbplugin)

// 修改用户名
rolename, _ := wraper.Getter[*pb.RoleName](wp)	 // 参看前述 要求DB的一级pb字段类型不能重复
rolename.Name = "王小明"
wp.Setter(a.wraper, rolename)

// 脏字段写回
if wp.HasDirty() {
	wp.Update(dbplugin)
}
```


## 熔断管理

## 服务网关

## 监控管理

## 链路追踪

## 日志管理

## 容灾&&扩缩容

## 压测

## 灰度

## 异步任务

## 热更新

## 集群部署

