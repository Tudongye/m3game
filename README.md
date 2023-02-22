# m3game

A game framework using GO language and Grpc

With the code generation capabilities of protobuf and grpc, developers only need to focus on specific game logic.The framework supports three common game service modes: multi-thread, single-thread asynchronous and Actor model(sync in actor).The services of the framework communicate with each other through grpc, and provide multiple routing (P2P, Random, Hash, Single) policies. Plug-in development interface is provided, which can be customized according to the infrastructure.

Runtime: The core of the framework, including the launcher, plug-in manager and Grpc service

App: program entity. Each service contains one app and several servers

Server: A business entity that carries specific game logic. There are three types of servers, Mutil, Async and Actor, which can be inherited and used according to specific business needs.

Mesh: service mesh, including load_ Balance and a Router plug-in based on Consul.

DB: DB plug-in, including a cache database for testing.

Resource: Business layer configuration

Config: configuration management

Demo: A group of demos includes complete business functions

# m3game

一个基于Golang和Grpc的游戏后端框架。

使用标准的protobuf和grpc来生成脚手架代码，业务开发只需要关注于具体的游戏逻辑。框架支持三种常见的服务模型，并发多线程，单线程异步 和 Actor模式（actor内同步）。各服务捅过grpc进行通讯，服务网格支持P2P，Random，Hash，Single路由模式。框架支持插件式开发，可以根据具体使用的基础设施进行定制化开发。

Runtime: 框架核心，包括了启动器，插件管理器 和 Grpcser。全部采用单例。

App: 程序实体，框架中以App作为路由实体，每个服务包含一个App和若干个Server。

Server: 业务实体，包含实际的游戏逻辑。这里有三种服务基类，Mutil，Async 和 Actor，可以根据具体的业务需求进行继承使用

Mesh: 服务网格，包含专用的load_Balance算法 和 一个基于Consul的路由插件

DB: DB插件，包含一个测试用的内存数据库

Resource： 资源管理（业务配置管理）

Config: 配置管理

Demo： 一组Demo，包括了完整的业务功能。


![未命名文件 (1)](https://user-images.githubusercontent.com/16680818/220582821-4def39e5-550f-4784-bc40-49779038c71e.png)


TODO:

1、Route Support BroadCast, MutilCast

路由策略支持 广播，多播

2、Metric, Trace, Distribute Lock Plugin

增加对监控，追踪，分布式锁的插件支持

3、Actor automatic cross-program migration

Actor模式下，Actor自动跨进程迁移

