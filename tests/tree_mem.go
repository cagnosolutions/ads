package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/cagnosolutions/ads/bpt"

	"github.com/cagnosolutions/ads/rbt"
)

const N = 100000 // 100,000

func main() {
	rbt_alloc()
	pause()
}

func pause() { fmt.Println("\nPress any key to continue..."); var n int; fmt.Scanln(&n) }

func bpt_alloc() {
	t1 := time.Now().UnixNano()
	t := bpt.NewTree()
	for i := 0; i < N; i++ {
		x := bpt.UUID()
		t.Insert(x, x)
	}
	t2 := time.Now().UnixNano()
	fmt.Printf("inserted %d elements in %dms", t.Size(), ((t2-t1)/1024)/1024)
}

type key []byte

func (k key) LessThan(v interface{}) bool {
	return bytes.Compare(k, v.(key)) == -1
}

func rbt_alloc() {
	t1 := time.Now().UnixNano()
	t := rbt.NewTree()
	for i := 0; i < N; i++ {
		x := bpt.UUID()
		t.Insert(key(x), x)
	}
	t2 := time.Now().UnixNano()
	fmt.Printf("inserted %d elements in %dms", t.Size(), ((t2-t1)/1024)/1024)
}
