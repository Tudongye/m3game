package mesh

import (
	"m3game/meta/errs"
	"math/rand"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/serialx/hashring"
)

var (
	_cfg MeshCfg
)

type MeshCfg struct {
	WatcherInterSecond int `mapstructure:"WatcherInterSecond"  validate:"gt=0"`
}

func Init(c map[string]interface{}) error {
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return errs.MeshInitFail.Wrap(err, "Decode MeshCfg")
	}
	validate := validator.New()
	if err := validate.Struct(&_cfg); err != nil {
		return err
	}
	return nil
}

func NewRouteHelper() *RouteHelper {
	return &RouteHelper{
		appIdSet: mapset.NewSet[string](),
	}
}

type RouteHelper struct {
	appIdSet  mapset.Set[string]
	hashRing  *hashring.HashRing
	minappid  string
	size      int
	appIdslic []string
}

func (r *RouteHelper) Add(appid string) {
	r.appIdSet.Add(appid)
}

func (r *RouteHelper) Compress() {
	r.hashRing = hashring.New(r.appIdSet.ToSlice())
	for _, appid := range r.appIdSet.ToSlice() {
		if r.minappid == "" {
			r.minappid = appid
		}
		if appid < r.minappid {
			r.minappid = appid
		}
	}
	r.size = r.appIdSet.Cardinality()
	r.appIdslic = r.appIdSet.ToSlice()
}

func (r *RouteHelper) RouteP2P(dstappid string) (string, error) {
	if r.size == 0 {
		return "", errs.MeshNoAvalibleDstApp.New("RouteP2P No Avalible For %s", dstappid)
	}
	if r.appIdSet.Contains(dstappid) {
		return dstappid, nil
	}
	return "", errs.MeshNoAvalibleDstApp.New("RouteP2P No DstApp For %s", dstappid)
}

func (r *RouteHelper) RouteRandom() (string, error) {
	if r.size == 0 {
		return "", errs.MeshNoAvalibleDstApp.New("RouteRandom No Avalible DstApp ")
	}
	return r.appIdslic[rand.Intn(r.size)], nil
}

func (r *RouteHelper) RouteHash(key string) (string, error) {
	if r.size == 0 {
		return "", errs.MeshNoAvalibleDstApp.New("RouteHash No Avalible DstApp ")
	}
	if dstappid, ok := r.hashRing.GetNode(key); !ok {
		return "", errs.MeshNoAvalibleDstApp.New("RouteHash No Avalible DstApp For HashKey %s", key)
	} else {
		return dstappid, nil
	}
}

func (r *RouteHelper) RouteSingle() (string, error) {
	if r.size == 0 {
		return "", errs.MeshNoAvalibleDstApp.New("RouteSingle No Avalible DstApp ")
	}
	return r.minappid, nil
}
