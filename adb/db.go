package adb

import "sync"

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
	if err := st.Add(key, val); err != nil {
		return false
	}
	return true
}

func (db *DB) Set(store, key string, val interface{}) bool {
	st, ok := db.namespace(store)
	if !ok {
		return false
	}
	if err := st.Set(key, val); err != nil {
		return false
	}
	return true
}

func (db *DB) Get(store, key string, ptr interface{}) bool {
	st, ok := db.namespace(store)
	if !ok {
		return false
	}
	if err := st.Get(key, ptr); err != nil {
		return false
	}
	return true
}

func (db *DB) All(store string, ptr interface{}) bool {
	st, ok := db.namespace(store)
	if !ok {
		return false
	}
	if err := st.All(ptr); err != nil {
		return false
	}
	return true
}

func (db *DB) Del(store, key string) bool {
	st, ok := db.namespace(store)
	if !ok {
		return false
	}
	st.Del(key)
	return true
}

func (db *DB) Match(store, qry string, ptr interface{}) bool {
	st, ok := db.namespace(store)
	if !ok {
		return false
	}
	if err := st.Match(qry, ptr); err != nil {
		return false
	}
	return true
}

func (db *DB) Close() {
	db.Lock()
	defer db.Unlock()
	for _, st := range db.stores {
		st.index.Close()
	}
}
