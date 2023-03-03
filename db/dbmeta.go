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
type DBCreater func() proto.Message

/*
DBMeta represents the format of the game entity in the DB, defined by the game logic, and used as the initialization parameter of the DB plug-in
DBMeta表示游戏实体在DB中的格式，由游戏逻辑定义，作为DB插件的初始化参数
*/
type DBMeta struct {
	Table     string
	Keyfield  string
	Allfields []string
	Creater   DBCreater
	fieldds   map[string]protoreflect.FieldDescriptor // field's pb options
}

/*
Create DBMeta
创建DBMeta
*/
func NewMeta(table string, creater DBCreater) *DBMeta {
	meta := &DBMeta{
		Table:   table,
		Creater: creater,
		fieldds: make(map[string]protoreflect.FieldDescriptor),
	}

	// game entity => proto.message.descriptor
	// 游戏实体 到 pb.message.descriptor 用于读取pb的自定义Option
	messaged := creater().ProtoReflect().Descriptor()
	messagename := messaged.Name()
	messageopts := messaged.Options()
	if messageopts == nil {
		panic(fmt.Sprintf("DB %s not have MessageOptions", messagename))
	}

	// game entity pb.message must have option E_DbPrimaryKey
	if v := proto.GetExtension(messageopts, pb.E_DbPrimaryKey); v == nil {
		panic(fmt.Sprintf("DB %s not have E_DbPrimaryKey", messagename))
	} else if keyfield, ok := v.(string); !ok {
		panic(fmt.Sprintf("DB %s E_DbPrimaryKey type err", messagename))
	} else {
		log.Fatal("DB %s KeyField => %s", messagename, keyfield)
		meta.Keyfield = keyfield
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
		meta.Allfields = append(meta.Allfields, fieldname)
		meta.fieldds[fieldname] = fieldd
		log.Fatal("DB %s AllField => %s", messagename, fieldname)
	}
	return meta
}

// Assignment based on string field name, pbReflect
func (meta *DBMeta) Setter(msg proto.Message, field string, buf []byte) error {
	fieldd, ok := meta.fieldds[field]
	if !ok {
		return nil
	}
	if fieldd.Kind() == protoreflect.StringKind {
		msg.ProtoReflect().Set(fieldd, protoreflect.ValueOfString(string(buf)))
	}
	if fieldd.Kind() == protoreflect.MessageKind {
		var m proto.Message
		if err := proto.Unmarshal(buf, m); err != nil {
			return err
		}
		msg.ProtoReflect().Set(fieldd, protoreflect.ValueOfMessage(m.ProtoReflect()))
	}
	return nil
}

// Read from string field name, pbReflect
func (meta *DBMeta) Getter(msg proto.Message, field string) ([]byte, error) {
	fieldd, ok := meta.fieldds[field]
	if !ok {
		return nil, nil
	}
	if fieldd.Kind() == protoreflect.StringKind {
		s := msg.ProtoReflect().Get(fieldd).String()
		return []byte(s), nil
	}
	if fieldd.Kind() == protoreflect.MessageKind {
		v := msg.ProtoReflect().Get(fieldd).Message().Interface()
		m, ok := v.(proto.Message)
		if !ok {
			return nil, nil
		}
		return proto.Marshal(m)
	}
	return nil, nil
}
