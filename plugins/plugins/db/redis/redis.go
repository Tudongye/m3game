// a test cache-db
package redis

import (
	"context"
	"fmt"
	"m3game/plugins/db"
	"m3game/plugins/log"
	"m3game/runtime/plugin"

	"github.com/go-playground/validator/v10"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/mitchellh/mapstructure"
	"google.golang.org/protobuf/proto"
)

var (
	_        db.DB            = (*RedisDB)(nil)
	_        plugin.Factory   = (*Factory)(nil)
	_        plugin.PluginIns = (*RedisDB)(nil)
	_redisdb *RedisDB
	_factory = &Factory{}
	_codec   = &db.CacheDBCodec{}
)

const (
	_name = "db_redis"
)

func init() {
	plugin.RegisterFactory(_factory)
}

type RedisCfg struct {
	Host      string `mapstructure:"Host" validate:"required"`
	Port      int    `mapstructure:"Port" validate:"gt=0"`
	Auth      string `mapstructure:"Auth"`
	MaxIdle   int    `mapstructure:"MaxIdle" validate:"gt=0"`
	MaxActive int    `mapstructure:"MaxActive" validate:"gt=0"`
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.DB
}

func (f *Factory) Name() string {
	return _name
}

func (f *Factory) Setup(ctx context.Context, c map[string]interface{}) (plugin.PluginIns, error) {
	if _redisdb != nil {
		return _redisdb, nil
	}
	var cfg RedisCfg
	if err := mapstructure.Decode(c, &cfg); err != nil {
		return nil, errors.Wrap(err, "RedisDB Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, err
	}
	_redisdb = &RedisDB{
		cfg: cfg,
		pool: &redis.Pool{
			MaxIdle:   cfg.MaxIdle,
			MaxActive: cfg.MaxActive,
			Dial: func() (redis.Conn, error) {
				addrStr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
				c, err := redis.Dial("tcp", addrStr)
				if err != nil {
					log.Fatal("error:%s", err.Error())
				}
				if cfg.Auth != "" {
					if _, err := c.Do("AUTH", cfg.Auth); err != nil {
						c.Close()
						log.Fatal("AUTH error:%s", err.Error())
					}
				}
				return c, err
			},
		}}
	if _, err := db.New(_redisdb); err != nil {
		return nil, err
	}
	log.Info("RedisDB...")
	return _redisdb, nil
}

func (f *Factory) Destroy(plugin.PluginIns) error {
	return nil
}

func (f *Factory) Reload(plugin.PluginIns, map[string]interface{}) error {
	return nil
}

func (f *Factory) CanUnload(plugin.PluginIns) bool {
	return false
}

type RedisDB struct {
	cfg  RedisCfg
	pool *redis.Pool
}

func (c *RedisDB) Factory() plugin.Factory {
	return _factory
}

func (c *RedisDB) Read(ctx context.Context, meta db.DBMetaInter, key interface{}, flags ...int32) (proto.Message, error) {
	obj := meta.New()
	fieldname := genCacheKey(key, meta.Table(), meta.FlagName(meta.KeyFlag()))
	rc := c.pool.Get()
	defer rc.Close()
	if v, err := rc.Do("GET", fieldname); err == redis.ErrNil || v == nil || len(v.([]byte)) == 0 {
		return nil, db.Err_KeyNotFound
	} else if err != nil {
		return nil, err
	}
	if len(flags) == 0 {
		flags = meta.AllFlags()
	}
	var args []interface{}
	for _, flag := range flags {
		fieldname := genCacheKey(key, meta.Table(), meta.FlagName(flag))
		args = append(args, fieldname)
	}
	values, err := redis.Values(rc.Do("MGET", redis.Args{}.AddFlat(args)...))
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	for i, v := range values {
		flag := flags[i]
		value, _ := _codec.Decode(meta.FlagKind(flag), v)
		meta.Setter(obj, flag, value)
	}
	return obj, nil
}

func (c *RedisDB) Update(ctx context.Context, meta db.DBMetaInter, key interface{}, obj proto.Message, flags ...int32) error {
	fieldname := genCacheKey(key, meta.Table(), meta.FlagName(meta.KeyFlag()))
	rc := c.pool.Get()
	defer rc.Close()
	if v, err := rc.Do("GET", fieldname); err == redis.ErrNil || v == nil {
		return db.Err_KeyNotFound
	} else if err != nil {
		return err
	}
	if len(flags) == 0 {
		flags = meta.AllFlags()
	}
	var args []interface{}
	for _, flag := range flags {
		fieldname := genCacheKey(key, meta.Table(), meta.FlagName(flag))
		v := meta.Getter(obj, flag)
		value, _ := _codec.Encode(meta.FlagKind(flag), v)
		args = append(args, fieldname, value)
	}
	_, err := rc.Do("MSET", args...)
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisDB) Create(ctx context.Context, meta db.DBMetaInter, key interface{}, obj proto.Message) error {
	fieldname := genCacheKey(key, meta.Table(), meta.FlagName(meta.KeyFlag()))
	rc := c.pool.Get()
	defer rc.Close()
	if v, err := rc.Do("GET", fieldname); v != nil && err != redis.ErrNil {
		return db.Err_DuplicateEntry
	}
	flags := meta.AllFlags()
	var args []interface{}
	for _, flag := range flags {
		fieldname := genCacheKey(key, meta.Table(), meta.FlagName(flag))
		v := meta.Getter(obj, flag)
		value, _ := _codec.Encode(meta.FlagKind(flag), v)
		args = append(args, fieldname, value)
	}
	_, err := rc.Do("MSET", args...)
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisDB) Delete(ctx context.Context, meta db.DBMetaInter, key interface{}) error {
	fieldname := genCacheKey(key, meta.Table(), meta.FlagName(meta.KeyFlag()))
	rc := c.pool.Get()
	defer rc.Close()
	if v, err := rc.Do("GET", fieldname); err == redis.ErrNil || v == nil {
		return db.Err_KeyNotFound
	} else if err != nil {
		return err
	}

	var args []interface{}
	for _, flag := range meta.AllFlags() {
		fieldname := genCacheKey(key, meta.Table(), meta.FlagName(flag))
		args = append(args, fieldname)
	}
	count, err := redis.Int(rc.Do("DEL", args...))
	if err != nil {
		return err
	}
	if count != len(args) {
		return fmt.Errorf("Del %d but want %d", count, len(args))
	}
	return nil
}

func (c *RedisDB) ReadMany(ctx context.Context, meta db.DBMetaInter, filters interface{}, flags ...int32) ([]proto.Message, error) {
	return nil, nil
}

func genCacheKey(key interface{}, table string, field string) string {
	return fmt.Sprintf("%v-%s-%s", key, table, field)
}
