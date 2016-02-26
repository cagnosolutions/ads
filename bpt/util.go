package bpt

// utility to give the height of the tree
func height(root *node) int {
	var h int
	var c *node = root
	for !c.is_leaf {
		c = c.ptrs[0].(*node)
		h++
	}
	return h
}

// utility function to give the length in edges
// for the path from any node to the root
func path_to_root(root, child *node) int {
	var length int
	var c *node = child
	for c != root {
		c = c.parent
		length++
	}
	return length
}

// find first leaf
func find_first_leaf(root *node) *node {
	if root == nil {
		return root
	}
	var c *node = root
	for !c.is_leaf {
		c = c.ptrs[0].(*node)
	}
	return c
}

// get tree size
func (t *tree) Size() int {
	if t.root == nil {
		return -1
	}
	var c *node = t.root
	for !c.is_leaf {
		c = c.ptrs[0].(*node)
	}
	var size int
	for {
		size += c.num_keys
		if c.ptrs[ORDER-1] != nil {
			c = c.ptrs[ORDER-1].(*node)
		} else {
			break
		}
	}
	return size
}
