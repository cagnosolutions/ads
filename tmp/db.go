package tmp

import (
	"encoding/json"
	"reflect"
	"sync"
)

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
	if _, ok := db.getstore(s); !ok {
		db.Lock()
		db.stores[s] = NewStore(db.path, s)
		db.Unlock()
		return nil
	}
	return StoreExists
}

func (db *DB) DelStore(s string) error {
	if _, ok := db.getstore(s); ok {
		delete(db.stores, s)
		return nil
	}
	return StoreNotFound
}

func (db *DB) Has(s string, k []byte) bool {
	if st, ok := db.getstore(s); ok {
		return st.Has(k)
	}
	return false
}

func (db *DB) Add(s string, k []byte, v interface{}) error {
	if st, ok := db.getstore(s); ok {
		doc, err := document(k, v)
		if err != nil {
			return err
		}
		return st.Add(k, doc)
	}
	return StoreNotFound
}

// create and check document based on the provided key and value
func document(k []byte, v interface{}) ([]byte, error) {
	// make sure we received a pointer value
	//if reflect.ValueOf(v).Kind() != reflect.Ptr {
	//	return NonPtrValue
	//}
	// marshal value into json data...
	var doc []byte
	if doc, err := json.Marshal(v); err != nil {
		return err
	}
	// append delim and json data to key
	k = append(k, '|')
	k = append(k, doc...)
	// check k (which is now full doc) in order
	// to ensure that it's not over the max size

	//
}

func (db *DB) Set(s string, k []byte, v interface{}) error {
	if st, ok := db.getstore(s); ok {
		return st.Set(k, v)
	}
	return StoreNotFound
}

func (db *DB) Get(s string, k []byte, v interface{}) error {
	if st, ok := db.getstore(s); ok {
		if !isPtr(v) {
			return NonPtrValue
		}
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
		if !isPtr(v) {
			return NonPtrValue
		}
		return st.All(v)
	}
	return StoreNotFound
}

func isPtr(v interface{}) bool {
	return reflect.ValueOf(v).Kind() == reflect.Ptr
}
