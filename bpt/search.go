package bpt

import "bytes"

// find record for a given key
// NOTE: best
func (t *tree) Find(key []byte) *record {
	n := find_leaf(t.root, key)
	if n == nil {
		return nil
	}
	var i int
	for i = 0; i < n.num_keys; i++ {
		if bytes.Equal(n.keys[i], key) {
			break
		}
	}
	if i == n.num_keys {
		return nil
	}
	return n.ptrs[i].(*record)
}

// find leaf type node for a given key
func find_leaf(n *node, key []byte) *node {
	if n == nil {
		return n
	}
	for !n.is_leaf {
		n = n.ptrs[search(n, key)].(*node)
	}
	return n
}

func search(n *node, key []byte) int {
	lo, hi := 0, n.num_keys-1
	for lo <= hi {
		md := (lo + hi) >> 1
		switch cmp := bytes.Compare(key, n.keys[md]); {
		case cmp > 0:
			lo = md + 1
		case cmp == 0:
			return md
		default:
			hi = md - 1
		}
	}
	return lo
}
