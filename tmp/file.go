package tmp

import (
	"errors"
	"os"
	"syscall"
)

var (
	SYS_PAGE = syscall.Getpagesize()
	NIL_PAGE = make([]byte, SYS_PAGE, SYS_PAGE)
)

const (
	MIN_MMAP = 1 << 24 // 16 MB
	MAX_MMAP = 1 << 31 //  2 GB
)

type MappedFile struct {
	path string
	file *os.File
	size int
	page int
	data Data
	indx *MappedIndx 
}

// open a mapped file, or create if needed and align the
// size to the minimum memory mapped file size (ie. 16 MB)
func OpenMappedFile(path string) (*MappedFile, bool) {
	file, path, size := OpenFile(path + ".dat")
	var iznu bool
	if size == 0 {
		size = resize(file.Fd(), MIN_MMAP)
		iznu = true
	}
	data := Mmap(file, 0, size)
	mapf := &MappedFile{
		path: path + ".dat",
		file: file,
		size: size,
		data: data,
		indx: OpenMappedIndx(path + ".idx"),
	}
	return mapf, iznu
}

// validate
func valid(n int) bool {
	if !(n < SYS_PAGE) {
		Log(errors.New("error: document is too large"))
		return false
	}
	return true
}

// inserts a new block
func (mf *MappedFile) Add(b []byte) int {
	if !valid(len(b)) {
		return -1
	}
	n, ok := mf.indx.Add()
	if !ok {
		mf.Grow()
	}
	copy(mf.data[n:], b)
	return n
}

// updates existing or inserts new block at offset n
func (mf *MappedFile) Set(n int, b []byte) int {
	if !valid(len(b)) {
		return -1
	}
	n, ok := mf.indx.Set(n)
	if !ok {
		mf.Grow()
	}
	copy(mf.data[n:], b)
	return n
}

// returns block at offset n
func (mf *MappedFile) Get(n int) []byte {
	if n > -1 && mf.data[n] != 0x00 {
		return mf.data[n : n+SYS_PAGE]
	}
	return nil
}

// removes block at offset n
func (mf *MappedFile) Del(n int) {
	mf.indx.Del(n)
	copy(mf.data[n:], NIL_PAGE)
}

// returns all non-empty blocks
func (mf *MappedFile) All() [][]byte {
	return nil
}

// closes the mapped file
func (mf *MappedFile) CloseMappedFile() {
	mf.data.Sync()
	mf.data.Munmap()
	Log(mf.file.Close())
}

func (mf *MappedFile) Grow() {
	return
}
