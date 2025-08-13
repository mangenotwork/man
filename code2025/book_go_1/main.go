package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime/pprof"
)

// golang学习 《Head First Go语言程序设计》

func main() {
	//case1()

	//case2()

	//case3()

	//case4()

	case5()
}

// 指针
func case1() {
	var a int = 10
	ap := &a
	println(ap)
	println(*ap)
	*ap = 20
	println(ap)
	println(a)
	a = 30
	println(a)
	pf(&a)     // 值被更新了
	println(a) // 31
	pf2(a)     // 值没有被更新
	println(a) // 31
	// 小结: 想在函数没有返回值没有重新分配内存空间下传入参数使用指针类型
	b := float64(*ap)
	fmt.Printf("%T", b)

}

func pf(a *int) {
	*a += 1
	println(a) // 指针指向传入的内存
}

func pf2(a int) {
	a += 1
	println(&a) // 重新分配了内存
}

func case2() {
	type a struct {
		a1 int
		a2 string
	}

	s := &a{1, "a"}

	s1 := "s1"

	fmt.Printf("%#v\n", s)
	fmt.Printf("%#v\n", s1)
}

// 判断变量是否为结构体类型
func isStruct(v interface{}) bool {
	t := reflect.TypeOf(v)

	// 处理指针情况，获取指针指向的元素类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Kind() == reflect.Struct
}

type Number int

// 值类型  接收参数 n 是不会被修改的
func (n Number) Tow() {
	n *= 2
}

// 指针类型 接收参数 n指针 会被修改
func (n *Number) TowP() {
	*n *= 2
}

func case3() {
	nt1 := Number(10)
	nt1.Tow()
	println(nt1) // 原本值不会被修改 不会翻倍

	nt2 := Number(10)
	nt2.TowP()
	println(nt2) // 传入参数是指针 会翻倍
}

// 注：  为了一致性，你所有的类型函数接受值类型，或者都接受指针类型，但是你应该避免混用的情况

type data1 struct {
	a int
}

// 接受值类型不会被修改，是拷贝不是原值;
func (d data1) set1(a int) {
	d.a = a
}

// 接受指针类型会修改原值,指向原值并修改
func (d *data1) set2(a int) {
	d.a = a
}

func case4() {
	d1 := data1{1}
	d1.set1(2)
	println(d1.a) // 1

	d1.set2(3)
	println(d1.a) // 3
}

// 将程序中的数据隐藏在一部分代码中而对另一部分不可见的方法称为封装

// case5  手动生成 profile 文件（适合短期程序）
// 对于一些短期运行的程序（如脚本、定时任务），没办法通过 HTTP 接口动态采集数据，这时可以手动生成 profile 文件。
func case5() {

	heavyTask := func() {
		// 模拟CPU密集型任务
		sum := 0
		for i := 0; i < 1e8; i++ {
			sum += i
		}
	}

	// 生成CPU profile文件
	cpuFile, _ := os.Create("cpu.pprof")

	defer cpuFile.Close()

	_ = pprof.StartCPUProfile(cpuFile) // 开始CPU采样

	defer pprof.StopCPUProfile() // 程序结束时停止采样

	// 生成堆内存profile文件

	heapFile, _ := os.Create("heap.pprof")

	defer heapFile.Close()

	defer pprof.WriteHeapProfile(heapFile) // 程序结束前写入堆内存数据

	// 你的业务逻辑，比如一段耗时的计算
	heavyTask()

}

//  go tool pprof 工具使用
//  安装 Graphviz 并配置环境变量
// go tool pprof cpu.pprof
// web
// 火焰图go-torch   需要安装 FlameGraph 还有perl

// AES-128 加密字符串到文件
func case6() {

	os.Open("")

	// 示例使用
	secretKey := []byte("mysecretkey1234") // 16字节密钥
	originalText := "这是一段需要加密保存的敏感信息: 123456"
	outputFile := "encrypted_data.txt"

	// 加密并保存到文件
	err := Aes128EncryptToFile(originalText, secretKey, outputFile)
	if err != nil {
		fmt.Printf("加密失败: %v\n", err)
		return
	}
	fmt.Printf("成功将加密数据保存到 %s\n", outputFile)

	// 从文件读取并解密
	decryptedText, err := Aes128DecryptFromFile(secretKey, outputFile)
	if err != nil {
		fmt.Printf("解密失败: %v\n", err)
		return
	}
	fmt.Println("解密后的内容:")
	fmt.Println(decryptedText)
}

// Aes128EncryptToFile 加密字符串并将结果保存到文件
// plaintext: 待加密的字符串
// key: 16字节的AES-128密钥
// filename: 保存加密结果的文件名
func Aes128EncryptToFile(plaintext string, key []byte, filename string) error {
	// 验证密钥长度
	if len(key) != 16 {
		return errors.New("AES-128 要求密钥长度必须为16字节")
	}

	// 创建AES密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("创建密码块失败: %v", err)
	}

	// 对明文进行PKCS#7填充
	paddingText := pkcs7Padding([]byte(plaintext), aes.BlockSize)

	// 生成随机IV(16字节)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return fmt.Errorf("生成IV失败: %v", err)
	}

	// 创建CBC加密模式
	mode := cipher.NewCBCEncrypter(block, iv)

	// 执行加密
	ciphertext := make([]byte, len(paddingText))
	mode.CryptBlocks(ciphertext, paddingText)

	// 组合IV和密文(IV需要用于解密，存储在密文前)
	result := append(iv, ciphertext...)

	// 将加密结果写入文件
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	// 写入Base64编码后的结果，方便存储和查看
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(result)))
	base64.StdEncoding.Encode(encoded, result)

	if _, err := file.Write(encoded); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// Aes128DecryptFromFile 从文件读取加密内容并解密
// key: 16字节的AES-128密钥
// filename: 存储加密结果的文件名
func Aes128DecryptFromFile(key []byte, filename string) (string, error) {
	// 验证密钥长度
	if len(key) != 16 {
		return "", errors.New("AES-128 要求密钥长度必须为16字节")
	}

	// 读取文件内容
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	// Base64解码
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(decoded, data)
	if err != nil {
		return "", fmt.Errorf("Base64解码失败: %v", err)
	}
	decoded = decoded[:n]

	// 分离IV和密文
	if len(decoded) < aes.BlockSize {
		return "", errors.New("加密数据长度不足")
	}
	iv := decoded[:aes.BlockSize]
	ciphertext := decoded[aes.BlockSize:]

	// 创建AES密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建密码块失败: %v", err)
	}

	// 创建CBC解密模式
	mode := cipher.NewCBCDecrypter(block, iv)

	// 执行解密
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除PKCS#7填充
	plaintext, err = pkcs7Unpadding(plaintext)
	if err != nil {
		return "", fmt.Errorf("去除填充失败: %v", err)
	}

	return string(plaintext), nil
}

// PKCS#7填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// PKCS#7去填充
func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("数据长度为0")
	}
	padding := int(data[length-1])
	if padding > length {
		return nil, errors.New("无效的填充数据")
	}
	return data[:length-padding], nil
}

// 生成安全的16字节AES密钥
func generateSecureKey() []byte {
	key := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err)
	}
	return key
}
