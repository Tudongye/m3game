package club

import (
	"context"
	"errors"
	"fmt"
	"m3game/demo/proto/pb"
	"m3game/meta/errs"
	"m3game/plugins/db"
	"m3game/plugins/log"
	"m3game/runtime/server/actor"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	ClubSlotNum   = 4
	ClubIdMetaKey = "clubslotid"
)

func SlotCreater(slotid string) actor.Actor {
	return &Slot{
		ActorBase: actor.ActorBaseCreator(slotid),
		clubs:     make(map[int64]*Club),
	}
}

func ConvertSlot(ctx context.Context) *Slot {
	a := actor.ParseActor(ctx)
	if a == nil {
		return nil
	}
	return a.(*Slot)
}

type Slot struct {
	*actor.ActorBase
	slotid int64
	clubs  map[int64]*Club
	logp   log.LogPlus
}

func (a *Slot) OnInit() error {
	log.InfoP(a.logp, "OnInit")
	// 加载slot数据
	if slotid, err := strconv.ParseInt(a.ID(), 10, 64); err != nil {
		return err
	} else {
		a.slotid = slotid
	}
	dbplugin := db.Instance()
	if clubids, err := dbplugin.ReadMany(context.TODO(), _clubdbmeta, bson.M{_clubdbmeta.FlagName(int32(pb.CFlag_CSlotId)): a.slotid}, int32(pb.CFlag_CClubId)); err != nil {
		return err
	} else {
		for _, clubid := range clubids {
			if club, err := newClub(context.TODO(), clubid.(*pb.ClubDB).ClubId); err != nil {
				continue
			} else {
				a.clubs[clubid.(*pb.ClubDB).ClubId] = club
			}
		}
	}
	return nil
}

func (a *Slot) OnTick() error {
	return nil
}

func (a *Slot) OnExit() error {
	log.InfoP(a.logp, "OnExit")
	return nil
}

func (a *Slot) OnSave() error {
	return nil
}

func (a *Slot) Club(clubid int64) *Club {
	return a.clubs[clubid]
}

func (a *Slot) CreateClub(ctx context.Context, clubid int64, roleid int64) error {
	dbplugin := db.Instance()
	crw := _clubrolewrapermeta.New(roleid)
	if err := crw.Read(ctx, dbplugin); err != nil {
		if errs.DBKeyNotFound.Is(err) {
			if err := crw.Create(ctx, dbplugin); err != nil {
				log.Error("%s", err.Error())
				return err
			}
		} else {
			log.Error("%s", err.Error())
			return err
		}
	}
	cw := _clubwrapermeta.New(clubid)
	cw.Set(pb.CFlag_COwnerId, roleid)
	if err := cw.Create(ctx, dbplugin); err != nil {
		return err
	}
	crw.Set(pb.CRFlag_CRClubId, clubid)
	if err := crw.Update(ctx, dbplugin); err != nil {
		return err
	}

	if club, err := newClub(ctx, clubid); err != nil {
		return err
	} else {
		a.clubs[clubid] = club
	}
	return nil
}

func (a *Slot) DeleteClub(ctx context.Context, clubid int64, roleid int64) error {
	club := a.clubs[clubid]
	if club == nil {
		return fmt.Errorf("Cant't find ClubId %d", clubid)
	}
	if roleid != club.wraper.Obj().OwnerId {
		return errors.New("Not Club Owner")
	}
	if len(club.clubroles) > 1 {
		return errors.New("Club Member > 1")
	}
	dbplugin := db.Instance()
	clubrole := club.clubroles[roleid]
	clubrole.wraper.Set(pb.CRFlag_CRClubId, int64(0))
	clubrole.wraper.Update(ctx, dbplugin)
	club.wraper.Delete(ctx, dbplugin)
	delete(a.clubs, clubid)
	return nil
}

func Club2SlotId(clubid int64) int64 {
	return clubid % ClubSlotNum
}
