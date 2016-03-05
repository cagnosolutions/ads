package main

import (
	"fmt"

	"github.com/cagnosolutions/ads/rbt"
)

type key int

func (n key) LessThan(b interface{}) bool {
	value, _ := b.(key)
	return n < value
}

/*
func Test_Preorder(t *testing.T) {
	tree := rbt.NewTree()
	tree.Insert(key(1), "123")
	tree.Insert(key(3), "234")
	tree.Insert(key(4), "dfa3")
	tree.Insert(key(6), "sd4")
	tree.Insert(key(5), "jcd4")
	tree.Insert(key(2), "bcd4")
	if tree.Size() != 6 {
		t.Error("Error size")
		return
	}
	tree.Preorder()
}

func Test_Find(t *testing.T) {
	tree := rbt.NewTree()
	tree.Insert(key(1), "123")
	tree.Insert(key(3), "234")
	tree.Insert(key(4), "dfa3")
	tree.Insert(key(6), "sd4")
	tree.Insert(key(5), "jcd4")
	tree.Insert(key(2), "bcd4")

	n := tree.FindIt(key(4))
	if n.Value != "dfa3" {
		t.Error("Error value")
		return
	}
	n.Value = "bdsf"
	if n.Value != "bdsf" {
		t.Error("Error value modify")
		return
	}
	value := tree.Find(key(5)).(string)
	if value != "jcd4" {
		t.Error("Error value after modifyed other node")
		return
	}
}
*/

func main() {
	fmt.Println("inserted: 1, 3, 4, 6, 5, 2")
	tree := rbt.NewTree()
	tree.Insert(key(1), "123")
	tree.Insert(key(3), "234")
	tree.Insert(key(4), "dfa3")
	tree.Insert(key(6), "sd4")
	tree.Insert(key(5), "jcd4")
	tree.Insert(key(2), "bcd4")
	it1 := tree.Iterator()
	for it1 != nil {
		fmt.Println(it1)
		it1 = it1.Next()
	}
	tree.Preorder()
	it2 := tree.Iterator()
	for it2 != nil {
		fmt.Println(it2)
		it2 = it2.Next()
	}
}

/*
func Test_Delete(t *testing.T) {
	tree := rbt.NewTree()
	tree.Insert(key(1), "123")
	tree.Insert(key(3), "234")
	tree.Insert(key(4), "dfa3")
	tree.Insert(key(6), "sd4")
	tree.Insert(key(5), "jcd4")
	tree.Insert(key(2), "bcd4")
	for i := 1; i <= 6; i++ {
		tree.Delete(key(i))
		if tree.Size() != 6-i {
			t.Error("Delete Error")
		}
	}
	tree.Insert(key(1), "bcd4")
	tree.Clear()
	tree.Preorder()
	if tree.Find(key(1)) != nil {
		t.Error("Can't clear")
		return
	}
}
*/
