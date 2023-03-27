// a test cache-db
package cache

import (
	"context"
	"fmt"
	"m3game/plugins/db"
	"m3game/plugins/log"
	"m3game/runtime/plugin"
	"sync"

	"google.golang.org/protobuf/proto"
)

var (
	_        db.DB            = (*CacheDB)(nil)
	_        plugin.Factory   = (*Factory)(nil)
	_        plugin.PluginIns = (*CacheDB)(nil)
	_cachedb *CacheDB
	_factory = &Factory{}
	_codec   = &db.CacheDBCodec{}
)

const (
	_name = "db_cache"
)

func init() {
	plugin.RegisterFactory(_factory)
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
	if _cachedb != nil {
		return _cachedb, nil
	}
	_cachedb = &CacheDB{
		cache: make(map[string]interface{}),
	}
	if _, err := db.New(_cachedb); err != nil {
		return nil, err
	}
	return _cachedb, nil
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

type CacheDB struct {
	cache map[string]interface{}
	lock  sync.RWMutex
}

func (c *CacheDB) Factory() plugin.Factory {
	return _factory
}

func (c *CacheDB) Read(ctx context.Context, meta db.DBMetaInter, key interface{}, flags ...int32) (proto.Message, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	obj := meta.New()
	fieldname := genCacheKey(key, meta.Table(), meta.FlagName(meta.KeyFlag()))
	log.Debug("Read %v", fieldname)
	if _, ok := c.cache[fieldname]; !ok {
		return nil, db.Err_KeyNotFound
	}
	if len(flags) == 0 {
		flags = meta.AllFlags()
	}
	for _, flag := range flags {
		fieldname := genCacheKey(key, meta.Table(), meta.FlagName(flag))
		if v, ok := c.cache[fieldname]; ok {
			value, _ := _codec.Decode(meta.FlagKind(flag), v)
			log.Debug("Read %v", value)
			meta.Setter(obj, flag, value)
		}
	}
	return obj, nil
}

func (c *CacheDB) Update(ctx context.Context, meta db.DBMetaInter, key interface{}, obj proto.Message, flags ...int32) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	fieldname := genCacheKey(key, meta.Table(), meta.FlagName(meta.KeyFlag()))
	log.Debug("Read %v", fieldname)
	if _, ok := c.cache[fieldname]; !ok {
		return db.Err_KeyNotFound
	}
	if len(flags) == 0 {
		flags = meta.AllFlags()
	}
	for _, flag := range flags {
		fieldname := genCacheKey(key, meta.Table(), meta.FlagName(flag))
		v := meta.Getter(obj, flag)
		value, _ := _codec.Encode(meta.FlagKind(flag), v)
		c.cache[fieldname] = value
	}
	return nil
}

func (c *CacheDB) Create(ctx context.Context, meta db.DBMetaInter, key interface{}, obj proto.Message) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	fieldname := genCacheKey(key, meta.Table(), meta.FlagName(meta.KeyFlag()))
	if _, ok := c.cache[fieldname]; ok {
		return db.Err_DuplicateEntry
	}
	flags := meta.AllFlags()
	log.Debug("%v %v %v", key, flags, meta.AllFlags())
	for _, flag := range flags {
		fieldname := genCacheKey(key, meta.Table(), meta.FlagName(flag))
		v := meta.Getter(obj, flag)
		value, _ := _codec.Encode(meta.FlagKind(flag), v)
		log.Debug("Set %v %v", fieldname, value)
		c.cache[fieldname] = value
	}
	return nil
}

func (c *CacheDB) Delete(ctx context.Context, meta db.DBMetaInter, key interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	fieldname := genCacheKey(key, meta.Table(), meta.FlagName(meta.KeyFlag()))
	if _, ok := c.cache[fieldname]; !ok {
		return db.Err_KeyNotFound
	}
	for _, flag := range meta.AllFlags() {
		fieldname := genCacheKey(key, meta.Table(), meta.FlagName(flag))
		delete(c.cache, fieldname)
	}
	return nil
}

func (c *CacheDB) ReadMany(ctx context.Context, meta db.DBMetaInter, filters interface{}, fields ...int32) ([]proto.Message, error) {
	return nil, nil
}

func genCacheKey(key interface{}, table string, field string) string {
	return fmt.Sprintf("%v-%s-%s", key, table, field)
}
