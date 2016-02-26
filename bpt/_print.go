package bpt

import "fmt"

var queue *node = nil

// helper function for printing the
// tree out. (see print_tree)
func enqueue(new_node *node) {
	var c *node
	if queue == nil {
		queue = new_node
		queue.next = nil
	} else {
		c = queue
		for c.next != nil {
			c = c.next
		}
		c.next = new_node
		new_node.next = nil
	}
}

// helper function for printing the
// tree out. (see print_tree)
func dequeue() *node {
	var n *node = queue
	queue = queue.next
	n.next = nil
	return n
}

// prints the bottom row of keys of the tree
func print_leaves(root *node) {
	var i int
	var c *node = root
	if root == nil {
		fmt.Printf("Empty tree.\n")
		return
	}
	for !c.is_leaf {
		c = c.ptrs[0].(*node)
	}
	for {
		for i = 0; i < c.num_keys; i++ {
			fmt.Printf("%s ", c.keys[i])
		}
		if c.ptrs[ORDER-1] != nil {
			fmt.Printf(" | ")
			c = c.ptrs[ORDER-1].(*node)
		} else {
			break
		}
	}
	fmt.Printf("\n")
}

// print tree out
func print_tree(root *node) {
	var i, rank, new_rank int
	if root == nil {
		fmt.Printf("Empty tree.\n")
		return
	}
	queue = nil
	enqueue(root)
	for queue != nil {
		n := dequeue()
		if n.parent != nil && n == n.parent.ptrs[0] {
			new_rank = path_to_root(root, n)
			if new_rank != rank {
				rank = new_rank
				fmt.Printf("\n")
			}
		}
		for i = 0; i < n.num_keys; i++ {
			fmt.Printf("%s ", n.keys[i])
		}
		if !n.is_leaf {
			for i = 0; i <= n.num_keys; i++ {
				enqueue(n.ptrs[i].(*node))
			}
		}
		fmt.Printf("| ")
	}
	fmt.Printf("\n")
}
