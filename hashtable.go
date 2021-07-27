// 使用GoLang的切片实现的hash table,hash方法是mod求余。
// 处理hash冲突的方式是使用拉链法，即在冲突位置使用链表来链接冲突结果
// hash table的slot初始值为常量值，固定为10

// 程序包含两个struct
// HashNode结构体定义val以及next两个属性，Val用于存储hash value，而next用于链接hash冲突时的下一个HashNode
// HashTable结构体，定义了指针数组Array以及Length两个属性，Array数组用于存放每一个hash value，
// 而Length则定义了value长度，即hash table中存放了多少value

// 程序初始化时。首先创建HashTable结构体，定义数组以及长度，初始数组切片为空，容量和长度使用常量Size定义，存放value长度Length为0
// 初始化之后HashTable可以插入value，插入值首先判断是否存在，不存在则插入，遇到冲突，则插入头部，这与redis的方式一致，存在则覆盖。

// Show() function则循环输出hash table的value，遇到冲突则循环输出next value

// MovNode() 移除一个value，找到value在数组中位置，如果不为空，则在链表中删除hashnode

package main

import (
	"fmt"
)

const Size = 10

type HashNode struct {
	Val  int
	Next *HashNode
}

type HashTable struct {
	Array []*HashNode
	Length int
}

func (n *HashTable) HashTableInit() {
	n.Array = make([]*HashNode, Size)
	n.Length = 0
}

func Hash(key int) int {
	return key % Size
}

func (n *HashTable) Search(value int) *HashNode {
	index := Hash(value)
	node := n.Array[index]
	for node != nil {
		if node.Val == value {
			return node
		} else {
			node = node.Next
		}
	}
	return nil
}

func (n *HashTable) InsertValue(value int) bool {
	node := n.Search(value)
	if node == nil {
		index := Hash(value)
		newNode := new(HashNode)
		newNode.Val = value
		newNode.Next = n.Array[index]
		n.Array[index] = newNode
		n.Length++
		return true
	}
	node.Val = value
	return true
}

// 找到链表位置，删除节点
func (n *HashTable) MovNode(value int) int {
	index := Hash(value)
	var cur, pre *HashNode
	cur = n.Array[index]
	//不为空，寻找value
	for cur != nil {
		if cur.Val == value {
			//index第一个节点value，则index位置为next
			if pre == nil {
				n.Array[index] = cur.Next
			} else {
				//非第一个节点，去掉节点就可以
				pre.Next = cur.Next
			}
			n.Length--
			return 0
		}
		//继续寻找value
		pre = cur
		cur = cur.Next
	}
	//等于nil则无此值
	return -1
}

func (n *HashTable) Show() {
	for i := 0; i < Size; i++ {
		cur := n.Array[i]
		for cur != nil {
			fmt.Printf("%v ", cur.Val)
			cur = cur.Next
		}
	}
}

func main() {
	hashTable := new(HashTable)
	hashTable.HashTableInit()
	hashTable.InsertValue(1)
	hashTable.InsertValue(11)
	hashTable.InsertValue(21)
	hashTable.InsertValue(31)
	hashTable.InsertValue(41)
	hashTable.InsertValue(2)
	hashTable.InsertValue(3)
	hashTable.InsertValue(13)
	hashTable.InsertValue(4)
	hashTable.InsertValue(7)
	hashTable.InsertValue(17)
	hashTable.InsertValue(27)
	hashTable.InsertValue(9)
	hashTable.Show()
	fmt.Println()
	hashTable.MovNode(2)
	hashTable.Show()
	fmt.Println()
	fmt.Println(hashTable.Search(9))
}
