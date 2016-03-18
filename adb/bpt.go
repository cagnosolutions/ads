package bpt

// ORDER is defined as the maximum number of pointers in any given node
// MIN_ORDER <= ORDER <= MAX_ORDER
// internal node min ptrs = ORDER/2 round up
// internal node max ptrs = ORDER
// leaf node min ptrs (ORDER-1)/ round up
// leaf node max ptrs ORDER-1

const ORDER = 32

type node struct {
	num_keys int
	keys     [ORDER - 1][]byte
	ptrs     [ORDER]interface{}
	parent   *node
	is_leaf  bool
}

type record struct {
	Val []byte
}

func cut(length int) int {
	if length%2 == 0 {
		return length / 2
	}
	return length/2 + 1
}

type tree struct {
	root *node
}

func (t *tree) Delete(key []byte) {
	record := t.Find(key)
	leaf := find_leaf(t.root, key)
	if record != nil && leaf != nil {
		t.root = delete_entry(t.root, leaf, key, record)
		record = nil // free key_record
	}
}
