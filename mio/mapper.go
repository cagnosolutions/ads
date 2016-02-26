package mio

import (
	"os"
	"syscall"
	"unsafe"
)

var SYS_PAGE = int64(syscall.Getpagesize())
var DEL_FLAG = []byte{0x00}

type mapper struct {
	path string
	file *os.File
	data []byte
	docs int64
	gaps []int
}

func Map(path string) *mapper {
	fd, err := os.OpenFile(path,
		syscall.O_RDWR|syscall.O_CREAT|syscall.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	fi, err := fd.Stat()
	if err != nil {
		panic(err)
	}
	size := fi.Size()
	if size == 0 {
		size = resizeFile(int(fd.Fd()), toPageSize(16*(1<<20)))
	}
	b, err := syscall.Mmap(int(fd.Fd()), 0, int(size),
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	return &mapper{
		path: path,
		file: fd,
		data: b,
		gaps: make([]int, 0),
	}
}

func (m *mapper) Set(data []byte, offset int) {
	pos := int64(offset) * SYS_PAGE
	copy(m.data[pos:pos+SYS_PAGE], data)
	m.docs++
}

func (m *mapper) Get(offset int) []byte {
	pos := int64(offset) * SYS_PAGE
	return m.data[pos : pos+SYS_PAGE]
}

func (m *mapper) Del(offset int) {
	pos := int64(offset) * SYS_PAGE
	copy(m.data[pos:pos+SYS_PAGE], DEL_FLAG)
}

func (m *mapper) Flush() {
	_, _, err := syscall.Syscall(syscall.SYS_MSYNC,
		uintptr(unsafe.Pointer(&m.data[0])), uintptr(len(m.data)),
		uintptr(syscall.MS_ASYNC))
	if err != 0 {
		panic(err)
	}
}

func (m *mapper) Unmap() {
	err := syscall.Munmap(m.data)
	if err != nil {
		panic(err)
	}
	m.data = nil
}

func toPageSize(size int64) int64 {
	if size > 0 {
		return (size + SYS_PAGE) &^ (SYS_PAGE)
	}
	return SYS_PAGE
}

func resizeFile(fd int, size int64) int64 {
	err := syscall.Ftruncate(fd, size)
	if err != nil {
		panic(err)
	}
	return size
	/*
		if err := fd.Truncate(size); err != nil {
			panic(err)
		}
		if err := fd.Sync(); err != nil {
			panic(err)
		}
	*/
}
