package db

import (
	"context"
	"m3game/meta/metapb"
	"m3game/plugins/log"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Flag interface {
	comparable
	protoreflect.Enum
}

type WraperMeta[TM proto.Message, TF Flag] struct {
	dbmeta *DBMeta[TM]
}

func NewWraperMeta[TM proto.Message, TF Flag](meta *DBMeta[TM]) *WraperMeta[TM, TF] {
	wrapermeta := &WraperMeta[TM, TF]{
		dbmeta: meta,
	}
	// 收集所有flag
	var tf TF
	flags := make(map[protoreflect.EnumNumber]protoreflect.Name)
	for i := 0; i < tf.Descriptor().Values().Len(); i++ {
		name := tf.Descriptor().Values().Get(i).Name()
		number := tf.Descriptor().Values().Get(i).Number()
		if number == 0 {
			continue
		}
		if meta.FlagField(int32(number)) == nil {
			log.Fatal("Flag %d not find in dbmeta", number)
		}
		flags[number] = name
	}
	for _, flag := range meta.AllFlags() {
		if _, ok := flags[protoreflect.EnumNumber(flag)]; !ok {
			log.Fatal("Flag %d not find in Enum", flag)
		}
		fieldd := meta.FlagField(flag)
		dbfieldoption := proto.GetExtension(fieldd.Options(), metapb.E_DbfieldOption).(*metapb.M3DBFieldOption)

		if flags[protoreflect.EnumNumber(flag)] != protoreflect.Name(dbfieldoption.Flag) {
			log.Fatal("Flag %d FieldFlag %s RealdFlag %s", flag,
				flags[protoreflect.EnumNumber(flag)], protoreflect.Name(dbfieldoption.Flag))
		}
	}
	return wrapermeta

}

func (w *WraperMeta[TM, TF]) Set(msg TM, flag TF, data interface{}) {
	w.dbmeta.Setter(msg, int32(flag.Number()), data)
}

func (w *WraperMeta[TM, TF]) Get(msg TM, flag TF) interface{} {
	return w.dbmeta.Getter(msg, int32(flag.Number()))
}

func (w *WraperMeta[TM, TF]) New(key interface{}) *Wraper[TM, TF] {
	wraper := &Wraper[TM, TF]{
		meta:   w,
		key:    key,
		dirtys: make(map[TF]bool),
		obj:    w.dbmeta.New().(TM),
	}
	w.dbmeta.Setter(wraper.obj, w.dbmeta.KeyFlag(), wraper.key)
	return wraper
}

type Wraper[TM proto.Message, TF Flag] struct {
	meta   *WraperMeta[TM, TF] // Meta
	key    interface{}         // 主键值
	obj    TM                  // 原始数据
	dirtys map[TF]bool         // 脏标记
}

func (w *Wraper[TM, TF]) Obj() TM {
	return w.obj
}

// 写入数据
func (w *Wraper[TM, TF]) Set(flag TF, value interface{}) {
	w.dirtys[flag] = true
	w.meta.Set(w.obj, flag, value)
}

// 获得数据
func (w *Wraper[TM, TF]) Get(flag TF) interface{} {
	return w.meta.Get(w.obj, flag)
}

// 是否置脏
func (w *Wraper[TM, TF]) IsDirty() bool {
	return len(w.dirtys) > 0
}

// 读取数据
func (w *Wraper[TM, TF]) Read(ctx context.Context, db DB) error {
	if value, err := db.Read(ctx, w.meta.dbmeta, w.key); err != nil {
		return err
	} else if v, ok := value.(TM); !ok {
		return nil
	} else {
		w.obj = v
		return nil
	}
}

// 新建数据
func (w *Wraper[TM, TF]) Create(ctx context.Context, db DB) error {
	w.dirtys = make(map[TF]bool)
	return db.Create(ctx, w.meta.dbmeta, w.key, w.obj)
}

// 更新数据
func (w *Wraper[TM, TF]) Update(ctx context.Context, db DB) error {
	var fields []int32
	for k := range w.dirtys {
		fields = append(fields, int32(k.Number()))
	}
	w.dirtys = make(map[TF]bool)
	return db.Update(ctx, w.meta.dbmeta, w.key, w.obj, fields...)
}

// 删除数据
func (w *Wraper[TM, TF]) Delete(ctx context.Context, db DB) error {
	return db.Delete(ctx, w.meta.dbmeta, w.key)
}
