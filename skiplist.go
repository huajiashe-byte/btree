// 参考埃默里大学课程实践，对skiplist进行了重构
// 原课程java代码
// http://www.mathcs.emory.edu/~cheung/Courses/253/Syllabus/Map/Progs/SkipList/SkipList.java

// 实现跳表数据结构过程讲解
// http://www.mathcs.emory.edu/~cheung/Courses/253/Syllabus/Map/skip-list-impl.html

// 在实现跳表时，首先定义了跳表节点的数据结构
// skiplit data struct：
//
//                 data: key，value 存放k,v键值对
//                 left: 指向在同一层中的前一个结点（每一层最左的节点，left值为nil，即没有前一节点）
//                right: 指向在同一层中的后一个结点（每一层最右的节点，right值为nil，即没有后一节点）
//                   up: 指向上一层有相同key的节点（最上层的节点，up值为nil，即没有上一层节点）
//                 down: 指向上一层有相同key的节点（最底层的节点，down值为nil，即没有下一层节点）

//  e.g.
//  min <--> 28 <------------------> 40 <--------------------------------------------------> max 
//  min <--> 28 <------------------> 40 <--------------------------------------------------> max
//  min <--> 28 <------------------> 40 <--------------------------> 81 <------------------> max
//  min <--> 28 <------------------> 40 <--------------------------> 81 <------------------> max
//  min <--> 28 <--> 31 <--> 37 <--> 40 <--> 41 <--> 47 <--> 66 <--> 81 <--> 87 <--> 88 <--> max 

// 在实现上，节点的key value使用int类型
// 最左（IntMin）以及最右节点（IntMax）使用了8字节长度int类型，取值范围：-9223372036854775808到9223372036854775807

// SkipList()函数
// 程序开始初始化一个空的skiplist，高度，长度为0，head为min节点，tail为max节点

// findKey()函数
// 在put\remove\get时，调用findKey()函数，可以获取key在skplist中的节点

// put()函数
// 调用findKey()获取到节点，key相同时修改value，不同时，插入到获取的节点右侧，修改上左右节点的指针
// 并且随机确定是否增加新插入节点的层数（在已有的层数上增加），并在其每一层的节点修改上下左右指针

// get()函数
// 调用findKey()函数，在skiplist从高层开始查找匹配节点，获取后判断key值是否相等，相等则匹配成功并返回。

// remove()函数
// 调用findKey()函数，删除掉key值匹配成功的节点，并循环逐层修改节点上下前后节点的指针

package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

const IntMax = int(^uint(0) >> 1)
const IntMin = ^IntMax

type SkipListNode struct {
	key   int
	value int
	right *SkipListNode
	left  *SkipListNode
	down  *SkipListNode
	up    *SkipListNode
	pos   int           // 用来协助打印输出，没有实际意义
}

/* ----------------------------------------------
   Constructor: empty skiplist

                        null        null
                         ^           ^
                         |           |
   head --->  null <-- -min <----> +max --> null
                         |           |
                         v           v
                        null        null
   ---------------------------------------------- */
type SkipList struct {
	n int // number of entries in the Skip list
	h int // Height

	head *SkipListNode // First element of the top level
	tail *SkipListNode // Last element of the top level
	r    int           // Coin toss
}

// Default constructor...
func (s *SkipList) SkipList() {
	s.n = 0
	s.h = 0
	s.head = &SkipListNode{IntMin, IntMin, nil, nil, nil, nil, 0}
	s.tail = &SkipListNode{IntMax, IntMax, nil, nil, nil, nil, 0}
	s.head.right = s.tail
	s.tail.left = s.head
	s.r = rand.Int()
}

/* ------------------------------------------------------
     findKey(k): find the largest key x <= k
		   on the LOWEST level of the Skip List
     ------------------------------------------------------ */
func (s *SkipList) findKey(k int) *SkipListNode {
	/* -----------------
	   Start at "head"
	   ----------------- */
	var p *SkipListNode
	p = s.head

	for {
		/* ------------------------------------------------
		           Search RIGHT until you find a LARGER entry

		           E.g.: k = 34

		                     10 ---> 20 ---> 30 ---> 40
		                                      ^
		                                      |
		                                      p must stop here
				   p.right.key = 40
		           ------------------------------------------------ */
		for {
			if p.right.key != IntMax && p.right.key <= k {
				p = p.right
			} else {
				break
			}
		}
		/* ---------------------------------
		   Go down one level if you can...
		   --------------------------------- */
		if p.down != nil {
			p = p.down
		} else {
			break
		}
	}
	return p // p.key <= k
}

func (s *SkipList) get(key int) *int {
	var p *SkipListNode
	p = s.findKey(key)
	if p.key == key {
		return &p.value
	}
	return nil
}

