package main

import (
	"testing"

	"github.com/cagnosolutions/ads/bpt"
)

func BenchmarkInsert(b *testing.B) {
	t := bpt.NewTree()
	var c int
	for i := 0; i < b.N; i++ {
		x := bpt.UUID()
		t.Insert(x, x)
		c++
	}
	if t.Size() != c {
		b.Errorf("t.Size(): %d, c: %d\n", t.Size(), c)
	}
}
