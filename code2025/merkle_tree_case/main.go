/*

Merkle Tree（默克尔树，又称哈希树）是一种基于哈希值的数据结构，由根节点、中间节点和叶子节点组成，叶子节点存储原始数据的哈希，非叶子节点存储其子节点哈希的组合哈希。其核心特性是高效验证数据完整性和快速定位数据差异，以下是其典型使用场景：
1. 区块链与加密货币
这是 Merkle Tree 最广为人知的应用场景：

交易验证：区块链中每个区块包含多笔交易，将所有交易哈希作为叶子节点构建 Merkle Tree，根哈希（Merkle Root）存入区块头。轻节点（如手机钱包）无需下载完整区块数据，只需通过 Merkle 路径（从目标交易哈希到根哈希的路径）即可验证某笔交易是否属于该区块，大幅减少数据传输量。
数据一致性：全节点之间同步区块时，通过比对 Merkle Root 可快速判断区块数据是否一致，若不一致可通过树结构定位差异交易。
2. 分布式文件系统与存储（如 IPFS、BitTorrent）
文件完整性校验：大型文件被分割为多个小块，每个块的哈希作为叶子节点构建 Merkle Tree。下载文件时，通过校验块的哈希与 Merkle Tree 中的对应值，可快速检测是否下载完整或被篡改。
高效同步：分布式节点之间同步文件时，通过比对 Merkle Root 确定是否需要更新；若需更新，通过树结构定位差异块，仅传输变化的部分，减少带宽消耗。
3. 数据库与数据备份
数据一致性检查：数据库分片存储或主从复制时，对各分片数据构建 Merkle Tree，通过比对根哈希验证主从数据是否一致，若不一致可快速定位差异分片。
增量备份：备份系统中，对文件或数据块构建 Merkle Tree，每次备份仅需传输哈希值变化的节点对应的原始数据，提高备份效率。
4. 分布式版本控制系统（如 Git）
Git 中使用类似 Merkle Tree 的结构（结合哈希链）管理文件版本：

每个文件的哈希、目录结构的哈希逐级向上聚合，最终形成一个唯一的 commit 哈希（类似 Merkle Root）。
通过比对 commit 哈希可快速判断两个版本是否一致，通过树结构可高效查找不同版本间的文件差异。
5. 零知识证明（ZKP）与隐私计算
在零知识证明中，Merkle Tree 可用于证明 “某个数据属于一个集合” 而不泄露集合的全部内容。例如，证明 “某笔交易在区块链中” 但不公开其他交易信息，只需提供该交易的 Merkle 路径和根哈希即可。
隐私计算场景中，多方数据联合计算时，通过 Merkle Tree 可验证参与方提供的数据是否未经篡改，同时保护数据隐私。
6. 日志校验与审计
系统日志或审计日志可构建 Merkle Tree，每个日志条目对应一个叶子节点。管理员通过根哈希验证日志是否完整（未被篡改或删除），若日志被修改，对应的叶子节点哈希变化会导致整个路径的哈希变化，最终反映在根哈希上。
定位篡改位置时，通过对比正常与异常 Merkle Tree 的结构，可快速找到被修改的日志条目。
总结
Merkle Tree 的核心价值在于用较小的计算和存储开销实现大规模数据的完整性验证与差异定位，尤其适合分布式系统、需要高频校验数据一致性或保护隐私的场景。其设计巧妙地将 “整体验证” 转化为 “路径验证”，既保证了安全性，又提升了效率。


*/

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
)

// MerkleTree 表示默克尔树结构
type MerkleTree struct {
	Root         *MerkleNode         // 根节点
	Leafs        []*MerkleNode       // 叶子节点
	hashStrategy func([]byte) []byte // 哈希函数策略，返回[]byte类型
}

// MerkleNode 表示默克尔树的节点
type MerkleNode struct {
	Left  *MerkleNode // 左子节点
	Right *MerkleNode // 右子节点
	Data  []byte      // 节点存储的数据（哈希值）
}

// NewMerkleTree 创建一个新的默克尔树
func NewMerkleTree(data [][]byte) *MerkleTree {
	var leafs []*MerkleNode

	// 包装sha256函数，使其返回[]byte而不是[32]byte
	hashStrategy := func(b []byte) []byte {
		hash := sha256.Sum256(b)
		return hash[:]
	}

	// 创建叶子节点
	for _, d := range data {
		hash := hashStrategy(d)
		leafs = append(leafs, &MerkleNode{Left: nil, Right: nil, Data: hash})
	}

	// 构建树
	root := buildTree(leafs, hashStrategy)

	return &MerkleTree{
		Root:         root,
		Leafs:        leafs,
		hashStrategy: hashStrategy,
	}
}

