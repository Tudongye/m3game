package db

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

func CreateDBMeta(table string, keyfield string, allfields []string, creater DBCreater, setter DBSetter, getter DBGetter) *DBMeta {
	return &DBMeta{
		Table:     table,
		Keyfield:  keyfield,
		Allfields: allfields,
		Creater:   creater,
		Setter:    setter,
		Getter:    getter,
	}
}

type DBCreater func() proto.Message
type DBSetter func(proto.Message, string, interface{})
type DBGetter func(proto.Message, string) interface{}

type DBMeta struct {
	Table     string
	Keyfield  string
	Allfields []string
	Creater   DBCreater
	Setter    DBSetter
	Getter    DBGetter
}

type DB interface {
	Get(meta *DBMeta, key string) (proto.Message, error)
	Update(meta *DBMeta, key string, obj proto.Message) error
	Insert(meta *DBMeta, key string, obj proto.Message) error
	Delete(meta *DBMeta, key string) error
}

var (
	Err_DB_notfind     = fmt.Errorf("Err_DB_notfind")
	Err_DB_repeatedkey = fmt.Errorf("Err_DB_notfind")
)
