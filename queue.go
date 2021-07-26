package main

import "fmt"

type Queue interface {
	push(e *Node)
	pop() *Node
	clear() bool
	size() int
	isEmpty() bool
}

type sliceEntry struct {
	//element []Element
	element []*Node
}

func newQueue() *sliceEntry {
	return &sliceEntry{}
}

//向队列中添加元素
func (entry *sliceEntry) push(e *Node) {
	entry.element = append(entry.element, e)
}

//移除队列中最前面的额元素
func (entry *sliceEntry) pop() *Node {
	if entry.isEmpty() {
		fmt.Println("queue is empty!")
		return nil
	}

	firstElement := entry.element[0]
	entry.element = entry.element[1:]
	return firstElement
}

//清空队列
func (entry *sliceEntry) clear() bool {
	if entry.isEmpty() {
		fmt.Println("queue is empty!")
		return false
	}
	for i := 0; i < entry.size(); i++ {
		entry.element[i] = nil
	}
	entry.element = nil
	return true
}

//队列长度
func (entry *sliceEntry) size() int {
	return len(entry.element)
}

//队列为空判断
func (entry *sliceEntry) isEmpty() bool {
	if len(entry.element) == 0 {
		return true
	}
	return false
}
