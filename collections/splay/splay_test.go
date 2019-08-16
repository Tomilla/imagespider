package splay

import (
	// "fmt"
	"testing"
)

func Test(t *testing.T) {
	tree := New(func(a, b interface{}) bool {
		return a.(string) < b.(string)
	})

	tree.Add("d")
	tree.Add("b")
	tree.Add("a")
	tree.Add("c")

	if tree.Len() != 4 {
		t.Errorf("expecting len 4")
	}

	tree.Remove("b")

	if tree.Len() != 3 {
		t.Errorf("expecting len 3")
	}
}
