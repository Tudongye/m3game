# m3game

A game framework using GO language and Grpc

# m3game

一个基于Golang和Grpc的游戏后端框架。

框架分为GameLogic, Frame-Runtime, Custom-Plugin三层。Frame-Runtime为框架驱动层，负责消息分发，服务网格，插件管理等核心驱动工作。Custom-Plugin为自定义插件层，框架层将第三方服务抽象为多种自定义插件接口，插件层根据实际的基础设施来进行实现。GameLogic为游戏逻辑层，用于承载实际的业务逻辑。框架使用protobuf来生成脚手架，通过引入自定义Option等方式将业务逻辑自动注入到框架层中。

优势：

1，简单但不简陋。框架包含了一个重度游戏后端的完备功能，支持绝大部分的业务场景

2、自动化的逻辑注入。借助pb的自定义选项，业务逻辑只需要很少的代码，就可以自动的注入到框架层

3、拒绝定制化工具。框架的代码生成和逻辑注入只依赖原生的protobuf和grpc，不需要额外安装任何定制化工具

4、

使用标准的protobuf和grpc来生成脚手架代码，业务开发只需要关注于具体的游戏逻辑。框架支持三种常见的服务模型，并发多线程，单线程异步 和 Actor模式（actor内同步）。各服务捅过grpc进行通讯，服务网格支持P2P，Random，Hash，Single路由模式。框架支持插件式开发，可以根据具体使用的基础设施进行定制化开发。框架基于Broker实现的Grpc兼容BrokerSer，用于处理BroadCast和MutilCast。可以直接复用Grpc接口代码

Runtime: 框架核心，包括了启动器，插件管理器 和 Grpcser。全部采用单例。

App: 程序实体，框架中以App作为路由实体，每个服务包含一个App和若干个Server。

Server: 业务实体，包含实际的游戏逻辑。这里有三种服务基类，Mutil，Async 和 Actor，可以根据具体的业务需求进行继承使用

Mesh: 服务网格，包含专用的load_Balance算法 和 一个基于Consul的路由插件

DB: DB插件，包含一个测试用的内存数据库

Broker： Broker插件，包含一个Nats实现

Resource： 资源管理（业务配置管理）

Config: 配置管理

Demo： 一组Demo，包括了完整的业务功能。


![未命名文件 (2)](https://user-images.githubusercontent.com/16680818/222721483-8f14f7f2-7bb9-4eb2-8688-1367a67ed2ac.png)

TODO:

1、Route Support BroadCast, MutilCast  Done

路由策略支持 广播，多播

2、Metric, Trace, Distribute Lock Plugin

增加对监控，追踪，分布式锁的插件支持

3、Actor automatic cross-program migration

Actor模式下，Actor自动跨进程迁移

