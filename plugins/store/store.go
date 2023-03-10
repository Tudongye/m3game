package store

import (
	"sync"
	"time"
)

var (
	_store        Store
	_defaultstore = newDefaultStore()
)

func Get() Store {
	if _store == nil {
		return _defaultstore
	}
	return _store
}

func Set(s Store) {
	if _store != nil {
		panic("Store Only One")
	}
	_store = s
}

type Session interface {
	Get(k string) string
	Set(k string, v string)
	Name() string
}

type Store interface {
	NewSession(name string, livetimes int64) Session
	GetSession(name string) Session
	FreeSession(name string)
}

func NewSession(name string, livetimes int64) Session {
	return Get().NewSession(name, livetimes)
}

func FreeSession(name string) {
	Get().FreeSession(name)
}

func GetSession(name string) Session {
	return Get().GetSession(name)
}

func newDefaultStore() *defaultStore {
	return &defaultStore{}
}

type dsession struct {
	name    string
	kv      map[string]string
	timeout int64
}

func (s *dsession) Get(k string) string {
	return s.kv[k]
}

func (s *dsession) Set(k string, v string) {
	s.kv[k] = v
}

func (s *dsession) Name() string {
	return s.name
}

type defaultStore struct {
	sync.Map
}

func (s *defaultStore) NewSession(name string, livetimes int64) Session {
	session := &dsession{
		name:    name,
		kv:      make(map[string]string),
		timeout: time.Now().Unix() + livetimes,
	}
	a, _ := s.Map.LoadOrStore(name, session)
	ns := a.(*dsession)
	if ns.timeout < session.timeout {
		ns.timeout = session.timeout
		ns.kv = session.kv
	}
	return ns
}

func (s *defaultStore) GetSession(name string) Session {
	if v, ok := s.Map.Load(name); v != nil && ok {
		ns := v.(*dsession)
		if ns.timeout < time.Now().Unix() {
			s.Map.Delete(name)
			return ns
		}
		return ns
	}
	return nil
}

func (s *defaultStore) FreeSession(name string) {
	s.Map.Delete(name)
}
