package roleserver

import (
	"m3game/db"
	"m3game/demo/proto/pb"
)

var (
	rolemeta *db.DBMeta[*pb.RoleDB]
)

func init() {
	rolemeta = db.NewMeta("role_table", roleDBCreater)
}

func roleDBCreater() *pb.RoleDB {
	return &pb.RoleDB{
		RoleID:       "",
		RoleName:     &pb.RoleName{},
		LocationInfo: &pb.LocationInfo{},
	}
}
