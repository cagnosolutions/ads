package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/cagnosolutions/ads/mio"
)

var THREE_SECONDS = time.Duration(3) * time.Second
var NEW_FILE = true

func main() {

	// open a new mapping (mapping will grow as needed)
	log.Printf("(0)[+] Opening mamory mapped file...\n")
	ts := time.Now().UnixNano()
	m := mio.Map("m.db")
	log.Printf("(0)[-] Mapping file took: %d ms\n", (time.Now().UnixNano()-ts)/1000/1000)
	log.Printf("(0)[=] Mapping Stats\n%s\n\n", m)
	time.Sleep(THREE_SECONDS)

	// if file already has data in it, lets read some of it
	if m.Pages() > 0 {
		NEW_FILE = false
		// reading data (using get) from disk, fifty pseudo-randomized records
		log.Printf("(1)[+] Reading 50 pseudo-random records from disk...\n")
		ts = time.Now().UnixNano()
		j, k := 0, (4096 * 63)
		for i := 0; i < 50; i++ {
			rand.Seed(int64(i))
			j = rand.Intn(k)
			r, ok := m.Get(j)
			if !ok {
				log.Panicf("[%d]: %s\n", i, "something bad happened.")
			} else {
				fmt.Printf("[%d]:\t%s\n", j, r)
			}
		}
		log.Printf("(1)[-] Reading 50 pseudo-random records took: %d ms\n", (time.Now().UnixNano()-ts)/1000/1000)
		log.Printf("(1)[=] Reading Stats\n%s\n\n", m)
		time.Sleep(THREE_SECONDS)
	}

	// if there are any gaps in the file, lets try to fill them
	if !NEW_FILE && m.Gaps() > 0 {
		// write data (using add) to disk, just enough records to fill up the
		// remaining gaps in file (this assumes the records are the same size)
		log.Printf("(2)[+] Filling gaps, writing %d records to disk...\n")
		ts = time.Now().UnixNano()
		for i := 0; i < m.Gaps(); i++ {
			dc := fmt.Sprintf(`{"id":"F-%d","desc":"this is a new record, used to fill any gaps\n"}`, i, i)
			ok := m.Add([]byte(dc))
			if !ok {
				log.Panicf("[%d]: %s\n", i, "something bad happened.")
			}
		}
		log.Printf("(2)[-] Filling gaps took: %d ms\n", (time.Now().UnixNano()-ts)/1000/1000)
		log.Printf("(2)[=] Filling Stats\n%s\n\n", m)
		time.Sleep(THREE_SECONDS)
	}

	// [START] NEW FILE CASE
	if NEW_FILE {
		// write data (using add) to disk, ~a quarter million records
		log.Printf("(1)[+] Writing 1GB of data, aprox %d records to disk...\n", (4096*64)-1)
		ts = time.Now().UnixNano()
		for i := 0; i < (4096*64)-1; i++ {
			dc := fmt.Sprintf(`{"id":%d,"desc":"some description for id-%d"}`, i, i)
			ok := m.Add([]byte(dc))
			if !ok {
				log.Panicf("[%d]: %s\n", i, "something bad happened.")
			}
		}
		log.Printf("(1)[-] Writing took: %d ms\n", (time.Now().UnixNano()-ts)/1000/1000)
		log.Printf("(1)[=] Writing Stats\n%s\n\n", m)
		time.Sleep(THREE_SECONDS)

		// modify data (using set) on disk, five hundred records starting at record number 3
		log.Printf("(2)[+] Modifying 500 records starting at record 3 (offset 12,288)\n")
		ts = time.Now().UnixNano()
		for i := 4096 * 3; i < (4096*3)+500; i++ {
			dc := fmt.Sprintf(`{"id":%d,"desc":"this record [id-%d] has now been modified"}`, i, i)
			ok := m.Set([]byte(dc), i)
			if !ok {
				log.Panicf("[%d]: %s\n", i, "something bad happened.")
			}
		}
		log.Printf("(2)[-] Modifying took: %d ms\n", (time.Now().UnixNano()-ts)/1000/1000)
		log.Printf("(2)[=] Modifying Stats\n%s\n\n", m)
		time.Sleep(THREE_SECONDS)

		// delete data (using del) on disk, seven hundred and fifty randomized records
		log.Printf("(3)[+] Deleting 750 pseudo-random records from disk...\n")
		ts = time.Now().UnixNano()
		j, k := 0, (4096 * 63)
		for i := 0; i < 750; i++ {
			rand.Seed(int64(i))
			j = rand.Intn(k)
			ok := m.Del(j)
			if !ok {
				log.Panicf("[%d]: %s\n", i, "something bad happened.")
			}
		}
		log.Printf("(3)[-] Deleting 750 pseudo-random numbers took: %d ms\n", (time.Now().UnixNano()-ts)/1000/1000)
		log.Printf("(3)[=] Deleting Stats\n%s\n\n", m)
		time.Sleep(THREE_SECONDS)

		// re-write data (using add) to disk, one thousand records (check for reclaimed holes)
		log.Printf("(4)[+] Re-Writing 1,000 records to disk...\n")
		ts = time.Now().UnixNano()
		for i := 0; i < 1000; i++ {
			dc := fmt.Sprintf(`{"id":"N-%d","desc":"this is a new record, part of the 1,000 (%d)\n"}`, i, i)
			ok := m.Add([]byte(dc))
			if !ok {
				log.Panicf("[%d]: %s\n", i, "something bad happened.")
			}
		}
		log.Printf("(4)[-] Re-Writing 1,000 took: %d ms\n", (time.Now().UnixNano()-ts)/1000/1000)
		log.Printf("(4)[=] Re-Writing Stats\n%s\n\n", m)
		time.Sleep(THREE_SECONDS)
	} // [END] NEW FILE CASE

	// unmap the mapping
	log.Printf("(5)[+] Unmapping file")
	m.Unmap()
	log.Printf("(5)[-] Finished\n")

	// pause program for system monitoring/analysis
	pause()
}

func pause() {
	fmt.Println("Press any key to continue...")
	var n byte
	fmt.Scanln(&n)
}
