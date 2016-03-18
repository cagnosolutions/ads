package bpt

import "bytes"

// ORDER is defined as the maximum number of pointers in any given node
// MIN_ORDER <= ORDER <= MAX_ORDER
// internal node min ptrs = ORDER/2 round up
// internal node max ptrs = ORDER
// leaf node min ptrs (ORDER-1)/ round up
// leaf node max ptrs ORDER-1

// ORDER maximum number of children a non-leaf node can hold
const ORDER = 32

type node struct {
	numKeys int
	keys    [ORDER - 1][]byte
	ptrs    [ORDER]interface{}
	parent  *node
	isLeaf  bool
}

// Record ...
type Record struct {
	Val []byte
}

// Tree ...
type Tree struct {
	root *node
}

// NewTree creates and returns a new tree
func NewTree() *Tree {
	return &Tree{}
}

// Add ...
func (t *Tree) Add(key, val []byte) {
	if t.Get(key) != nil {
		return
	}
	t.Set(key, val)
}

// Set ...
func (t *Tree) Set(key []byte, val []byte) {
	// ignore duplicates: if a value can be found for
	// given key, simply return without inserting
	if r := t.Get(key); r != nil {
		r.Val = val
		return
	}
	// create record ptr for given value
	ptr := &Record{val}
	// if the tree is empty, start a new one
	if t.root == nil {
		t.root = startNewTree(key, ptr)
		return
	}
	// tree already exists, and ready to insert a non
	// duplicate value. find proper leaf to insert into
	leaf := findLeaf(t.root, key)
	// if the leaf has room, then insert key and record
	if leaf.numKeys < ORDER-1 {
		insertIntoLeaf(leaf, key, ptr)
		return
	}
	// otherwise, insert, split, and balance... returning updated root
	t.root = insertIntoLeafAfterSplitting(t.root, leaf, key, ptr)
}

// Get find record for a given key
func (t *Tree) Get(key []byte) *Record {
	n := findLeaf(t.root, key)
	if n == nil {
		return nil
	}
	var i int
	for i = 0; i < n.numKeys; i++ {
		if bytes.Equal(n.keys[i], key) {
			break
		}
	}
	if i == n.numKeys {
		return nil
	}
	return n.ptrs[i].(*Record)
}

// Del ...
func (t *Tree) Del(key []byte) {
	record := t.Get(key)
	leaf := findLeaf(t.root, key)
	if record != nil && leaf != nil {
		t.root = deleteEntry(t.root, leaf, key, record)
	}
}
