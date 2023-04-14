package club

import (
	"context"
	"m3game/demo/proto/pb"
	"m3game/plugins/db"
	"m3game/plugins/lease"
)

var (
	_clubdbmeta         *db.DBMeta[*pb.ClubDB]
	_clubroledbmeta     *db.DBMeta[*pb.ClubRoleDB]
	_clubwrapermeta     *db.WraperMeta[*pb.ClubDB, pb.CFlag]
	_clubrolewrapermeta *db.WraperMeta[*pb.ClubRoleDB, pb.CRFlag]

	_clubleasemeta *lease.LeaseMeta
)

func init() {
	_clubdbmeta = db.NewMeta[*pb.ClubDB]("club_table")
	_clubroledbmeta = db.NewMeta[*pb.ClubRoleDB]("clubrole_table")
	_clubwrapermeta = db.NewWraperMeta[*pb.ClubDB, pb.CFlag](_clubdbmeta)
	_clubrolewrapermeta = db.NewWraperMeta[*pb.ClubRoleDB, pb.CRFlag](_clubroledbmeta)
}

func newClubDB(ctx context.Context, clubid int64, ownroleid int64) error {
	dbplugin := db.Instance()
	w := _clubwrapermeta.New(clubid)
	w.Set(pb.CFlag_CClubId, clubid)
	w.Set(pb.CFlag_COwnerId, ownroleid)
	return w.Create(ctx, dbplugin)
}

func newClubRoleDB(ctx context.Context, roleid int64, clubid int64) error {
	dbplugin := db.Instance()
	w := _clubrolewrapermeta.New(roleid)
	w.Set(pb.CRFlag_CRRoleId, roleid)
	w.Set(pb.CRFlag_CRClubId, clubid)
	return w.Create(ctx, dbplugin)
}

func NewLeaseMeta(lm *lease.LeaseMeta) {
	_clubleasemeta = lm
}

func LeaseMeta() *lease.LeaseMeta {
	return _clubleasemeta
}
