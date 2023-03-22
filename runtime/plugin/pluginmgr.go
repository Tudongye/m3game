package plugin

import (
	"fmt"
	"m3game/plugins/log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func init() {
	_pluginMgr = &PluginMgr{
		insMap: make(map[Type]map[string]PluginIns),
	}
}

type PluginMgr struct {
	insMap map[Type]map[string]PluginIns // type name tag
}

var (
	_pluginMgr *PluginMgr
)

func Reload(v viper.Viper) error {
	var cfg config
	if err := v.Unmarshal(&cfg); err != nil {
		return errors.Wrap(err, "Unmarshal PluginCfg")
	}
	for _, typ := range _pluginserial {
		for name, p := range _pluginMgr.insMap[typ] {
			if err := p.Factory().Reload(p, cfg.Plugin[string(typ)][name]); err != nil {
				return err
			}
		}
	}
	return nil
}

func Destroy() error {
	for i := len(_pluginserial) - 1; i >= 0; i-- {
		typ := _pluginserial[i]
		for name, p := range _pluginMgr.insMap[typ] {
			log.Info("Destory  %s.%s", typ, name)
			if err := p.Factory().Destroy(p); err != nil {
				log.Error("Destory type %s name %s Fail %s", typ, name, err.Error())
			}
		}
	}
	return nil
}

func registerPluginIns(typ Type, name string, p PluginIns) error {
	if _, ok := _pluginMgr.insMap[typ]; !ok {
		_pluginMgr.insMap[typ] = make(map[string]PluginIns)
	}
	if _, ok := _pluginMgr.insMap[typ][name]; !ok {
		_pluginMgr.insMap[typ][name] = p
	} else {
		return fmt.Errorf("Repeated plugin type %s name %s", typ, name)
	}
	return nil
}

func getPluginByType(typ Type) []PluginIns {
	var inss []PluginIns
	if tm, ok := _pluginMgr.insMap[typ]; !ok {
		return nil
	} else {
		for _, p := range tm {
			inss = append(inss, p)
		}
	}
	return inss
}

func getPluginByName(typ Type, name string) PluginIns {
	if tm, ok := _pluginMgr.insMap[typ]; !ok {
		return nil
	} else {
		if p, ok := tm[name]; ok {
			return p
		}
		return nil
	}
}
