package tmp

import "os"

const IDX_SIZE = 1 << 19

type MappedIndx struct {
	path string
	file *os.File
	indx Data
}

func OpenMappedIndx(path string) *MappedIndx {
	file, path, size := OpenFile(path)
	if size == 0 {
		size = resize(file.Fd(), IDX_SIZE/WS)
	}
	indx := Mmap(file, 0, size)
	mapx := &MappedIndx{
		path: path,
		file: file,
		indx: indx,
	}
	return mapx
}

func (mx *MappedIndx) Has(v int) bool {
	return mx.indx.Has(v)
}

func (mx *MappedIndx) Add() (int, bool) {
	p := mx.indx.Next()
	if p == -1 {
		return -1, false
	}
	mx.indx.Add(p)
	return p, true
}

func (mx *MappedIndx) Set(n int) {
	mx.indx.Add(n)
}

func (mx *MappedIndx) Del(n int) {
	mx.indx.Del(n)
}
