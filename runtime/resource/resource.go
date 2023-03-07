package resource

import (
	"context"
	"fmt"
	"sync"

	"github.com/mitchellh/mapstructure"
)

var (
	_loadercreaters []ResLoaderCreater
	_cfg            ResourceCfg
	_mgr            *resourceMgr
	_loaderflag     = "_loaderflag"
)

type ResLoader interface {
	Load(ctx context.Context, cfgpath string) error // 资源更新
	Name() string
}

type ResLoaderCreater interface {
	NewLoader() ResLoader
	Name() string
}

type ResourceCfg struct {
	CfgPath string
}

func Init(c map[string]interface{}) error {
	if _mgr != nil {
		return nil
	}
	_mgr = &resourceMgr{
		vaildidx: 0,
	}
	return ReLoad(c)
}

func ReLoad(c map[string]interface{}) error {
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return err
	}
	return _mgr.reLoad()
}

func RegisterResource(creater ResLoaderCreater) {
	_loadercreaters = append(_loadercreaters, creater)
}

func GetLoader[T ResLoader](ctx context.Context) T {
	var t T
	name := t.Name()
	loader := _mgr.getLoader(ctx, name)
	if loader == nil {
		panic(fmt.Sprintf("ResLoader %s not find", name))
	}
	if l, ok := loader.(T); !ok {
		panic(fmt.Sprintf("ResLoader %s type not match", name))
	} else {
		return l
	}
}

func newLoadingCtx() context.Context {
	return context.WithValue(context.Background(), _loaderflag, true)
}

func isLoadingCtx(ctx context.Context) bool {
	if v := ctx.Value(_loaderflag); v == nil {
		return false
	} else if bv, ok := v.(bool); !ok {
		return false
	} else {
		return bv
	}
}

type resourceMgr struct {
	resloaders [2]map[string]ResLoader
	vaildidx   int
	lock       sync.RWMutex
}

func (rm *resourceMgr) reLoad() error {
	ctx := newLoadingCtx()
	newloaders := make(map[string]ResLoader)
	for _, creater := range _loadercreaters {
		loader := creater.NewLoader()
		if _, ok := newloaders[loader.Name()]; ok {
			panic(fmt.Sprintf("ResLoader %s repeated", loader.Name()))
		}
		if loader.Name() != creater.Name() {
			panic(fmt.Sprintf("ResLoader %s and Creater Name %s not match", loader.Name(), creater.Name()))
		}
		if err := loader.Load(ctx, _cfg.CfgPath); err != nil {
			return fmt.Errorf("Load %s err %s", loader.Name(), err.Error())
		}
		newloaders[loader.Name()] = loader
	}
	rm.lock.Lock()
	defer rm.lock.Unlock()
	rm.vaildidx = (rm.vaildidx + 1) % 2
	rm.resloaders[rm.vaildidx] = newloaders
	return nil
}

func (rm *resourceMgr) getLoader(ctx context.Context, name string) ResLoader {
	isloading := isLoadingCtx(ctx)
	rm.lock.RLock()
	defer rm.lock.RUnlock()
	if isloading {
		return rm.resloaders[(rm.vaildidx+1)%2][name]
	} else {
		return rm.resloaders[rm.vaildidx][name]
	}
}
