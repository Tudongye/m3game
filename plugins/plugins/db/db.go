package db

import (
	"context"
	"fmt"
	"m3game/plugins/log"
	"m3game/runtime/plugin"

	"google.golang.org/protobuf/proto"
)

var (
	_db DB
)

// db-plguin interface
type DB interface {
	plugin.PluginIns
	Read(ctx context.Context, meta DBMetaInter, key interface{}, flags ...int32) (proto.Message, error)
	Update(ctx context.Context, meta DBMetaInter, key interface{}, obj proto.Message, flags ...int32) error
	Create(ctx context.Context, meta DBMetaInter, key interface{}, obj proto.Message) error
	Delete(ctx context.Context, meta DBMetaInter, key interface{}) error

	ReadMany(ctx context.Context, meta DBMetaInter, filters interface{}, flags ...int32) ([]proto.Message, error)
}

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
