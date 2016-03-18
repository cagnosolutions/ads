package tmp

import "os"

const (
	IDX_SIZE = 1 << 19
	WS       = 8
)

var TBL = [16]byte{0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4}

type MappedMeta struct {
	path string
	file *os.File
	meta Data
}

func OpenMappedMeta(path string) *MappedMeta {
	file, path, size := OpenFile(path + ".idx")
	if size == 0 {
		size = resize(file.Fd(), IDX_SIZE/WS)
	}
	return &MappedMeta{
		path: path + ".idx",
		file: file,
		meta: Mmap(file, 0, size),
	}
}

func (mx *MappedMeta) Has(k int) bool {
	return (mx.meta[k/WS] & (1 << (uint(k % WS)))) != 0
}

func (mx *MappedMeta) Add() int {
	if k := mx.Next(); k != -1 {
		mx.Set(k) // add
		return k
	}
	return -1
}

func (mx *MappedMeta) Set(k int) {
	// flip the n-th bit on; add/set
	mx.meta[k/WS] |= (1 << uint(k%WS))
}

func (mx *MappedMeta) Del(k int) {
	// flip the k-th bit off; delete
	mx.meta[k/WS] &= ^(1 << uint(k%WS))
}

func (mx *MappedMeta) bits(n byte) int {
	return int(TBL[n>>4] + TBL[n&0x0f])
}

func (mx *MappedMeta) Next() int {
	for i := 0; i < len(mx.meta); i++ {
		if mx.bits(mx.meta[i]) < 8 {
			for j := 0; j < 8; j++ {
				cur := (i * WS) + j
				if !mx.Has(cur) {
					return cur
				}
			}
		}
	}
	return -1
}

func (mx *MappedMeta) Size() int {
	var used int
	for i := 0; i < len(mx.meta); i++ {
		if mx.bits(mx.meta[i]) < 8 {
			for j := 0; j < 8; j++ {
				if mx.Has((i * WS) + j) {
					used++
				}
			}
		}
	}
	return used
}
