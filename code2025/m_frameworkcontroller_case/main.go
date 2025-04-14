package main

import (
	"bytes"
	"cmp"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows/registry"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/klog/v2"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	case1()

	case2()

	case3()

	case6()

	case9()
}

// case1 通过使用 k8s.io/apimachinery/pkg/types 库，你可以更方便地处理 Kubernetes 资源的命名空间、名称和唯一标识符，
// 从而更好地与 Kubernetes API 进行交互
func case1() {
	// 创建一个 NamespacedName 实例
	nn := types.NamespacedName{
		Namespace: "default",
		Name:      "my-pod",
	}

	// 打印 NamespacedName 的信息
	fmt.Printf("Namespace: %s, Name: %s\n", nn.Namespace, nn.Name)
}

// case2 通过运用 k8s.io/apimachinery/pkg/util/uuid 库，你能够方便地生成符合标准的 UUID，以此来满足 Kubernetes
// 系统里对资源唯一标识的需求
func case2() {
	// 生成一个新的 UUID
	newUUID := uuid.NewUUID()
	fmt.Printf("生成的 UUID: %s\n", newUUID)
}

// case3 k8s.io/klog 是 Kubernetes 项目使用的日志记录库，它提供了一套灵活且强大的日志功能，方便开发者在 Kubernetes
// 相关的应用程序、控制器、插件等中进行日志记录和调试
func case3() {
	// 初始化 klog
	klog.InitFlags(nil)

	// 记录不同级别的日志
	klog.Info("这是一条信息级别的日志")
	klog.Warning("这是一条警告级别的日志")
	klog.Error("这是一条错误级别的日志")

	// 详细日志级别，通过 --v 标志控制
	klog.V(2).Info("这是详细级别 2 的日志")

	// 程序结束时，刷新日志缓冲区
	klog.Flush()
}

// case4 通过使用 k8s.io/apimachinery/pkg/util/net 库，你可以更方便地处理 Kubernetes 中的网络操作，提高代码的健壮性和可维护性

// case5 通过使用 k8s.io/client-go/tools/cache 库，你可以更高效地与 Kubernetes API 服务器进行交互，实现对 Kubernetes 资源的实时监控和处理。

// case6 github.com/pkg/errors是 Go 语言中一个用于处理错误的库，它提供了丰富的功能来增强错误处理的能力，比如错误的包装、格式化以及堆栈跟踪等
func case6() {
	// 定义一个原始错误
	err := errors.New("原始错误信息")
	// 包装原始错误，并添加额外的上下文信息
	newErr := errors.Wrap(err, "这是包装后的错误信息")
	fmt.Println(newErr)

	// 创建格式化的错误信息
	err = errors.Errorf("发生错误：%s", "这是一个格式化的错误")
	fmt.Println(err)

	err = someFunction()
	// 使用errors.WithStack获取堆栈跟踪信息
	stackErr := errors.WithStack(err)
	fmt.Println(stackErr)

	newErr = errors.Wrap(newErr, "这是包装后的错误信息2")
	fmt.Println(newErr)
	fmt.Println(GetErrorCause(newErr))
}
func someFunction() error {
	// 模拟一个错误
	return errors.New("函数内部发生错误")
}

// case7 一些实用方法

// Quote 引用
func Quote(s string) string {
	return `"` + s + `"`
}

// SecToDuration 将秒数转换为 time.Duration 类型
func SecToDuration(sec int64) time.Duration {
	return time.Duration(sec) * time.Second
}

// InitRandSeed 初始化随机数种子
func InitRandSeed() {
	rand.NewSource(time.Now().UTC().UnixNano())
}

