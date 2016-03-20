package adb

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
	Key []byte
	Val []byte
	Idx int
}

// Doc ...
func Doc(k, v []byte, p int) *Record {
	return &Record{k, v, p}
}

// Tree ...
type Tree struct {
	root *node
}

// NewTree creates and returns a new tree
func NewTree() *Tree {
	return &Tree{}
}

// Has ...
func (t *Tree) Has(key []byte) bool {
	return t.Get(key) != nil
}

// Add ...
func (t *Tree) Add(rec *Record) {
	// ignore duplicates: if a value
	// can be found for a given key,
	// simply return, don't insert
	if t.Get(rec.Key) != nil {
		return
	}
	// otherwise simply call set
	t.Set(rec)
}

// Set ...
func (t *Tree) Set(rec *Record) {
	// don't ignore duplicates: if
	// a value can be found for a
	// given key, simply update the
	// record value and return
	if r := t.Get(rec.Key); r != nil {
		r.Val = rec.Val
		return
	}
	// otherwise...
	// create record ptr for given value
	ptr := rec
	// if the tree is empty, start a new one
	if t.root == nil {
		t.root = startNewTree(rec.Key, ptr)
		return
	}
	// tree already exists, and ready to insert a non
	// duplicate value. find proper leaf to insert into
	leaf := findLeaf(t.root, rec.Key)
	// if the leaf has room, then insert key and record
	if leaf.numKeys < ORDER-1 {
		insertIntoLeaf(leaf, rec.Key, ptr)
		return
	}
	// otherwise, insert, split, and balance... returning updated root
	t.root = insertIntoLeafAfterSplitting(t.root, leaf, rec.Key, ptr)
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

// All ...
func (t *Tree) All() []*Record {
	leaf := findFirstLeaf(t.root)
	var r []*Record
	for {
		for i := 0; i < leaf.numKeys; i++ {
			if leaf.ptrs[i] != nil {
				r = append(r, leaf.ptrs[i].(*Record))
			}
		}
		if leaf.ptrs[ORDER-1] == nil {
			break
		}
		leaf = leaf.ptrs[ORDER-1].(*node)
	}

	return r
}
