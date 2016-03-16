package tmp

import (
	"os"
	"strings"
	"syscall"
)

func sanitize(path string) string {
	if path[len(path)-1] == '/' {
		return path[:len(path)-1]
	}
	if x := strings.Index(path, "."); x != -1 {
		return path[:x]
	}
	return path
}

func mkdirs(path string) string {
	path = sanitize(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			Log(err)
		}
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

func PutInt32(b []byte, v int32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}

func Int32(b []byte) int32 {
	return int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16 | int32(b[3])<<24
}
