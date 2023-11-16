package main

/*
树排序是一种在线排序算法。它使用二叉搜索树数据结构来存储元素。通过按顺序遍历二叉搜索树，可以按排序顺序检索元素。
由于它是一种在线排序算法，因此插入的元素始终按排序顺序进行维护。

二、实现逻辑
假设使用一组未排序的数组 array 包含 n 个元素。

算法主体的步骤：
1. 通过在二叉搜索树中插入数组中的元素来构建二进制搜索树；
2. 在树上执行顺序遍历，以使元素按排序顺序返回。

插入排序的步骤：
1. 创建一个BST节点，其值等于数组元素 array[i]；
2. Insert(node, key): 如果 root == null，那么返回新形成的节点；如果 root->data < key，那么
root->right = insert(root->right, key)；如果 root->data > key，那么 root->left = insert(root->left, key)；
3. 返回指向原始根结点的游标。

顺序遍历操作：遍历左子树 → 访问根结点 → 遍历右子树。
*/

// definition of a bst node
type node struct {
	val   int
	left  *node
	right *node
}

// definition of a node
type btree struct {
	root *node
}

// allocating a new node
func newNode(val int) *node {
	return &node{val, nil, nil}
}

// insert nodes into a binary search tree
func insert(root *node, val int) *node {
	if root == nil {
		return newNode(val)
	}
	if val < root.val {
		root.left = insert(root.left, val)
	} else {
		root.right = insert(root.right, val)
	}
	return root
}

// inorder traversal algorithm
// Copies the elements of the bst to the array in sorted order
func inorderCopy(n *node, array []int, index *int) {
	if n != nil {
		inorderCopy(n.left, array, index)
		array[*index] = n.val
		*index++
		inorderCopy(n.right, array, index)
	}
}

func TreeSort(array []int, tree *btree) {
	// build the binary search tree
	for _, element := range array {
		tree.root = insert(tree.root, element)
	}
	index := 0
	// perform inorder traversal to get the elements in sorted order
	inorderCopy(tree.root, array, &index)
}
