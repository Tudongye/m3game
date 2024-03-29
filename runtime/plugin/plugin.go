package plugin

import (
	"context"
	"m3game/meta/errs"
	"m3game/plugins/log"

	"github.com/spf13/viper"
)

type Type string

const (
	DB     Type = "db"     // 存储
	Router Type = "router" // 服务发现
	Trace  Type = "trace"  // 链路追踪
	Metric Type = "metric" // 监控
	Broker Type = "broker" // 消息队列
	Log    Type = "log"    // 日志
	Shape  Type = "shape"  // 流量管理
	Gate   Type = "gate"   // CS连接
	Lease  Type = "lease"  // 租约
	Lua    Type = "lua"    // Lua
	Trans  Type = "trans"  // 内部通信
)

var (
	_pluginserial = []Type{Log, Trans, Broker, Router, Trace, Metric, DB, Lease, Shape, Gate, Lua} // Plugin加载顺序
	_factoryMap   = make(map[string]Factory)
)

func RegisterFactory(f Factory) {
	if _, ok := _factoryMap[f.Name()]; ok {
		log.Fatal("RegisterFactory factory name repeatad %s", f.Name())
	}
	_factoryMap[f.Name()] = f
}

type PluginIns interface {
	Factory() Factory
}

type Factory interface {
	Type() Type
	Name() string
	Setup(context.Context, map[string]interface{}) (PluginIns, error)
	Destroy(PluginIns) error
	Reload(PluginIns, map[string]interface{}) error
	CanUnload(PluginIns) bool
}

type config struct {
	Plugin map[string]map[string]map[string]interface{} `toml:"Plugin"`
}

func InitPlugins(ctx context.Context, v viper.Viper) error {
	var cfg config
	if err := v.Unmarshal(&cfg); err != nil {
		return errs.PluginsInitFail.Wrap(err, "Unmarshal PluginCfg")
	}
	for _, typ := range _pluginserial {
		for name, nm := range cfg.Plugin[string(typ)] {
			factory, ok := _factoryMap[name]
			if !ok {
				return errs.PluginsFactoryNotRegister.New("Factory not find %s", name)
			}
			log.Info("Plugin Setup %s", name)
			pluginIns, err := factory.Setup(ctx, nm)
			if err != nil {
				return errs.PluginsSetupFail.Wrap(err, "Factory %s", name)
			}
			if err := addPluginIns(factory.Type(), name, pluginIns); err != nil {
				return errs.PluginsAddPluginFail.Wrap(err, "Factory %s", name)
			}
		}
	}
	return nil
}
