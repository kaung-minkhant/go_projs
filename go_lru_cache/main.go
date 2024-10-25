package main

import "fmt"

type Node struct {
	Left  *Node
	Right *Node
	Value string
}

type Queue struct {
	Head   *Node
	Tail   *Node
	Length int
}

type Cache struct {
	Queue Queue
	Hash  Hash
	Size  int
}

type Hash map[string]*Node

func NewQueue() Queue {
	head := &Node{}
	tail := &Node{}

	head.Right = tail
	tail.Left = head

	return Queue{
		Head: head,
		Tail: tail,
	}
}

func (q *Queue) Display() {
	node := q.Head.Right
	fmt.Printf("%d - [", q.Length)
	for i := 0; i < q.Length; i++ {
		fmt.Printf("{%s}", node.Value)
		if i < q.Length-1 {
			fmt.Printf("<-->")
		}
		node = node.Right
	}
	fmt.Println("]")
}

func NewCache(size int) Cache {
	return Cache{
		Queue: NewQueue(),
		Hash:  Hash{},
		Size:  size,
	}
}

func (c *Cache) Remove(nodeToRemove *Node) *Node {
	fmt.Printf("Remove: %s\n", nodeToRemove.Value)

	left := nodeToRemove.Left
	right := nodeToRemove.Right

	left.Right = right
	right.Left = left

	c.Queue.Length -= 1
	delete(c.Hash, nodeToRemove.Value)
	return nodeToRemove
}

func (c *Cache) Add(nodeToAdd *Node) {
	fmt.Printf("Add: %s\n", nodeToAdd.Value)
	nodeToAdd.Right = c.Queue.Head.Right
	nodeToAdd.Left = c.Queue.Head
	c.Queue.Head.Right = nodeToAdd
	nodeToAdd.Right.Left = nodeToAdd

	c.Hash[nodeToAdd.Value] = nodeToAdd
	c.Queue.Length += 1
	if c.Queue.Length > c.Size {
		c.Remove(c.Queue.Tail.Left)
	}
}

func (c *Cache) Check(value string) {
	node := &Node{}

	if nodeWithValue, ok := c.Hash[value]; ok {
		node = c.Remove(nodeWithValue)
	} else {
		node = &Node{
			Value: value,
		}
	}

	c.Add(node)
	c.Hash[value] = node
}

func (c *Cache) Display() {
	c.Queue.Display()
}

func main() {
	fmt.Println("START CACHE")
	cache := NewCache(5)

	for _, word := range []string{"apple", "orange", "parrot", "beach", "orange", "food", "apple"} {
		cache.Check(word)
		cache.Display()
	}
}
