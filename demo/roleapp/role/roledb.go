package role

import (
	"m3game/demo/proto/pb"
	"m3game/plugins/db"
)

var (
	_roledbmeta     *db.DBMeta[*pb.RoleDB]
	_rolewrapermeta *db.WraperMeta[*pb.RoleDB, pb.RFlag]
)

func init() {
	_roledbmeta = db.NewMeta[*pb.RoleDB]("role_table")
	_rolewrapermeta = db.NewWraperMeta[*pb.RoleDB, pb.RFlag](_roledbmeta)
}
