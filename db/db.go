package db

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

// db-plguin interface
type DB interface {
	Get(meta *DBMeta, key string) (proto.Message, error)
	Update(meta *DBMeta, key string, obj proto.Message) error
	Insert(meta *DBMeta, key string, obj proto.Message) error
	Delete(meta *DBMeta, key string) error
}

var (
	Err_DB_notfindkey  = fmt.Errorf("Err_DB_notfindkey")
	Err_DB_repeatedkey = fmt.Errorf("Err_DB_repeatedkey")
)
