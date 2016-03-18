package bpt

func NewTree() *tree {
	return &tree{}
}

/*
var Usage = `Enter any of the following commands after the prompt >
	i [k] --insert <k>
	f [k] --find <k>
	d [k] --delete <k>
	t 	  --print tree
	l 	  --print leaf values
	q 	  --quit`
*/

/*func (t *tree) Set(key []byte, val []byte) {
	t.insert(key, val)
}

func (t *tree) Get(key []byte) []byte {
	return t.find(key).value
}

func (t *tree) Del(key []byte) {
	t.delete(key)
}

func (t *tree) Close() {
	destroy_tree(t.root)
}

func (t *tree) Print_Tree() {
	print_tree(t.root)
}

func (t *tree) Print_Leaves() {
	print_leaves(t.root)
}
*/
