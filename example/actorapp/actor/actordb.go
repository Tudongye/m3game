package actor

import (
	"fmt"
	"m3game/example/proto/pb"
	"m3game/plugins/db"
	"m3game/plugins/db/wraper"
	"m3game/plugins/log"
)

var (
	actormeta *db.DBMeta[*pb.ActorDB]
)

func init() {
	actormeta = db.NewMeta("actor_table", actorDBCreater)
}

func actorDBCreater() *pb.ActorDB {
	return &pb.ActorDB{
		ActorID:   "",
		ActorName: &pb.ActorName{},
		ActorInfo: &pb.ActorInfo{},
	}
}

func Register(name string) (string, error) {
	actorid := fmt.Sprintf("ActorID-%s", name)
	dbplugin := db.Get()
	if dbplugin == nil {
		return "", _err_actor_dbplugin
	}
	w := wraper.New(actormeta, actorid)
	if err := wraper.KeySetter(w, actorid); err != nil {
		log.Error(err.Error())
		return "", _err_actor_dbplugin
	}
	if err := wraper.Setter(w, &pb.ActorName{
		Name: name,
	}); err != nil {
		log.Error(err.Error())
		return "", _err_actor_dbplugin
	}
	if err := wraper.Setter(w, &pb.ActorInfo{
		Level: 0,
	}); err != nil {
		log.Error(err.Error())
		return "", _err_actor_dbplugin
	}
	return actorid, w.Create(dbplugin)
}
