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
)

func (st *Store) isOk(k, v []byte) error {
	if (len(k) + 1 + len(v)) > SYS_PAGE {
		return ErrTooLarge
	}
	if st.index.Size()+1 > MAX_DOCS {
		return ErrStoreFull
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

func (st *Store) Add(k, v []byte) error {
	if err := st.isOk(k, v); err != nil {
		return err
	}
	st.Lock()
	st.index.Add(Doc(k, v, -1))
	st.Unlock()
	return nil
}

func (st *Store) Set(k, v []byte) error {
	if err := st.isOk(k, v); err != nil {
		return err
	}
	st.Lock()
	st.index.Set(Doc(k, v, -1))
	st.Unlock()
	return nil
}

func (st *Store) Get(k []byte) []byte {
	st.RLock()
	defer st.RUnlock()
	rec := st.index.Get(k)
	if rec != nil {
		return rec.Val
	}
	return nil
}

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

func (st *Store) Del(k []byte) {
	st.Lock()
	st.index.Del(k)
	st.Unlock()
}
