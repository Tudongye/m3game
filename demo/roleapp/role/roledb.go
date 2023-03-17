package role

import (
	"m3game/demo/proto/pb"
	"m3game/plugins/db"
)

var (
	rolemeta *db.DBMeta[*pb.RoleDB]
)

func init() {
	rolemeta = db.NewMeta("role_table", roledbCreater)
}

func roledbCreater() *pb.RoleDB {
	return &pb.RoleDB{
		RoleId: "",
		Name:   &pb.RoleName{},
		Power:  &pb.RolePower{},
	}
}
