package tmp

import (
	"log"
)

var (
	WS = 8
	SZ = 65536
)

type BitIdx []byte

func NewBitIdx() BitIdx {
	log.Printf("Bit vector of %d provides %d indexes (%dKB)\n", SZ, WS*SZ, SZ/1024)
	return make([]byte, SZ, SZ)
}

func (idx BitIdx) Has(k int) bool {
	return (idx[k/WS] & (1 << (uint(k % WS)))) != 0
}

func (idx BitIdx) Add(k int) {
	idx[k/WS] |= (1 << uint(k%WS))
}

func (idx BitIdx) Del(k int) {
	idx[k/WS] &= ^(1 << uint(k%WS))
}
