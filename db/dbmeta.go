package db

import (
	"fmt"
	"m3game/proto/pb"
	"m3game/util/log"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type DBCreater func() proto.Message

type DBMeta struct {
	Table            string
	Keyfield         string
	Allfields        []string
	Creater          DBCreater
	fieldDescriptors map[string]protoreflect.FieldDescriptor
}

func NewMeta(table string, creater DBCreater) *DBMeta {
	meta := &DBMeta{
		Table:            table,
		Creater:          creater,
		fieldDescriptors: make(map[string]protoreflect.FieldDescriptor),
	}
	msg := creater()
	md := msg.ProtoReflect()
	mdescriptor := md.Descriptor()
	mname := mdescriptor.Name()
	moptions := mdescriptor.Options()
	if moptions == nil {
		panic(fmt.Sprintf("NewMeta %s not have MessageOptions", mname))
	}
	if v := proto.GetExtension(moptions, pb.E_DbPrimaryKey); v == nil {
		panic(fmt.Sprintf("NewMeta %s not have E_DbPrimaryKey", mname))
	} else if keyfield, ok := v.(string); !ok {
		panic(fmt.Sprintf("NewMeta %s E_DbPrimaryKey type err", mname))
	} else {
		log.Fatal("DBMeta %s KeyField => %s", mname, keyfield)
		meta.Keyfield = keyfield
	}
	for i := 0; i < mdescriptor.Fields().Len(); i++ {
		fdescriptor := mdescriptor.Fields().Get(i)
		fname := string(fdescriptor.Name())
		if fdescriptor.Kind() != protoreflect.StringKind &&
			fdescriptor.Kind() != protoreflect.MessageKind {
			panic(fmt.Sprintf("NewMeta %s field %s Kind not vaild", mname, fname))
		}
		meta.Allfields = append(meta.Allfields, fname)
		meta.fieldDescriptors[fname] = fdescriptor
		log.Fatal("DBMeta %s AllField => %s", mname, fname)
	}
	return meta
}

func (meta *DBMeta) Setter(msg proto.Message, field string, buf []byte) error {
	fdescriptor, ok := meta.fieldDescriptors[field]
	if !ok {
		return nil
	}
	if fdescriptor.Kind() == protoreflect.StringKind {
		msg.ProtoReflect().Set(fdescriptor, protoreflect.ValueOfString(string(buf)))
	}
	if fdescriptor.Kind() == protoreflect.MessageKind {
		var m proto.Message
		if err := proto.Unmarshal(buf, m); err != nil {
			return err
		}
		msg.ProtoReflect().Set(fdescriptor, protoreflect.ValueOfMessage(m.ProtoReflect()))
	}
	return nil
}
func (meta *DBMeta) Getter(msg proto.Message, field string) ([]byte, error) {
	fdescriptor, ok := meta.fieldDescriptors[field]
	if !ok {
		return nil, nil
	}
	if fdescriptor.Kind() == protoreflect.StringKind {
		s := msg.ProtoReflect().Get(fdescriptor).String()
		return []byte(s), nil
	}
	if fdescriptor.Kind() == protoreflect.MessageKind {
		v := msg.ProtoReflect().Get(fdescriptor).Message().Interface()
		m, ok := v.(proto.Message)
		if !ok {
			return nil, nil
		}
		return proto.Marshal(m)
	}
	return nil, nil

}
