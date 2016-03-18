package tmp

import (
	"encoding/json"
	"sync"
)

type Store struct {
	path string
	indx *tree
	meta *MappedMeta
	file *MappedFile
	sync.RWMutex
}

func NewStore(path, store string) *Store {
	st := &Store{}
	st.Lock()
	st.path = path + "/" + store
	st.indx = NewTree()
	st.meta = OpenMappedMeta(st.path)
	st.file = OpenMappedFile(st.path, st.meta.Size())
	st.Unlock()
	return st

}

func (st *Store) Has(k []byte) bool {
	st.RLock()
	ok := st.indx.Has(k)
	st.RUnlock()
	return ok
}

// **ADD TO STORE**
//
// 1) Run helper method (canAdd()), this will ensure
// that the data you are about to enter does not exceed
// the maximum size, that the store is not full, and that
// the key does not already exist; you cannot add duplicate
// keys. 2) Run the add method from the meta index which
// returns an available page id, as well as marks that as
// used. 3) Write data to file at the supplied page id; note
// that we are creating a document using a helper function
// (called document()) as we pass it to the file's Set()
// function. 4) Add the key and page id into the index tree.
// > Return any errors encountered.
func (st *Store) Add(k, v []byte) error {
	// run add checks, if successful...
	if err := st.canAdd(k, v); err != nil {
		return err
	}
	st.Lock()
	defer st.Unlock()
	// get page from meta index; set bit in meta
	if p := st.meta.Add(); p != -1 {
		// write to mapped file
		st.file.Set(p, document(k, v))
		// add to tree index
		st.indx.Set(k, p)
		return nil
	}
	Log(ErrUnknown)
	return ErrUnknown
}

// **SET INTO STORE / UPDATE**
//
func (st *Store) Set(k, v []byte) error {
	st.Lock()
	defer st.Unlock()
	return nil
}

// **GET FROM STORE**
//
// 1) Return page id from the index tree. 2) Make a new
// slice to copy the found data. 3) Get data from the
// file, and copy into fresh slice; in this way, we are
// not mutating the underlying mapped file. 4) Unmarshal
// copied data into the provided value; note we run the
// data through the strip() function helper to strip
// all of the null bytes out of the data.
// > Return any errors encountered.
func (st *Store) Get(k []byte, v interface{}) error {
	st.RLock()
	defer st.RUnlock()
	// try to get page from index...
	if p := st.indx.Get(k); p != -1 {
		// make []byte to copy data into
		d := make([]byte, SYS_PAGE)
		// copy data at page p into d
		copy(d, st.file.Get(p))
		// unmarshal into supplied value, if err return err
		if err := json.Unmarshal(strip(d), v); err != nil {
			return err
		}
		// success, return nil
		return nil
	}
	// not found
	return DocNotFound
}

// **DELETE FROM STORE**
//
// 1) Return page id from the index tree. 2) Delete page
// id from index tree. 3) Remove the data from the file.
// 4) Unset the indexed page bit in the meta index.
// > Return any errors encountered.
func (st *Store) Del(k []byte) error {
	st.Lock()
	defer st.Unlock()
	// try to get page from index, then remove
	if p := st.indx.Get(k); p != -1 {
		// delete from index tree
		st.indx.Del(k)
		// delete from file
		st.file.Del(p)
		// unset bit in meta
		st.meta.Del(p)
		return nil
	}
	// not found
	return DocNotFound
}

func (st *Store) All(v interface{}) error {
	st.RLock()
	defer st.RUnlock()
	return nil
}

// This ensures that the data you are about to enter
// does not exceed the maximum size, that the store is
// not full, and that the key does not already exist;
// you cannot add duplicate keys.
func (st *Store) canAdd(k, v []byte) error {
	st.RLock()
	defer st.RUnlock()
	if docTooLarge(k, v) {
		return DocTooLarge
	}
	if storeIsFull(st.indx.size()) {
		return StoreIsFull
	}
	if st.indx.Has(k) {
		return DocExists
	}
	return nil
}

/*
func document(k, v []byte) []byte {
	kl, vl := len(k), len(v)
	b := make([]byte, kl+vl+1, kl+vl+1)
	copy(b[0:kl], k)
	copy(b[kl+1:kl+1+vl], v)
	b[kl+1] = byte('|')
	return b
}
*/

func docTooLarge(k, v []byte) bool {
	return (len(k) + 1 + len(v)) > SYS_PAGE
}

func storeIsFull(records int) bool {
	return records > IDX_SIZE/WS
}

func strip(b []byte) []byte {
	for i, j := 0, len(b)-1; i <= j; i, j = i+1, j-1 {
		if b[i] == 0x00 {
			return b[:i]
		}
		if b[j] != 0x00 {
			return b[:j+1]
		}
	}
	return b
}
