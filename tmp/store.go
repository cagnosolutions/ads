package tmp

import "sync"

type Store struct {
	path string
	indx *tree
	data *MappedFile
	sync.RWMutex
}

func NewStore(path string) *Store {
	st := &Store{
		path: path,
		indx: NewTree(),
	}
	file, _ := OpenMappedFile(path)
	st.data = file
	return st
}

func (st *Store) Has() bool {
	st.RLock()
	defer st.RUnlock()
}

func (st *Store) Add() {
	st.Lock()
	defer st.Unlock()
}

func (st *Store) Set() {
	st.Lock()
	defer st.Unlock()
}

func (st *Store) Get() {
	st.RLock()
	defer st.RUnlock()
}

func (st *Store) Del() {
	st.Lock()
	defer st.Unlock()
}

func (st *Store) All() {
	st.RLock()
	defer st.RUnlock()
}
