package xdb

import "sync"

// NOTE: this could be used as the "store manager"; containing helper
//methods to enable a cleaner approach from the store's perspective.

// get disk ptr offset from index
func (st *store) get_offset(k []byte) int /* <- type may change */ {
	// locks go in here ??
	return 0
}

func (st *store) get_record(k []byte) []byte {
	_ = st.get_offset(k)
	// get record from disk
	return nil
}

type mgr struct {
	idx *tree
	dat *mapper
	sync.RWMutex
}

func (m *mgr) has(k []byte) bool {
	return false
}

func (m *mgr) add(k, v []byte) bool {
	return false
}

/*func (m *mgr) set(k, v []byte) bool {

	return false
}*/

func (m *mgr) get(k []byte) []byte {
	m.RLock()
	if n, ok := m.idx.val(k); ok {
		return m.dat.get(n)
		m.RUnlock()
	}
	m.RUnlock()
	return nil
}

func (m *mgr) del(k []byte) bool {
	return false
}

func (m *mgr) all() []byte {
	return nil
}

func (m *mgr) getIdx(k []byte) (int, bool) {
	m.RLock()
	n, ok := m.idx.val(k)
	m.RUnlock()
	return n, ok
}

func (m *mgr) getDat(n int) []byte {
	m.RLock()
	//b := m.dat.get(n)
	m.RUnlock()
	//return decdoc(b)
	return nil
}

func (m *mgr) set(k, v []byte) bool {
	m.RLock()
	i, ok := m.idx.val(k)
	m.RUnlock()
	var ni int
	m.Lock()
	if !ok {
		ni, ok = m.dat.add(encdoc(k, v))
	} else {
		ni, ok = m.dat.set(encdoc(k, v), i)
	}
	if ok && ni != i {
		m.idx.set(k, ni)
	}
	m.Unlock()
	return ok
}
