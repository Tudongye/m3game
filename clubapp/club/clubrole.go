package club

import (
	"context"
	"m3game/demo/proto/pb"
	"m3game/plugins/db"
	"m3game/plugins/log"
)

func newClubRole(ctx context.Context, roleid int64) (*ClubRole, error) {
	clubrole := &ClubRole{
		wraper: _clubrolewrapermeta.New(roleid),
	}
	dbplugin := db.Instance()
	if err := clubrole.wraper.Read(ctx, dbplugin); err != nil {
		return nil, err
	}
	return clubrole, nil
}

type ClubRole struct {
	wraper *db.Wraper[*pb.ClubRoleDB, pb.CRFlag]
	logp   log.LogPlus
}

func (cr *ClubRole) Obj() *pb.ClubRoleDB {
	return cr.wraper.Obj()
}

func (cr *ClubRole) Exit(ctx context.Context) error {
	cr.wraper.Set(pb.CRFlag_CRClubId, int64(0))
	cr.wraper.Update(ctx, db.Instance())
	return nil
}
