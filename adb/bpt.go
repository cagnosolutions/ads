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

type Encoder func(interface{}) ([]byte, error)
type Decoder func([]byte, interface{}) error

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
	//Val []byte
	Idx int
}

/*
func enc(r *Record) []byte {
	return append(append(r.Key, byte('|')), r.Val...)
}

func dec(d []byte, p int) *Record {
	b := bytes.Split(d, []byte{'|'})
	return &Record{b[0], b[1], p}
}
*/

// Doc ...
func Doc(k, v []byte, p int) *Record {
	return &Record{k, v, p}
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
		Logger(json.Unmarshal(t.file.Get(p), &data).Error())
		t.Put(&Record{[]byte(data.Key), p})
	}
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

// Put is mainly used for re-indexing
// as it assumes the data to already
// be contained on the disk, so it's
// policy is to just forcefully put it
// in the btree index. it will overwrite
// duplicate keys, as it does not check
// to see if the key already exists...
func (t *Tree) Put(rec *Record) {
	// create record ptr for given value
	ptr := rec // NOTE: THIS MAY BE UN-NESSICARY
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

// Set ...
func (t *Tree) Set(rec *Record) {
	// don't ignore duplicates: if
	// a value can be found for a
	// given key, simply update the
	// record value and return
	if r := t.Get(rec.Key); r != nil {
		// update in tree
		r.Val = rec.Val
		// encode...
		b, err := t.enc(r)
		if err != nil {
			Logger(err.Error())
		}
		if len(b) > SYS_PAGE {
			Logger("document is larger than system page size; cannot persist")
		}
		// update on disk
		t.file.Set(r.Idx, b)
		return
	}
	// otherwise insert new data
	if rec.Idx == -1 {
		// if brand new record, find a page in meta index
		if p := t.meta.Add(); p != -1 { // NOTE: this should never be -1
			// change record index to found page
			rec.Idx = p
			// encode...
			b, err := t.enc(rec)
			if err != nil {
				Logger(err.Error())
			}
			if len(b) > SYS_PAGE {
				Logger("document is larger than system page size; cannot persist")
			}
			// indexed in meta, add to page on disk
			t.file.Set(rec.Idx, b)
		}
	}
	// increase size
	t.size++
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
		// wipe page in file
		t.file.Del(record.Idx)
		// remove from meta index
		t.meta.Del(record.Idx)
		// remove from tree, and rebalance
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
