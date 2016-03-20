package adb

import "errors"

var (
	ErrDupKey    = errors.New("found duplicate key; no duplicate keys allowed")
	ErrStoreFull = errors.New("could not find space to insert; store appears to be full")
)

type Store struct {
	name string
	indx *Tree
	meta *MappedMeta
	file *MappedFile
}

func NewStore(name string) *Store {
	st := &Store{}
	st.name = name
	st.indx = NewTree()
	st.meta = OpenMappedMeta(name)
	st.file = OpenMappedFile(name, st.meta.Size())
	return st
}

func (st *Store) Add(k, v []byte) error {
	if st.indx.Has(k) {
		return ErrDupKey
	}
	if p := st.meta.Add(); p != -1 {
		doc := Doc(k, v, p)
		b, err := Enc(doc)
		if err != nil {
			return err
		}
		st.file.Set(p, b)
		st.indx.Add(doc)
	}
	return ErrStoreFull
}

func (st *Store) Set(k, v []byte) {

}

func (st *Store) Get(k []byte, v interface{}) {

}

func (st *Store) All(v interface{}) {

}

func (st *Store) Del(key []byte) {

}
