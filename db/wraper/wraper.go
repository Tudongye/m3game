package wraper

import (
	"fmt"
	"m3game/db"
	"m3game/log"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func New[T proto.Message](meta *db.DBMeta[T], key string) *Wraper[T] {
	w := &Wraper[T]{
		obj:    meta.ObjCreater(),
		key:    key,
		meta:   meta,
		dirtys: make(map[string]bool),
	}
	return w
}

type Wraper[T proto.Message] struct {
	key    string          // 实体Key
	obj    T               // 实体pb.Message
	meta   *db.DBMeta[T]   // Meta
	dirtys map[string]bool // 置脏标记
}

func (w *Wraper[T]) Meta() *db.DBMeta[T] {
	return w.meta
}

func (w *Wraper[T]) Obj() proto.Message {
	return w.obj
}

func (w *Wraper[T]) HasDirty() bool {
	return len(w.dirtys) > 0
}

func (w *Wraper[T]) SetDirty(field string) {
	w.dirtys[field] = true
}

func (w *Wraper[T]) GetDirtyList() []string {
	var dirtys []string
	for field := range w.dirtys {
		dirtys = append(dirtys, field)
	}
	return dirtys
}

func (w *Wraper[T]) ClearDirty() {
	w.dirtys = make(map[string]bool)
}

func (w *Wraper[T]) Update(db db.DB) error {
	if !w.HasDirty() {
		return nil
	}
	dirts := w.GetDirtyList()
	log.Debug("Update %v", dirts)
	if err := db.Update(w.meta, w.key, w.obj, dirts...); err != nil {
		return errors.Wrapf(err, "DB Obj %s Update Fail field %v", w.Meta().ObjName(), dirts)
	}
	w.ClearDirty()
	return nil
}

func (w *Wraper[T]) Create(db db.DB) error {
	log.Debug("%v", w.obj)
	if err := db.Create(w.meta, w.key, w.obj); err != nil {
		return err
	}
	w.ClearDirty()
	return nil
}

func (w *Wraper[T]) Delete(db db.DB) error {
	return db.Delete(w.meta, w.key)
}

func (w *Wraper[T]) Read(db db.DB) error {
	if v, err := db.Read(w.meta, w.key); err != nil {
		return err
	} else if value, ok := v.(T); !ok {
		return nil
	} else {
		w.obj = value
	}
	return nil
}

func KeySetter[T proto.Message](wraper *Wraper[T], value string) error {
	fieldd, ok := wraper.Meta().Fieldds()[wraper.Meta().KeyField()]
	if !ok {
		return fmt.Errorf("DB field %s not find in %s Meta", wraper.Meta().KeyField(), wraper.Meta().ObjName())
	}
	wraper.Obj().ProtoReflect().Set(fieldd, protoreflect.ValueOfString(value))
	wraper.SetDirty(wraper.Meta().KeyField())
	return nil
}

func KeyGetter[T proto.Message](wraper *Wraper[T]) (string, error) {
	fieldd, ok := wraper.Meta().Fieldds()[wraper.Meta().KeyField()]
	if !ok {
		return "", fmt.Errorf("DB field %s not find in %s Meta", wraper.Meta().KeyField(), wraper.Meta().ObjName())
	}
	return wraper.Obj().ProtoReflect().Get(fieldd).String(), nil
}

func Setter[P, T proto.Message](wraper *Wraper[T], value P) error {
	var p P
	pbname := string(p.ProtoReflect().Descriptor().FullName())
	if fieldname, ok := wraper.Meta().AllPBName()[pbname]; !ok {
		log.Debug("%v", wraper.Meta().AllPBName())
		return fmt.Errorf("DB Obj %s not have field type %s", wraper.Meta().ObjName(), pbname)
	} else {
		fieldd, ok := wraper.Meta().Fieldds()[fieldname]
		if !ok {
			return fmt.Errorf("DB field %s not find in %s Meta", fieldname, wraper.Meta().ObjName())
		}
		wraper.Obj().ProtoReflect().Set(fieldd, protoreflect.ValueOfMessage(value.ProtoReflect()))
		wraper.SetDirty(fieldname)
		return nil
	}
}

func Getter[P, T proto.Message](wraper *Wraper[T]) (P, error) {
	var p P
	pbname := string(p.ProtoReflect().Descriptor().FullName())
	if fieldname, ok := wraper.Meta().AllPBName()[pbname]; !ok {
		return p, fmt.Errorf("DB Obj %s not have field type %s", wraper.Meta().ObjName(), pbname)
	} else {
		fieldd, ok := wraper.Meta().Fieldds()[fieldname]
		if !ok {
			return p, fmt.Errorf("DB field %s not find in %s Meta", fieldname, wraper.Meta().ObjName())
		}
		if v := wraper.Obj().ProtoReflect().Get(fieldd).Message().Interface(); v == nil {
			return p, fmt.Errorf("DB Obj %s Field %s is nil", wraper.Meta().ObjName(), fieldname)
		} else if value, ok := v.(P); !ok {
			return p, fmt.Errorf("DB Obj %s Field %s is not %s", wraper.Meta().ObjName(), fieldname, pbname)
		} else {
			return value, nil
		}
	}
}
