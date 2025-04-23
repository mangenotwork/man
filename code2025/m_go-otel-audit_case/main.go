package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/exp/constraints"
	"io"
	"log"
	"net"
	"net/netip"
	"sync"
	"unsafe"
)

func main() {
	fmt.Println("------------ case1_1 ------------")
	case1_1()
	fmt.Println("\n\n------------ case1_2 ------------")
	case1_2()
	fmt.Println("\n\n------------ case1_3 ------------")
	case1_3()
	fmt.Println("\n\n------------ case1_4 ------------")
	case1_4()
	fmt.Println("\n\n------------ case1_5 ------------")
	case1_5()
	fmt.Println("\n\n------------ case1_6 ------------")
	case1_6()
	fmt.Println("\n\n------------ case1_7 ------------")
	case1_7()
	fmt.Println("\n\n------------ case1_8 ------------")
	case1_8()

	fmt.Println("\n\n------------ case2 ------------")
	case2()

	fmt.Println("\n\n------------ case3 ------------")
	case3()

	fmt.Println("\n\n------------ case4 ------------")
	case4()

	fmt.Println("\n\n------------ case5 ------------")
	case5()

	fmt.Println("\n\n------------ case6 ------------")
	case6()

	fmt.Println("\n\n------------ case7 ------------")
	case7()

	fmt.Println("\n\n------------ case8 ------------")
	case8()
}

// case1 net/netip 包使用场景
// Go 1.18 引入了 net/netip 包，这是一个新的网络地址和前缀处理包，设计上比传统的 net.IP 和 net.IPNet 更高效、更安全。
//它的主要目标是提供一个现代化的 API 来处理 IP 地址（包括 IPv4 和 IPv6）以及相关的网络操作。

// 1. 高性能的 IP 地址处理 net/netip 使用值类型（struct）而不是指针类型，因此在内存分配和性能上更加高效。
// 它避免了 net.IP 的动态内存分配问题，适合需要频繁操作 IP 地址的场景
func case1_1() {
	ip := netip.MustParseAddr("192.168.1.1")
	fmt.Println("IP:", ip)
	fmt.Println("Is IPv4:", ip.Is4())
	fmt.Println("Is IPv6:", ip.Is6())
}

// 2. IP 地址验证
func case1_2() {
	addr, err := netip.ParseAddr("256.256.256.256") // 无效地址
	if err != nil {
		log.Println("Invalid IP:", err)
		return
	}
	fmt.Println("Valid IP:", addr)
}

// 3. IP 前缀和子网操作
func case1_3() {
	prefix := netip.MustParsePrefix("192.168.1.0/24")
	fmt.Println("Prefix:", prefix)
	fmt.Println("Network address:", prefix.Addr())
	fmt.Println("Mask length:", prefix.Bits())
}

// 4. IP 地址比较
func case1_4() {
	ip1 := netip.MustParseAddr("192.168.1.1")
	ip2 := netip.MustParseAddr("192.168.1.2")

	fmt.Println("IP1 < IP2:", ip1.Less(ip2))
	fmt.Println("IP1 == IP2:", ip1 == ip2)
}

// 5. IPv4 和 IPv6 的统一处理
func case1_5() {
	ipv4 := netip.MustParseAddr("192.168.1.1")
	ipv6 := netip.MustParseAddr("::1")

	fmt.Println("IPv4:", ipv4)
	fmt.Println("IPv6:", ipv6)
	fmt.Println("IPv4 is IPv6-mapped?", ipv4.Is4In6())
	fmt.Println("IPv6 is loopback?", ipv6.IsLoopback())
}

// 6. 范围检查
func case1_6() {
	start := netip.MustParseAddr("192.168.1.1")
	end := netip.MustParseAddr("192.168.1.10")
	ip := netip.MustParseAddr("192.168.1.5")

	fmt.Println("InRange:", !ip.Less(start) && !end.Less(ip))
}

// 7. 与现有 net 包的互操作性
func case1_7() {
	ip := net.ParseIP("192.168.1.1")
	addr, _ := netip.AddrFromSlice(ip)
	fmt.Println("netip.Addr:", addr)

	stdIP := addr.AsSlice()
	fmt.Println("net.IP:", net.IP(stdIP))
}

