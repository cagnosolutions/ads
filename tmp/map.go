package tmp

import (
	"os"
	"syscall"
	"unsafe"
)

type mmap []byte

func open(f *os.File, offset, length int) (mmap, error) {
	return syscall.Mmap(int(f.Fd()), int64(offset), length, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
}

func (m mmap) lock() error {
	return syscall.Mlock(m)
}

func (m mmap) unlock() error {
	return syscall.Munlock(m)
}

func (m mmap) close() error {
	err := syscall.Munmap(m)
	m = nil
	return err
}

func (m mmap) sync() error {
	_, _, err := syscall.Syscall(syscall.SYS_MSYNC,
		uintptr(unsafe.Pointer(&m[0])), uintptr(len(m)),
		uintptr(syscall.MS_ASYNC))
	if err != 0 {
		return err
	}
	return nil
}