// RandInt64 生成一个随机整数
func RandInt64(min int64, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Compress 使用 gzip 压缩字符串
func Compress(rawStr string) ([]byte, error) {
	compressedBuffer := &bytes.Buffer{}
	compressor := gzip.NewWriter(compressedBuffer)
	if _, err := compressor.Write([]byte(rawStr)); err != nil {
		return nil, fmt.Errorf(
			"Failed to compress %#v when writing: %v",
			rawStr, err)
	} else {
		if err := compressor.Close(); err != nil {
			return nil, fmt.Errorf(
				"Failed to compress %#v when closing: %v",
				rawStr, err)
		} else {
			return compressedBuffer.Bytes(), nil
		}
	}
}

// Decompress 使用 gzip 解压缩字符串
func Decompress(compressedBytes []byte) (string, error) {
	compressedReader := bytes.NewReader(compressedBytes)
	if decompressor, err := gzip.NewReader(compressedReader); err != nil {
		return "", fmt.Errorf(
			"Failed to decompress %#v when initializing: %v",
			compressedBytes, err)
	} else {
		if rawBytes, err := ioutil.ReadAll(decompressor); err != nil {
			return "", fmt.Errorf(
				"Failed to decompress %#v when reading: %v",
				compressedBytes, err)
		} else {
			return string(rawBytes), nil
		}
	}
}

// GetErrorCause 通过 errors.Cause 方法，从最终的包装错误 finalErr 中获取到了最初的原始错误 originalErr，并打印出了原始错误信息。
// 这样在错误处理过程中，无论错误被包装了多少次，都可以方便地找到错误的根源。
func GetErrorCause(err error) error {
	causeErr := errors.Cause(err)
	if causeErr == nil {
		return err
	} else {
		return causeErr
	}
}

// Retry 运行函数，直到成功或达到尝试限制。
func Retry(attempts int, f func() error) error {
	var i = 0
	for ; i < attempts; i++ {
		if attempts > 1 {
			fmt.Printf("---- Running attempt %v of %v...\n", i+1, attempts)
		}
		err := f()
		if err != nil {
			if i+1 < attempts {
				fmt.Printf("---- Attempt failed with error: %v\n", err)
				continue
			}
			fmt.Printf("---- Final attempt failed.\n")
			return err
		}
		break
	}
	fmt.Printf("---- Successful on attempt %v of %v.\n", i+1, attempts)
	return nil
}

// RunCmdMultiWriter runs a command and outputs the stdout to multiple [io.Writer].
// The writers are closed after the command completes.
func RunCmdMultiWriter(cmdline []string, stdout ...io.Writer) (err error) {
	c := exec.Command(cmdline[0], cmdline[1:]...)
	c.Stdout = io.MultiWriter(stdout...)
	c.Stderr = os.Stderr

	fmt.Printf("---- Running command: %v\n", c.Args)
	return c.Run()
}

// copyFile copies src to dst, creating dst's directory if necessary. Handles errors robustly,
// see https://github.com/golang/go/blob/c3458e35f4/src/cmd/internal/archive/archive_test.go#L57
// Doesn't copy file permissions.
func copyFile(dst, src string) (err error) {
	err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return err
	}
	var s, d *os.File
	s, err = os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err = os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if e := d.Close(); err == nil {
			err = e
		}
	}()
	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}
	return nil
}

// case8 golang.org/x/sys/windows/registry 是 Go 语言中用于操作 Windows 注册表的库。它提供了一组函数和类型，
// 允许 Go 程序与 Windows 注册表进行交互，例如读取、写入和删除注册表项和值
//
// 常见使用场景
// 应用程序配置：应用程序可以将配置信息存储在注册表中，例如应用程序的设置、用户偏好等。在程序启动时，可以从注册表中读取这些配置信息，
// 以便根据用户的设置来初始化应用程序的行为。
//
// 系统服务管理：对于需要与系统服务进行交互的 Go 程序，可以使用该库来读取和修改服务相关的注册表项，例如服务的启动参数、依赖关系等。
// 还可以通过注册表来注册或注销系统服务。
//
// 软件安装与卸载：在安装软件时，通常需要在注册表中记录软件的安装路径、版本信息等。卸载软件时，需要从注册表中删除相应的记录。通过操作注册表，
// 可以实现软件安装和卸载过程中的相关配置管理。
//
// 系统设置调整：可以使用该库来修改一些系统级别的设置，如网络设置、显示设置等。但需要注意的是，对系统设置的修改可能会影响系统的稳定性和其他应
// 用程序的正常运行，因此需要谨慎操作。
func case8() {
	// 打开注册表键
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		// 处理错误
	}
	defer key.Close()

	// 读取注册表值
	value, _, err := key.GetStringValue("ProgramFileDir")
	if err != nil {
		// 处理错误
	}
	// value就是读取到的值
	fmt.Println(value)

	// 写入注册表值
	key, err = registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		// 处理错误
	}
	defer key.Close()

	err = key.SetStringValue("NewApp", "NewAppValue")
	if err != nil {
		// 处理错误
	}

	// 创建注册表键
	newKey, has, err := registry.CreateKey(registry.CURRENT_USER, `Software\MyCompany\MyApp`, registry.ALL_ACCESS)
	if err != nil {
		// 处理错误
	}
	defer newKey.Close()
	fmt.Println(has)

	// 删除注册表值和键
	delKey, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		// 处理错误
	}
	defer delKey.Close()

	err = delKey.DeleteValue("SomeApp")
	if err != nil {
		// 处理错误
	}

	err = registry.DeleteKey(registry.CURRENT_USER, `Software\MyCompany\MyApp`)
	if err != nil {
		// 处理错误
	}
}

func copyFile2(dst, src string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	return cmp.Or(copyToFile(dst, f), f.Close())
}

func copyToFile(path string, r io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o777); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return cmp.Or(err, f.Close()) // cmp.Or 或者返回不等于零值的第一个参数。如果没有参数为非零值，则返回零值。
	// 如果 err == nil 则执行 f.Close()
}

func case9() {

	WriteSHA256ChecksumFile("./go.mod")
}

func WriteSHA256ChecksumFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	checksum := sha256.New()
	if _, err = io.Copy(checksum, file); err != nil {
		return err
	}
	// Write the checksum in a format that "sha256sum -c" can work with. Use the base path of the
	// tarball (not full path, not relative path) because then "sha256sum -c" automatically works
	// when the file and the checksum file are downloaded to the same directory.
	content := fmt.Sprintf("%v  %v\n", hex.EncodeToString(checksum.Sum(nil)), filepath.Base(path))
	outputPath := path + ".sha256"
	if err := os.WriteFile(outputPath, []byte(content), 0o666); err != nil {
		return err
	}
	fmt.Printf("Wrote checksum file %q with content: %v", outputPath, content)
	return nil
}
