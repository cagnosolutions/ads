package tmp

import (
	"log"
)

var (
	WS = 8
	SZ = 1 // 65536 index to index 2GB worth of pages
	LU = [16]byte{0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4}
)

func clean(idx int) {
    if idx < WS {
		idx = WS
	} else if idx%WS != 0 {
		idx++
	}
	SZ = idx / WS    
}

type BitIdx []byte

func NewBitIdx(idx int) BitIdx {
	clean(idx)	
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

func (idx BitIdx) Bits(n byte) int {
	return int(LU[n>>4] + LU[n&0x0f])
}
