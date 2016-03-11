package tmp

import (
	"encoding/binary"
	"errors"
	"os"
)

const IDX_SIZE = 1 << 12 // 4 KB

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

func (mx *MappedIndx) Add() (int, bool) {
	return -1, false
}

func (mx *MappedIndx) Set(n int) (int, bool) {
	return -1, false
}

func (mx *MappedIndx) Del(n int) {
	return
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
	// account for page count offset (8)
	mx.indx[n+8] = 0x00
}

func (mx *MappedIndx) DelHole(n int) {
	// account for page count offset (8)
	mx.indx[n+8] = 0x01
}

func (mx *MappedIndx) GetHole() int {
	var i int
	for i = 8; i < len(mx.indx); i++ {
		if mx.indx[i] == 0x00 {
			return i
		}
	}
	return -1
}
