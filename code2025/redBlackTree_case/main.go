/*

	该用例是学习红黑树概念后，由AI进行生成，辅助学习

*/

package main

import (
	"fmt"
)

// 定义颜色类型
// 在红黑树中，节点颜色只有红色和黑色两种，使用布尔类型来表示，true 代表红色，false 代表黑色
const (
	RED   = true
	BLACK = false
)

// 定义红黑树节点结构
// 每个节点包含键（key）、值（value）、颜色（color）、左子节点（left）、右子节点（right）和父节点（parent）
type Node struct {
	key    int
	value  interface{}
	color  bool
	left   *Node
	right  *Node
	parent *Node
}

// 定义红黑树结构
// 红黑树由一个根节点（root）表示
type RedBlackTree struct {
	root *Node
}

// 新建一个节点
// 根据传入的键和值创建一个新的节点，默认颜色为红色，左右子节点和父节点初始化为 nil
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
// 左旋是红黑树中用于维持树平衡的一种操作，以节点 x 为支点进行左旋
func (t *RedBlackTree) leftRotate(x *Node) {
	// 取出 x 的右子节点 y
	y := x.right
	// 将 y 的左子节点作为 x 的右子节点
	x.right = y.left
	// 如果 y 的左子节点不为空，更新其 parent 指针指向 x
	if y.left != nil {
		y.left.parent = x
	}
	// 将 y 的 parent 指针指向 x 的父节点
	y.parent = x.parent
	// 如果 x 是根节点，将 y 设置为新的根节点
	if x.parent == nil {
		t.root = y
	} else if x == x.parent.left {
		// 如果 x 是其父节点的左子节点，将 y 设置为其父节点的左子节点
		x.parent.left = y
	} else {
		// 如果 x 是其父节点的右子节点，将 y 设置为其父节点的右子节点
		x.parent.right = y
	}
	// 将 x 作为 y 的左子节点
	y.left = x
	// 更新 x 的 parent 指针指向 y
	x.parent = y
}

// 右旋操作
// 右旋是红黑树中用于维持树平衡的一种操作，以节点 y 为支点进行右旋
func (t *RedBlackTree) rightRotate(y *Node) {
	// 取出 y 的左子节点 x
	x := y.left
	// 将 x 的右子节点作为 y 的左子节点
	y.left = x.right
	// 如果 x 的右子节点不为空，更新其 parent 指针指向 y
	if x.right != nil {
		x.right.parent = y
	}
	// 将 x 的 parent 指针指向 y 的父节点
	x.parent = y.parent
	// 如果 y 是根节点，将 x 设置为新的根节点
	if y.parent == nil {
		t.root = x
	} else if y == y.parent.right {
		// 如果 y 是其父节点的右子节点，将 x 设置为其父节点的右子节点
		y.parent.right = x
	} else {
		// 如果 y 是其父节点的左子节点，将 x 设置为其父节点的左子节点
		y.parent.left = x
	}
	// 将 y 作为 x 的右子节点
	x.right = y
	// 更新 y 的 parent 指针指向 x
	y.parent = x
}

// 插入修复操作
// 在向红黑树中插入新节点后，可能会破坏红黑树的性质，需要进行修复操作
func (t *RedBlackTree) insertFixup(z *Node) {
	// 当 z 的父节点不为空且父节点颜色为红色时，需要进行调整
	for z.parent != nil && z.parent.color == RED {
		if z.parent == z.parent.parent.left {
			// 如果 z 的父节点是其祖父节点的左子节点
			y := z.parent.parent.right
			if y != nil && y.color == RED {
				// 情况 1：z 的叔叔节点 y 为红色
				// 将 z 的父节点和叔叔节点颜色设置为黑色，祖父节点颜色设置为红色
				z.parent.color = BLACK
				y.color = BLACK
				z.parent.parent.color = RED
				// 将 z 指针移动到祖父节点，继续向上检查
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					// 情况 2：z 是其父节点的右子节点
					// 将 z 指针移动到其父节点，然后以 z 为支点进行左旋
					z = z.parent
					t.leftRotate(z)
				}
				// 情况 3：z 是其父节点的左子节点
				// 将 z 的父节点颜色设置为黑色，祖父节点颜色设置为红色
				z.parent.color = BLACK
				z.parent.parent.color = RED
				// 以 z 的祖父节点为支点进行右旋
				t.rightRotate(z.parent.parent)
			}
		} else {
			// 如果 z 的父节点是其祖父节点的右子节点
			y := z.parent.parent.left
			if y != nil && y.color == RED {
				// 情况 1：z 的叔叔节点 y 为红色
				// 将 z 的父节点和叔叔节点颜色设置为黑色，祖父节点颜色设置为红色
				z.parent.color = BLACK
				y.color = BLACK
				z.parent.parent.color = RED
				// 将 z 指针移动到祖父节点，继续向上检查
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					// 情况 2：z 是其父节点的左子节点
					// 将 z 指针移动到其父节点，然后以 z 为支点进行右旋
					z = z.parent
					t.rightRotate(z)
				}
				// 情况 3：z 是其父节点的右子节点
				// 将 z 的父节点颜色设置为黑色，祖父节点颜色设置为红色
				z.parent.color = BLACK
				z.parent.parent.color = RED
				// 以 z 的祖父节点为支点进行左旋
				t.leftRotate(z.parent.parent)
			}
		}
	}
	// 确保根节点颜色为黑色
	t.root.color = BLACK
}

