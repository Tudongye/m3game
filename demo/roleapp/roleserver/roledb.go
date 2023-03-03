package roleserver

import (
	"m3game/db"
	"m3game/demo/proto/pb"

	"google.golang.org/protobuf/proto"
)

var (
	rolemeta *db.DBMeta
)

func init() {
	rolemeta = db.NewMeta("role_table", RoleDBCreater)
}

func RoleDBCreater() proto.Message {
	return roleDBCreater()
}

func roleDBCreater() *pb.RoleDB {
	return &pb.RoleDB{}
}
