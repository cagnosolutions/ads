package adb

// find first leaf
func findFirstLeaf(root *node) *node {
	if root == nil {
		return root
	}
	c := root
	for !c.isLeaf {
		c = c.ptrs[0].(*node)
	}
	return c
}

// Size ...
func (t *Tree) Size() int {
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

func cut(length int) int {
	if length%2 == 0 {
		return length / 2
	}
	return length/2 + 1
}
