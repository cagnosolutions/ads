package tmp

import (
	"os"
	"syscall"
)

func sanitize(path string) string {
	if path[len(path)-1] == '/' {
		return path[:len(path)-1]
	}
	return path
}

func OpenFile(path string) (*os.File, string, int) {
	fd, err := os.OpenFile(path, syscall.O_RDWR|syscall.O_CREAT|syscall.O_APPEND, 0644)
	if err != nil {
		Log(err)
	}
	fi, err := fd.Stat()
	if err != nil {
		Log(err)
	}
	return fd, sanitize(fi.Name()), int(fi.Size())
}

func align(size int) int {
	if size > 0 {
		return (size + SYS_PAGE - 1) &^ (SYS_PAGE - 1)
	}
	return SYS_PAGE
}

func resize(fd uintptr, size int) int {
	err := syscall.Ftruncate(int(fd), int64(align(size)))
	if err != nil {
		Log(err)
	}
	return size
}
