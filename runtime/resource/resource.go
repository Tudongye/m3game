package resource

import (
	"context"
	"m3game/meta/errs"
	"m3game/plugins/log"
	"sync"
	"sync/atomic"

	"github.com/mitchellh/mapstructure"
)

var (
	_loadercreaters []ResLoaderBuilder
	_mgr            *ResourceMgr
)

// 资源加载器
type ResLoader interface {
	Load(ctx context.Context, cfgpath string) error // 资源更新
	Name() string
}

type ResLoaderBuilder interface {
	NewLoader() ResLoader
	Name() string
}

func RegisterResource(creater ResLoaderBuilder) {
	_loadercreaters = append(_loadercreaters, creater)
}

type ResourceCfg struct {
	CfgPath string
}

func GetLoader[T ResLoader](ctx context.Context) T {
	var t T
	name := t.Name()
	loader := _mgr.getLoader(ctx, name)
	if loader == nil {
		log.Fatal("ResLoader %s not find", name)
	}
	if l, ok := loader.(T); !ok {
		log.Fatal("ResLoader %s type not match", name)
	} else {
		return l
	}
	return t
}

func New(c map[string]interface{}) (*ResourceMgr, error) {
	if _mgr != nil {
		return _mgr, nil
	}
	_mgr = &ResourceMgr{
		vaildidx: 0,
	}
	if err := mapstructure.Decode(c, &_mgr.cfg); err != nil {
		return nil, err
	}
	return _mgr, nil
}

type ResourceMgr struct {
	cfg        ResourceCfg
	resloaders [2]*sync.Map
	vaildidx   int32
	lock       sync.RWMutex
}

func (rm *ResourceMgr) reLoad() error {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	ctx := rm.newLoadingCtx()
	var newloaders sync.Map
	newidx := (rm.vaildidx + 1) % 2
	rm.resloaders[newidx] = &newloaders
	for _, creater := range _loadercreaters {
		loader := creater.NewLoader()
		if _, ok := newloaders.Load(loader.Name()); ok {
			log.Fatal("ResLoader %s repeated", loader.Name())
		}
		if loader.Name() != creater.Name() {
			log.Fatal("ResLoader %s and Creater Name %s not match", loader.Name(), creater.Name())
		}
		if err := loader.Load(ctx, rm.cfg.CfgPath); err != nil {
			return errs.ResourceLoadFail.New("Load %s err %s", loader.Name(), err.Error())
		}
		newloaders.Store(loader.Name(), loader)
	}
	atomic.StoreInt32(&rm.vaildidx, int32(newidx))
	return nil
}

func (rm *ResourceMgr) getLoader(ctx context.Context, name string) ResLoader {
	isloading := rm.isLoadingCtx(ctx)
	if isloading {
		vaildidx := atomic.LoadInt32(&rm.vaildidx)
		if v, ok := rm.resloaders[(vaildidx+1)%2].Load(name); !ok {
			return nil
		} else {
			return v.(ResLoader)
		}
	} else {
		vaildidx := atomic.LoadInt32(&rm.vaildidx)
		if v, ok := rm.resloaders[vaildidx].Load(name); !ok {
			return nil
		} else {
			return v.(ResLoader)
		}
	}
}

func (rm *ResourceMgr) ReLoad(c map[string]interface{}) error {
	if err := mapstructure.Decode(c, &rm.cfg); err != nil {
		return err
	}
	return rm.reLoad()
}

func (rm *ResourceMgr) newLoadingCtx() context.Context {
	return context.WithValue(context.Background(), rm.loadFlag(), true)
}

func (rm *ResourceMgr) isLoadingCtx(ctx context.Context) bool {
	if v := ctx.Value(rm.loadFlag()); v == nil {
		return false
	} else if bv, ok := v.(bool); !ok {
		return false
	} else {
		return bv
	}
}

func (rm *ResourceMgr) loadFlag() string {
	return "_loaderflag"
}
