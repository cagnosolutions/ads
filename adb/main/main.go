package main

import (
	"fmt"

	"github.com/cagnosolutions/ads/adb"
)

func main() {

	// create a new tree instance
	t := adb.NewTree("users")

	// add 255 records....
	for i := 0; i < 255; i++ {
		t.Add(adb.Doc(adb.UUID(), []byte{byte(i)}))
	}

	// print total record count
	fmt.Println(t.Size())

	// range all records in order
	for _, r := range t.All() {
		fmt.Printf("doc-> k:%x, v:%d\n", r.Key, r.Val)
	}

	// wait... press any key to continue
	pause()
}

func pause() {
	var n int
	fmt.Println("Press any key to continue...")
	fmt.Scanln(&n)
}
