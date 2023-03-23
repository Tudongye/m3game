package db

import (
	"fmt"
	"m3game/plugins/log"
	"m3game/runtime/plugin"

	"google.golang.org/protobuf/proto"
)

// db-plguin interface
type DB interface {
	plugin.PluginIns
	Read(meta DBMetaInter, key string, filters ...string) (proto.Message, error)
	Update(meta DBMetaInter, key string, obj proto.Message, filters ...string) error
	Create(meta DBMetaInter, key string, obj proto.Message, filters ...string) error
	Delete(meta DBMetaInter, key string) error
}

var (
	_db DB
)

func New(db DB) (DB, error) {
	if _db != nil {
		return nil, fmt.Errorf("db is newed %s", _db.Factory().Name())
	}
	_db = db
	return _db, nil
}

func Instance() DB {
	if _db == nil {
		log.Fatal("DB not newd")
		return nil
	}
	return _db
}

var (
	Err_KeyNotFound    = fmt.Errorf("Err_KeyNotFound")
	Err_DuplicateEntry = fmt.Errorf("Err_DuplicateEntry")
)

func IsErrKeyNotFound(e error) bool {
	if e == Err_KeyNotFound {
		return true
	}
	return false
}

func IsErrDuplicateEntry(e error) bool {
	if e == Err_DuplicateEntry {
		return true
	}
	return false
}
