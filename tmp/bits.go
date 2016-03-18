package tmp

/*
const (
	WS = 8
	SZ = 1 // 65536 index to index 2GB worth of pages
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

func (idx Data) Has(k int) bool {
	return (idx[k/WS] & (1 << (uint(k % WS)))) != 0
}

func (idx Data) Add(k int) {
	idx[k/WS] |= (1 << uint(k%WS))
}

func (idx Data) Del(k int) {
	idx[k/WS] &= ^(1 << uint(k%WS))
}

func (idx Data) Bits(n byte) int {
	tbl := [16]byte{0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4}
	return int(tbl[n>>4] + tbl[n&0x0f])
}

func (idx Data) Next() int {
	for i := 0; i < len(idx); i++ {
		if idx.Bits(idx[i]) < 8 {
			for j := 0; j < 8; j++ {
				cur := (i * WS) + j
				if !idx.Has(cur) {
					return cur
				}
			}
		}
	}
	return -1
}

func atMax(i int) bool {
	return i == -1
}
*/
