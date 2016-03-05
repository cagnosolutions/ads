package main

import (
	"bytes"
	"fmt"

	"github.com/cagnosolutions/ads/rbt"
)

type key []byte

func (k key) LessThan(v interface{}) bool {
	return bytes.Compare(k, v.(key)) == -1
}

func main() {
	/*
		ts1 := time.Now().UnixNano()
		fmt.Println("Inserting 10 key/value pairs...")
	*/
	t := rbt.NewTree()
	t.Insert(key(`dk-8`), []byte(`{"id":"dv-8","desc":"dv-8 data goes in here"}`))
	t.Insert(key(`dk-1`), []byte(`{"id":"dv-1","desc":"dv-1 data goes in here"}`))
	t.Insert(key(`dk-2`), []byte(`{"id":"dv-2","desc":"dv-2 data goes in here"}`))
	t.Insert(key(`dk-0`), []byte(`{"id":"dv-0","desc":"dv-0 data goes in here"}`))
	t.Insert(key(`dk-9`), []byte(`{"id":"dv-9","desc":"dv-9 data goes in here"}`))
	t.Insert(key(`dk-3`), []byte(`{"id":"dv-3","desc":"dv-3 data goes in here"}`))
	t.Insert(key(`dk-4`), []byte(`{"id":"dv-4","desc":"dv-4 data goes in here"}`))
	t.Insert(key(`dk-6`), []byte(`{"id":"dv-6","desc":"dv-6 data goes in here"}`))
	t.Insert(key(`dk-5`), []byte(`{"id":"dv-5","desc":"dv-5 data goes in here"}`))
	t.Insert(key(`dk-7`), []byte(`{"id":"dv-7","desc":"dv-7 data goes in here"}`))
	/*
		ts2 := time.Now().UnixNano()
		fmt.Printf("Took %dms\n", ((ts2-ts1)/1024)/1024)
		fmt.Printf("t.Size(): %d\n", t.Size())
		fmt.Printf("Tree...\n")
		t.Preorder()
		fmt.Println()
		fmt.Printf("%#v\n", t.BFS())
	*/
	it := t.Iterator()
	for it != nil {
		fmt.Println(it)
		it = it.Next()
	}
}
