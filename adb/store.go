package adb

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"
)

var (
	ErrTooLarge  = errors.New("key and value data is too large; maximum limit of 4KB")
	ErrStoreFull = errors.New("maximum number of records was reached; store is full")
	ErrNotFound  = errors.New("could not locate; not found")
	ErrNonPtrVal = errors.New("expected pointer to value, not value")
)

func encode(k string, v interface{}) ([]byte, error) {
	data := struct {
		Key string
		Val interface{}
	}{k, v}
	b, err := json.Marshal(data)
	if err != nil {
		Logger(err.Error())
		return nil, err
	}
	if len(b) > SYS_PAGE {
		Logger(ErrTooLarge.Error())
		return nil, ErrTooLarge
	}
	return b, nil
}

func decode(b []byte, v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		Logger(ErrNonPtrVal.Error())
		return ErrNonPtrVal
	}
	data := struct {
		Key string
		Val interface{}
	}{
		Val: v, // ptr to value passed in by user
	}
	err := json.Unmarshal(b, &data)
	if err != nil {
		Logger(err.Error())
		return err
	}
	return nil
}

type Store struct {
	index *Tree
	sync.RWMutex
}

func NewStore(name string) *Store {
	return &Store{
		index: NewTree(name),
	}
}

func (st *Store) Add(k string, v interface{}) error {
	doc, err := encode(k, v)
	if err != nil {
		return err
	}
	st.Lock()
	st.index.Add([]byte(k), doc)
	st.Unlock()
	return nil
}

func (st *Store) Set(k string, v interface{}) error {
	doc, err := encode(k, v)
	if err != nil {
		return err
	}
	st.Lock()
	st.index.Set([]byte(k), doc)
	st.Unlock()
	return nil
}

func (st *Store) Get(k string, v interface{}) error {
	st.RLock()
	defer st.RUnlock()
	if doc := st.index.GetDoc([]byte(k)); doc != nil {
		if err := decode(doc, v); err != nil {
			return err
		}
	}
	return ErrNotFound
}

/*
func (st *Store) All() []byte {
	st.RLock()
	size := st.index.Size()
	recs := make([][]byte, size)
	for i, rec := range st.index.All() {
		if i == 0 {
			recs[i] = append([]byte{'['}, rec.Val...)
		}
		if i == size-1 {
			recs[i] = append(rec.Val, byte(']'))
		}
		recs[i] = rec.Val
	}
	st.RUnlock()
	return bytes.Join(recs, []byte{','})
}
*/

func (st *Store) Del(k string) {
	st.Lock()
	st.index.Del([]byte(k))
	st.Unlock()
}
