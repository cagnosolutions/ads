package bpt

import (
	"bytes"
	"log"
)

// master delete
func (t *tree) Delete(key []byte) {
	record := t.Find(key)
	leaf := find_leaf(t.root, key)
	if record != nil && leaf != nil {
		t.root = delete_entry(t.root, leaf, key, record)
		record = nil // free key_record
	}
}

// helper for delete methods... returns index of
// a nodes nearest sibling to the left if one exists
func get_neighbor_index(n *node) int {
	for i := 0; i <= n.parent.num_keys; i++ {
		if n.parent.ptrs[i] == n {
			return i - 1
		}
	}
	log.Fatalf("Search for nonexistent ptr to node in parent.\nNode: %p\n", n)
	return 1
}

func remove_entry_from_node(n *node, key []byte, ptr interface{}) *node {
	var i, num_ptrs int
	// remove key and shift over keys accordingly
	for !bytes.Equal(n.keys[i], key) {
		i++
	}
	for i += 1; i < n.num_keys; i++ {
		n.keys[i-1] = n.keys[i]
	}
	// remove ptr and shift other ptrs accordingly
	// first determine the number of ptrs
	if n.is_leaf {
		num_ptrs = n.num_keys
	} else {
		num_ptrs = n.num_keys + 1
	}
	i = 0
	for n.ptrs[i] != ptr {
		i++
	}

	//for n.ptrs[i].(*node) != ptr {
	//	i++
	//}
	for i += 1; i < num_ptrs; i++ {
		n.ptrs[i-1] = n.ptrs[i]
	}
	// one key has been removed
	n.num_keys--
	// set other ptrs to nil for tidiness; remember leaf
	// nodes use the last ptr to point to the next leaf
	if n.is_leaf {
		for i := n.num_keys; i < ORDER-1; i++ {
			n.ptrs[i] = nil
		}
	} else {
		for i := n.num_keys + 1; i < ORDER; i++ {
			n.ptrs[i] = nil
		}
	}
	return n
}

// deletes an entry from the tree; removes record, key, and ptr from leaf and rebalances tree
func delete_entry(root, n *node, key []byte, ptr interface{}) *node {
	var k_prime_index, capacity int
	var neighbor *node
	var k_prime []byte

	// remove key, ptr from node
	n = remove_entry_from_node(n, key, ptr)
	//switch ptr.(type) {
	//case *node:
	//	n = remove_entry_from_node(n, key, ptr.(*node))
	//case *record:
	//	n = remove_entry_from_node(n, key, ptr.(*record))
	//}
	if n == root {
		return adjust_root(root)
	}

	var min_keys int
	// case: delete from inner node
	if n.is_leaf {
		min_keys = cut(ORDER - 1)
	} else {
		min_keys = cut(ORDER) - 1
	}
	// case: node stays at or above min order
	if n.num_keys >= min_keys {
		return root
	}

	// case: node is bellow min order; coalescence or redistribute
	neighbor_index := get_neighbor_index(n)
	if neighbor_index == -1 {
		k_prime_index = 0
	} else {
		k_prime_index = neighbor_index
	}
	k_prime = n.parent.keys[k_prime_index]
	if neighbor_index == -1 {
		neighbor = n.parent.ptrs[1].(*node)
	} else {
		neighbor = n.parent.ptrs[neighbor_index].(*node)
	}
	if n.is_leaf {
		capacity = ORDER
	} else {
		capacity = ORDER - 1
	}

	// coalescence
	if neighbor.num_keys+n.num_keys < capacity {
		return coalesce_nodes(root, n, neighbor, neighbor_index, k_prime)
	}
	return redistribute_nodes(root, n, neighbor, neighbor_index, k_prime_index, k_prime)
}

func adjust_root(root *node) *node {
	// if non-empty root key and ptr
	// have already been deleted, so
	// nothing to be done here
	if root.num_keys > 0 {
		return root
	}
	var new_root *node
	// if root is empty and has a child
	// promote first (only) child as the
	// new root node. If it's a leaf then
	// the while tree is empty...
	if !root.is_leaf {
		new_root = root.ptrs[0].(*node)
		new_root.parent = nil
	} else {
		new_root = nil
	}
	root = nil // free root
	return new_root
}

