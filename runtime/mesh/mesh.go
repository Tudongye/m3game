package mesh

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
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
