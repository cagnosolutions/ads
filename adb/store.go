package adb

import (
	"bytes"
	"errors"
	"sync"
)

var (
	ErrTooLarge  = errors.New("key and value data is too large; maximum limit of 4KB")
	ErrStoreFull = errors.New("maximum number of records was reached; store is full")
	ErrNotFound  = errors.New("could not locate; not found")
	ErrNonPtrVal = errors.New("expected pointer to value, not value")
)

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

func (st *Store) Get(k string, ptr interface{}) error {
	st.RLock()
	defer st.RUnlock()
	if doc := st.index.GetDoc([]byte(k)); doc != nil {
		if err := decode(doc, ptr); err != nil {
			return err
		}
		return nil
	}
	return ErrNotFound
}

func (st *Store) Del(k string) {
	st.Lock()
	st.index.Del([]byte(k))
	st.Unlock()
}

func (st *Store) All(ptr interface{}) error {
	st.RLock()
	defer st.RUnlock()
	docs := bytes.Join(st.index.All(), []byte{','})
	docs = append([]byte{'['}, append(docs, byte(']'))...)
	if err := decode(docs, ptr); err != nil {
		return err
	}
	return nil
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
