package main

import (
	"fmt"

	"github.com/cagnosolutions/ads/adb"
)

type R struct {
	Id   string `json:"id"`
	Desc string `json:"desc"`
}

func main() {

	// create a new tree instance
	t := adb.NewTree("users")

	// add 255 records....
	for i := 0; i < 524289; i++ {
		x := adb.UUID()
		v := fmt.Sprintf(`{"id":%x,"desc":"this is record %d"}`, x, i+1)
		t.Add(adb.Doc(x, []byte(v), -1))
	}
	// range all records in order
	/*for _, r := range t.All() {
		fmt.Printf("doc-> k:%x, v:%s\n", r.Key, r.Val)
	}*/

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