// 8. 线程安全性
// 由于 net/netip 使用的是值类型，所有的操作都是无状态的，天然具有线程安全性。这使得它非常适合在并发环境中使用。
func case1_8() {
	var wg sync.WaitGroup
	ip := netip.MustParseAddr("192.168.1.1")

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(ip)
		}()
	}

	wg.Wait()
}

// 总结
/*
net/netip 是一个现代化的网络地址处理包，提供了以下优势：

高效：值类型设计，减少内存分配。
安全：内置验证，避免无效地址。
统一：简化 IPv4 和 IPv6 的处理。
易用：直观的 API 和丰富的功能。
适合的使用场景包括：

高性能网络服务。
IP 地址验证和解析。
子网和范围操作。
并发环境下的线程安全操作。
*/

// case2
// golang.org/x/exp/constraints 是 Go 语言实验性扩展库中的一个包，属于 golang.org/x/exp 项目。
//它为 Go 的泛型编程提供了通用的约束（constraints）定义，用于简化泛型代码的编写
/*
主要功能
constraints 包定义了一些常用的类型约束，这些约束可以直接用于泛型代码中，避免重复定义。例如：

constraints.Integer：表示所有整数类型（如 int, int8, uint16 等）。
constraints.Float：表示所有浮点数类型（如 float32, float64）。
constraints.Complex：表示所有复数类型（如 complex64, complex128）。
constraints.Ordered：表示所有可以比较大小的类型（如整数、浮点数、字符串）。
constraints.Signed：表示所有有符号整数类型（如 int, int8, int16 等）。
constraints.Unsigned：表示所有无符号整数类型（如 uint, uint8, uint16 等）。
这些约束通过组合 Go 的接口和类型集（type set）来实现。
*/

// 定义一个泛型函数，支持所有整数类型
func Max[T constraints.Integer](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// 定义一个泛型函数，支持所有有序类型
func Compare[T constraints.Ordered](a, b T) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

// 定义一个自定义约束，支持整数和浮点数
type Numeric interface {
	constraints.Integer | constraints.Float
}

// 泛型函数，支持整数和浮点数
func Sum[T Numeric](values ...T) T {
	var total T
	for _, v := range values {
		total += v
	}
	return total
}

func case2() {
	fmt.Println(Max(10, 20))     // 输出: 20
	fmt.Println(Max(int8(5), 3)) // 输出: 5

	fmt.Println(Compare(10, 20))            // 输出: -1
	fmt.Println(Compare(3.14, 2.71))        // 输出: 1
	fmt.Println(Compare("apple", "banana")) // 输出: -1

	fmt.Println(Sum(1, 2, 3))       // 输出: 6
	fmt.Println(Sum(1.5, 2.5, 3.0)) // 输出: 7.0

}

// case3
// 将字节片转换为没有副本的字符串（即没有分配）。
// 这使用了unsafe，因此仅在您知道字节片不会被修改的情况下使用。
func bytesToStr(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
func case3() {
	fmt.Println(bytesToStr([]byte("adssadsad")))
}

// case4 格式化动词 %q 用于以双引号括起来的字符串形式输出数据，并对字符串中的特殊字符进行转义处理
func case4() {
	// 字符串示例
	fmt.Printf("%q\n", "Hello, World!")
	// 输出: "Hello, World!"

	// 包含特殊字符的字符串示例
	fmt.Printf("%q\n", "Hello\nWorld!")
	// 输出: "Hello\nWorld!"

	// 非字符串类型示例
	fmt.Printf("%q\n", 65)
	// 输出: ASCII 值

	// 字符类型示例
	fmt.Printf("%q\n", 'A')
	// 输出: 'A' (注意：对于单个字符，它不会加上双引号而是直接显示字符)
}

// case5 PBKDF2（Password-Based Key Derivation Function 2）是一种密钥导出函数，常用于从密码生成加密密钥。
// 它通过使用哈希函数、加盐（salt）以及多次迭代来增加暴力破解的难度，从而提高密码的安全性
func case5() {
	password := "my-password" // 用户提供的密码
	salt := make([]byte, 16)  // 生成一个随机的盐值
	rand.Read(salt)           // 使用 crypto/rand 生成安全的随机数

	// 使用 pbkdf2 生成密钥
	key := pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)

	fmt.Printf("Salt: %x\n", salt)
	fmt.Printf("Key: %x\n", key)
}

// case6 使用 pbkdf2 对密码加密 解密
// 要实现加密和解密，通常需要结合 PBKDF2 和其他加密算法（如 AES）。以下是一个完整的示例，
// 展示如何使用 pbkdf2 生成密钥，并结合 AES 实现加密和解密
const (
	keyLength  = 32   // AES-256 密钥长度为 32 字节
	saltLength = 16   // 盐值长度
	iterations = 4096 // 迭代次数
)

// 使用 PBKDF2 根据密码和盐值生成密钥
func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, iterations, keyLength, sha256.New)
}