// 递归构建默克尔树
func buildTree(nodes []*MerkleNode, hashStrategy func([]byte) []byte) *MerkleNode {
	// 如果只有一个节点，直接返回（根节点）
	if len(nodes) == 1 {
		return nodes[0]
	}

	var newLevel []*MerkleNode

	// 两两组合节点，计算父节点的哈希
	for i := 0; i < len(nodes); i += 2 {
		// 如果是奇数个节点，最后一个节点与自身组合
		j := i + 1
		if j >= len(nodes) {
			j = i
		}

		// 合并两个子节点的哈希
		combinedHash := append(nodes[i].Data, nodes[j].Data...)
		hash := hashStrategy(combinedHash)

		// 创建父节点
		parent := &MerkleNode{
			Left:  nodes[i],
			Right: nodes[j],
			Data:  hash,
		}

		newLevel = append(newLevel, parent)
	}

	// 递归构建上一层
	return buildTree(newLevel, hashStrategy)
}

// RootHash 返回默克尔树的根哈希（十六进制字符串）
func (t *MerkleTree) RootHash() string {
	if t.Root == nil {
		return ""
	}
	return hex.EncodeToString(t.Root.Data)
}

// GetMerklePath 获取数据在默克尔树中的验证路径
func (t *MerkleTree) GetMerklePath(data []byte) ([][]byte, []int, error) {
	// 计算目标数据的哈希
	targetHash := t.hashStrategy(data)

	// 查找目标叶子节点
	var targetNode *MerkleNode
	for _, leaf := range t.Leafs {
		if hex.EncodeToString(leaf.Data) == hex.EncodeToString(targetHash) {
			targetNode = leaf
			// 移除未使用的index变量，改为直接break
			break
		}
	}

	if targetNode == nil {
		return nil, nil, fmt.Errorf("数据不在默克尔树中")
	}

	// 收集路径和位置（0表示左节点，1表示右节点）
	var path [][]byte
	var positions []int
	currentNode := targetNode
	// 从叶子节点向上遍历到根节点
	for currentNode != t.Root {
		parent, isLeft := findParent(t.Root, currentNode)
		if parent == nil {
			return nil, nil, fmt.Errorf("无法找到父节点")
		}

		// 记录兄弟节点的哈希和当前节点的位置
		if isLeft {
			path = append(path, parent.Right.Data)
			positions = append(positions, 0) // 当前节点是左节点
		} else {
			path = append(path, parent.Left.Data)
			positions = append(positions, 1) // 当前节点是右节点
		}

		currentNode = parent
	}

	return path, positions, nil
}

// 辅助函数：查找节点的父节点并判断是否为左子节点
func findParent(root, node *MerkleNode) (*MerkleNode, bool) {
	if root.Left == nil || root.Right == nil {
		return nil, false
	}

	if root.Left == node {
		return root, true
	}

	if root.Right == node {
		return root, false
	}

	// 递归查找左子树
	if parent, isLeft := findParent(root.Left, node); parent != nil {
		return parent, isLeft
	}

	// 递归查找右子树
	return findParent(root.Right, node)
}

// Verify 验证数据是否属于默克尔树且未被篡改
func (t *MerkleTree) Verify(data []byte, path [][]byte, positions []int) bool {
	// 计算数据的哈希
	currentHash := t.hashStrategy(data)

	// 沿着路径向上计算哈希，直到得到根哈希
	for i, siblingHash := range path {
		var combined []byte
		if positions[i] == 0 {
			// 当前哈希在左，兄弟哈希在右
			combined = append(currentHash, siblingHash...)
		} else {
			// 当前哈希在右，兄弟哈希在左
			combined = append(siblingHash, currentHash...)
		}
		currentHash = t.hashStrategy(combined)
	}

	// 对比计算得到的根哈希与树的根哈希
	return hex.EncodeToString(currentHash) == t.RootHash()
}

func main() {
	// 示例数据：一组交易哈希（模拟区块链中的交易）
	data := [][]byte{
		[]byte("transaction1"),
		[]byte("transaction2"),
		[]byte("transaction3"),
		[]byte("transaction4"),
		[]byte("transaction5"),
	}

	// 创建默克尔树
	merkleTree := NewMerkleTree(data)
	fmt.Printf("默克尔树根哈希: %s\n", merkleTree.RootHash())

	// 验证数据
	targetData := []byte("transaction3")
	path, positions, err := merkleTree.GetMerklePath(targetData)
	if err != nil {
		log.Fatalf("获取验证路径失败: %v", err)
	}

	// 验证数据有效性
	isValid := merkleTree.Verify(targetData, path, positions)
	fmt.Printf("数据 '%s' 验证结果: %v\n", string(targetData), isValid)

	// 验证被篡改的数据
	tamperedData := []byte("transaction3_tampered")
	isValid = merkleTree.Verify(tamperedData, path, positions)
	fmt.Printf("被篡改数据 '%s' 验证结果: %v\n", string(tamperedData), isValid)
}
