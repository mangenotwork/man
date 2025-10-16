package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/zeebo/blake3" // 主流的Blake3 Go实现库
)

// 支持的哈希算法
const (
	AlgorithmMD5    = "md5"
	AlgorithmSHA1   = "sha1"
	AlgorithmSHA256 = "sha256"
)

// CalculateFileHash 流式计算文件的哈希值
// 特点：不会一次性加载整个文件到内存，适合大文件处理
func CalculateFileHash(filePath string, algorithm string) (string, error) {
	// 输入参数验证
	if filePath == "" {
		return "", errors.New("文件路径不能为空")
	}
	if algorithm == "" {
		return "", errors.New("哈希算法不能为空")
	}

	// 检查文件是否存在且可访问
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close() // 确保文件最终会被关闭

	// 验证文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("无法获取文件信息: %w", err)
	}
	if fileInfo.IsDir() {
		return "", errors.New("指定路径是目录，不是文件")
	}
	if fileInfo.Size() == 0 {
		// 空文件也能计算哈希，这里仅作为信息提示，不返回错误
		fmt.Printf("警告: 计算空文件 %s 的哈希值\n", filePath)
	}

	// 创建对应的哈希计算器
	var hasher hash.Hash
	switch algorithm {
	case AlgorithmMD5:
		hasher = md5.New()
	case AlgorithmSHA1:
		hasher = sha1.New()
	case AlgorithmSHA256:
		hasher = sha256.New()
	default:
		return "", fmt.Errorf("不支持的哈希算法: %s，支持的算法有 %s, %s, %s",
			algorithm, AlgorithmMD5, AlgorithmSHA1, AlgorithmSHA256)
	}

	// 流式处理：分块读取文件并更新哈希
	// io.Copy 使用io.Copy实现标准流式处理，内部会自动管理缓冲区 内部会使用默认缓冲区(通常是32KB)，比手动实现更高效
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("计算哈希过程出错: %w", err)
	}

	// 生成最终哈希值的十六进制表示
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

// -------------- 扩展学习   blake3算法 是什么?
/*

Blake3 本质上是一款由 Jack O'Connor 等人开发的加密哈希算法，基于 Blake2 算法改进而来，2020 年正式发布，主打 “高性能 + 高安全性” 的平衡，
尤其适合需要高效处理大文件或海量数据的场景。

Blake3 的核心特点  极致速度  高安全性  功能灵活
Blake3 的典型应用场景   文件校验   密码存储   数据同步  区块链 / 加密货币

*/

// CalculateFileBlake3 流式计算文件的Blake3哈希值
// 特点：分块读取文件，内存占用恒定，适合GB级大文件
func CalculateFileBlake3(filePath string) (string, error) {
	// 输入参数验证
	if filePath == "" {
		return "", errors.New("文件路径不能为空")
	}

	// 打开文件（只读模式）
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close() // 确保文件最终关闭，释放资源

	// 验证文件属性
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败: %w", err)
	}
	if fileInfo.IsDir() {
		return "", errors.New("指定路径是目录，不是文件")
	}

	// 初始化Blake3哈希器
	hasher := blake3.New()

	// 流式处理：使用io.Copy高效分块读写
	// 内部缓冲区默认32KB，可自动适配不同大小文件
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("计算哈希过程出错: %w", err)
	}

	// 生成最终哈希值（默认32字节，转为64位十六进制字符串）
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

// ------------------------- 哈希运行需要一次性取出整个文件吗？
/*

不需要一次性取出整个文件，主流哈希算法（如 Blake3、SHA-256、MD5）都支持 “流式计算”，只需分块读取文件数据逐步更新哈希状态，
内存占用始终保持在固定大小（仅缓冲区或算法内部状态大小），完全适配大文件场景。

*/
