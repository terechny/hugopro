package binary

import (
	"fmt"
	"time"
)

type Node struct {
	Key    int
	Height int
	Left   *Node
	Right  *Node
}

type AVLTree struct {
	Root *Node
}

var ResultString string

func (t *AVLTree) PrintThree() {

	ParseLeaf(t.Root, 0)
}

func ParseLeaf(n *Node, parent int) {

	if parent != 0 {
		v := fmt.Sprintf("%d --> %d \n", parent, n.Key)
		ResultString += v
	}

	if n.Left != nil {
		ParseLeaf(n.Left, n.Key)
	}

	if n.Right != nil {
		ParseLeaf(n.Right, n.Key)
	}
}

func NewNode(key int) *Node {
	return &Node{Key: key, Height: 1}
}

func (t *AVLTree) Insert(key int) {
	t.Root = insert(t.Root, key)
}

func (t *AVLTree) ToMermaid() string {
	return ""
}

func height(node *Node) int {
	if node == nil {
		return 0
	}
	return node.Height
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func updateHeight(node *Node) {
	node.Height = 1 + max(height(node.Left), height(node.Right))
}

func getBalance(node *Node) int {
	if node == nil {
		return 0
	}
	return height(node.Left) - height(node.Right)
}

func leftRotate(x *Node) *Node {
	y := x.Right
	x.Right = y.Left
	y.Left = x
	updateHeight(x)
	updateHeight(y)
	return y
}

func rightRotate(y *Node) *Node {
	x := y.Left
	y.Left = x.Right
	x.Right = y
	updateHeight(y)
	updateHeight(x)
	return x
}

func insert(node *Node, key int) *Node {
	if node == nil {
		return NewNode(key)
	}
	if key < node.Key {
		node.Left = insert(node.Left, key)
	} else if key > node.Key {
		node.Right = insert(node.Right, key)
	} else {
		return node
	}
	updateHeight(node)
	balance := getBalance(node)
	if balance > 1 && key < node.Left.Key {
		return rightRotate(node)
	}
	if balance < -1 && key > node.Right.Key {
		return leftRotate(node)
	}
	if balance > 1 && key > node.Left.Key {
		node.Left = leftRotate(node.Left)
		return rightRotate(node)
	}
	if balance < -1 && key < node.Right.Key {
		node.Right = rightRotate(node.Right)
		return leftRotate(node)
	}
	return node
}

func GenerateTree(count int) *AVLTree {

	tree := &AVLTree{}

	go func() {

		cnt := 0
		for i := 1; i < count; i++ {

			tree.Insert(i)

			if i > 2 {

				time.Sleep(2 * time.Second)
			}

			cnt++

			if cnt > 10 {

				tree.Root.Left = nil
				tree.Root.Right = nil
				//tree = &AVLTree{}
				cnt = 0
			}
		}

	}()

	return tree
}
