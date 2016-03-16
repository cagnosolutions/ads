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

func (st *Store) Has(k []byte) bool {
	st.RLock()
	ok := st.indx.Has(k)
	st.RUnlock()
	return ok
}

func (st *Store) Add(k, v []byte) error {
	st.Lock()
	defer st.Unlock()
	if err := st.CanInsert(k, v); err != nil {
		return err
	}
	if st.indx.Has(k) {
		return DocExists
	}
	st.indx.Add(k, v)
}

func (st *Store) Set(k, v []byte) error {
	st.Lock()
	defer st.Unlock()
}

func (st *Store) Get(k []byte, v interface{}) error {
	st.RLock()
	defer st.RUnlock()
}

func (st *Store) Del(k []byte) error {
	st.Lock()
	defer st.Unlock()
}

func (st *Store) All(v interface{}) error {
	st.RLock()
	defer st.RUnlock()
}

func (st *Store) CanInsert(k, v, []byte) error {
	if docTooLarge(k, v) {
		return DocTooLarge
	}
	if storeIsFull(st.indx.size()) {
		return StoreIsFull
	}
}

func docTooLarge(k, v []byte) bool {
	return (len(k) + 1 + len(v)) > SYS_PAGE
}

func storeIsFull(records int) bool {
	return records > IDX_SIZE/WS
}
