package db

import (
	"fmt"
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

func Get() DB {
	return _db
}

func Set(d DB) {
	_db = d
}

var (
	Err_DB_notfindkey  = fmt.Errorf("Err_DB_notfindkey")
	Err_DB_repeatedkey = fmt.Errorf("Err_DB_repeatedkey")
)

func IsErrDBNotFindKey(e error) bool {
	if e == Err_DB_notfindkey {
		return true
	}
	return false
}

func IsErrDBRepeatedKey(e error) bool {
	if e == Err_DB_repeatedkey {
		return true
	}
	return false
}
