package main

import (
	"fmt"
	"strconv"
)

var MAX int
var root *Node

type Node struct {
	isLeaf bool  // 是否为叶子节点
	parent *Node // 父节点
	key    []int // 关键字
	size   int
	ptr    []*Node // 节点
	// friend BPTree
}

// insert函数负责进行插入
func insert(x int) {
	if root == nil {
		// 根节点为空则创建一个根节点
		root = &Node{
			true,
			nil,
			make([]int, MAX+1),
			1,
			make([]*Node, MAX+1),
		}
		root.key[0] = x
	} else {
		// 如果不为根节点则要找到插入位置（到叶节点才停止）同时记录插入节点的父节点
		cursor := root
		var parent *Node

		for {
			if cursor != nil && !cursor.isLeaf {
				parent = cursor
				for i := 0; i < cursor.size; i++ {
					if x < cursor.key[i] {
						cursor = cursor.ptr[i]
						break
					}
					if i == (cursor.size - 1) {
						cursor = cursor.ptr[i+1]
						break
					}
				}
			}
			break
		}

		if cursor != nil && cursor.size < MAX {
			// 如果插入节点满足关键字个数<MAX( 就是M-1) 我们就可以直接插入
			insertVal(x, cursor)
			cursor.parent = parent
			cursor.ptr[cursor.size] = cursor.ptr[cursor.size-1]
			cursor.ptr[cursor.size-1] = nil
		} else {
			// 否则需要split
			split(x, parent, cursor)
		}
	}
}

// insertVal函数负责找到插入的位置并返回
func insertVal(x int, cursor *Node) int {
	i := 0
	for i = 0; i < cursor.size; i++ {
		if x > cursor.key[i] {
			break
		}
	}
	for j := cursor.size; j > i; j-- {
		cursor.key[j] = cursor.key[j-1]
	}
	cursor.key[i] = x
	cursor.size = cursor.size + 1
	return i
}

// 需要split的插入
func split(x int, parent *Node, cursor *Node) {
	LLeaf := &Node{
		true,
		nil,
		make([]int, MAX+1),
		(MAX + 1) / 2,
		make([]*Node, MAX+1),
	}
	RLeaf := &Node{
		true,
		nil,
		make([]int, MAX+1),
		(MAX + 1) - (MAX+1)/2,
		make([]*Node, MAX+1),
	}

	insertVal(x, cursor)

	for i := 0; i < MAX+1; i++ {
		LLeaf.ptr[i] = cursor.ptr[i]
	}

	LLeaf.ptr[LLeaf.size] = RLeaf
	RLeaf.ptr[RLeaf.size] = LLeaf.ptr[MAX]
	LLeaf.ptr[MAX] = nil

	for i := 0; i < LLeaf.size; i++ {
		LLeaf.key[i] = cursor.key[i]
	}

	j := LLeaf.size
	for i := 0; i < RLeaf.size; i++ {
		RLeaf.key[i] = cursor.key[j]
		j++
	}

	if cursor == root {
		// 叶子节点拆分之后，提上去的节点为根节点
		newRoot := &Node{
			false,
			nil,
			make([]int, (MAX+1)/2),
			1,
			make([]*Node, MAX+1),
		}
		newRoot.key[0] = RLeaf.key[0]
		newRoot.ptr[0] = LLeaf
		newRoot.ptr[1] = RLeaf
		newRoot.isLeaf = false
		newRoot.size = 1
		root = newRoot
		LLeaf.parent = newRoot
		RLeaf.parent = newRoot
	} else {
		// 否则需要调用insertInternal函数
		insertInternal(RLeaf.key[0], parent, LLeaf, RLeaf)
	}
}

// insertInternal插入的实现
func insertInternal(x int, cursor *Node, LLeaf *Node, RRLeaf *Node) {
	if cursor.size < MAX {
		// 如果由于拆分之后提上去的节点不会再产生拆分则直接插入
		i := insertVal(x, cursor)
		for j := cursor.size; j > i+1; j-- {
			cursor.ptr[j] = cursor.ptr[j-1]
		}
		cursor.ptr[i] = LLeaf
		cursor.ptr[i+1] = RRLeaf
	} else {
		// 拆分如果提到根节点则创建新的根节点
		newLchild := &Node{
			true,
			nil,
			make([]int, (MAX+1)/2),
			(MAX + 1) / 2,
			nil,
		}
		newRchild := &Node{
			true,
			nil,
			make([]int, (MAX+1)/2),
			(MAX + 1) / 2,
			nil,
		}
		var virtualPtr = make([]*Node, MAX+2)
		for i := 0; i < MAX+1; i++ {
			virtualPtr[i] = cursor.ptr[i]
		}
		i := insertVal(x, cursor)
		for j := MAX + 2; j > i+1; j-- {
			virtualPtr[j] = virtualPtr[j-1]
		}
		virtualPtr[i] = LLeaf
		virtualPtr[i+1] = RRLeaf
		newLchild.isLeaf = false
		newRchild.isLeaf = false
		//这里和叶子结点上有区别的
		newLchild.size = (MAX + 1) / 2
		newRchild.size = MAX - (MAX+1)/2
		for i := 0; i < newLchild.size; i++ {
			newLchild.key[i] = cursor.key[i]
		}

		j1 := newLchild.size + 1
		for i := 0; i < newRchild.size; i++ {
			newRchild.key[i] = cursor.key[j1]
			j1++
		}

		for i := 0; i < LLeaf.size+1; i++ {
			newLchild.ptr[i] = virtualPtr[i]
		}

		j2 := LLeaf.size + 1
		for i := 0; i < RRLeaf.size+1; i++ {
			newRchild.ptr[i] = virtualPtr[j2]
			j2++
		}

		if cursor == root {
			newRoot := &Node{
				false,
				nil,
				make([]int, (MAX+1)/2),
				(MAX + 1) / 2,
				nil,
			}
			newRoot.key[0] = cursor.key[newLchild.size]
			newRoot.ptr[0] = newLchild
			newRoot.ptr[1] = newRchild
			newRoot.isLeaf = false
			newRoot.size = 1
			root = newRoot
			newLchild.parent = newRoot
			newRchild.parent = newRoot
		} else {
			// 否则就继续调用insertInternal
			insertInternal(cursor.key[newLchild.size], cursor.parent, newLchild, newRchild)
		}
	}
}

func Display() {
	queue := newQueue()
	queue.push(root)
	for {
		if !queue.isEmpty() {
			fmt.Println(queue.size())
			for sizeT := queue.size(); sizeT >= 0; sizeT-- {
				t := queue.element[queue.size()-1]
				for i := 0; i < t.size+1; i++ {
					if !t.isLeaf {
						queue.push(t.ptr[i])
					}
				}
				for i :=0;i<t.size;i++{
					fmt.Println(strconv.Itoa(t.key[i])+",")
				}
				fmt.Println("  ")
				queue.pop()
			}
		}
		break
	}
}

func main() {
	MAX = 3
	insert(5)
	insert(8)
	insert(10)
	insert(15)
	insert(16)
	insert(20)
	insert(19)
	Display()
}
