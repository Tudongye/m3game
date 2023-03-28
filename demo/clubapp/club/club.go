package club

import (
	"context"
	"errors"
	"m3game/demo/proto/pb"
	"m3game/plugins/db"
	"m3game/plugins/log"

	"go.mongodb.org/mongo-driver/bson"
)

func newClub(ctx context.Context, clubid int64) (*Club, error) {
	club := &Club{
		wraper:    _clubwrapermeta.New(clubid),
		clubroles: make(map[int64]*ClubRole),
	}
	dbplugin := db.Instance()
	if err := club.wraper.Read(ctx, dbplugin); err != nil {
		return nil, err
	}

	if roleids, err := dbplugin.ReadMany(context.TODO(), _clubroledbmeta, bson.M{_clubroledbmeta.FlagName(int32(pb.CRFlag_CRClubId)): clubid}, int32(pb.CRFlag_CRRoleId)); err != nil {
		return nil, err
	} else {
		for _, roldid := range roleids {
			if clubrole, err := newClubRole(context.TODO(), roldid.(*pb.ClubRoleDB).RoleId); err != nil {
				continue
			} else {
				club.clubroles[roldid.(*pb.ClubRoleDB).RoleId] = clubrole
			}
		}
	}
	return club, nil
}

type Club struct {
	wraper    *db.Wraper[*pb.ClubDB, pb.CFlag]
	clubroles map[int64]*ClubRole
	logp      log.LogPlus
}

func (c *Club) Obj() *pb.ClubDB {
	return c.wraper.Obj()
}

func (c *Club) ClubRole(roleid int64) *ClubRole {
	return c.clubroles[roleid]
}

func (c *Club) Join(ctx context.Context, roleid int64) error {
	crw := _clubrolewrapermeta.New(roleid)
	dbplugin := db.Instance()
	crw.Read(ctx, dbplugin)
	if crw.Obj().ClubId != 0 {
		return errors.New("")
	}
	crw.Set(pb.CRFlag_CRClubId, c.wraper.Obj().ClubId)
	crw.Update(ctx, dbplugin)
	if clubrole, err := newClubRole(context.TODO(), roleid); err != nil {
		return err
	} else {
		c.clubroles[roleid] = clubrole
	}
	return nil
}

func (c *Club) Exit(ctx context.Context, roleid int64) error {
	clubrole := c.clubroles[roleid]
	if clubrole == nil {
		return errors.New("")
	}
	if roleid == c.Obj().OwnerId {
		return errors.New("")
	}
	if err := clubrole.Exit(ctx); err != nil {
		return err
	}
	delete(c.clubroles, roleid)
	return nil
}
