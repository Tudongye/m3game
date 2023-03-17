package mesh

import (
	"fmt"
	"math/rand"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/serialx/hashring"
)

var (
	_cfg MeshCfg
)

type MeshCfg struct {
	WatcherInterSecond uint32 `mapstructure:"WatcherInterSecond"`
}

func (c *MeshCfg) checkvaild() error {
	if c.WatcherInterSecond <= 0 {
		return fmt.Errorf("WatcherInterSecond %d invaild", c.WatcherInterSecond)
	}
	return nil
}

func Init(c map[string]interface{}) error {
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return errors.Wrap(err, "decode cfg")
	}
	if err := _cfg.checkvaild(); err != nil {
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
	appIdSet mapset.Set[string]
	hashRing *hashring.HashRing
	minappid string
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
}

func (r *RouteHelper) Size() int {
	return r.appIdSet.Cardinality()
}

func (r *RouteHelper) RouteP2P(dstappid string) (string, error) {
	if r.Size() == 0 {
		return "", errors.New("no avalible dst")
	}
	if r.appIdSet.Contains(dstappid) {
		return dstappid, nil
	}
	return "", errors.New("no avalible dst")
}

func (r *RouteHelper) RouteRandom() (string, error) {
	if r.Size() == 0 {
		return "", errors.New("no avalible dst")
	}
	return r.appIdSet.ToSlice()[rand.Int()%r.Size()], nil
}

func (r *RouteHelper) RouteHash(key string) (string, error) {
	if r.Size() == 0 {
		return "", errors.New("no avalible dst")
	}
	if dstappid, ok := r.hashRing.GetNode(key); !ok {
		return "", errors.New("no avalible dst")
	} else {
		return dstappid, nil
	}
}

func (r *RouteHelper) RouteSingle() (string, error) {
	if r.Size() == 0 {
		return "", errors.New("no avalible dst")
	}
	return r.minappid, nil
}
