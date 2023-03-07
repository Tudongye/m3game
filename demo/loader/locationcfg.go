package loader

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"m3game/runtime/resource"
)

var (
	_                resource.ResLoader = (*LocationCfgLoader)(nil)
	_cfgname                            = "LocationCfg.json"
	_locationcfgflag                    = "locationcfg"
)

func RegisterLocationCfg() {
	resource.RegisterResource(&locationCfgLoaderCreater{})
}

type locationCfgLoaderCreater struct {
}

func (l *locationCfgLoaderCreater) Name() string {
	return _locationcfgflag
}
func (l *locationCfgLoaderCreater) NewLoader() resource.ResLoader {
	return &LocationCfgLoader{
		cfgs: make(map[int32]LocationCfg),
	}
}

type LocationCfg struct {
	Distance   int32  `json:"Distance"`
	LocateName string `json:"LocateName"`
}

type LocationCfgFile struct {
	LocationCfgs []LocationCfg `json:"LocationCfgs"`
}

type LocationCfgLoader struct {
	cfgs map[int32]LocationCfg
}

func (l *LocationCfgLoader) Name() string {
	return _locationcfgflag
}

func (l *LocationCfgLoader) Load(ctx context.Context, cfgpath string) error {
	if data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", cfgpath, _cfgname)); err != nil {
		return err
	} else {
		var cfgfile LocationCfgFile
		if err := json.Unmarshal(data, &cfgfile); err != nil {
			return err
		}
		for _, cfg := range cfgfile.LocationCfgs {
			if _, ok := l.cfgs[cfg.Distance]; ok {
				return fmt.Errorf("LocationCfg Distance %d repeated", cfg.Distance)
			}
			l.cfgs[cfg.Distance] = cfg
		}
	}
	return nil
}

func (l *LocationCfgLoader) GetNameByDistance(distance int32) string {
	maxdis := -1
	name := "start"
	for _, locationcfg := range l.cfgs {
		if distance >= locationcfg.Distance && maxdis < int(locationcfg.Distance) {
			maxdis = int(locationcfg.Distance)
			name = locationcfg.LocateName
		}
	}
	return name
}
