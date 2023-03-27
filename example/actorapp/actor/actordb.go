package actor

import (
	"context"
	"m3game/example/proto/pb"
	"m3game/plugins/db"
	"m3game/plugins/log"
)

var (
	actordbmeta      *db.DBMeta[*pb.ActorDB]
	actorwrapermeata *db.WraperMeta[*pb.ActorDB, pb.AcFlag]
)

func init() {
	actordbmeta = db.NewMeta[*pb.ActorDB]("actor_table")
	actorwrapermeata = db.NewWraperMeta[*pb.ActorDB, pb.AcFlag](actordbmeta)
}

func Register(ctx context.Context, playerid string, name string) (string, error) {
	dbplugin := db.Instance()
	if dbplugin == nil {
		return "", _err_actor_dbplugin
	}
	log.Debug(playerid)
	w := actorwrapermeata.New(playerid)
	w.Set(pb.AcFlag_FActorID, playerid)
	w.Set(pb.AcFlag_FActorName, name)
	w.Set(pb.AcFlag_FActorLevel, int32(0))
	return playerid, w.Create(ctx, dbplugin)
}
