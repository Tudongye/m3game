// a test cache-db
package redis

import (
	"fmt"
	"m3game/plugins/db"
	"m3game/plugins/log"
	"m3game/runtime/plugin"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/mitchellh/mapstructure"
	"google.golang.org/protobuf/proto"
)

var (
	_         db.DB            = (*RedisDB)(nil)
	_         plugin.Factory   = (*Factory)(nil)
	_         plugin.PluginIns = (*RedisDB)(nil)
	_cfg                       = RedisCfg{}
	_instance *RedisDB
	_factory  = &Factory{}
)

const (
	_factoryname = "db_redis"
)

func init() {
	plugin.RegisterFactory(_factory)
}

type RedisCfg struct {
	Host      string `mapstructure:"Host"`
	Port      int    `mapstructure:"Port"`
	Auth      string `mapstructure:"Auth"`
	MaxIdle   int    `mapstructure:"MaxIdle"`
	MaxActive int    `mapstructure:"MaxActive"`
}

func (c RedisCfg) CheckVaild() error {
	if c.Host == "" {
		return errors.New("Host cant be space")
	}
	if c.Port == 0 {
		return errors.New("Port cant be 0")
	}
	if c.MaxIdle == 0 {
		return errors.New("MaxIdle cant be 0")
	}
	if c.MaxActive == 0 {
		return errors.New("MaxActive cant be 0")
	}
	return nil
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.DB
}
func (f *Factory) Name() string {
	return _factoryname
}

func (f *Factory) Setup(c map[string]interface{}) (plugin.PluginIns, error) {
	if _instance != nil {
		return _instance, nil
	}
	if err := mapstructure.Decode(c, &_cfg); err != nil {
		return nil, errors.Wrap(err, "RedisDB Decode Cfg")
	}
	if err := _cfg.CheckVaild(); err != nil {
		return nil, err
	}
	_instance = &RedisDB{pool: &redis.Pool{
		MaxIdle:   _cfg.MaxIdle,
		MaxActive: _cfg.MaxActive,
		Dial: func() (redis.Conn, error) {
			addrStr := fmt.Sprintf("%s:%d", _cfg.Host, _cfg.Port)
			c, err := redis.Dial("tcp", addrStr)
			if err != nil {
				log.Fatal("error:%s", err.Error())
			}
			if _cfg.Auth != "" {
				if _, err := c.Do("AUTH", _cfg.Auth); err != nil {
					c.Close()
					log.Fatal("AUTH error:%s", err.Error())
				}
			}
			return c, err
		},
	}}
	db.Set(_instance)
	log.Info("RedisDB...")
	return _instance, nil
}

func (f *Factory) Destroy(plugin.PluginIns) error {
	return nil
}

func (f *Factory) Reload(plugin.PluginIns, map[string]interface{}) error {
	return nil
}

func (f *Factory) CanDelete(plugin.PluginIns) bool {
	return false
}

type RedisDB struct {
	pool *redis.Pool
}

func (c *RedisDB) Factory() plugin.Factory {
	return _factory
}

func (c *RedisDB) Read(meta db.DBMetaInter, key string, filters ...string) (proto.Message, error) {
	obj := meta.Creater()()
	fieldname := genCacheKey(key, meta.Table(), meta.KeyField())
	rc := c.pool.Get()
	defer rc.Close()
	if v, err := rc.Do("GET", fieldname); err == redis.ErrNil || v == nil {
		return nil, db.Err_DB_notfindkey
	} else if err != nil {
		return nil, err
	}

	var fields []string
	if len(filters) == 0 {
		fields = meta.AllFields()
	} else {
		for _, field := range filters {
			if !meta.HasField(field) {
				log.Error("Obj %s not have field %s in filter", meta.ObjName, field)
			}
			fields = append(fields, field)
		}
	}
	for _, field := range fields {
		fieldname := genCacheKey(key, meta.Table(), field)
		if v, err := redis.Bytes(rc.Do("GET", fieldname)); err != nil {
			log.Error(err.Error())
			continue
		} else if err := meta.Decode(obj, field, v); err != nil {
			return nil, err
		}
	}
	return obj, nil
}

func (c *RedisDB) Update(meta db.DBMetaInter, key string, obj proto.Message, filters ...string) error {
	fieldname := genCacheKey(key, meta.Table(), meta.KeyField())
	rc := c.pool.Get()
	defer rc.Close()
	if v, err := rc.Do("GET", fieldname); err == redis.ErrNil || v == nil {
		return db.Err_DB_notfindkey
	} else if err != nil {
		return err
	}
	var fields []string
	if len(filters) == 0 {
		fields = meta.AllFields()
	} else {
		for _, field := range filters {
			if !meta.HasField(field) {
				log.Error("Obj %s not have field %s in filter", meta.ObjName, field)
			}
			fields = append(fields, field)
		}
	}
	for _, field := range fields {
		if v, err := meta.Encode(obj, field); err != nil {
			return err
		} else {
			fieldname := genCacheKey(key, meta.Table(), field)
			if _, err := rc.Do("Set", fieldname, v); err != nil {
				log.Error("Redis Set %s = %s Fail %s", fieldname, v, err.Error())
			}
		}
	}
	return nil
}
func (c *RedisDB) Create(meta db.DBMetaInter, key string, obj proto.Message, filters ...string) error {
	fieldname := genCacheKey(key, meta.Table(), meta.KeyField())
	rc := c.pool.Get()
	defer rc.Close()
	if v, err := rc.Do("GET", fieldname); v != nil && err != redis.ErrNil {
		return db.Err_DB_repeatedkey
	}

	var fields []string
	if len(filters) == 0 {
		fields = meta.AllFields()
	} else {
		for _, field := range filters {
			if !meta.HasField(field) {
				log.Error("Obj %s not have field %s in filter", meta.ObjName, field)
			}
			fields = append(fields, field)
		}
	}

	for _, field := range fields {
		if v, err := meta.Encode(obj, field); err != nil {
			return err
		} else {
			fieldname := genCacheKey(key, meta.Table(), field)
			if _, err := rc.Do("Set", fieldname, v); err != nil {
				log.Error("Redis Set %s = %s Fail %s", fieldname, v, err.Error())
			}
		}
	}
	return nil
}
func (c *RedisDB) Delete(meta db.DBMetaInter, key string) error {
	fieldname := genCacheKey(key, meta.Table(), meta.KeyField())
	rc := c.pool.Get()
	defer rc.Close()
	if v, err := rc.Do("GET", fieldname); err == redis.ErrNil || v == nil {
		return db.Err_DB_notfindkey
	} else if err != nil {
		return err
	}

	for _, field := range meta.AllFields() {
		fieldname := genCacheKey(key, meta.Table(), field)
		if _, err := rc.Do("Del", fieldname); err != nil {
			log.Error("Redis Del %s  Fail %s", fieldname, err.Error())
		}
	}
	return nil
}

func genCacheKey(key string, table string, field string) string {
	return fmt.Sprintf("%s-%s-%s", key, table, field)
}
