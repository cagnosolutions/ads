package main

import (
	"fmt"
	"time"

	"github.com/cagnosolutions/ads/adb"
)

func main() {

	// create a new tree instance
	t := adb.NewTree("users")
	fmt.Printf("tree size: %d\n", t.Size())

	time.Sleep(time.Duration(5) * time.Second)

	// add 255 records....
	for i := 0; i < 255; i++ {
		x := adb.UUID()
		v := fmt.Sprintf(`{"id":%x,"desc":"this is record %x"}`, x, x)
		t.Add(adb.Doc(x, []byte(v), -1))
	}

	// range all records in order
	for _, r := range t.All() {
		fmt.Printf("doc-> k:%x, v:%s\n", r.Key, r.Val)
	}

	// close
	t.Close()

	// wait... press any key to continue
	pause()
}

func pause() {
	var n int
	fmt.Println("Press any key to continue...")
	fmt.Scanln(&n)
}
