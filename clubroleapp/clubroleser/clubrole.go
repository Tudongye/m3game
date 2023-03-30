package clubroleser

import (
	"m3game/demo/proto/pb"
	"m3game/plugins/db"
)

var (
	_clubroledbmeta     *db.DBMeta[*pb.ClubRoleDB]
	_clubrolewrapermeta *db.WraperMeta[*pb.ClubRoleDB, pb.CRFlag]
)

func init() {
	_clubroledbmeta = db.NewMeta[*pb.ClubRoleDB]("clubrole_table")
	_clubrolewrapermeta = db.NewWraperMeta[*pb.ClubRoleDB, pb.CRFlag](_clubroledbmeta)
}
