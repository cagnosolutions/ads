package main

import (
	"bytes"
	"testing"

	"github.com/cagnosolutions/ads/bpt"
	"github.com/cagnosolutions/ads/rbt"
)

func Benchmark_bpt(b *testing.B) {
	t := bpt.NewTree()
	for i := 0; i < b.N; i++ {
		x := bpt.UUID()
		t.Insert(x, x)
	}
}

type key []byte

func (k key) LessThan(v interface{}) bool {
	return bytes.Compare(k, v.(key)) == -1
}

func Benchmark_rbt(b *testing.B) {
	t := rbt.NewTree()
	for i := 0; i < b.N; i++ {
		x := bpt.UUID()
		t.Insert(key(x), x)
	}
}
