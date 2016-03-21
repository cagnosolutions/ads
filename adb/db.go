package adb

import (
	"encoding/json"
	"log"
	"sync"
)

type DB struct {
	stores map[string]*Store
	sync.RWMutex
}

func NewDB() *DB {
	return &DB{
		stores: make(map[string]*Store, 0),
	}
}

func (db *DB) namespace(store string) (*Store, bool) {
	db.RLock()
	st, ok := db.stores[store]
	db.RUnlock()
	return st, ok
}

func (db *DB) AddStore(store string) {
	if _, ok := db.namespace(store); !ok {
		db.Lock()
		db.stores[store] = NewStore(store)
		db.Unlock()
	}
}

func (db *DB) DelStore(store string) {
	if _, ok := db.namespace(store); ok {
		db.Lock()
		delete(db.stores, store)
		db.Unlock()
	}
}

func (db *DB) Add(store, key string, val interface{}) bool {
	st, ok := db.namespace(store)
	if !ok {
		return false
	}
	b, err := json.Marshal(val)
	if err != nil {
		return logger(err)
	}
	err = st.Add([]byte(key), b)
	if err != nil {
		return logger(err)
	}
	return true
}

func (db *DB) Set(store, key string, val interface{}) bool {
	st, ok := db.namespace(store)
	if !ok {
		return false
	}
	b, err := json.Marshal(val)
	if err != nil {
		return logger(err)
	}
	err = st.Set([]byte(key), b)
	if err != nil {
		return logger(err)
	}
	return true
}

func (db *DB) Get(store, key string, ptr interface{}) bool {
	st, ok := db.namespace(store)
	if !ok {
		return false
	}
	b := st.Get([]byte(key))
	if err := json.Unmarshal(b, ptr); err != nil {
		return logger(err)
	}
	return true
}

func (db *DB) Del(store, key string) bool {
	st, ok := db.namespace(store)
	if !ok {
		return false
	}
	st.Del([]byte(key))
	return true
}

func (db *DB) All(store string, ptr interface{}) bool {
	st, ok := db.namespace(store)
	if !ok {
		return false
	}
	if err := json.Unmarshal(st.All(), ptr); err != nil {
		return logger(err)
	}
	return true
}

func logger(err error) bool {
	// log
	log.Printf("ERROR: %v\n", err)
	return false
}
