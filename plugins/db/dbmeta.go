/*
	package db is defined dbmeta and db-plugin interface
	db包用来定义dbmeta 和 db插件的接口
*/

package db

import (
	"m3game/meta/metapb"
	"m3game/plugins/log"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Create DB Object

type DBMetaInter interface {
	Setter(msg proto.Message, flag int32, data interface{}) // 赋值
	Getter(msg proto.Message, flag int32) interface{}       // 读取
	FlagKind(flag int32) protoreflect.Kind                  // 获取字段类型
	FlagName(flag int32) string                             // 获取字段类型
	KeyFlag() int32                                         // 主键字段
	AllFlags() []int32                                      // 所有字段名
	New() proto.Message
	Table() string
}

/*
DBMeta represents the format of the game entity in the DB, defined by the game logic, and used as the initialization parameter of the DB plug-in
DBMeta表示游戏实体在DB中的格式，由游戏逻辑定义，作为DB插件的初始化参数
TM 内存对象 TD 存储结构 TF 置脏标记
*/
type DBMeta[TM proto.Message] struct {
	name     string                                 // 数据名
	table    string                                 // DB表名
	keyflag  []int32                                // 主键flag
	allflags []int32                                // 所有flag
	fieldds  map[int32]protoreflect.FieldDescriptor // flag到field映射
}

var (
	_ DBMetaInter = (*DBMeta[proto.Message])(nil)
)

/*
Create DBMeta
创建DBMeta
*/
func NewMeta[TM proto.Message](table string) *DBMeta[TM] {
	meta := &DBMeta[TM]{
		table:   table,
		fieldds: make(map[int32]protoreflect.FieldDescriptor),
	}
	var tm TM
	messaged := tm.ProtoReflect().Descriptor()
	meta.name = string(messaged.Name())

	// 收集所有字段
	for i := 0; i < messaged.Fields().Len(); i++ {
		fieldd := messaged.Fields().Get(i)
		name := fieldd.Name()
		number := fieldd.Number()
		fieldopts := fieldd.Options()
		if fieldopts == nil {
			return nil
		}
		var dbfieldoption *metapb.M3DBFieldOption
		if v := proto.GetExtension(fieldopts, metapb.E_DbfieldOption); v == nil {
			log.Fatal("DB %s not have E_DbPrimaryKey", meta.name)
		} else if v, ok := v.(*metapb.M3DBFieldOption); !ok {
			log.Fatal("DB %s E_DbPrimaryKey type err", meta.name)
		} else {
			dbfieldoption = v
		}
		meta.fieldds[int32(number)] = fieldd
		if dbfieldoption.Primary {
			// 主键
			meta.keyflag = append(meta.keyflag, int32(number))
			log.Info("DB %s KeyField => %v", meta.name, meta.keyflag)
		}
		meta.allflags = append(meta.allflags, int32(number))
		log.Info("DB %s Field => %s", meta.name, name)
	}
	if len(meta.keyflag) != 1 {
		log.Fatal("DB %s not find KeyField", meta.name)
	}
	return meta
}

func (d *DBMeta[TM]) FlagKind(flag int32) protoreflect.Kind {
	return d.fieldds[flag].Kind()
}

func (d *DBMeta[TM]) FlagName(flag int32) string {
	return string(d.fieldds[flag].Name())
}

func (d *DBMeta[TM]) FlagField(flag int32) protoreflect.FieldDescriptor {
	return d.fieldds[flag]
}

func (d *DBMeta[TM]) KeyFlag() int32 {
	return d.keyflag[0]
}

func (d *DBMeta[TM]) AllFlags() []int32 {
	return d.allflags
}

func (d *DBMeta[TM]) Setter(msg proto.Message, flag int32, data interface{}) {
	msg.ProtoReflect().Set(d.fieldds[flag], getReflectValue(d.FlagKind(flag), data))
}

func (d *DBMeta[TM]) Getter(msg proto.Message, flag int32) interface{} {
	return getRealValue(d.FlagKind(flag), msg.ProtoReflect().Get(d.fieldds[flag]))
}

func (d *DBMeta[TM]) New() proto.Message {
	var tm TM
	refType := reflect.TypeOf(tm).Elem()
	value := reflect.New(refType)
	return value.Interface().(proto.Message)
}

func (d *DBMeta[TM]) Table() string {
	return d.table
}

func getReflectValue(kind protoreflect.Kind, data interface{}) protoreflect.Value {
	switch kind {
	case protoreflect.BoolKind:
		return protoreflect.ValueOfBool(data.(bool))
	case protoreflect.StringKind:
		return protoreflect.ValueOfString(data.(string))
	case protoreflect.BytesKind:
		return protoreflect.ValueOfBytes(data.([]byte))
	case protoreflect.FloatKind:
		return protoreflect.ValueOfFloat32(data.(float32))
	case protoreflect.DoubleKind:
		return protoreflect.ValueOfFloat64(data.(float64))
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return protoreflect.ValueOfInt32(data.(int32))
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return protoreflect.ValueOfUint32(data.(uint32))
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return protoreflect.ValueOfInt64(data.(int64))
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return protoreflect.ValueOfUint64(data.(uint64))
	case protoreflect.EnumKind:
		return protoreflect.ValueOfEnum(data.(protoreflect.EnumNumber))
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return protoreflect.ValueOfMessage(data.(protoreflect.ProtoMessage).ProtoReflect())
	default:
		// 其它类型，返回空类型
		log.Fatal("Unknow ReflectType %v", kind)
		return protoreflect.Value{}
	}
}

func getRealValue(kind protoreflect.Kind, data protoreflect.Value) interface{} {
	switch kind {
	case protoreflect.BoolKind:
		return data.Bool()
	case protoreflect.StringKind:
		return data.String()
	case protoreflect.BytesKind:
		return data.Bytes()
	case protoreflect.FloatKind:
		return float32(data.Float())
	case protoreflect.DoubleKind:
		return data.Float()
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return int32(data.Int())
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return uint32(data.Uint())
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return data.Int()
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return data.Uint()
	case protoreflect.EnumKind:
		return data.Enum()
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return data.Message().Interface()
	default:
		// 其它类型，返回空类型
		log.Fatal("Unknow ReflectType %v", kind)
		return nil
	}
}
