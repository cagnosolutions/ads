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
		k := fmt.Sprintf(`record-%.3d`, i)
		v := fmt.Sprintf(`{"id":%.3d,"desc":"this is record #%.3d"}`, i, i)
		t.Add(adb.Doc([]byte(k), []byte(v), -1))
	}

	// print total record count
	fmt.Println(t.Size())
	fmt.Println(t.Count())

	// range all records in order
	for _, r := range t.All() {
		fmt.Printf("doc-> k:%x, v:%s\n", r.Key, r.Val)
	}

	// wait... press any key to continue
	pause()
}

func pause() {
	var n int
	fmt.Println("Press any key to continue...")
	fmt.Scanln(&n)
}
