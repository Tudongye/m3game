package onlineser

import (
	"context"
	"errors"
	"fmt"
	"m3game/config"
	"m3game/demo/proto"
	"m3game/demo/proto/pb"
	"m3game/meta"
	"m3game/meta/errs"
	"m3game/plugins/db"
	"m3game/plugins/log"
	"m3game/plugins/router"
	"sync"
	"time"

	"github.com/bluele/gcache"
)

var (
	_onlineroledbmeta     *db.DBMeta[*pb.OnlineRoleDB]
	_onlinerolewrapermeta *db.WraperMeta[*pb.OnlineRoleDB, pb.ORFlag]
	_onlinepool           *OnlinePool
)

func init() {
	_onlineroledbmeta = db.NewMeta[*pb.OnlineRoleDB]("onlinerole_table")
	_onlinerolewrapermeta = db.NewWraperMeta[*pb.OnlineRoleDB, pb.ORFlag](_onlineroledbmeta)
}

func newPool() *OnlinePool {
	if _onlinepool != nil {
		return _onlinepool
	}
	_onlinepool = &OnlinePool{
		isopen: false,
	}
	return _onlinepool
}

func newCache() gcache.Cache {
	return gcache.New(_cfg.CachePoolSize).LRU().
		LoaderFunc(func(key interface{}) (interface{}, error) {
			dbp := db.Instance()
			roleid := key.(int64)
			w := _onlinerolewrapermeta.New(roleid)
			if err := w.Read(context.TODO(), dbp); err == nil {
				if app := w.Obj().GetOnlineApp(); app != nil {
					return app, nil
				} else {
					// 异常
					return nil, errors.New("OnlinePool is Err, app is nil")
				}
			} else {
				return nil, err
			}
		}).Build()
}

func Pool() *OnlinePool {
	return _onlinepool
}

type OnlinePool struct {
	isopen bool
	mu     sync.Mutex

	rolecache gcache.Cache
	appcache  sync.Map
}

type AppCache struct {
	Ver            string
	LastUpdateTime int64
}

func (u *OnlinePool) Close() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.isopen = false
}

func (u *OnlinePool) Open() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.isopen = true
	u.rolecache = newCache()
}

func (u *OnlinePool) IsOpen() bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.isopen
}

func (u *OnlinePool) OnlineCreate(ctx context.Context, roleid int64, appid string) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if !u.isopen {
		return errors.New("OnlinePool is Close")
	}
	// 查缓存
	created := false
	if v, err := u.rolecache.Get(roleid); err == nil {
		roleapp := v.(*pb.OnlineApp)
		if roleapp.AppId != appid {
			if v, ok := u.appcache.Load(roleapp.AppId); ok {
				appcache := v.(*AppCache)
				if appcache.Ver == roleapp.Ver && appcache.LastUpdateTime+int64(_cfg.AppAliveTimeOut) > time.Now().Unix() {
					return fmt.Errorf("RoleId %d have online in %s:%s", roleid, roleapp.AppId, roleapp.Ver)
				}
			}
		}
		created = true
	} else if !errs.DBKeyNotFound.Is(err) {
		return err
	}
	// 查App
	var appcache *AppCache
	if v, ok := u.appcache.Load(appid); ok {
		appcache = v.(*AppCache)
		if appcache.LastUpdateTime+int64(_cfg.AppAliveTimeOut) < time.Now().Unix() {
			return fmt.Errorf("App %s not alive", appid)

		}
	} else {
		return fmt.Errorf("App %s not alive", appid)
	}
	// 写入DB
	dbp := db.Instance()
	w := _onlinerolewrapermeta.New(roleid)
	onlineapp := &pb.OnlineApp{AppId: appid, Ver: appcache.Ver}
	w.Set(pb.ORFlag_OROnlineApp, onlineapp)
	if created {
		if err := w.Update(ctx, dbp); err != nil {
			log.Error("%s", err.Error())
			return err
		}
	} else {
		if err := w.Create(ctx, dbp); err != nil {
			log.Error("%s", err.Error())
			return err
		}
	}
	// 返回
	u.rolecache.Set(roleid, onlineapp)
	return nil
}

func (u *OnlinePool) OnlineRead(roleid int64) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if !u.isopen {
		return "", errors.New("OnlinePool is Close")
	}
	// 查缓存
	if v, err := u.rolecache.Get(roleid); err == nil {
		roleapp := v.(*pb.OnlineApp)
		if v, ok := u.appcache.Load(roleapp.AppId); ok {
			appcache := v.(*AppCache)
			if appcache.Ver == roleapp.Ver && appcache.LastUpdateTime+int64(_cfg.AppAliveTimeOut) > time.Now().Unix() {
				return roleapp.AppId, nil
			}
		}
		return "", nil
	} else if errs.DBKeyNotFound.Is(err) {
		return "", nil
	} else {
		return "", err
	}
}

func (u *OnlinePool) OnlineDelete(ctx context.Context, roleid int64, appid string) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if !u.isopen {
		return errors.New("OnlinePool is Close")
	}
	// 查缓存
	created := false
	if v, err := u.rolecache.Get(roleid); err == nil {
		roleapp := v.(*pb.OnlineApp)
		if v, ok := u.appcache.Load(roleapp.AppId); ok {
			appcache := v.(*AppCache)
			if appcache.Ver == roleapp.Ver && appcache.LastUpdateTime+int64(_cfg.AppAliveTimeOut) > time.Now().Unix() {
				return fmt.Errorf("RoleId %d have online in %s:%s", roleid, roleapp.AppId, roleapp.Ver)
			}
		}
		created = true
	} else if !errs.DBKeyNotFound.Is(err) {
		return err
	}
	u.rolecache.Remove(roleid)
	if created {
		dbp := db.Instance()
		w := _onlinerolewrapermeta.New(roleid)
		if err := w.Delete(ctx, dbp); err != nil {
			return err
		}
	}
	return nil
}

func (u *OnlinePool) LoadAppCache() error {
	// 加载App在线信息
	env, world, _ := config.GetWorldID().Parse()
	rolesvcid := meta.GenRouteSvc(env, world, proto.RoleFuncID)
	timenow := time.Now().Unix()
	if inns, err := router.GetAllInstances(rolesvcid.String()); err != nil {
		return err
	} else {
		for _, ins := range inns {
			ver, ok := ins.GetMeta(string(meta.M3AppVer))
			if !ok {
				continue
			}
			u.appcache.Store(ins.GetIDStr(), &AppCache{Ver: ver, LastUpdateTime: timenow})
		}
	}
	return nil
}
