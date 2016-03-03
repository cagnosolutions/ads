package main

import (
	"fmt"
	"log"
	"time"

	"github.com/cagnosolutions/ads/mio"
)

func main() {
	ts := time.Now().UnixNano()
	m := mio.Map("m.db")
	for i := 0; i < (4096*64)-1; i++ {
		dc := fmt.Sprintf(`{"id":%d,"desc":"some description for id-%d"}`, i, i)
		ok := m.Set([]byte(dc), i)
		if !ok {
			log.Panicf("[%d]: %s\n", i, "something bad happened.")
		}
	}
	m.Unmap()
	fmt.Printf("Took: %d ms\n", (time.Now().UnixNano()-ts)/1000/1000)
}