// 加密函数
func encrypt(plaintext []byte, password string) (ciphertext []byte, salt []byte, err error) {
	// 生成随机盐值
	salt = make([]byte, saltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, nil, err
	}

	// 派生密钥
	key := deriveKey(password, salt)

	// 创建 AES 密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	// 生成随机 IV（初始化向量）
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}

	// 创建加密器
	stream := cipher.NewCFBEncrypter(block, iv)

	// 加密数据
	ciphertext = make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)

	// 将 IV 附加到密文前面
	ciphertext = append(iv, ciphertext...)

	return ciphertext, salt, nil
}

// 解密函数
func decrypt(ciphertext []byte, password string, salt []byte) (plaintext []byte, err error) {
	// 派生密钥
	key := deriveKey(password, salt)

	// 提取 IV（初始化向量）
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// 创建 AES 密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建解密器
	stream := cipher.NewCFBDecrypter(block, iv)

	// 解密数据
	plaintext = make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

func case6() {
	password := "my-strong-password"
	plaintext := []byte("This is a secret message!")

	// 加密
	ciphertext, salt, err := encrypt(plaintext, password)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return
	}
	fmt.Printf("Salt: %x\n", salt)
	fmt.Printf("Ciphertext: %x\n", ciphertext)

	// 解密
	decrypted, err := decrypt(ciphertext, password, salt)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}
	fmt.Printf("Decrypted text: %s\n", decrypted)
}

// case7 使用pbkdf2实现生成密钥验证密钥

// 生成随机盐值
func generateSalt() ([]byte, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// 生成并存储密钥（模拟注册过程）
func generateAndStoreKey(password string) (string, []byte, error) {
	// 生成随机盐值
	salt, err := generateSalt()
	if err != nil {
		return "", nil, err
	}

	// 派生密钥
	key := deriveKey(password, salt)

	// 将密钥转换为十六进制字符串以便存储
	hexKey := hex.EncodeToString(key)
	return hexKey, salt, nil
}

// 验证用户输入的密码是否正确
func verifyKey(storedKey string, salt []byte, password string) bool {
	// 派生密钥
	key := deriveKey(password, salt)

	// 将派生的密钥转换为十六进制字符串
	hexKey := hex.EncodeToString(key)

	// 比较生成的密钥与存储的密钥
	return hexKey == storedKey
}

func case7() {
	// 模拟用户注册
	password := "my-strong-password"
	storedKey, salt, err := generateAndStoreKey(password)
	if err != nil {
		fmt.Println("Error generating key:", err)
		return
	}
	fmt.Printf("Stored Key: %s\n", storedKey)
	fmt.Printf("Salt: %x\n", salt)

	// 模拟用户登录
	loginPassword := "my-strong-password"
	isValid := verifyKey(storedKey, salt, loginPassword)
	if isValid {
		fmt.Println("times 1 Password is valid!")
	} else {
		fmt.Println("times 1 Invalid password!")
	}

	// 模拟用户登录2
	loginPassword2 := "my-strong-password1"
	isValid2 := verifyKey(storedKey, salt, loginPassword2)
	if isValid2 {
		fmt.Println("times 2 Password is valid!")
	} else {
		fmt.Println("times 2 Invalid password!")
	}

}

// case8 随机返回几个bits
func case8() {
	randBits := func(n uint) ([]byte, error) {
		b := make([]byte, n)
		_, err := rand.Read(b)
		return b, err
	}
	for i := 0; i < 10; i++ {
		data, _ := randBits(10)
		fmt.Printf("Hex String: %x\n", data)
	}

}