// 插入操作
// 向红黑树中插入一个新的键值对
func (t *RedBlackTree) Insert(key int, value interface{}) {
	// 创建一个新的节点
	z := newNode(key, value)
	y := (*Node)(nil)
	x := t.root
	// 找到新节点应该插入的位置
	for x != nil {
		y = x
		if z.key < x.key {
			x = x.left
		} else {
			x = x.right
		}
	}
	// 设置新节点的父节点
	z.parent = y
	if y == nil {
		// 如果树为空，将新节点设置为根节点
		t.root = z
	} else if z.key < y.key {
		// 如果新节点的键小于其父节点的键，将新节点作为其父节点的左子节点
		y.left = z
	} else {
		// 如果新节点的键大于其父节点的键，将新节点作为其父节点的右子节点
		y.right = z
	}
	// 插入新节点后，进行修复操作以维持红黑树的性质
	t.insertFixup(z)
}

// 查找操作
// 根据键在红黑树中查找对应的节点
func (t *RedBlackTree) Search(key int) (*Node, bool) {
	x := t.root
	for x != nil {
		if key == x.key {
			// 找到匹配的键，返回节点和 true
			return x, true
		} else if key < x.key {
			// 如果键小于当前节点的键，继续在左子树中查找
			x = x.left
		} else {
			// 如果键大于当前节点的键，继续在右子树中查找
			x = x.right
		}
	}
	// 未找到匹配的键，返回 nil 和 false
	return nil, false
}

// 中序遍历
// 中序遍历红黑树，按照左子树、根节点、右子树的顺序访问节点
func inorderTraversal(node *Node) {
	if node != nil {
		// 递归遍历左子树
		inorderTraversal(node.left)
		// 输出当前节点的键、值和颜色
		fmt.Printf("Key: %d, Value: %v, Color: %v\n", node.key, node.value, node.color)
		// 递归遍历右子树
		inorderTraversal(node.right)
	}
}

// 打印红黑树（中序遍历）
// 调用中序遍历函数，打印红黑树的所有节点信息
func (t *RedBlackTree) PrintTree() {
	inorderTraversal(t.root)
}

