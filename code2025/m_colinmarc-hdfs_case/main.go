package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

func main() {
	case1()
	case2()
	case3()
	case4()
	case5()
	case6()
}

// case1
// file.Seek 是 os.File 类型的一个方法，其用途是设置文件的读写偏移量，也就是确定接下来读写操作从文件的哪个位置开始。
// func (f *File) Seek(offset int64, whence int) (ret int64, err error)
// offset：这是一个 int64 类型的值，代表要移动的字节数。正值表示向文件末尾方向移动，负值表示向文件开头方向移动。
// whence：这是一个 int 类型的值，用于指定偏移量的起始位置，有以下三种取值：
// io.SeekStart（值为 0）：从文件开头开始计算偏移量。
// io.SeekCurrent（值为 1）：从当前文件读写位置开始计算偏移量。
// io.SeekEnd（值为 2）：从文件末尾开始计算偏移量。
func case1() {
	// 打开文件
	file, err := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("打开文件出错: %v\n", err)
		return
	}
	defer file.Close()

	// 写入一些数据
	data := []byte("Hello, World!")
	_, err = file.Write(data)
	if err != nil {
		fmt.Printf("写入文件出错: %v\n", err)
		return
	}

	// 从文件开头偏移 7 个字节
	newOffset, err := file.Seek(7, io.SeekStart)
	if err != nil {
		fmt.Printf("Seek 操作出错: %v\n", err)
		return
	}
	fmt.Printf("新的偏移量（从文件开头）: %d\n", newOffset)

	// 从当前位置偏移 2 个字节
	newOffset, err = file.Seek(2, io.SeekCurrent)
	if err != nil {
		fmt.Printf("Seek 操作出错: %v\n", err)
		return
	}
	fmt.Printf("新的偏移量（从文件开头）: %d\n", newOffset)

	// 从文件末尾偏移 -6 个字节
	newOffset, err = file.Seek(-6, io.SeekEnd)
	if err != nil {
		fmt.Printf("Seek 操作出错: %v\n", err)
		return
	}
	fmt.Printf("新的偏移量（从文件开头）: %d\n", newOffset)

	// 读取从当前偏移位置开始的数据
	buffer := make([]byte, 5)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Printf("读取文件出错: %v\n", err)
		return
	}
	fmt.Printf("读取的数据: %s\n", string(buffer[:n]))
}

// case2
// os.IsNotExist 是 os 包提供的一个实用函数，其主要用途是判断一个错误是否由文件或目录不存在所导致
func case2() {
	// 尝试打开一个不存在的文件
	_, err := os.Open("nonexistent_file.txt")
	if os.IsNotExist(err) {
		fmt.Println("文件不存在")
	} else if err != nil {
		fmt.Printf("打开文件时出现其他错误: %v\n", err)
	} else {
		fmt.Println("文件打开成功")
	}

	// 尝试获取一个不存在的目录的信息
	_, err = os.Stat("nonexistent_directory")
	if os.IsNotExist(err) {
		fmt.Println("目录不存在")
	} else if err != nil {
		fmt.Printf("获取目录信息时出现其他错误: %v\n", err)
	} else {
		fmt.Println("目录存在")
	}
}

// case3  字节转字符串
func formatBytes(i uint64) string {
	switch {
	case i > (1024 * 1024 * 1024 * 1024):
		return fmt.Sprintf("%#.1fT", float64(i)/1024/1024/1024/1024)
	case i > (1024 * 1024 * 1024):
		return fmt.Sprintf("%#.1fG", float64(i)/1024/1024/1024)
	case i > (1024 * 1024):
		return fmt.Sprintf("%#.1fM", float64(i)/1024/1024)
	case i > 1024:
		return fmt.Sprintf("%#.1fK", float64(i)/1024)
	default:
		return fmt.Sprintf("%dB", i)
	}
}
func case3() {
	fmt.Println(formatBytes(51161112))
	fmt.Println(formatBytes(84545))
	fmt.Println(formatBytes(123))
}

//case4
/*
使用 io.ReadFull 读取数据：调用 io.ReadFull 函数从 reader 中读取数据到 buffer 中，根据返回的 n 和 err 进行相应的处理。
若 err 为 io.ErrUnexpectedEOF，表示没有读满缓冲区就遇到了文件结束。
若 err 不为 nil 且不是 io.ErrUnexpectedEOF，表示读取过程中出现了其他错误。
若 err 为 nil，表示成功将缓冲区填满。

适用场景
固定长度数据读取：当你需要读取固定长度的数据时，使用 io.ReadFull 可以确保读取到足够的数据，避免数据不足的问题。
数据完整性检查：在某些场景下，需要确保读取的数据长度符合预期，io.ReadFull 可以帮助你进行数据完整性检查。
*/
func case4() {
	// 创建一个字节切片作为数据源
	data := []byte("Hello, World!")
	// 将字节切片包装成 io.Reader
	reader := bytes.NewReader(data)

	// 创建一个缓冲区，大小为 5 字节
	buffer := make([]byte, 50)

	// 使用 io.ReadFull 读取数据
	n, err := io.ReadFull(reader, buffer)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			fmt.Printf("未读满缓冲区，实际读取了 %d 字节\n", n)
		} else {
			fmt.Printf("读取数据时出错: %v\n", err)
		}
	} else {
		fmt.Printf("成功读取 %d 字节数据: %s\n", n, string(buffer))
	}
}

// case5
/*
binary.BigEndian.PutUint32 和 binary.BigEndian.PutUint16 是用于将无符号整数转换为大端字节序字节数据的实用方法，
在处理网络通信、文件格式等需要特定字节序的场景中非常有用。
*/
func case5() {
	var num uint32 = 0x12345678
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, num)
	fmt.Printf("转换后的字节数据: %x\n", buf)

	var num16 uint16 = 0x1234
	buf16 := make([]byte, 2)
	binary.BigEndian.PutUint16(buf16, num16)
	fmt.Printf("转换后的字节数据: %x\n", buf16)
}

// case6
func case6() {
	for i := 0; i < 10000; i++ {
		fmt.Println(string(newClientID()))
		time.Sleep(100 * time.Millisecond)
	}
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func newClientID() []byte {
	id := make([]byte, 128)

	rand.NewSource(time.Now().UTC().UnixNano())
	for i := range id {
		id[i] = chars[rand.Intn(len(chars))]
	}

	return id
}
