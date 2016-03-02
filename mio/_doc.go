package mio

import "fmt"

var DOC_HEAD = 0x05

type data []byte

func align(kn, vn int) int {
	return ((DOC_HEAD + kn + vn) + int(SYS_PAGE) - 1) &^ (int(SYS_PAGE) - 1)
}

func Doc(k, v []byte) data {
	kn, vn := len(k), len(v)
	// return nil if key exceeds 255 B or if value exceeds 65,535 B
	if kn > 0xff || vn > 0xffff {
		return nil
	}
	// set current offset and align bytes for key and value
	co, sz := DOC_HEAD, align(kn, vn)
	// create doc using the aligned byte size
	dc := make(data, sz, sz)
	// create header and write to doc
	copy(dc[0:co+1], []byte{0x01, byte(sz / int(SYS_PAGE)), byte(kn), byte(vn), byte(vn) >> 8})
	// write key to doc
	copy(dc[co:co+kn], k)
	// update current offset
	co += kn
	// write val to doc
	copy(dc[co:co+vn], v)
	return dc
}

func (d data) empty() bool {
	return d[0] == 0x00
}

func (d data) pages() int {
	return int(d[1])
}

func (d data) key() []byte {
	kn := int(d[2])
	return d[DOC_HEAD : DOC_HEAD+kn]
}

func (d data) val() []byte {
	kn, vn := int(d[2]), int(uint16(d[3])|uint16(d[3]>>8))
	return d[DOC_HEAD+kn : DOC_HEAD+kn+vn]
}

func (d data) remove() ([]byte, int) {
	p := int(d[1])
	for i := 0; i < p; i++ {
		fmt.Printf("d[%d:%d]\n", i*int(SYS_PAGE), (i+1)*int(SYS_PAGE))
		copy(d[i*int(SYS_PAGE):(i+1)*int(SYS_PAGE)], NIL_PAGE)
	}
	return d, p
}
