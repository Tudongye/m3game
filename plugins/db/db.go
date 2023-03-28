package db

import (
	"context"
	"m3game/meta/errs"
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
		return nil, errs.DBInsHasNewed.New("db is newed %s", _db.Factory().Name())
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
