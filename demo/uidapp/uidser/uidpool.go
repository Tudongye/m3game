package uidser

import (
	"context"
	"errors"
	"m3game/config"
	"m3game/demo/proto/pb"
	"m3game/meta/errs"
	"m3game/plugins/db"
	"m3game/plugins/log"
	"sync"

	"github.com/bluele/gcache"
)

var (
	_uidclubiddbmeta *db.DBMeta[*pb.UidClubIdDB]
	_uidroleiddbmeta *db.DBMeta[*pb.UidRoleIdDB]
	_uidmetadbmeta   *db.DBMeta[*pb.UidMetaDB]

	_uidclubiddbwrapermeta *db.WraperMeta[*pb.UidClubIdDB, pb.UCFlag]
	_uidroleiddbwrapermeta *db.WraperMeta[*pb.UidRoleIdDB, pb.URFlag]
	_uidmetadwraperbmeta   *db.WraperMeta[*pb.UidMetaDB, pb.UMFlag]

	_uidpool *UidPool
)

func init() {
	_uidclubiddbmeta = db.NewMeta[*pb.UidClubIdDB]("uidclubid_table")
	_uidroleiddbmeta = db.NewMeta[*pb.UidRoleIdDB]("uidroleid_table")
	_uidmetadbmeta = db.NewMeta[*pb.UidMetaDB]("uidmeta_table")

	_uidclubiddbwrapermeta = db.NewWraperMeta[*pb.UidClubIdDB, pb.UCFlag](_uidclubiddbmeta)
	_uidroleiddbwrapermeta = db.NewWraperMeta[*pb.UidRoleIdDB, pb.URFlag](_uidroleiddbmeta)
	_uidmetadwraperbmeta = db.NewWraperMeta[*pb.UidMetaDB, pb.UMFlag](_uidmetadbmeta)
}

func newPool() *UidPool {
	if _uidpool != nil {
		return _uidpool
	}
	_uidpool = &UidPool{
		isopen: false,
	}
	return _uidpool
}

func newCache() gcache.Cache {
	return gcache.New(_cfg.CachePoolSize).LRU().
		LoaderFunc(func(key interface{}) (interface{}, error) {
			dbp := db.Instance()
			openid := key.(string)
			w := _uidroleiddbwrapermeta.New(openid)
			if err := w.Read(context.TODO(), dbp); err == nil {
				return w.Obj().GetRoleId(), nil
			} else {
				return nil, err
			}
		}).Build()
}

func Pool() *UidPool {
	return _uidpool
}

type UidPool struct {
	isopen bool
	mu     sync.RWMutex

	metawraper *db.Wraper[*pb.UidMetaDB, pb.UMFlag]
	cache      gcache.Cache
}

func (u *UidPool) Close() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.isopen = false
}

func (u *UidPool) Open() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.isopen = true
	u.cache = newCache()
	u.metawraper = _uidmetadwraperbmeta.New(config.GetWorldID().String())
}

func (u *UidPool) IsOpen() bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.isopen
}

func (u *UidPool) GetMetaWraper(ctx context.Context) (*db.Wraper[*pb.UidMetaDB, pb.UMFlag], error) {
	dbp := db.Instance()
	if err := u.metawraper.Read(ctx, dbp); err == nil {
		return u.metawraper, nil
	} else if !errs.DBKeyNotFound.Is(err) {
		return nil, err
	} else {
		if err := u.metawraper.Create(ctx, dbp); err != nil {
			return nil, err
		} else {
			return u.metawraper, nil
		}
	}
}

func (u *UidPool) AllocRoleId(ctx context.Context, openid string) (int64, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if !u.isopen {
		return 0, errors.New("UidPool is Close")
	}
	// 先查缓存
	if value, err := u.cache.Get(openid); err == nil {
		log.Debug("Cache Get %s %v", openid, value)
		return value.(int64), nil
	} else if !errs.DBKeyNotFound.Is(err) {
		return 0, err
	}
	dbp := db.Instance()
	var roleid int64
	// 分配新RoleId
	if metaw, err := u.GetMetaWraper(ctx); err != nil {
		return 0, err
	} else {
		roleid = metaw.Obj().CurRoleId
		metaw.Set(pb.UMFlag_UMCurRoleId, roleid+1)
		if err := metaw.Update(ctx, dbp); err != nil {
			return 0, err
		}
	}
	// 写入DB
	w := _uidroleiddbwrapermeta.New(openid)
	w.Set(pb.URFlag_URRoleId, roleid)
	if err := w.Create(ctx, dbp); err != nil {
		return 0, err
	}
	// 返回
	u.cache.Set(openid, roleid)
	log.Debug("Alloc New %s %d", openid, roleid)
	return roleid, nil
}

func (u *UidPool) AllocClubId(ctx context.Context, roleid int64) (int64, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if !u.isopen {
		return 0, errors.New("UidPool is Close")
	}
	dbp := db.Instance()
	// 分配新ClubId
	var clubid int64
	if metaw, err := u.GetMetaWraper(ctx); err != nil {
		return 0, err
	} else {
		clubid = metaw.Obj().CurClubId
		metaw.Set(pb.UMFlag_UMCurClubId, clubid+1)
		if err := metaw.Update(ctx, dbp); err != nil {
			return 0, err
		}
	}
	// 写入DB
	w := _uidclubiddbwrapermeta.New(clubid)
	w.Set(pb.UCFlag_UCOwnerId, roleid)
	if err := w.Create(ctx, dbp); err != nil {
		return 0, err
	}
	// 返回
	return clubid, nil
}
