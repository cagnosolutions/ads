package xdb

import (
	"os"
	"sync"
)

type engine struct {
	mp string // main path
	st map[string]*store
	//*log.Logger
	sync.RWMutex
}

// return a new db engine instance
func NewDB(path string) *engine {
	// check if path exists...
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path does not exists, so let's initialize
		if err := os.MkdirAll(path, 0755); err != nil {
			panic(err)
		}
	}
	return &engine{
		mp: path,
		st: make(map[string]*store),
	}
}

// returns a boolean displaying if the named
// store exists or not
func (e *engine) HasStore(name string) bool {
	e.RLock()
	_, ok := e.st[name]
	e.RUnlock()
	return ok
}

// add a new store if it doesn't exist...
// if it exists, return false
func (e *engine) AddStore(name string) bool {
	e.Lock()
	_, ok := e.st[name]
	if !ok {
		e.st[name] = NewStore(mp, name)
	}
	e.Unlock()
	return ok
}

// deletes a store by name; if the store
// doesn't exist, returns false
func (e *engine) DelStore(name string) bool {
	e.Lock()
	_, ok := e.st[name]
	delete(e.st, name)
	e.Unlock()
	return ok
}

// has returns weather the store/key exists
func (e *eingine) Has(st string, k []byte) bool {
	return true
}

// add inserts new data, returning false if it already exists
func (e *eingine) Add(st string, k, v []byte) bool {
	return true
}

// set volatiles' overrites/updates the data
func (e *eingine) Set(st string, k, v []byte) bool {
	return true
}

// returns single value matching given store and key
func (e *eingine) Get(st string, k []byte) []byte {
	return nil
}

// delets a record based on the store/key given...
// returning false if the store or key doesn't exist
func (e *eingine) Del(st string, k []byte) bool {
	return true
}

// returns all the data from a given store...
// as well as how many records it retreived
func (e *eingine) All(st string) ([][]byte, int) {
	return true
}
