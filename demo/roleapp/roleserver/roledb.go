package roleserver

import (
	"m3game/db"
	"m3game/demo/proto/pb"
	"m3game/util/log"

	"google.golang.org/protobuf/proto"
)

var (
	rolemeta *db.DBMeta
)

func init() {
	rolemeta = db.NewMeta("role_table", "RoleID", []string{"RoleID", "Name"}, RoleDBCreater, RoleSetter, RoleGetter)
}

func RoleDBCreater() proto.Message {
	return roleDBCreater()
}

func RoleSetter(obj proto.Message, field string, value interface{}) {
	if roledb, ok := obj.(*pb.RoleDB); !ok {
		log.Error("unknow message type, want roledb")
	} else {
		switch field {
		case "RoleID":
			FillValue(&roledb.RoleID, value, field)
		case "Name":
			FillValue(&roledb.Name, value, field)
		default:
			log.Error("unknow field %s for roledb", field)
		}
	}
}
func RoleGetter(obj proto.Message, field string) interface{} {
	if roledb, ok := obj.(*pb.RoleDB); !ok {
		log.Error("unknow message type, want roledb")
		return nil
	} else {
		switch field {
		case "RoleID":
			return roledb.RoleID
		case "Name":
			return roledb.Name
		default:
			log.Error("unknow field %s for roledb", field)
			return nil
		}
	}
}

func roleDBCreater() *pb.RoleDB {
	return &pb.RoleDB{}
}

func FillValue[T any](p *T, v interface{}, field string) {
	if value, ok := v.(T); ok {
		*p = value
	} else {
		log.Error("unknow type for field %s", field)
	}
}
