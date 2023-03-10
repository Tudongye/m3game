package plugin

import (
	"fmt"
	"m3game/plugins/log"

	"github.com/pkg/errors"
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
	Shape  Type = "Shape"  // 流量管理
	Gate   Type = "Gate"   // CS连接
)

const (
	_defaulttag = "default"
)

var (
	_factoryMap map[string]Factory
)

func init() {
	_factoryMap = make(map[string]Factory)
}

func RegisterFactory(f Factory) {
	if _, ok := _factoryMap[f.Name()]; ok {
		panic(fmt.Sprintf("RegisterFactory factory name repeatad %s", f.Name()))
	}
	_factoryMap[f.Name()] = f
}

type PluginIns interface {
	Factory() Factory
}

type Factory interface {
	Type() Type
	Name() string
	Setup(map[string]interface{}) (PluginIns, error)
	Destroy(PluginIns) error
	Reload(PluginIns, map[string]interface{}) error
	CanDelete(PluginIns) bool
}

type config struct {
	Plugin map[string]map[string]map[string]interface{} `toml:"Plugin"`
}

func InitPlugins(v viper.Viper) error {
	var cfg config
	if err := v.Unmarshal(&cfg); err != nil {
		return errors.Wrap(err, "Unmarshal PluginCfg")
	}
	for _, tm := range cfg.Plugin {
		for name, nm := range tm {
			factory, ok := _factoryMap[name]
			if !ok {
				return fmt.Errorf("Factory not find %s", name)
			}
			log.Info("Plugin Setup %s", name)
			pluginIns, err := factory.Setup(nm)
			if err != nil {
				return errors.Wrapf(err, "Factory %s", name)
			}
			if err := registerPluginIns(factory.Type(), name, getPluginTag(nm), pluginIns); err != nil {
				return errors.Wrapf(err, "Factory %s", name)
			}
		}
	}
	return nil
}

func getPluginTag(m map[string]interface{}) string {
	if v, ok := m["tag"]; !ok {
		return _defaulttag
	} else if tag, ok := v.(string); !ok {
		return _defaulttag
	} else {
		return tag
	}
}

func Gett[T PluginIns]() T {
	var t T
	v := getPluginByType(t.Factory().Type())
	return v.(T)
}
