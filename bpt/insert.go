package bpt

// master insert function
func (t *tree) Insert(key []byte, val []byte) {
	// ignore duplicates: if a value can be found for
	// given key, simply return without inserting
	if t.Find(key) != nil {
		return
	}
	// create record ptr for given value
	ptr := &record{val}
	// if the tree is empty, start a new one
	if t.root == nil {
		t.root = start_new_tree(key, ptr)
		return
	}
	// tree already exists, and ready to insert a non
	// duplicate value. find proper leaf to insert into
	leaf := find_leaf(t.root, key)
	// if the leaf has room, then insert key and record
	if leaf.num_keys < ORDER-1 {
		/*leaf =*/ insert_into_leaf(leaf, key, ptr)
		return
	}
	// otherwise, insert, split, and balance... returning updated root
	t.root = insert_into_leaf_after_splitting(t.root, leaf, key, ptr)
}

// first insertion, start a new tree
func start_new_tree(key []byte, ptr *record) *node {
	root := &node{is_leaf: true}
	root.keys[0] = key
	root.ptrs[0] = ptr
	root.ptrs[ORDER-1] = nil
	root.parent = nil
	root.num_keys++
	return root
}

// creates a new root for two sub-trees and inserts the key into the new root
func insert_into_new_root(left *node, key []byte, right *node) *node {
	root := &node{}
	root.keys[0] = key
	root.ptrs[0] = left
	root.ptrs[1] = right
	root.num_keys++
	root.parent = nil
	left.parent = root
	right.parent = root
	return root
}

// insert a new node (leaf or internal) into tree, return root of tree
func insert_into_parent(root, left *node, key []byte, right *node) *node {
	var left_index int
	var parent *node
	parent = left.parent
	if parent == nil {
		return insert_into_new_root(left, key, right)
	}
	left_index = get_left_index(parent, left)
	if parent.num_keys < ORDER-1 {
		return insert_into_node(root, parent, left_index, key, right)
	}
	return insert_into_node_after_splitting(root, parent, left_index, key, right)
}

// helper->insert_into_parent
// used to find index of the parent's ptr to the
// node to the left of the key to be inserted
// NOTE: best
func get_left_index(parent, left *node) int {
	var left_index int
	for left_index <= parent.num_keys && parent.ptrs[left_index] != left {
		left_index++
	}
	return left_index
}
