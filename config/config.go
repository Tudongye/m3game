package config

import (
	"flag"
	"fmt"
	"m3game/util/log"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	_rconf   *m3Config
	_lock    sync.RWMutex
	_kvCache = make(map[string]interface{})
)

var (
	_defaultCfgPath = "m3game.toml"
)

type m3Config struct {
	config    *viper.Viper
	configBak *viper.Viper
	isUseBak  bool
	cfgPath   string
	idstr     string
	envmap    mapFlag
}

type mapFlag map[string]string

func (f mapFlag) String() string {
	return fmt.Sprintf("%v", map[string]string(f))
}

func (f mapFlag) Set(value string) error {
	split := strings.SplitN(value, "=", 2)
	f[split[0]] = split[1]
	return nil
}

func (c *m3Config) getBakConfig() *viper.Viper {
	if c.isUseBak {
		return c.config
	}
	return c.configBak
}

func (c *m3Config) loadconfig() error {
	v := c.getBakConfig()
	v.SetConfigFile(c.cfgPath)
	v.SetConfigType("toml")
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	c.isUseBak = !c.isUseBak
	clearCache()
	return nil
}

func GetConfig() *viper.Viper {
	if _rconf.isUseBak {
		return _rconf.configBak
	}
	return _rconf.config
}

func GetIDStr() string {
	return _rconf.idstr
}

func GetEnv(key string) string {
	return _rconf.envmap[key]
}

func init() {
	_rconf = &m3Config{
		config:    viper.New(),
		configBak: viper.New(),
		isUseBak:  false,
		envmap:    make(map[string]string),
	}
	flag.StringVar(&_rconf.cfgPath, "conf", _defaultCfgPath, "server config path")
	flag.StringVar(&_rconf.idstr, "idstr", "", "app idstr env.world.func.ins")
	flag.Var(&_rconf.envmap, "env", "other config")
	flag.Parse()
	log.Fatal("CfgPath:%s", _rconf.cfgPath)
	if _rconf.cfgPath == "" {
		_rconf.cfgPath = _defaultCfgPath
	}
	if err := _rconf.loadconfig(); err != nil {
		panic(fmt.Sprintf("LoadConfig Fail %s", err.Error()))
	}
}

func Reload() error {
	if err := _rconf.loadconfig(); err != nil {
		return err
	}
	return nil
}

func GetString(k string) string {
	if v := getFromCache(k); v != nil {
		if realVal, ok := v.(string); ok {
			return realVal
		}
	}
	c := GetConfig()
	v := c.GetString(k)
	insertToCache(k, v)
	return v
}

func GetInt(k string) int {
	if v := getFromCache(k); v != nil {
		if realVal, ok := v.(int); ok {
			return realVal
		}
	}
	c := GetConfig()
	v := c.GetInt(k)
	insertToCache(k, v)
	return v
}

func GetBool(k string) bool {
	if v := getFromCache(k); v != nil {
		if realVal, ok := v.(bool); ok {
			return realVal
		}
	}
	c := GetConfig()
	v := c.GetBool(k)
	insertToCache(k, v)
	return v
}

func getFromCache(k string) interface{} {
	_lock.RLock()
	defer _lock.RUnlock()
	if v, ok := _kvCache[k]; ok {
		return v
	}
	return nil
}

func clearCache() {
	_lock.Lock()
	defer _lock.Unlock()
	_kvCache = make(map[string]interface{})
}

func insertToCache(k string, v interface{}) {
	_lock.Lock()
	defer _lock.Unlock()
	_kvCache[k] = v
}
