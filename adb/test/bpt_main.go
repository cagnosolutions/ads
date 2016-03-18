package main

import (
	"fmt"

	"github.com/cagnosolutions/ads/bpt"
)

func main() {
	t := bpt.NewTree()
	for i := 0; i < 100000; i++ {
		x := bpt.UUID()
		t.Insert(x, x)
	}
	pause()
}

func pause() {
	var n int
	fmt.Println("Press any key to continue...")
	fmt.Scanln(&n)
}
