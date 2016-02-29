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
		k := fmt.Sprintf(`foo-%d`, i)
		v := fmt.Sprintf(`{"id":%d,"desc":"some description for id-%d"}`, i, i)

		if err := m.Set([]byte(k), []byte(v), i); err != nil {
			log.Panicf("[0]: %s\n", err)
		}

		/*
			if err := m.Add([]byte(k), []byte(v)); err != nil {
				log.Panicf("[%d]: %s\n", err)
			}
		*/
	}
	m.Unmap()
	fmt.Printf("Took: %d ms\n", (time.Now().UnixNano()-ts)/1000/1000)
}
