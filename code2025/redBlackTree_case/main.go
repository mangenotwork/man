/*

	该用例是学习红黑树概念后，由AI进行生成，辅助学习

*/

package main

import (
	"fmt"
)

// 定义颜色类型
const (
	RED   = true
	BLACK = false
)

// 定义红黑树节点结构
type Node struct {
	key    int
	value  interface{}
	color  bool
	left   *Node
	right  *Node
	parent *Node
}

// 定义红黑树结构
type RedBlackTree struct {
	root *Node
}

// 新建一个节点
func newNode(key int, value interface{}) *Node {
	return &Node{
		key:    key,
		value:  value,
		color:  RED,
		left:   nil,
		right:  nil,
		parent: nil,
	}
}

// 左旋操作
func (t *RedBlackTree) leftRotate(x *Node) {
	y := x.right
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		t.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

// 右旋操作
func (t *RedBlackTree) rightRotate(y *Node) {
	x := y.left
	y.left = x.right
	if x.right != nil {
		x.right.parent = y
	}
	x.parent = y.parent
	if y.parent == nil {
		t.root = x
	} else if y == y.parent.right {
		y.parent.right = x
	} else {
		y.parent.left = x
	}
	x.right = y
	y.parent = x
}

// 插入修复操作
func (t *RedBlackTree) insertFixup(z *Node) {
	for z.parent != nil && z.parent.color == RED {
		if z.parent == z.parent.parent.left {
			y := z.parent.parent.right
			if y != nil && y.color == RED {
				z.parent.color = BLACK
				y.color = BLACK
				z.parent.parent.color = RED
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					z = z.parent
					t.leftRotate(z)
				}
				z.parent.color = BLACK
				z.parent.parent.color = RED
				t.rightRotate(z.parent.parent)
			}
		} else {
			y := z.parent.parent.left
			if y != nil && y.color == RED {
				z.parent.color = BLACK
				y.color = BLACK
				z.parent.parent.color = RED
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					z = z.parent
					t.rightRotate(z)
				}
				z.parent.color = BLACK
				z.parent.parent.color = RED
				t.leftRotate(z.parent.parent)
			}
		}
	}
	t.root.color = BLACK
}

// 插入操作
func (t *RedBlackTree) Insert(key int, value interface{}) {
	z := newNode(key, value)
	y := (*Node)(nil)
	x := t.root
	for x != nil {
		y = x
		if z.key < x.key {
			x = x.left
		} else {
			x = x.right
		}
	}
	z.parent = y
	if y == nil {
		t.root = z
	} else if z.key < y.key {
		y.left = z
	} else {
		y.right = z
	}
	t.insertFixup(z)
}

// 查找操作
func (t *RedBlackTree) Search(key int) (*Node, bool) {
	x := t.root
	for x != nil {
		if key == x.key {
			return x, true
		} else if key < x.key {
			x = x.left
		} else {
			x = x.right
		}
	}
	return nil, false
}

// 中序遍历
func inorderTraversal(node *Node) {
	if node != nil {
		inorderTraversal(node.left)
		fmt.Printf("Key: %d, Value: %v, Color: %v\n", node.key, node.value, node.color)
		inorderTraversal(node.right)
	}
}

// 打印红黑树（中序遍历）
func (t *RedBlackTree) PrintTree() {
	inorderTraversal(t.root)
}

// 打印红黑树结构图
func printTreeStructure(node *Node, prefix string, isLeft bool) {
	if node == nil {
		return
	}
	// 输出节点信息
	colorStr := "R"
	if !node.color {
		colorStr = "B"
	}
	fmt.Print(prefix)
	if isLeft {
		fmt.Print("├── ")
		prefix += "│   "
	} else {
		fmt.Print("└── ")
		prefix += "    "
	}
	fmt.Printf("%d(%s)\n", node.key, colorStr)
	// 递归打印左子树
	printTreeStructure(node.left, prefix, true)
	// 递归打印右子树
	printTreeStructure(node.right, prefix, false)

}

