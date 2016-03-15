package tmp

import (
	"encoding/binary"
	"errors"
	"os"
)

<<<<<<< HEAD
const (
	IDX_SIZE = 1 << 12   // 4 KB
	MAX_PAGS = 1<<19 - 1 // 524,287 PAGES
	NIL_HOLE = make([]byte, 4, 4)
)
=======
// NOTE: this should probably be 1 << 19
const IDX_SIZE = 1 << 16 // 64KB
>>>>>>> 5c5d0ad7afcdc9169c4d17acc15bada652278f44

type MappedIndx struct {
	path string
	file *os.File
	indx Data
}

func OpenMappedIndx(path string) *MappedIndx {
	file, path, size := OpenFile(path)
	if size == 0 {
		size = resize(file.Fd(), IDX_SIZE)
	}
	indx := Mmap(file, 0, size)
	mapx := &MappedIndx{
		path: path,
		file: file,
		indx: indx,
	}
	return mapx
}

func (mx *MappedIndx) Init() {
	copy(mx.indx, make([]byte, IDX_SIZE, IDX_SIZE))
}

func (mx *MappedIndx) CanInsert() bool {
	return mx.GetPage() < MAX_PAGS
}

func (mx *MappedIndx) Add() (int, bool) {
	if !mx.CanInsert() {
		return -1, false
	}
	p, ok := mx.GetHole()
	if !ok {
		p = mx.GetPage() * SYS_PAGE
	}
	mx.AddPage()
	return p, true
}

func (mx *MappedIndx) Set(n int) (int, bool) {
	return -1, false
}

func (mx *MappedIndx) Del(n int) {
	mx.DelPage()
	mx.AddHole(n)
}

func (mx *MappedIndx) AddPage() {
	mx.AddPages(1)
}

func (mx *MappedIndx) GetPage() int {
	p, n := binary.Varint(mx.indx[0:8])
	if n < 1 {
		Log(errors.New("error:extracting int64 from mapped index"))
	}
	return int(p)
}

func (mx *MappedIndx) DelPage() {
	mx.DelPages(1)
}

func (mx *MappedIndx) AddPages(c int) {
	p, n := binary.Varint(mx.indx[0:8])
	if n < 1 {
		Log(errors.New("error:extracting int64 from mapped index"))
	}
	p += int64(c)
	binary.PutVarint(mx.indx[0:8], p)
}

func (mx *MappedIndx) DelPages(c int) {
	p, n := binary.Varint(mx.indx[0:8])
	if n < 1 {
		Log(errors.New("error:extracting int64 from mapped index"))
	}
	p -= int64(c)
	binary.PutVarint(mx.indx[0:8], p)
}

func (mx *MappedIndx) AddHole(n int) {
	if mx.GetPage() < MAX_PAGES {

	}
	// account for page count offset (8)
	mx.indx[n+8] = 0x00
}

/*
func (mx *MappedIndx) DelHole(n int) {
	// account for page count offset (8)
	i := n * SYS_PAGE
	copy(m.indx[i+8:i+8+4], )
	mx.indx[n+8] = 0x01
}
*/

// finds and returns hole/gap; if one is found and
// returned successfully GetHole() also removes it.
func (mx *MappedIndx) GetHole() (int, bool) {
	var i int
	for i = 8; i < len(mx.indx)-4; i += 4 {
		if mx.indx[i] != 0x00 {
			p := Int32(mx.indx[i : i+4])
			copy(m.indx[i:i+4], NIL_HOLE)
			return p, true
		}
	}
	return -1, false
}
