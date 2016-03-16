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

func (db *DB) AddStore(s string) (*Store, error) {
	db.Lock()
	defer db.Unlock()
	if st, ok := db.stores[s]; ok {
		return st, StoreExists
	}
	return NewStore(db.path), nil
}

func (db *DB) GetStore(s string) (*Store, error) {
	db.Lock()
	defer db.Unlock()
	if st, ok := db.stores[s]; ok {
		return st, nil
	}
	return nil, StoreNotFound
}

func (db *DB) DelStore(s string) error {
	db.Lock()
	defer db.Unlock()
	if st, ok := db.stores[s]; ok {
		delete(db.stores, s)
		return nil
	}
	return StoreNotFound
}
