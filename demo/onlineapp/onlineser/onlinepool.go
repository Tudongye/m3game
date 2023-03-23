package onlineser

import (
	"errors"
	"fmt"
	"m3game/config"
	"m3game/demo/proto"
	"m3game/demo/proto/pb"
	"m3game/meta"
	"m3game/plugins/db"
	"m3game/plugins/db/wraper"
	"m3game/plugins/router"
	"sync"
	"time"

	"github.com/bluele/gcache"
)

var (
	_onlineroledbmeta *db.DBMeta[*pb.OnlineRoleDB]
	_onlinepool       *OnlinePool
)

func init() {
	_onlineroledbmeta = db.NewMeta("onlinerole_table", onlineroleCreater)
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
			openid := key.(string)
			w := wraper.New(_onlineroledbmeta, openid)
			if err := w.Read(dbp); err == nil {
				if app := w.TObj().GetApp(); app != nil {
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

func onlineroleCreater() *pb.OnlineRoleDB {
	return &pb.OnlineRoleDB{
		RoleId: "",
		App:    &pb.OnlineRoleApp{},
	}
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
	return u.isopen
}

func (u *OnlinePool) OnlineCreate(roleid string, appid string) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if !u.IsOpen() {
		return errors.New("OnlinePool is Close")
	}
	// 查缓存
	created := false
	if v, err := u.rolecache.Get(roleid); err == nil {
		roleapp := v.(*pb.OnlineRoleApp)
		if roleapp.AppId != appid {
			if v, ok := u.appcache.Load(roleapp.AppId); ok {
				appcache := v.(*AppCache)
				if appcache.Ver == roleapp.Ver && appcache.LastUpdateTime+int64(_cfg.AppAliveTimeOut) > time.Now().Unix() {
					return fmt.Errorf("RoleId %s have online in %s:%s", roleid, roleapp.AppId, roleapp.Ver)
				}
			}
		}
		created = true
	} else if !db.IsErrKeyNotFound(err) {
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
	w := wraper.New(_onlineroledbmeta, roleid)
	roleapp := &pb.OnlineRoleApp{AppId: appid, Ver: appcache.Ver}
	if err := wraper.Setter(w, roleapp); err != nil {
		return err
	}
	if created {
		if err := w.Update(dbp); err != nil {
			return err
		}
	} else {
		if err := w.Create(dbp); err != nil {
			return err
		}
	}
	// 返回
	u.rolecache.Set(roleid, roleapp)
	return nil
}

func (u *OnlinePool) OnlineRead(roleid string) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if !u.IsOpen() {
		return "", errors.New("OnlinePool is Close")
	}
	// 查缓存
	if v, err := u.rolecache.Get(roleid); err == nil {
		roleapp := v.(*pb.OnlineRoleApp)
		if v, ok := u.appcache.Load(roleapp.AppId); ok {
			appcache := v.(*AppCache)
			if appcache.Ver == roleapp.Ver && appcache.LastUpdateTime+int64(_cfg.AppAliveTimeOut) > time.Now().Unix() {
				return roleapp.AppId, nil
			}
		}
		return "", nil
	} else if db.IsErrKeyNotFound(err) {
		return "", nil
	} else {
		return "", err
	}
}

func (u *OnlinePool) OnlineDelete(roleid string, appid string) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if !u.IsOpen() {
		return errors.New("OnlinePool is Close")
	}
	// 查缓存
	created := false
	if v, err := u.rolecache.Get(roleid); err == nil {
		roleapp := v.(*pb.OnlineRoleApp)
		if v, ok := u.appcache.Load(roleapp.AppId); ok {
			appcache := v.(*AppCache)
			if appcache.Ver == roleapp.Ver && appcache.LastUpdateTime+int64(_cfg.AppAliveTimeOut) > time.Now().Unix() {
				return fmt.Errorf("RoleId %s have online in %s:%s", roleid, roleapp.AppId, roleapp.Ver)
			}
		}
		created = true
	} else if !db.IsErrKeyNotFound(err) {
		return err
	}
	u.rolecache.Remove(roleid)
	if created {
		dbp := db.Instance()
		w := wraper.New(_onlineroledbmeta, roleid)
		if err := w.Delete(dbp); err != nil {
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
