package uidser

import (
	"errors"
	"fmt"
	"m3game/config"
	"m3game/demo/proto/pb"
	"m3game/plugins/db"
	"m3game/plugins/db/wraper"
	"sync"

	"github.com/bluele/gcache"
)

var (
	_uidclubiddbmeta *db.DBMeta[*pb.UidClubIdDB]
	_uidroleiddbmeta *db.DBMeta[*pb.UidRoleIdDB]
	_uidmetadbmeta   *db.DBMeta[*pb.UidMetaDB]
	_uidpool         *UidPool
)

func init() {
	_uidclubiddbmeta = db.NewMeta("uidclubid_table", uidclubidCreater)
	_uidroleiddbmeta = db.NewMeta("uidroleid_table", uidroleidCreater)
	_uidmetadbmeta = db.NewMeta("uidmeta_table", uidmetaCreater)
	_uidpool = &UidPool{
		isopen: false,
	}
}

func newCache() gcache.Cache {
	return gcache.New(_cfg.CachePoolSize).LRU().
		LoaderFunc(func(key interface{}) (interface{}, error) {
			dbp := db.Get()
			openid := key.(string)
			w := wraper.New(_uidroleiddbmeta, openid)
			if err := w.Read(dbp); err == nil {
				if roleid := w.TObj().GetRoleId(); roleid != nil {
					return roleid.Value, nil
				} else {
					// 异常
					return nil, errors.New("UidPool is Err, roldid is nil")
				}
			} else {
				return nil, err
			}
		}).Build()
}

func uidclubidCreater() *pb.UidClubIdDB {
	return &pb.UidClubIdDB{
		ClubId: "",
		RoleId: &pb.UidClubIdRoleId{},
	}
}

func uidroleidCreater() *pb.UidRoleIdDB {
	return &pb.UidRoleIdDB{
		OpenId: "",
		RoleId: &pb.UidRoleIdRoleId{},
	}
}

func uidmetaCreater() *pb.UidMetaDB {
	return &pb.UidMetaDB{
		WorldId:   "",
		CurRoleId: &pb.UidMetaCurRoleId{},
		CurClubId: &pb.UidMetaCurClubId{},
	}
}

func Pool() *UidPool {
	return _uidpool
}

type UidPool struct {
	isopen bool
	mu     sync.Mutex

	metawraper *wraper.Wraper[*pb.UidMetaDB]
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
	u.metawraper = wraper.New(_uidmetadbmeta, config.GetWorldID().String())
}

func (u *UidPool) IsOpen() bool {
	return u.isopen
}

func (u *UidPool) GetMetaWraper() (*wraper.Wraper[*pb.UidMetaDB], error) {
	dbp := db.Get()
	if err := u.metawraper.Read(dbp); err == nil {
		return u.metawraper, nil
	} else if !db.IsErrDBNotFindKey(err) {
		return nil, err
	} else {
		if err := u.metawraper.Create(dbp); err != nil {
			return nil, err
		} else {
			return u.metawraper, nil
		}
	}
}

func (u *UidPool) AllocRoleId(openid string) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	var roldid string
	roldid = ""
	if !u.IsOpen() {
		return roldid, errors.New("UidPool is Close")
	}
	// 先查缓存
	if value, err := u.cache.Get(openid); err == nil {
		return value.(string), nil
	} else if !db.IsErrDBNotFindKey(err) {
		return roldid, err
	}
	dbp := db.Get()
	// 分配新RoleId
	if metaw, err := u.GetMetaWraper(); err != nil {
		return roldid, err
	} else {
		var curroleid *pb.UidMetaCurRoleId
		if curroleid = metaw.TObj().GetCurRoleId(); curroleid == nil {
			if err := wraper.Setter(metaw, &pb.UidMetaCurRoleId{Value: 1}); err != nil {
				return roldid, err
			}
			curroleid = metaw.TObj().GetCurRoleId()
		}
		roldid = fmt.Sprintf("%d", curroleid.Value)
		curroleid.Value += 1
		wraper.Setter(metaw, curroleid)
		if err := metaw.Update(dbp); err != nil {
			return roldid, err
		}
	}
	// 写入DB
	w := wraper.New(_uidroleiddbmeta, openid)
	if err := wraper.Setter(w, &pb.UidRoleIdRoleId{Value: roldid}); err != nil {
		return roldid, err
	}
	if err := w.Create(dbp); err != nil {
		return roldid, err
	}
	// 返回
	u.cache.Set(openid, roldid)
	return roldid, nil
}

func (u *UidPool) AllocClubId(roldid string) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	var clubid string
	clubid = ""
	if !u.IsOpen() {
		return clubid, errors.New("UidPool is Close")
	}
	dbp := db.Get()
	// 分配新ClubId
	if metaw, err := u.GetMetaWraper(); err != nil {
		return clubid, err
	} else {
		var curclubid *pb.UidMetaCurClubId
		if curclubid = metaw.TObj().GetCurClubId(); curclubid == nil {
			if err := wraper.Setter(metaw, &pb.UidMetaCurClubId{Value: 1}); err != nil {
				return clubid, err
			}
			curclubid = metaw.TObj().GetCurClubId()
		}
		clubid = fmt.Sprintf("%d", curclubid.Value)
		curclubid.Value += 1
		wraper.Setter(metaw, curclubid)
		if err := metaw.Update(dbp); err != nil {
			return clubid, err
		}
	}
	// 写入DB
	w := wraper.New(_uidclubiddbmeta, clubid)
	if err := wraper.Setter(w, &pb.UidClubIdRoleId{Value: roldid}); err != nil {
		return clubid, err
	}
	if err := w.Create(dbp); err != nil {
		return clubid, err
	}
	// 返回
	return clubid, nil
}
