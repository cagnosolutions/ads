package bpt

// ORDER is defined as the maximum number of pointers in any given node
// MIN_ORDER <= ORDER <= MAX_ORDER
// internal node min ptrs = ORDER/2 round up
// internal node max ptrs = ORDER
// leaf node min ptrs (ORDER-1)/ round up
// leaf node max ptrs ORDER-1

const ORDER = 32

type tree struct {
	root *node
}

type node struct {
	num_keys int
	keys     [ORDER - 1][]byte
	ptrs     [ORDER]interface{}
	parent   *node
	//next     *node
	is_leaf bool
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

/*
func destroy_tree(n *node) {
	destroy_tree_nodes(n)
}

func destroy_tree_nodes(n *node) {
	if n.is_leaf {
		for i := 0; i < n.num_keys; i++ {
			n.ptrs[i] = nil
		}
	} else {
		for i := 0; i < n.num_keys+1; i++ {
			destroy_tree_nodes(n.ptrs[i].(*node))
		}
	}
	n = nil // free
}
*/
