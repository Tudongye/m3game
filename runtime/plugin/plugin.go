package plugin

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Type string

const (
	DB     Type = "db"
	Router Type = "router" // 只可以有一个
	Trace  Type = "trace"
	Metric Type = "metric"
	Broker Type = "broker"
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

type PluginIns interface {
	Name() string
}

type Factory interface {
	Type() Type
	Name() string
	Setup(map[string]interface{}) (PluginIns, error)
	Destroy(PluginIns) error
	Reload(PluginIns, map[string]interface{}) error
	CanDelete(PluginIns) bool
}

func RegisterPluginFactory(f Factory) {
	if _, ok := _factoryMap[f.Name()]; ok {
		log.Panicf("RegisterPluginFactory factory name repeatad %s", f.Name())
	}
	_factoryMap[f.Name()] = f
}

type config struct {
	Plugin map[string]map[string]map[string]interface{} `toml:"Plugin"`
}

func InitPlugins(v viper.Viper) error {
	var cfg config
	if err := v.Unmarshal(&cfg); err != nil {
		return err
	}
	for _, tm := range cfg.Plugin {
		for name, nm := range tm {
			if factory, ok := _factoryMap[name]; !ok {
				return fmt.Errorf("Factory not find %s", name)
			} else {
				log.Printf("Plugin Setup %s\n", name)
				if p, err := factory.Setup(nm); err != nil {
					return err
				} else {
					if err := registerPluginIns(factory.Type(), name, getPluginTag(nm), p); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func getPluginTag(m map[string]interface{}) string {
	if v, ok := m["tag"]; !ok {
		return _defaulttag
	} else {
		return v.(string)
	}
}
