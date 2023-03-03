# m3game

A game framework using GO language and Grpc

# m3game

一个基于Golang和Grpc的游戏后端框架。

框架分为GameLogic, Frame-Runtime, Custom-Plugin三层。Frame-Runtime为框架驱动层，负责消息驱动，服务网格，插件管理等核心驱动工作。Custom-Plugin为自定义插件层，框架层将第三方服务抽象为多种自定义插件接口，插件层根据实际的基础设施来进行实现。GameLogic为游戏逻辑层，用于承载实际的业务逻辑。框架使用protobuf来生成脚手架，通过引入自定义Option等方式将业务逻辑自动注入到框架层中。

优势：

1，简单但不简陋。框架包含了一个重度游戏后端的完备功能，囊括了大部分的业务场景。

2、自动化的逻辑注入。借助pb的自定义选项，业务逻辑只需要很少的代码，就可以自动的注入到框架层

3、拒绝定制工具。框架的代码生成和逻辑注入只依赖原生的protobuf和grpc，不需要额外安装任何定制工具

![未命名文件 (2)](https://user-images.githubusercontent.com/16680818/222721483-8f14f7f2-7bb9-4eb2-8688-1367a67ed2ac.png)

Mutil,Async,Actor-Server: 游戏后台常见的业务模式，分别对应并发，单线程异步，Actor模式

App: 用于承载业务逻辑的服务实体，是服务网格中的独立个体，由“环境ID.区服ID.功能ID.实例ID”唯一标识。一个App可以承载一个或多个Server

Client：RPC客户端，由服务提供方编写，包含一些参数校验，和路由规则

ResourceLoader: 可线上热更新的资源加载器，一般用于GameLogic Config的管理

Runtime: 驱动器

Transport: 提供服务之间问答式RPC调用能力，采用tcp/GrpcSer实现一对一传输

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
Step3、定义服务客户端Client
```
package dirclient

var (
	_client *Client
)

// 初始化
func Init(srcins *pb.RouteIns, opts ...grpc.CallOption) error {
	_client = &Client{
		Meta: client.NewMeta(
			dpb.File_dir_proto.Services().Get(0),  // 逻辑注入
			srcins,
			&pb.RouteSvc{
				EnvID:   srcins.EnvID,
				WorldID: srcins.WorldID,
				FuncID:  srcins.FuncID,
				IDStr:   util.SvcID2Str(srcins.EnvID, srcins.WorldID, dproto.DirAppFuncID),
			},
		),
		opts: opts,
	}
	var err error
	if _client.conn, err = grpc.Dial(
		fmt.Sprintf("router://%s", util.SvcID2Str(srcins.EnvID, srcins.WorldID, dproto.DirAppFuncID)),  // 自定义路由
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(client.SendInteror()),
	); err != nil {
		return err
	} else {
		_client.DirSerClient = dpb.NewDirSerClient(_client.conn)
		return nil
	}
}

type Client struct {
	client.Meta
	dpb.DirSerClient
	conn *grpc.ClientConn
	opts []grpc.CallOption
}

// 客户端参数检查
func Hello(ctx context.Context, hellostr string, opts ...grpc.CallOption) (string, error) {
	var in dpb.Hello_Req
	in.Req = hellostr
	out, err := client.RPCCallRandom(_client, _client.Hello, ctx, &in, append(opts, _client.opts...)...)
	if err != nil {
		return "", err
	} else {
		return out.Rsp, nil
	}
}
```


Step4、定义服务实体DirApp
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


## 消息驱动

## 三种业务模式

## 服务发现与路由

## RPC与广播

## 资源管理

## 数据存储