// 打印红黑树结构图
// 以树形结构打印红黑树，便于直观观察树的结构
func printTreeStructure(node *Node, prefix string, isLeft bool) {
	if node == nil {
		return
	}
	// 根据节点颜色设置颜色字符串
	colorStr := "R"
	if !node.color {
		colorStr = "B"
	}
	// 打印前缀和节点信息
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
// 检查树是否为空，如果为空则输出提示信息，否则调用 printTreeStructure 函数打印树的结构
func (t *RedBlackTree) PrintTreeStructure() {
	if t.root == nil {
		fmt.Println("Empty tree")
		return
	}
	printTreeStructure(t.root, "", false)
}

// 寻找最小节点
// 在以 x 为根的子树中寻找最小键的节点
func minimum(x *Node) *Node {
	for x.left != nil {
		x = x.left
	}
	return x
}

// 删除修复操作
// 在删除节点后，可能会破坏红黑树的性质，需要进行修复操作
func (t *RedBlackTree) deleteFixup(x *Node) {
	if x == nil {
		return
	}
	// 当 x 不是根节点且颜色为黑色时，需要进行调整
	for x != t.root && x.color == BLACK {
		if x == x.parent.left {
			// 如果 x 是其父节点的左子节点
			w := x.parent.right
			if w.color == RED {
				// 情况 1：x 的兄弟节点 w 为红色
				// 将 w 的颜色设置为黑色，x 的父节点颜色设置为红色，然后以 x 的父节点为支点进行左旋
				w.color = BLACK
				x.parent.color = RED
				t.leftRotate(x.parent)
				w = x.parent.right
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				// 情况 2：x 的兄弟节点 w 的左右子节点都为黑色
				// 将 w 的颜色设置为红色，x 指针移动到其父节点
				w.color = RED
				x = x.parent
			} else {
				if w.right.color == BLACK {
					// 情况 3：x 的兄弟节点 w 的右子节点为黑色
					// 将 w 的左子节点颜色设置为黑色，w 的颜色设置为红色，然后以 w 为支点进行右旋
					w.left.color = BLACK
					w.color = RED
					t.rightRotate(w)
					w = x.parent.right
				}
				// 情况 4：x 的兄弟节点 w 的右子节点为红色
				// 将 w 的颜色设置为 x 父节点的颜色，x 的父节点颜色设置为黑色，w 的右子节点颜色设置为黑色，然后以 x 的父节点为支点进行左旋
				w.color = x.parent.color
				x.parent.color = BLACK
				w.right.color = BLACK
				t.leftRotate(x.parent)
				x = t.root
			}
		} else {
			// 如果 x 是其父节点的右子节点
			w := x.parent.left
			if w.color == RED {
				// 情况 1：x 的兄弟节点 w 为红色
				// 将 w 的颜色设置为黑色，x 的父节点颜色设置为红色，然后以 x 的父节点为支点进行右旋
				w.color = BLACK
				x.parent.color = RED
				t.rightRotate(x.parent)
				w = x.parent.left
			}
			if w.right.color == BLACK && w.left.color == BLACK {
				// 情况 2：x 的兄弟节点 w 的左右子节点都为黑色
				// 将 w 的颜色设置为红色，x 指针移动到其父节点
				w.color = RED
				x = x.parent
			} else {
				if w.left.color == BLACK {
					// 情况 3：x 的兄弟节点 w 的左子节点为黑色
					// 将 w 的右子节点颜色设置为黑色，w 的颜色设置为红色，然后以 w 为支点进行左旋
					w.right.color = BLACK
					w.color = RED
					t.leftRotate(w)
					w = x.parent.left
				}
				// 情况 4：x 的兄弟节点 w 的左子节点为红色
				// 将 w 的颜色设置为 x 父节点的颜色，x 的父节点颜色设置为黑色，w 的左子节点颜色设置为黑色，然后以 x 的父节点为支点进行右旋
				w.color = x.parent.color
				x.parent.color = BLACK
				w.left.color = BLACK
				t.rightRotate(x.parent)
				x = t.root
			}
		}
	}
	// 将 x 的颜色设置为黑色
	x.color = BLACK
}

// 删除操作
// 根据键从红黑树中删除对应的节点
func (t *RedBlackTree) Delete(key int) bool {
	// 查找要删除的节点
	z, found := t.Search(key)
	if !found {
		// 未找到要删除的节点，返回 false
		return false
	}
	var y *Node
	var x *Node
	yOriginalColor := z.color
	if z.left == nil {
		// 如果 z 没有左子节点，用其右子节点替换 z
		x = z.right
		t.transplant(z, z.right)
	} else if z.right == nil {
		// 如果 z 没有右子节点，用其左子节点替换 z
		x = z.left
		t.transplant(z, z.left)
	} else {
		// 如果 z 有左右子节点，找到 z 右子树中的最小节点 y 替换 z
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
		// 如果被删除节点或替换节点的原始颜色为黑色，需要进行修复操作
		t.deleteFixup(x)
	}
	return true
}

// 移植操作
// 将节点 v 替换节点 u 在树中的位置
func (t *RedBlackTree) transplant(u, v *Node) {
	if u.parent == nil {
		// 如果 u 是根节点，将 v 设置为新的根节点
		t.root = v
	} else if u == u.parent.left {
		// 如果 u 是其父节点的左子节点，将 v 设置为其父节点的左子节点
		u.parent.left = v
	} else {
		// 如果 u 是其父节点的右子节点，将 v 设置为其父节点的右子节点
		u.parent.right = v
	}
	if v != nil {
		// 如果 v 不为空，更新 v 的父节点指针
		v.parent = u.parent
	}
}

// 修改节点值的函数
// 根据键查找节点，并将其值修改为新的值
func (t *RedBlackTree) Modify(key int, newValue interface{}) bool {
	// 查找要修改的节点
	node, found := t.Search(key)
	if !found {
		// 未找到要修改的节点，返回 false
		return false
	}
	// 修改节点的值
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