func (s *SkipList) put(key int, value int) *int {
	var p, q *SkipListNode
	p = s.findKey(key)

	/* ------------------------
	Check if key is found, replace the value
	------------------------ */
	if p.key == key {
		x := p.value
		p.value = value
		return &x
	}

	/* --------------------------------------------------------------
	        Insert q into the lowest level after SkipListEntry p:
	                         p       put q here           p        q
	                         |            |               |        |
		 	                 V            V               V        V        V
	        Lower level:    [ ] <------> [ ]      ==>    [ ] <--> [ ] <--> [ ]
	        --------------------------------------------------------------- */
	q = &SkipListNode{key, value, nil, nil, nil, nil, 0}
	q.left = p
	q.right = p.right
	p.right.left = q
	p.right = q

	i := 0 // Current level = 0
	for {
		// Coin flip success: make one more level....
		s.r = rand.Intn(10)
		if s.r < 5 {
			/* ---------------------------------------------
				   Check if height exceed current height.
			 	   If so, make a new EMPTY level
				   --------------------------------------------- */
			if i >= s.h {
				var p1, p2 *SkipListNode
				p1 = &SkipListNode{IntMin, IntMin, nil, nil, nil, nil, 0}
				p2 = &SkipListNode{IntMax, IntMax, nil, nil, nil, nil, 0}

				p1.right = p2
				p2.left = p1

				p1.down = s.head
				p2.down = s.tail

				s.head.up = p1
				s.tail.up = p2

				s.head = p1
				s.tail = p2

				s.h = s.h + 1
			}

			/* -------------------------
			   Scan backwards...
			   ------------------------- */
			for {
				if p.up == nil {
					p = p.left
				} else {
					break
				}
			}
			p = p.up

			/* ---------------------------------------------------
			       Add one more (k, nil) to the column
			               p <--> e(k, nil) <--> p.right
			                         ^
					          |
					          v
					          q
				   ---------------------------------------------------- */
			var e *SkipListNode
			// Don't need the value... so give a zero value
			e = &SkipListNode{key, 0, nil, nil, nil, nil, 0}
			e.key = key
			e.left = p
			e.right = p.right
			e.down = q

			/* ---------------------------------------
			   Change the neighboring links..
			   --------------------------------------- */
			p.right.left = e
			p.right = e
			q.up = e

			q = e     // Set q up for the next iteration
			i = i + 1 // Current level increased by 1
		} else {
			break
		}
	}
	s.n = s.n + 1
	return nil // No old value
}

func (s *SkipList) remove(key int) *int {
	var p, q *SkipListNode
	p = s.findKey(key)
	if p.key != key {
		return nil
	}

	x := p.value
	for {
		if p != nil {
			q = p.up
			p.left.right = p.right
			p.right.left = p.left
			p = q
		} else {
			break
		}
	}
	return &x
}

func (s *SkipList) getOneRow(p *SkipListNode) string {
	var str string
	var a, b, i int

	a = 0

	str = "" + strconv.Itoa(p.key)
	p = p.right

	for {
		if p != nil {
			var q *SkipListNode
			q = p

			for {
				if q.down != nil {
					q = q.down
				} else {
					break
				}
			}

			b = q.pos
			str = str + " <-"

			for i = a + 1; i < b; i++ {
				str = str + "------"
			}
			str = str + "> " + strconv.Itoa(p.key)
			a = b
			p = p.right
		} else {
			break
		}
	}
	return str
}

func (s *SkipList) printHorizontal() {
	var str = ""
	var i int

	var p *SkipListNode
	p = s.head
	for {
		if p.down != nil {
			p = p.down
		} else {
			break
		}
	}

	i = 0
	for {
		if p != nil {
			i++
			p.pos = i
			p = p.right
		} else {
			break
		}
	}

	p = s.head

	/* -------------------
	Print...
	------------------- */
	for {
		if p != nil {
			str = s.getOneRow(p)
			fmt.Println(str)
			p = p.down
		} else {

		}
	}
}

func (s *SkipList) printVertical() {
	var str = ""
	var p *SkipListNode
	p = s.head
	for {
		if p.down != nil {
			p = p.down
		} else {
			break
		}
	}
	for {
		if p != nil {
			str = s.getOneColumn(p)
			fmt.Println(str)
			p = p.right
		} else {
			break
		}
	}
}

func (s *SkipList) getOneColumn(p *SkipListNode) string {
	var str = ""
	for {
		if p != nil {
			str = str + " " + strconv.Itoa(p.key)
			p = p.up
		} else {
			break
		}
	}
	return str
}

func main() {
	s := new(SkipList)
        s.SkipList()
	for i := 0; i < 40; i++ {
	 	s.put(rand.Intn(100), rand.Intn(100))
	}
	// s.printHorizontal()
	s.printVertical()

	x1 := s.get(88)
	if x1 != nil {
		// fmt.Println(reflect.TypeOf(x1))
		fmt.Println(*x1)
	}

	s.remove(88)
	s.printVertical()

	x2 := s.get(88)
	if x2 != nil {
		fmt.Println(*x2)
	} else {
		fmt.Println("not found")
	}
}