// 对外暴露的打印红黑树结构图函数
func (t *RedBlackTree) PrintTreeStructure() {
	if t.root == nil {
		fmt.Println("Empty tree")
		return
	}
	printTreeStructure(t.root, "", false)
}

// 删除

// 寻找最小节点
func minimum(x *Node) *Node {
	for x.left != nil {
		x = x.left
	}
	return x
}

// 删除修复操作
func (t *RedBlackTree) deleteFixup(x *Node) {
	if x == nil {
		return
	}
	for x != t.root && x.color == BLACK {
		if x == x.parent.left {
			w := x.parent.right
			if w.color == RED {
				w.color = BLACK
				x.parent.color = RED
				t.leftRotate(x.parent)
				w = x.parent.right
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				x = x.parent
			} else {
				if w.right.color == BLACK {
					w.left.color = BLACK
					w.color = RED
					t.rightRotate(w)
					w = x.parent.right
				}
				w.color = x.parent.color
				x.parent.color = BLACK
				w.right.color = BLACK
				t.leftRotate(x.parent)
				x = t.root
			}
		} else {
			w := x.parent.left
			if w.color == RED {
				w.color = BLACK
				x.parent.color = RED
				t.rightRotate(x.parent)
				w = x.parent.left
			}
			if w.right.color == BLACK && w.left.color == BLACK {
				w.color = RED
				x = x.parent
			} else {
				if w.left.color == BLACK {
					w.right.color = BLACK
					w.color = RED
					t.leftRotate(w)
					w = x.parent.left
				}
				w.color = x.parent.color
				x.parent.color = BLACK
				w.left.color = BLACK
				t.rightRotate(x.parent)
				x = t.root
			}
		}
	}
	x.color = BLACK
}

// 删除操作
func (t *RedBlackTree) Delete(key int) bool {
	z, found := t.Search(key)
	if !found {
		return false
	}
	var y *Node
	var x *Node
	yOriginalColor := z.color
	if z.left == nil {
		x = z.right
		t.transplant(z, z.right)
	} else if z.right == nil {
		x = z.left
		t.transplant(z, z.left)
	} else {
		y = minimum(z.right)
		yOriginalColor = y.color
		x = y.right
		if y.parent == z {
			x.parent = y
		} else {
			t.transplant(y, y.right)
			y.right = z.right
			y.right.parent = y
		}
		t.transplant(z, y)
		y.left = z.left
		y.left.parent = y
		y.color = z.color
	}
	if yOriginalColor == BLACK {
		t.deleteFixup(x)
	}
	return true
}

// 移植操作
func (t *RedBlackTree) transplant(u, v *Node) {
	if u.parent == nil {
		t.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	if v != nil {
		v.parent = u.parent
	}
}

// 修改节点值的函数
func (t *RedBlackTree) Modify(key int, newValue interface{}) bool {
	node, found := t.Search(key)
	if !found {
		return false
	}
	node.value = newValue
	return true
}

func main() {
	tree := &RedBlackTree{}
	tree.Insert(1, "Value 1")
	tree.Insert(2, "Value 2")
	tree.Insert(3, "Value 3")
	tree.Insert(14, "Value 4")
	tree.Insert(15, "Value 5")
	tree.Insert(5, "Value 5")
	tree.Insert(7, "Value 7")
	tree.Insert(9, "Value 9")
	tree.Insert(12, "Value 12")
	tree.PrintTree()
	value, found := tree.Search(3)
	if found {
		fmt.Printf("Found key 3 with value: %v\n", value.value)
	} else {
		fmt.Println("Key 3 not found")
	}

	tree.PrintTreeStructure()

	deleted := tree.Delete(3)
	if deleted {
		fmt.Println("\nAfter deleting key 3:")
		tree.PrintTree()
		tree.PrintTreeStructure()
	} else {
		fmt.Println("Key 3 not found in the tree.")
	}

	tree.Modify(15, "Value 15")

	tree.PrintTree()

	deleted = tree.Delete(13)
	if deleted {
		fmt.Println("\nAfter deleting key 13:")
		tree.PrintTree()
		tree.PrintTreeStructure()
	} else {
		fmt.Println("Key 13 not found in the tree.")
	}
}
