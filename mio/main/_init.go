package main

import (
	"fmt"
	"log"

	"github.com/cagnosolutions/ads/mio"
)

func main() {
	fmt.Println("MAPPING")
	m := mio.Map("init.db")
	fmt.Println("ADDING 10 RECORDS")
	if e := m.Add([]byte("foo-0"), []byte("bar-0")); e != nil {
		log.Panicf("[0] -> %s\n", e)
	}
	if e := m.Add([]byte("foo-1"), []byte("bar-1")); e != nil {
		log.Panicf("[1] -> %s\n", e)
	}
	if e := m.Add([]byte("foo-2"), []byte("bar-2")); e != nil {
		log.Panicf("[2] -> %s\n", e)
	}
	if e := m.Add([]byte("foo-3"), []byte("bar-3")); e != nil {
		log.Panicf("[3] -> %s\n", e)
	}
	if e := m.Add([]byte("foo-4"), []byte("bar-4")); e != nil {
		log.Panicf("[4] -> %s\n", e)
	}
	if e := m.Add([]byte("foo-5"), []byte("bar-5")); e != nil {
		log.Panicf("[5] -> %s\n", e)
	}
	if e := m.Add([]byte("foo-6"), []byte("bar-6")); e != nil {
		log.Panicf("[6] -> %s\n", e)
	}
	if e := m.Add([]byte("foo-7"), []byte("bar-7")); e != nil {
		log.Panicf("[7] -> %s\n", e)
	}
	if e := m.Add([]byte("foo-8"), []byte("bar-8")); e != nil {
		log.Panicf("[8] -> %s\n", e)
	}
	if e := m.Add([]byte("foo-9"), []byte("bar-9")); e != nil {
		log.Panicf("[9] -> %s\n", e)
	}
	fmt.Println("DELETING RECORDS 1, 2, 5")
	m.Del(1)
	m.Del(2)
	m.Del(5)
	fmt.Println("UNMAPPING")
	m.Unmap()
	fmt.Println("REMAPPING")
	m = mio.Map("init.db")
	fmt.Println("ADDING 4 MORE RECORDS")
	if e := m.Add([]byte("foo-?1"), []byte("bar-?1")); e != nil {
		log.Panicf("[?1] -> %s\n", e)
	}
	if e := m.Add([]byte("foo-?2"), []byte("bar-?2")); e != nil {
		log.Panicf("[?2] -> %s\n", e)
	}
	if e := m.Add([]byte("foo-?3"), []byte("bar-?3")); e != nil {
		log.Panicf("[?3] -> %s\n", e)
	}
	if e := m.Add([]byte("foo-?4"), []byte("bar-?4")); e != nil {
		log.Panicf("[?4] -> %s\n", e)
	}
	fmt.Println("UNMAPPING")
	m.Unmap()
}
