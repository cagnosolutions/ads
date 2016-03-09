package xdb

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
