package resource

import (
	"fmt"
	"sync"

	"github.com/mitchellh/mapstructure"
)

var (
	_resloadercreaters = make(map[string]ResLoaderCreater)
	_mgr               *ResLoaderMgr
	_lock              sync.RWMutex
	_cfg               ResourceCfg
)

type ResourceCfg struct {
	CfgPath string
}

func Init(c map[string]interface{}) error {
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return err
	}
	return Load()
}

func Load() error {
	newmgr := &ResLoaderMgr{
		resloaders: make(map[string]ResLoader),
		cfgpath:    _cfg.CfgPath,
	}
	if err := newmgr.reload(); err != nil {
		return err
	}

	_lock.Lock()
	defer _lock.Unlock()
	_mgr = newmgr
	return nil
}

func GetResource(k string) ResLoader {
	_lock.RLock()
	defer _lock.RUnlock()
	return _mgr.getResource(k)
}

func RegisterResLoader(name string, f ResLoaderCreater) {
	if _, ok := _resloadercreaters[name]; ok {
		panic(fmt.Sprintf("Resource Name %s repeated", name))
	}
	_resloadercreaters[name] = f
}

type ResLoaderCreater func() ResLoader
type ResLoaderGetter func(string) ResLoader

type ResLoader interface {
	Load(cfgpath string) error
	Check(ResLoaderGetter) error
}

type ResLoaderMgr struct {
	resloaders map[string]ResLoader
	cfgpath    string
}

func (r *ResLoaderMgr) reload() error {
	r.resloaders = make(map[string]ResLoader)
	for name, creater := range _resloadercreaters {
		loader := creater()
		if err := loader.Load(r.cfgpath); err != nil {
			return fmt.Errorf("Load %s err %s", name, err.Error())
		}
		r.resloaders[name] = loader
	}
	for name, loader := range r.resloaders {
		if err := loader.Check(r.getResource); err != nil {
			return fmt.Errorf("Check %s err %s", name, err.Error())
		}
	}
	return nil
}

func (r *ResLoaderMgr) getResource(k string) ResLoader {
	return _mgr.resloaders[k]
}
