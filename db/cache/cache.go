// a test cache-db
package cache

import (
	"fmt"
	"m3game/db"
	"m3game/runtime/plugin"
	"sync"

	"google.golang.org/protobuf/proto"
)

var (
	_         db.DB            = (*CacheDB)(nil)
	_         plugin.Factory   = (*Factory)(nil)
	_         plugin.PluginIns = (*CacheDB)(nil)
	_instance *CacheDB
)

const (
	_factoryname = "db_cache"
)

func Name() string {
	return _factoryname
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
	_instance = &CacheDB{
		cache: make(map[string][]byte),
	}
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

type CacheDB struct {
	cache map[string][]byte
	lock  sync.RWMutex
}

func (c *CacheDB) Name() string {
	return _factoryname
}

func (c *CacheDB) Get(meta *db.DBMeta, key string) (proto.Message, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	obj := meta.Creater()
	fieldname := genCacheKey(key, meta.Table, meta.Keyfield)
	if _, ok := c.cache[fieldname]; !ok {
		return nil, db.Err_DB_notfindkey
	}
	for _, field := range meta.Allfields {
		fieldname := genCacheKey(key, meta.Table, field)
		if v, ok := c.cache[fieldname]; ok {
			if err := meta.Setter(obj, field, v); err != nil {
				return nil, err
			}
		}
	}
	return obj, nil
}

func (c *CacheDB) Update(meta *db.DBMeta, key string, obj proto.Message) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	fieldname := genCacheKey(key, meta.Table, meta.Keyfield)
	if _, ok := c.cache[fieldname]; !ok {
		return db.Err_DB_notfindkey
	}
	for _, field := range meta.Allfields {
		if v, err := meta.Getter(obj, field); err != nil {
			return err
		} else {
			fieldname := genCacheKey(key, meta.Table, field)
			c.cache[fieldname] = v
		}
	}
	return nil
}
func (c *CacheDB) Insert(meta *db.DBMeta, key string, obj proto.Message) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	fieldname := genCacheKey(key, meta.Table, meta.Keyfield)
	if _, ok := c.cache[fieldname]; ok {
		return db.Err_DB_repeatedkey
	}
	for _, field := range meta.Allfields {
		if v, err := meta.Getter(obj, field); err != nil {
			return err
		} else {
			fieldname := genCacheKey(key, meta.Table, field)
			c.cache[fieldname] = v
		}
	}
	return nil
}
func (c *CacheDB) Delete(meta *db.DBMeta, key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	fieldname := genCacheKey(key, meta.Table, meta.Keyfield)
	if _, ok := c.cache[fieldname]; ok {
		return db.Err_DB_notfindkey
	}
	for _, field := range meta.Allfields {
		fieldname := genCacheKey(key, meta.Table, field)
		delete(c.cache, fieldname)
	}
	return nil
}

func genCacheKey(key string, table string, field string) string {
	return fmt.Sprintf("%s-%s-%s", key, table, field)
}
