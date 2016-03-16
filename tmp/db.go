package tmp

import "sync"

type DB struct {
	path   string
	stores map[string]*Store
	sync.RWMutex
}

func NewDB(path string) *DB {
	return &DB{
		path:   mkdirs(path),
		stores: make(map[string]*Store, 0),
	}
}

/*
 * store level methods
 */

func (db *DB) getstore(s string) (*Store, bool) {
	db.Lock()
	defer db.Unlock()
	st, ok := db.stores[s]
	if !ok {
		Log(StoreNotFound)
	}
	return st, ok
}

func (db *DB) HasStore(s string) bool {
	if _, ok := db.getstore(s); ok {
		return true
	}
	return false
}

func (db *DB) AddStore(s string) error {
	if st, ok := db.getstore(s); !ok {
		st = NewStore(db.path + s)
		return nil
	}
	return StoreExists
}

func (db *DB) DelStore(s string) error {
	if st, ok := db.getstore(s); ok {
		delete(db.stores, s)
		return nil
	}
	return StoreNotFound
}

/*
 *	document level methods
 */

func (db *DB) Has(s string, k []byte) bool {
	if st, ok := db.getstore(s); ok {
		return st.Has(k)
	}
	return false
}

func (db *DB) Add(s string, k, v []byte) error {
	if st, ok := db.getstore(s); ok {
		return st.Add(k, v)
	}
	return StoreNotFound
}

func (db *DB) Set(s string, k, v []byte) error {
	if st, ok := db.getstore(s); ok {
		return st.Set(k, v)
	}
	return StoreNotFound
}

func (db *DB) Get(s string, k []byte, v interface{}) error {
	if st, ok := db.getstore(s); ok {
		return st.Get(k, v)
	}
	return StoreNotFound
}

func (db *DB) Del(s string, k []byte) error {
	if st, ok := db.getstore(s); ok {
		return st.Del(k)
	}
	return StoreNotFound
}

func (db *DB) All(s string, v interface{}) error {
	if st, ok := db.getstore(s); ok {
		return st.All(v)
	}
	return StoreNotFound
}
