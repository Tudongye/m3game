package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"m3game/resource"
)

var (
	_               resource.ResLoader = (*LocationCfgLoader)(nil)
	_cfgname                           = "LocationCfg.json"
	LocationCfgName                    = "locationcfg"
)

func RegisterLocationCfg() {
	resource.RegisterResLoader(LocationCfgName, locationCfgLoaderCreater)
}

func GetLocationCfgLoader() *LocationCfgLoader {
	return resource.GetResource(LocationCfgName).(*LocationCfgLoader)
}

func locationCfgLoaderCreater() resource.ResLoader {
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

func (l *LocationCfgLoader) Load(cfgpath string) error {
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

func (l *LocationCfgLoader) Check(f resource.ResLoaderGetter) error {
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
