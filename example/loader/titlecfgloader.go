package loader

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"m3game/runtime/resource"
)

var (
	_             resource.ResLoader = (*TitleCfgLoader)(nil)
	_cfgname                         = "titlecfg.json"
	_titlecfgflag                    = "titlecfg"
)

func RegisterTitleCfg() {
	resource.RegisterResource(&titleCfgLoaderCreater{})
}

type titleCfgLoaderCreater struct {
}

func (l *titleCfgLoaderCreater) Name() string {
	return _titlecfgflag
}
func (l *titleCfgLoaderCreater) NewLoader() resource.ResLoader {
	return &TitleCfgLoader{
		cfgs: make(map[int32]TitleCfg),
	}
}

type TitleCfg struct {
	Level int32  `json:"Level"`
	Title string `json:"Title"`
}

type TitleCfgFile struct {
	TitleCfgs []TitleCfg `json:"TitleCfgs"`
}

type TitleCfgLoader struct {
	cfgs map[int32]TitleCfg
}

func (l *TitleCfgLoader) Name() string {
	return _titlecfgflag
}

func (l *TitleCfgLoader) Load(ctx context.Context, cfgpath string) error {
	if data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", cfgpath, _cfgname)); err != nil {
		return err
	} else {
		var cfgfile TitleCfgFile
		if err := json.Unmarshal(data, &cfgfile); err != nil {
			return err
		}
		for _, cfg := range cfgfile.TitleCfgs {
			if _, ok := l.cfgs[cfg.Level]; ok {
				return fmt.Errorf("TitleCfgs Level %d repeated", cfg.Level)
			}
			l.cfgs[cfg.Level] = cfg
		}
	}
	return nil
}

func (l *TitleCfgLoader) GetTitleByLv(level int32) string {
	maxdis := -1
	name := "start"
	for _, cfg := range l.cfgs {
		if level >= cfg.Level && maxdis < int(cfg.Level) {
			maxdis = int(cfg.Level)
			name = cfg.Title
		}
	}
	return name
}
