/*
	package db is defined dbmeta and db-plugin interface
	db包用来定义dbmeta 和 db插件的接口
*/

package db

import (
	"fmt"
	"m3game/proto/pb"
	"m3game/util/log"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Create DB Object
type DBMetaInter interface {
	ObjName() string
	Table() string
	KeyField() string
	AllFields() []string
	AllPBName() map[string]string
	Creater() func() proto.Message
	Fieldds() map[string]protoreflect.FieldDescriptor

	Decode(msg proto.Message, field string, buf []byte) error
	Encode(msg proto.Message, field string) ([]byte, error)
	HasField(field string) bool
}

/*
DBMeta represents the format of the game entity in the DB, defined by the game logic, and used as the initialization parameter of the DB plug-in
DBMeta表示游戏实体在DB中的格式，由游戏逻辑定义，作为DB插件的初始化参数
*/
type DBMeta[T proto.Message] struct {
	objName   string
	table     string                                  // DB表名
	keyField  string                                  // 主键，强制为string
	allFields []string                                // 所有数据键
	allPBName map[string]string                       // 类型名到字段名映射
	creater   func() T                                // 游戏实体工场
	fieldds   map[string]protoreflect.FieldDescriptor // 游戏实体字段反射信息
}

var (
	_ DBMetaInter = (*DBMeta[proto.Message])(nil)
)

/*
Create DBMeta
创建DBMeta
*/
func NewMeta[T proto.Message](table string, creater func() T) *DBMeta[T] {
	meta := &DBMeta[T]{
		table:     table,
		creater:   creater,
		allPBName: make(map[string]string),
		fieldds:   make(map[string]protoreflect.FieldDescriptor),
	}

	// game entity => proto.message.descriptor
	// 游戏实体 到 pb.message.descriptor 用于读取pb的自定义Option
	messaged := creater().ProtoReflect().Descriptor()
	messagename := messaged.Name()
	messageopts := messaged.Options()
	if messageopts == nil {
		panic(fmt.Sprintf("DB %s not have MessageOptions", messagename))
	}
	meta.objName = string(messagename)
	// game entity pb.message must have option E_DbPrimaryKey
	keyfield := ""
	if v := proto.GetExtension(messageopts, pb.E_DbPrimaryKey); v == nil {
		panic(fmt.Sprintf("DB %s not have E_DbPrimaryKey", messagename))
	} else if v, ok := v.(string); !ok {
		panic(fmt.Sprintf("DB %s E_DbPrimaryKey type err", messagename))
	} else {
		keyfield = v
	}

	// collect all field
	for i := 0; i < messaged.Fields().Len(); i++ {
		fieldd := messaged.Fields().Get(i)
		fieldname := string(fieldd.Name())

		// Currently only supported string and pb.message in first layer
		// 当前一级结构仅支持string和message
		if fieldd.Kind() != protoreflect.StringKind &&
			fieldd.Kind() != protoreflect.MessageKind {
			panic(fmt.Sprintf("DB %s field %s Kind not vaild", messagename, fieldname))
		}
		if fieldname == keyfield {
			if fieldd.Kind() != protoreflect.StringKind {
				panic(fmt.Sprintf("DB %s KeyField %s Kind must be string", messagename, fieldname))
			}
			meta.keyField = keyfield
			log.Fatal("DB %s KeyField => %s", messagename, keyfield)

		} else {
			if fieldd.Kind() != protoreflect.MessageKind {
				panic(fmt.Sprintf("DB %s Field %s Kind must be pb.Message", messagename, fieldname))
			}
			if _, ok := meta.allPBName[string(fieldd.Message().FullName())]; ok {
				panic(fmt.Sprintf("DB %s field %s FullName %s repeated", messagename, fieldname, fieldd.Message().FullName()))
			}
			meta.allPBName[string(fieldd.Message().FullName())] = fieldname
			log.Fatal("DB %s PBField => %s", messagename, fieldname)
		}
		meta.allFields = append(meta.allFields, fieldname)
		meta.fieldds[fieldname] = fieldd
		log.Fatal("DB %s AllField => %s", messagename, fieldname)
	}
	if meta.keyField == "" {
		panic(fmt.Sprintf("DB %s not find KeyField", messagename))
	}
	return meta
}

func (meta *DBMeta[T]) Decode(msg proto.Message, field string, buf []byte) error {
	fieldd, ok := meta.fieldds[field]
	if !ok {
		return fmt.Errorf("DB field %s not find in %s Meta", field, meta.ObjName())
	}
	switch fieldd.Kind() {
	case protoreflect.StringKind:
		msg.ProtoReflect().Set(fieldd, protoreflect.ValueOfString(string(buf)))
		return nil
	case protoreflect.MessageKind:
		v := msg.ProtoReflect().Get(fieldd).Message().Interface()
		m, ok := v.(proto.Message)
		if !ok {
			return fmt.Errorf("DB Obj %s Field %s is not pb.message", meta.ObjName(), field)
		}
		if fieldd.Message().FullName() != m.ProtoReflect().Descriptor().FullName() {
			return fmt.Errorf("DB Obj %s Field %s FullName %s but UnMarshal FullName %s", meta.ObjName(), field,
				fieldd.Message().FullName(), m.ProtoReflect().Descriptor().FullName())
		}
		msg.ProtoReflect().Set(fieldd, protoreflect.ValueOfMessage(m.ProtoReflect()))
		return nil
	default:
		return fmt.Errorf("DB Obj %s Field %s Kind %d is invaild", meta.ObjName(), field, fieldd.Kind())
	}
}

func (meta *DBMeta[T]) Encode(msg proto.Message, field string) ([]byte, error) {
	fieldd, ok := meta.fieldds[field]
	if !ok {
		return nil, fmt.Errorf("DB field %s not find in %s Meta", field, meta.ObjName())
	}
	switch fieldd.Kind() {
	case protoreflect.StringKind:
		s := msg.ProtoReflect().Get(fieldd).String()
		return []byte(s), nil
	case protoreflect.MessageKind:
		v := msg.ProtoReflect().Get(fieldd).Message().Interface()
		m, ok := v.(proto.Message)
		if !ok {
			return nil, fmt.Errorf("DB Obj %s Field %s is not pb.message", meta.ObjName(), field)
		}
		b, e := proto.Marshal(m)
		return b, e
	default:
		return nil, fmt.Errorf("DB Obj %s Field %s Kind %d is invaild", meta.ObjName(), field, fieldd.Kind())
	}
}

func (meta *DBMeta[T]) HasField(field string) bool {
	_, ok := meta.fieldds[field]
	return ok
}

func (meta *DBMeta[T]) ObjName() string {
	return meta.objName
}
func (meta *DBMeta[T]) Table() string {
	return meta.table
}
func (meta *DBMeta[T]) KeyField() string {
	return meta.keyField
}
func (meta *DBMeta[T]) AllFields() []string {
	return meta.allFields
}

func (meta *DBMeta[T]) AllPBName() map[string]string {
	return meta.allPBName
}

func (meta *DBMeta[T]) Creater() func() proto.Message {
	return func() proto.Message {
		return meta.creater()
	}
}

func (meta *DBMeta[T]) ObjCreater() T {
	return meta.creater()
}

func (meta *DBMeta[T]) Fieldds() map[string]protoreflect.FieldDescriptor {
	return meta.fieldds
}
