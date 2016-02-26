package bpt

import "bytes"

// inserts a new key and *record into a leaf, then returns leaf
// NOTE: best
func insert_into_leaf(leaf *node, key []byte, ptr *record) /**node*/ {
	var i, insertion_point int
	for insertion_point < leaf.num_keys && bytes.Compare(leaf.keys[insertion_point], key) == -1 {
		insertion_point++
	}
	for i = leaf.num_keys; i > insertion_point; i-- {
		leaf.keys[i] = leaf.keys[i-1]
		leaf.ptrs[i] = leaf.ptrs[i-1]
	}
	leaf.keys[insertion_point] = key
	leaf.ptrs[insertion_point] = ptr
	leaf.num_keys++
	//return leaf
}

// inserts a new key and *record into a leaf, so as
// to exceed the order, causing the leaf to be split
func insert_into_leaf_after_splitting(root, leaf *node, key []byte, ptr *record) *node {

	// perform linear search to find index to insert new record
	var insertion_index int
	for insertion_index < ORDER-1 && bytes.Compare(leaf.keys[insertion_index], key) == -1 {
		insertion_index++
	}

	var tmp_keys [ORDER][]byte
	var tmp_ptrs [ORDER]interface{}
	var i, j int

	// copy leaf keys & ptrs to temp
	// reserve space at insertion index for new record
	for i < leaf.num_keys {
		if j == insertion_index {
			j++
		}
		tmp_keys[j] = leaf.keys[i]
		tmp_ptrs[j] = leaf.ptrs[i]
		i++
		j++
	}
	tmp_keys[insertion_index] = key
	tmp_ptrs[insertion_index] = ptr

	leaf.num_keys = 0

	// index where to split leaf
	split := cut(ORDER - 1)

	// over write original leaf up to split point
	for i = 0; i < split; i++ {
		leaf.ptrs[i] = tmp_ptrs[i]
		leaf.keys[i] = tmp_keys[i]
		leaf.num_keys++
	}

	// create new leaf
	new_leaf := &node{is_leaf: true}

	// writing to new leaf from split point to end of giginal leaf pre split
	j = 0
	for i = split; i < ORDER; i++ {
		new_leaf.ptrs[j] = tmp_ptrs[i]
		new_leaf.keys[j] = tmp_keys[i]
		new_leaf.num_keys++
		j++
	}

	// freeing tmps...
	for i = 0; i < ORDER; i++ {
		tmp_ptrs[i] = nil
		tmp_keys[i] = nil
	}

	new_leaf.ptrs[ORDER-1] = leaf.ptrs[ORDER-1]
	leaf.ptrs[ORDER-1] = new_leaf

	for i = leaf.num_keys; i < ORDER-1; i++ {
		leaf.ptrs[i] = nil
	}

	for i = new_leaf.num_keys; i < ORDER-1; i++ {
		new_leaf.ptrs[i] = nil
	}

	new_leaf.parent = leaf.parent
	new_key := new_leaf.keys[0]

	return insert_into_parent(root, leaf, new_key, new_leaf)
}
