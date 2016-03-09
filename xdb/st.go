package xdb

import "sync"

type store struct {
	sid string
	idx *tree
	dat *mapper
	sync.RWMutex
}

// returns boolean of keys existance
func (st *store) has(k []byte) bool {
	// return if key is in index
	return st.idx.has(k)
}

// inserts new data, if key exists return false
func (st *store) add(k, v []byte) bool {
	// if key is not in index
	if !st.idx.has(k) {
		// encdoc using k/v pair; if ok continue
		if r, ok := encdoc(k, v); ok {
			// add to disk, if ok continue
			if i, ok := st.dat.add(r); ok {
				// add key and record offset in index
				st.idx.set(k, uitob(i))
				return true
			}
		}
	}
	// else, something went wrong
	return false
}

// volatiles overwrite/update data
func (st *store) set(k, v []byte) bool {
	// if k/r exists in index...
	var i uint
	if r := st.idx.get(k); r != nil {
		// read ptr offset from index
		i = btoui(r.val)
		// del indexed val for clean update
		st.idx.del(k)
	}
	// encdoc using k/v pair; if ok continue
	if r, ok := encdoc(k, v); ok {
		// set on disk, if ok continue
		if i, ok := st.dat.set(r, int(i)); ok {
			// add key and record offset in index
			st.idx.set(k, uitob(i))
			return true
		}
	}
	// else, something went wrong
	return false
}

// return single value matching key
func (st *store) get(k []byte) []byte {
	// if k/r exists in index...
	if r := st.idx.get(k); r != nil {
		// read ptr offset from index
		i := btoui(r.val)
		// read data off disk; if ok continue
		if b, ok := st.dat.get(int(i)); ok {
			// discard disk header
			b = b[2:]
			// decode data; if ok return v
			if _, v, ok := decdoc(b); ok {
				return v
			}
		}
	}
	// else, something went wrong
	return nil
}

// delete a record matching key
func (st *store) del(k []byte) bool {
	// if k/r exists in index...
	if r := st.idx.get(k); r != nil {
		// read ptr offset from index
		i := btoui(r.val)
		// remove from disk; if ok continue
		if ok := st.dat.remove(int(i)); ok {
			// remove from index and return
			st.idx.del(k)
			return true
		}
	}
	// else, something went wrong
	return false
}

// return all records from store/namespace
func (st *store) all() [][]byte {
	// get all records from index...
	if recs := st.idx.all(); recs != nil {
		// iterate indexed records...
		var data [][]byte
		for j := 0; j < len(recs); j++ {
			// read ptr offset from record index
			i := btoui(recs[j].val)
			// read data off disk; if ok continue
			if b, ok := st.dat.get(int(i)); ok {
				// discard disk header
				b = b[2:]
				// decode data; if ok append v
				if _, v, ok := decdoc(b); ok {
					data = append(data, v)
				}
			}
		}
		return data
	}
	// else, something went wrong
	return nil
}

const (
	K_MAX = 0xff      // 255; 1 byte (uint8)
	V_MAX = 0xfffffff // 4294967295; 4 bytes (uint32)
)

// uint to bytes
func uitob(n uint) []byte {
	return []byte{
		byte(n),
		byte(n),
		byte(n >> 8),
		byte(n >> 24),
	}
}

// bytes to uint
func btoui(b []byte) uint {
	return uint(b[0]) | uint(b[1])<<8 | uint(b[2])<<16 | uint(b[3])<<24
}

// encode key/val to doc and header
func encdoc(k, v []byte) ([]byte, bool) {
	ks, vs := len(k), len(v)
	if ks > K_MAX || vs > V_MAX {
		return nil, false
	}
	ds := 5 + ks + vs
	d := make([]byte, ds, ds)
	d[0] = byte(ks)
	d[1] = byte(vs)
	d[2] = byte(vs >> 8)
	d[3] = byte(vs >> 16)
	d[4] = byte(vs >> 24)
	copy(d[5:5+ks], k)
	copy(d[5+ks:ds], v)
	return d, true
}

// decode doc and header back to key/val
func decdoc(d []byte) ([]byte, []byte, bool) {
	if len(d) < 4 {
		return nil, nil, false
	}
	ks := uint32(d[0])
	vs := uint32(d[1]) | uint32(d[2])<<8 | uint32(d[3])<<16 | uint32(d[4])<<24
	return d[5 : 5+ks], d[5+ks : 5+ks+vs], true
}
