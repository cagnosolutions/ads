package tmp

import (
	"errors"
	"syscall"
)

var SYS_PAGE = syscall.Getpagesize()
var MIN_MMAP = 0x1000000 // 16 MB
var ERR_DOC_OVERFLOW = errors.New("document too large")
var ERR_DOC_ABSENT = errors.New("document does not exist")

type store struct {
	*data
}

func newStore(path string) (*store, error) {

	return nil, nil
}

func (st *store) add(d []byte) (int, error) {

	return 0, nil
}

func (st *store) set(n int, d []byte) error {

	return nil
}

func (st *store) get(n int) []byte {

	return nil
}

func (st *store) del(n int) {

}

func (st *store) all(fn func(n int, d []byte) bool) {

}
