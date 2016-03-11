package tmp

type Store struct {
	*MappedFile
}

func newStore(path string) (*Store, error) {

	return nil, nil
}

func (st *Store) add(d []byte) (int, error) {

	return 0, nil
}

func (st *Store) set(n int, d []byte) error {

	return nil
}

func (st *Store) get(n int) []byte {

	return nil
}

func (st *Store) del(n int) {

}

func (st *Store) all(fn func(n int, d []byte) bool) {

}
