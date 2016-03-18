package main

import (
	"fmt"

	"github.com/cagnosolutions/ads/adb"
)

func main() {
	t := adb.NewTree()

	for i := 0; i < 255; i++ {
		//x := adb.UUID()
		t.Add(adb.UUID(), []byte{byte(i)})
	}
	fmt.Println(t.Size())
	r := t.All()
	for _, v := range r {
		fmt.Printf("%d\n", v.Val)
	}

	pause()
}

func pause() {
	var n int
	fmt.Println("Press any key to continue...")
	fmt.Scanln(&n)
}
