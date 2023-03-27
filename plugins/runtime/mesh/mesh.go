package mesh

import (
	"math/rand"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
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
		return errors.Wrap(err, "decode cfg")
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
		return "", errors.New("no avalible dst")
	}
	if r.appIdSet.Contains(dstappid) {
		return dstappid, nil
	}
	return "", errors.New("no avalible dst")
}

func (r *RouteHelper) RouteRandom() (string, error) {
	if r.size == 0 {
		return "", errors.New("no avalible dst")
	}
	return r.appIdslic[rand.Intn(r.size)], nil
}

func (r *RouteHelper) RouteHash(key string) (string, error) {
	if r.size == 0 {
		return "", errors.New("no avalible dst")
	}
	if dstappid, ok := r.hashRing.GetNode(key); !ok {
		return "", errors.New("no avalible dst")
	} else {
		return dstappid, nil
	}
}

func (r *RouteHelper) RouteSingle() (string, error) {
	if r.size == 0 {
		return "", errors.New("no avalible dst")
	}
	return r.minappid, nil
}
