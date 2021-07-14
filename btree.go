package main

import (
	"fmt"
	"strconv"
)

var broot *BTreeNode
var M int

func main() {
	InitBTree(3)

	//Insert(8)
	//Insert(9)

	for i := 10; i < 30; i++{
		Insert(i)
	}

	Show(broot)

	broot.Contain(20)
	broot.Contain(6)
}

type BTreeNode struct {
	n     int
	leaf  bool
	keys  []int
	child []*BTreeNode
}

// New btree
func InitBTree(t int) {
	broot = &BTreeNode{
		0,
		true,
		make([]int, 2*t, 2*t),
		make([]*BTreeNode, 2*t, 2*t),
	}
	M = 3
}

// Insert node
func Insert(k int) {
	var r *BTreeNode
	r = broot
	if r.n == 2*M-1 {
		var s *BTreeNode
		s = &BTreeNode{
			0,
			true,
			make([]int, 2*M, 2*M),
			make([]*BTreeNode, 2*M, 2*M),
		}
		broot = s
		s.leaf = false
		s.n = 0
		s.child[0] = r
		Split(s, 0, r)
		InsertValue(s, k)
	} else {
		InsertValue(r, k)
	}
}

// Insert key
func InsertValue(x *BTreeNode, k int) {
	if x.leaf {
		i := 0
		for i = x.n - 1; i >= 0 && k < x.keys[i]; i-- {
			x.keys[i+1] = x.keys[i]
		}
		x.n = x.n + 1
		x.keys[i+1] = k
	} else {
		i := 0
		for i = x.n - 1; i >= 0 && k < x.keys[i]; i-- {
		}
		i++
		tmp := x.child[i]
		if tmp.n == 2*M-1 {
			Split(x, i, tmp)
			if k > x.keys[i] {
				i++
			}
		}
		InsertValue(x.child[i], k)
	}
}

// Search key
func Search(x *BTreeNode, key int) *BTreeNode {
	i := 0
	if x == nil {
		return nil
	}
	for i = 0; i < x.n; i++ {
		if key < x.keys[i] {
			break
		}
		if key == x.keys[i] {
			return x
		}
	}
	if x.leaf {
		return nil
	} else {
		r := Search(x.child[i], key)
		return r
	}
}

// split the node
func Split(x *BTreeNode, pos int, y *BTreeNode) {
	var z *BTreeNode
	z = &BTreeNode{
		0,
		true,
		make([]int, 2*M, 2*M),
		make([]*BTreeNode, 2*M, 2*M),
	}
	z.leaf = y.leaf
	z.n = M - 1
	for j := 0; j < M-1; j++ {
		z.keys[j] = y.keys[j+M]
	}
	if !y.leaf {
		for j := 0; j < M; j++ {
			z.child[j] = y.child[j+M]
		}
	}
	y.n = M - 1
	for j := x.n; j >= pos+1; j-- {
		x.child[j+1] = x.child[j]
	}
	x.child[pos+1] = z
	for j := x.n - 1; j >= pos; j-- {
		x.keys[j+1] = x.keys[j]
	}
	x.keys[pos] = y.keys[M-1]
	x.n = x.n + 1
}

// show Node
func Show(x *BTreeNode) {
	for i := 0; i < x.n; i++ {
		fmt.Print(x.keys[i])
		fmt.Print(" ")
	}
	fmt.Println("\n")
	if !x.leaf {
		for i := 0; i < x.n+1; i++ {
			Show(x.child[i])
		}
	}
}

// Check if present
func (no *BTreeNode) Contain(k int) {
	x := Search(broot, k)
	if x != nil {
		fmt.Println(strconv.Itoa(k) + " is found")
	} else {
		fmt.Println(strconv.Itoa(k) + " is not found")
	}
}
