package plugin

import (
	"fmt"
	"m3game/db"
	"m3game/mesh/router"
)

func init() {
	_pluginMgr = &PluginMgr{
		insMap: make(map[Type]map[string]map[string]PluginIns),
	}
}

type PluginMgr struct {
	insMap map[Type]map[string]map[string]PluginIns // type name tag
}

var (
	_pluginMgr *PluginMgr
)

func registerPluginIns(typ Type, name string, tag string, p PluginIns) error {
	if _, ok := _pluginMgr.insMap[typ]; !ok {
		_pluginMgr.insMap[typ] = make(map[string]map[string]PluginIns)
	} else if typ == Router {
		return fmt.Errorf("Plugin Router only one")
	}
	if _, ok := _pluginMgr.insMap[typ][name]; !ok {
		_pluginMgr.insMap[typ][name] = make(map[string]PluginIns)
	}
	if _, ok := _pluginMgr.insMap[typ][name][tag]; ok {
		return fmt.Errorf("Plugin repeated type %s name %s tag %s", typ, name, tag)
	}
	_pluginMgr.insMap[typ][name][tag] = p
	return nil
}

func getPluginByType(typ Type) PluginIns {
	if tm, ok := _pluginMgr.insMap[typ]; !ok {
		return nil
	} else {
		for _, nm := range tm {
			for _, p := range nm {
				return p
			}
		}
	}
	return nil
}

func getPluginByName(typ Type, name string) PluginIns {
	if tm, ok := _pluginMgr.insMap[typ]; !ok {
		return nil
	} else {
		if nm, ok := tm[name]; !ok {
			return nil
		} else {
			for _, p := range nm {
				return p
			}
		}
	}
	return nil
}

func getPluginByTag(typ Type, name string, tag string) PluginIns {
	if tm, ok := _pluginMgr.insMap[typ]; !ok {
		return nil
	} else {
		if nm, ok := tm[name]; !ok {
			return nil
		} else {
			if p, ok := nm[tag]; !ok {
				return nil
			} else {
				return p
			}
		}
	}
}

func GetRouterPlugin() router.Router {
	p := getPluginByType(Router)
	if p == nil {
		return nil
	}
	return p.(router.Router)
}

func GetDBPlugin() db.DB {
	p := getPluginByType(DB)
	if p == nil {
		return nil
	}
	return p.(db.DB)
}

func GetDBPluginByName(name string) db.DB {
	p := getPluginByName(DB, name)
	if p == nil {
		return nil
	}
	return p.(db.DB)
}
