package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"sort"
	"strings"
	"sync"
)

// 一致性哈希环结构
type ConsistentHash struct {
	ring         map[uint64]string // 哈希值 -> 表名
	sortedHashes []uint64          // 排序后的哈希值，用于快速查找
	tables       map[string]bool   // 记录已添加的表（去重）
	virtualNodes int               // 每个真实表对应的虚拟节点数量
	mu           sync.RWMutex      // 并发安全锁
	metadata     map[string]string // 企业名称（标准化）-> 实际存储表名
}

// 创建一致性哈希实例
func NewConsistentHash(virtualNodes int) *ConsistentHash {
	return &ConsistentHash{
		ring:         make(map[uint64]string),
		tables:       make(map[string]bool),
		virtualNodes: virtualNodes,
		metadata:     make(map[string]string),
	}
}

// 添加表（支持动态增加，避免重复添加）
func (c *ConsistentHash) AddTable(tableName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 避免重复添加表
	if c.tables[tableName] {
		return fmt.Errorf("表 %s 已存在", tableName)
	}

	// 为每个真实表创建多个虚拟节点（使用固定规则生成虚拟节点，确保表名变化时映射可预期）
	for i := 0; i < c.virtualNodes; i++ {
		// 虚拟节点生成规则：表名+固定前缀+序号（增强稳定性）
		virtualNode := fmt.Sprintf("tbl_%s_%d", tableName, i)
		hash := c.hash(virtualNode)
		c.ring[hash] = tableName
	}

	c.tables[tableName] = true
	c.refreshSortedHashes()
	return nil
}

// 哈希计算（内部方法，确保一致性）
func (c *ConsistentHash) hash(key string) uint64 {
	h := sha256.Sum256([]byte(key))
	return binary.BigEndian.Uint64(h[:8])
}

// 刷新排序的哈希列表
func (c *ConsistentHash) refreshSortedHashes() {
	hashes := make([]uint64, 0, len(c.ring))
	for h := range c.ring {
		hashes = append(hashes, h)
	}
	sort.Slice(hashes, func(i, j int) bool {
		return hashes[i] < hashes[j]
	})
	c.sortedHashes = hashes
}

// 标准化企业名称
func (c *ConsistentHash) normalizeCompanyName(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}

// 获取数据存储表名（核心方法）
func (c *ConsistentHash) GetStorageTable(companyName string) string {
	normalizedName := c.normalizeCompanyName(companyName)
	if normalizedName == "" {
		return ""
	}

	// 优先从元数据获取（保证历史数据稳定）
	c.mu.RLock()
	if table, exists := c.metadata[normalizedName]; exists {
		c.mu.RUnlock()
		return table
	}
	c.mu.RUnlock()

	// 元数据无记录时，计算哈希并写入元数据
	c.mu.Lock()
	defer c.mu.Unlock()

	// 双重检查
	if table, exists := c.metadata[normalizedName]; exists {
		return table
	}

	// 计算哈希映射
	hash := c.hash(normalizedName)
	idx := sort.Search(len(c.sortedHashes), func(i int) bool {
		return c.sortedHashes[i] >= hash
	})
	if idx == len(c.sortedHashes) {
		idx = 0
	}
	tableName := c.ring[c.sortedHashes[idx]]

	// 写入元数据
	c.metadata[normalizedName] = tableName
	return tableName
}

// 更新元数据（用于数据迁移）
func (c *ConsistentHash) UpdateMetadata(companyName, newTable string) {
	normalizedName := c.normalizeCompanyName(companyName)
	if normalizedName == "" {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.metadata[normalizedName] = newTable
}

var (
	dbHost          = "10.0.40.9"
	port            = 3306
	user            = "root"
	password        = "123456"
	spiderBase      = "test"
	SpiderBaseDB, _ = gt.NewMysql(dbHost, port, user, password, spiderBase)
)

func main() {
	hashRing := NewConsistentHash(10)

	// 初始表：3个
	initTables := []string{"data_0", "data_1", "data_2"}
	for _, table := range initTables {
		hashRing.AddTable(table)
	}
	fmt.Println("初始表数量:", len(hashRing.tables))

	// 模拟存储一批旧数据（元数据会记录它们的存储位置）
	companies := []string{
		"阿里巴巴集团",
		"腾讯科技有限公司",
		"百度在线网络技术有限公司",
	}
	fmt.Println("\n=== 初始存储 ===")
	for _, company := range companies {
		table := hashRing.GetStorageTable(company)
		inData := map[string]interface{}{
			"name": company,
		}
		err := SpiderBaseDB.InsertAt(table, inData)
		if err != nil {
			gt.Error(err)
		}
		fmt.Printf("企业: %-20s 存储到: %s\n", company, table)
	}

	// 动态扩展：新增第4个表
	newTable := "data_3"
	hashRing.AddTable(newTable)
	fmt.Println("\n=== 新增表后（总表数:", len(hashRing.tables), "）===")

	// 1. 查询旧数据：仍指向原表（元数据保证）
	fmt.Println("\n旧数据查询:")
	for _, company := range companies {
		table := hashRing.GetStorageTable(company)
		fmt.Printf("企业: %-20s 实际存储表: %s\n", company, table)
		outData, err := SpiderBaseDB.Select(fmt.Sprintf("select * from %s where name=\"%s\"", table, company))
		gt.Info(outData, err)
	}

	// 2. 新增数据：按新哈希环映射（可能指向新表）
	newCompanies := []string{
		"华为技术有限公司",
		"字节跳动有限公司",
	}
	fmt.Println("\n新数据存储:")
	for _, company := range newCompanies {
		table := hashRing.GetStorageTable(company)
		fmt.Printf("企业: %-20s 存储到: %s\n", company, table)
		inData := map[string]interface{}{
			"name": company,
		}
		err := SpiderBaseDB.InsertAt(table, inData)
		if err != nil {
			gt.Error(err)
		}
	}
}
