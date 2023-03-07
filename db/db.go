package db

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

// db-plguin interface
type DB interface {
	Read(meta DBMetaInter, key string, filters ...string) (proto.Message, error)
	Update(meta DBMetaInter, key string, obj proto.Message, filters ...string) error
	Create(meta DBMetaInter, key string, obj proto.Message, filters ...string) error
	Delete(meta DBMetaInter, key string) error
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
