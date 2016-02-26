package bpt

// insert a new key, ptr to a node
func insert_into_node(root, n *node, left_index int, key []byte, right *node) *node {
	var i int
	for i = n.num_keys; i > left_index; i-- {
		n.ptrs[i+1] = n.ptrs[i]
		n.keys[i] = n.keys[i-1]
	}
	n.ptrs[left_index+1] = right
	n.keys[left_index] = key
	n.num_keys++
	return root
}

// insert a new key, ptr to a node causing node to split
func insert_into_node_after_splitting(root, old_node *node, left_index int, key []byte, right *node) *node {
	var i, j int
	var child *node
	var tmp_keys [ORDER][]byte
	var tmp_ptrs [ORDER + 1]interface{}
	var k_prime []byte

	for i < old_node.num_keys+1 {
		if j == left_index+1 {
			j++
		}
		tmp_ptrs[j] = old_node.ptrs[i]
		i++
		j++
	}

	i = 0
	j = 0

	for i < old_node.num_keys {
		if j == left_index {
			j++
		}
		tmp_keys[j] = old_node.keys[i]
		i++
		j++
	}

	tmp_ptrs[left_index+1] = right
	tmp_keys[left_index] = key

	split := cut(ORDER)
	new_node := &node{}
	old_node.num_keys = 0

	for i = 0; i < split-1; i++ {
		old_node.ptrs[i] = tmp_ptrs[i]
		old_node.keys[i] = tmp_keys[i]
		old_node.num_keys++
	}

	old_node.ptrs[i] = tmp_ptrs[i]
	k_prime = tmp_keys[split-1]

	j = 0
	for i += 1; i < ORDER; i++ {
		new_node.ptrs[j] = tmp_ptrs[i]
		new_node.keys[j] = tmp_keys[i]
		new_node.num_keys++
		j++
	}

	new_node.ptrs[j] = tmp_ptrs[i]

	// free tmps...
	for i = 0; i < ORDER; i++ {
		tmp_keys[i] = nil
		tmp_ptrs[i] = nil
	}
	tmp_ptrs[ORDER] = nil

	new_node.parent = old_node.parent

	for i = 0; i <= new_node.num_keys; i++ {
		child = new_node.ptrs[i].(*node)
		child.parent = new_node
	}
	return insert_into_parent(root, old_node, k_prime, new_node)
}
