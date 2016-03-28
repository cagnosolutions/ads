package adb

import (
	"bytes"
	"encoding/json"
)

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
	Idx int
	//Val []byte
}

// Doc ...
func Doc(k []byte, p int) *Record {
	return &Record{k, p}
}

// Tree ...
type Tree struct {
	root *node
	name string
	size int
	meta *MappedMeta
	file *MappedFile
}

// NewTree creates and returns a new tree
func NewTree(name string) *Tree {
	t := &Tree{}
	t.name = name
	t.meta = OpenMappedMeta(name)
	t.file = OpenMappedFile(name, t.meta.used)
	t.size = t.meta.used
	t.load()
	return t
}

func (t *Tree) load() {
	data := struct {
		Key string
		Val interface{}
	}{}
	for _, p := range t.meta.All() {
		err := json.Unmarshal(t.file.Get(p), &data)
		if err != nil {
			Logger(err.Error())
		}
		t.Put([]byte(data.Key), p)
	}
}

// Has ...
func (t *Tree) Has(key []byte) bool {
	return t.Get(key) != nil
}

// Add ...
//func (t *Tree) Add(rec *Record) {
func (t *Tree) Add(key, val []byte) {
	// ignore duplicates: if a value
	// can be found for a given key,
	// simply return, don't insert
	if t.Get(key) != nil {
		return
	}
	// otherwise simply call set
	t.Set(key, val)
}

// Put is mainly used for re-indexing
// as it assumes the data to already
// be contained on the disk, so it's
// policy is to just forcefully put it
// in the btree index. it will overwrite
// duplicate keys, as it does not check
// to see if the key already exists...
//func (t *Tree) Put(rec *Record) {
func (t *Tree) Put(key []byte, idx int) {
	// create record ptr for given value
	ptr := &Record{key, idx}
	// if the tree is empty, start a new one
	if t.root == nil {
		t.root = startNewTree(ptr.Key, ptr)
		return
	}
	// tree already exists, and ready to insert a non
	// duplicate value. find proper leaf to insert into
	leaf := findLeaf(t.root, ptr.Key)
	// if the leaf has room, then insert key and record
	if leaf.numKeys < ORDER-1 {
		insertIntoLeaf(leaf, ptr.Key, ptr)
		return
	}
	// otherwise, insert, split, and balance... returning updated root
	t.root = insertIntoLeafAfterSplitting(t.root, leaf, ptr.Key, ptr)
}

// Set ...
func (t *Tree) Set(key, doc []byte) {
	// don't ignore duplicates: if
	// a value can be found for a
	// given key, simply update the
	// record value and return
	if r := t.Get(key); r != nil {
		// update on disk
		t.file.Set(r.Idx, doc)
		return
	}
	// otherwise insert new data...
	// if brand new record, find a page in meta index
	page := t.meta.Add()
	if page == -1 {
		return
	}
	// indexed in meta, add to page on disk
	t.file.Set(page, doc)
	// increase size
	t.size++
	// create record ptr for given value
	ptr := &Record{key, page}
	// if the tree is empty, start a new one
	if t.root == nil {
		t.root = startNewTree(ptr.Key, ptr)
		return
	}
	// tree already exists, and ready to insert a non
	// duplicate value. find proper leaf to insert into
	leaf := findLeaf(t.root, ptr.Key)
	// if the leaf has room, then insert key and record
	if leaf.numKeys < ORDER-1 {
		insertIntoLeaf(leaf, ptr.Key, ptr)
		return
	}
	// otherwise, insert, split, and balance... returning updated root
	t.root = insertIntoLeafAfterSplitting(t.root, leaf, ptr.Key, ptr)
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

func (t *Tree) GetDoc(key []byte) []byte {
	if rec := t.Get(key); rec != nil {
		return t.file.GetDoc(rec.Idx, len(key))
	}
	return nil
}

// Del ...
func (t *Tree) Del(key []byte) {
	record := t.Get(key)
	leaf := findLeaf(t.root, key)
	if record != nil && leaf != nil {
		// wipe page in file
		t.file.Del(record.Idx)
		// remove from meta index
		t.meta.Del(record.Idx)
		// remove from tree, and rebalance
		t.root = deleteEntry(t.root, leaf, key, record)
	}
}

// All ...
func (t *Tree) All() [][]byte {
	leaf := findFirstLeaf(t.root)
	var docs [][]byte
	for {
		for i := 0; i < leaf.numKeys; i++ {
			if leaf.ptrs[i] != nil {
				// get record from leaf
				rec := leaf.ptrs[i].(*Record)
				// get doc and append to docs
				docs = append(docs, t.file.GetDoc(rec.Idx, len(rec.Key)))
			}
		}
		// we're at the end, no more leaves to iterate
		if leaf.ptrs[ORDER-1] == nil {
			break
		}
		// increment/follow pointer to next leaf node
		leaf = leaf.ptrs[ORDER-1].(*node)
	}
	return docs
}

func (t *Tree) Match(qry []byte) [][]byte {
	leaf := findFirstLeaf(t.root)
	var docs [][]byte
	for {
		for i := 0; i < leaf.numKeys; i++ {
			if leaf.ptrs[i] != nil {
				// get record from leaf
				rec := leaf.ptrs[i].(*Record)
				// get doc and append to docs
				if d := t.file.GetDocIfItContains(rec.Idx, len(rec.Key), qry); d != nil {
					docs = append(docs, d)
				}
			}
		}
		// we're at the end, no more leaves to iterate
		if leaf.ptrs[ORDER-1] == nil {
			break
		}
		// increment/follow pointer to next leaf node
		leaf = leaf.ptrs[ORDER-1].(*node)
	}
	return docs
}

// Count ...
func (t *Tree) Count() int {
	if t.root == nil {
		return -1
	}
	c := t.root
	for !c.isLeaf {
		c = c.ptrs[0].(*node)
	}
	var size int
	for {
		size += c.numKeys
		if c.ptrs[ORDER-1] != nil {
			c = c.ptrs[ORDER-1].(*node)
		} else {
			break
		}
	}
	return size
}

// Size ...
func (t *Tree) Size() int {
	return t.size
}

// Close ...
func (t *Tree) Close() {
	t.file.CloseMappedFile()
	t.meta.CloseMappedMeta()
}
