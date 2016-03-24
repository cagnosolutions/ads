package adb

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

// find first leaf
func findFirstLeaf(root *node) *node {
	if root == nil {
		return root
	}
	c := root
	for !c.isLeaf {
		c = c.ptrs[0].(*node)
	}
	return c
}

func cut(length int) int {
	if length%2 == 0 {
		return length / 2
	}
	return length/2 + 1
}

func OpenFile(path string) (*os.File, string, int) {
	fd, err := os.OpenFile(path, syscall.O_RDWR|syscall.O_CREAT|syscall.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	fi, err := fd.Stat()
	if err != nil {
		panic(err)
	}
	return fd, sanitize(fi.Name()), int(fi.Size())
}

func sanitize(path string) string {
	if path[len(path)-1] == '/' {
		return path[:len(path)-1]
	}
	if x := strings.Index(path, "."); x != -1 {
		return path[:x]
	}
	return path
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
		panic(err)
	}
	return size
}

func strip(b []byte) []byte {
	for i, j := 0, len(b)-1; i <= j; i, j = i+1, j-1 {
		if b[i] == 0x00 {
			return b[:i]
		}
		if b[j] != 0x00 {
			return b[:j+1]
		}
	}
	return b
}

/*

func Encode(k string, v interface{}) []byte {
	data := struct {
		Key string
		Val interface{}
	}{
		k,
		v,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	return b
}

func DecodeVal(b []byte, v interface{}) {
	data := struct {
		Key string
		Val interface{}
	}{
		Val: v,
	}
	json.Unmarshal(b, &data)
}

func DecodeRecord(b []byte) Record {
	data := struct {
		Key string
		Val interface{}
	}{}
	err := json.Unmarshal(b, &data)
	return Record{[]byte(data.Key), 4}
}
*/

func GetDocArr(d []byte, kLen int) []byte {
	start := kLen + 4
	var count int
	for i, j, set := start, len(d)-1, 1; i < j; i, j = i+1, j-1 {
		count++
		if d[i] == '[' {
			set++
		}
		if d[i] == ']' {
			set--
		}
		if set == 0 || d[j] == ']' {
			fmt.Printf("Took %d loops\n", count)
			if d[i] == ']' {
				return d[start:i]
			}
			return d[start:j]
		}
	}
	return d
}
