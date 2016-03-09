package xdb

import (
	"bytes"
	"encoding/json"
	"os"
	"sync"
)

type engine struct {
	p string
	s map[string]*store
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
		p: path,
		s: make(map[string]*store),
	}
}

// returns boolean displayiny if store exists
func (e *engine) HasStore(s string) bool {
	e.RLock()
	_, ok := e.s[s]
	e.RUnlock()
	return ok
}

// add store if it doesn't exist, else return false
func (e *engine) AddStore(s string) bool {
	e.Lock()
	_, ok := e.s[s]
	if !ok {
		e.s[s] = &store{sid: s}
	}
	e.Unlock()
	return ok
}

// deletes store if it exists, else return false
func (e *engine) DelStore(s string) bool {
	e.Lock()
	_, ok := e.s[s]
	delete(e.s, s)
	e.Unlock()
	return ok
}

// return store along with boolean if it exists (FOR INTERNAL USE ONLY)
func (e *engine) namespace(s string) (*store, bool) {
	e.RLock()
	defer e.RUnlock()
	st, ok := e.s[s]
	return st, ok
}

// returns boolean displaying the existance of the store/key
func (e *engine) Has(s string, k []byte) bool {
	if st, ok := e.namespace(s); ok {
		return st.has(k)
	}
	return false
}

// inserts new data, returning false if it already exists
func (e *engine) Add(s string, k []byte, v interface{}) bool {
	if st, ok := e.namespace(s); ok {
		if b, ok := encode(v); ok {
			return st.add(k, b)
		}
	}
	return false
}

// volatiles' overrites/updates the data
func (e *engine) Set(s string, k []byte, v interface{}) bool {
	if st, ok := e.namespace(s); ok {
		if b, ok := encode(v); ok {
			return st.set(k, b)
		}
	}
	return false
}

// returns single value matching given store/key
func (e *engine) Get(s string, k []byte, v interface{}) bool {
	if st, ok := e.namespace(s); ok {
		if b := st.get(k); b != nil {
			return decode(b, v)
		}
	}
	return false
}

// delets record using store/key given, else returns false
func (e *engine) Del(s string, k []byte) bool {
	if st, ok := e.namespace(s); ok {
		return st.del(k)
	}
	return false
}

// returns all the data from a given store
func (e *engine) All(s string, v interface{}) bool {
	if st, ok := e.namespace(s); ok {
		if bs := st.all(); bs != nil {
			b := bytes.Join(bs, []byte{','})
			return decode(b, v)
		}
	}
	return false
}

// obj to byte encoder
func encode(v interface{}) ([]byte, bool) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, false
	}
	return b, true
}

// byte to obj decoder
func decode(b []byte, v interface{}) bool {
	err := json.Unmarshal(b, v)
	if err != nil {
		return false
	}
	return true
}