// merge (underflow)
func coalesce_nodes(root, n, neighbor *node, neighbor_index int, k_prime []byte) *node {
	var i, j, neighbor_insertion_index, n_end int
	var tmp *node
	// swap neight with node if nod eis on the
	// extreme left and neighbor is to its right
	if neighbor_index == -1 {
		tmp = n
		n = neighbor
		neighbor = tmp
	}
	// starting index for merged pointers
	neighbor_insertion_index = neighbor.num_keys
	// case nonleaf node, append k_prime and the following ptr.
	// append all ptrs and keys for the neighbors.
	if !n.is_leaf {
		// append k_prime (key)
		neighbor.keys[neighbor_insertion_index] = k_prime
		neighbor.num_keys++
		n_end = n.num_keys
		i = neighbor_insertion_index + 1
		j = 0
		for j < n_end {
			neighbor.keys[i] = n.keys[j]
			neighbor.ptrs[i] = n.ptrs[j]
			neighbor.num_keys++
			n.num_keys--
			i++
			j++
		}
		neighbor.ptrs[i] = n.ptrs[j]
		for i = 0; i < neighbor.num_keys+1; i++ {
			tmp = neighbor.ptrs[i].(*node)
			tmp.parent = neighbor
		}
	} else {
		// in a leaf; append the keys and ptrs.
		i = neighbor_insertion_index
		j = 0
		for j < n.num_keys {
			neighbor.keys[i] = n.keys[j]
			neighbor.ptrs[i] = n.ptrs[j]
			neighbor.num_keys++
			i++
			j++
		}
		neighbor.ptrs[ORDER-1] = n.ptrs[ORDER-1]
	}
	root = delete_entry(root, n.parent, k_prime, n)
	n = nil // free n
	return root
}

// merge / redistribute
func redistribute_nodes(root, n, neighbor *node, neighbor_index, k_prime_index int, k_prime []byte) *node {
	var i int
	var tmp *node

	// case: node n has a neighnor to the left
	if neighbor_index != -1 {
		if !n.is_leaf {
			n.ptrs[n.num_keys+1] = n.ptrs[n.num_keys]
		}
		for i = n.num_keys; i > 0; i-- {
			n.keys[i] = n.keys[i-1]
			n.ptrs[i] = n.ptrs[i-1]
		}
		if !n.is_leaf {
			n.ptrs[0] = neighbor.ptrs[neighbor.num_keys]
			tmp = n.ptrs[0].(*node)
			tmp.parent = n
			neighbor.ptrs[neighbor.num_keys] = nil
			n.keys[0] = k_prime
			n.parent.keys[k_prime_index] = neighbor.keys[neighbor.num_keys-1]
		} else {
			n.ptrs[0] = neighbor.ptrs[neighbor.num_keys-1]
			neighbor.ptrs[neighbor.num_keys-1] = nil
			n.keys[0] = neighbor.keys[neighbor.num_keys-1]
			n.parent.keys[k_prime_index] = n.keys[0]
		}
	} else {
		// case: n is left most child (n has no left neighbor)
		if n.is_leaf {
			n.keys[n.num_keys] = neighbor.keys[0]
			n.ptrs[n.num_keys] = neighbor.ptrs[0]
			n.parent.keys[k_prime_index] = neighbor.keys[1]
		} else {
			n.keys[n.num_keys] = k_prime
			n.ptrs[n.num_keys+1] = neighbor.ptrs[0]
			tmp = n.ptrs[n.num_keys+1].(*node)
			tmp.parent = n
			n.parent.keys[k_prime_index] = neighbor.keys[0]
		}
		for i = 0; i < neighbor.num_keys-1; i++ {
			neighbor.keys[i] = neighbor.keys[i+1]
			neighbor.ptrs[i] = neighbor.ptrs[i+1]
		}
		if !n.is_leaf {
			neighbor.ptrs[i] = neighbor.ptrs[i+1]
		}
	}

	n.num_keys++
	neighbor.num_keys--
	return root
}
