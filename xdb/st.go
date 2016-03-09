package xdb

import "sync"

const (
	DELIM = 0x1f
)

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
		// add to disk, if ok continue
		// encdoc using k/v pair
		if i, ok := st.dat.add(encdoc(k, v)); ok {
			// add key and record offset in index
			st.idx.set(k, i)
			return true
		}
	}
	// else, something went wrong
	return false
}

// volatiles overwrite/update data
func (st *store) set(k, v []byte) bool {
	/*m.RLock()
	i, ok := m.idx.val(k)
	m.RUnlock()
	if !ok {

	}
	m.Lock()
	ni, ok := m.dat.set(encdoc(k, v), i)
	if ok && ni != i {
		m.idx.set(k, ni)
	}
	m.Unlock()
	return ok*/
	return false
}

// return single value matching key
func (st *store) get(k []byte) []byte {

	// if k/r exists in index...
	if i, ok := st.idx.val(k); ok {
		// read data off disk; if ok continue
		if b := st.dat.get(int(i)); b != nil {
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
	if i, ok := st.idx.val(k); ok {
		// remove from disk; if ok continue
		if ok := st.dat.del(i); ok {
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
			i := recs[j]
			// read data off disk; if ok continue
			if b := st.dat.get(int(i)); b != nil {
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
func encdoc(k, v []byte) []byte {

	return append(k, append([]byte{DELIM}, v...)...)

	/*ks, vs := len(k), len(v)
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
	return d, true*/
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
