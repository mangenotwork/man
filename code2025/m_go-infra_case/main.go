package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	case1()

	log.Println(preHyphenated("aaa", "", "vvvv", "124312"))

	case3()

	case4()
}

// case1 查找环境变量
func case1() {
	key := "MY_VARIABLE"
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("Env var not defined: %v", key)
	}
	if v == "" {
		log.Printf("Env var defined as empty string: %v", key)
	}
	log.Printf("Env var value: %v", v)
}

// case2 预连字符 用 -
func preHyphenated(s ...string) string {
	if len(s) == 0 {
		return ""
	}
	var b strings.Builder
	first := true
	for i := 0; i < len(s); i++ {
		if s[i] == "" {
			continue
		}
		if !first {
			b.WriteRune('-')
		}
		b.WriteString(s[i])
		first = false
	}
	return b.String()
}

// case3 exec.LookPath 是 Go 语言标准库 os/exec 包中的一个函数，其主要作用是在系统的 PATH 环境变量所指定的目录中查找可执行文件
// 使用场景
// - 脚本执行前的检查
// - 跨平台兼容性
// - 动态调用外部程序
func case3() {
	path, err := exec.LookPath("python")
	if err != nil {
		fmt.Println("未找到可执行文件:", err)
	} else {
		fmt.Println("找到可执行文件，路径为:", path)
	}
}

// case4 封装命令执行
func case4() {
	// Run a command and print the output.
	if err := Run(Dir(".", "python", "-c", "print('hello world')")); err != nil {
		log.Fatal(err)
	}

	log.Println(SpaceTrimmedCombinedOutput(Dir(".", "python", "-V")))
}

// Run sets up the command to log directly to our stdout/stderr streams, then runs it.
func Run(c *exec.Cmd) error {
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return RunQuiet(c)
}

// RunQuiet logs the command line and runs the given command, but sends the output to os.DevNull.
func RunQuiet(c *exec.Cmd) error {
	fmt.Printf("---- Running command: %v %v\n", c.Path, c.Args)
	return c.Run()
}

// CombinedOutput runs a command and returns the output string of c.CombinedOutput.
func CombinedOutput(c *exec.Cmd) (string, error) {
	fmt.Printf("---- Running command: %v %v\n", c.Path, c.Args)
	out, err := c.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// SpaceTrimmedCombinedOutput runs CombinedOutput and trims leading/trailing spaces from the result.
func SpaceTrimmedCombinedOutput(c *exec.Cmd) (string, error) {
	out, err := CombinedOutput(c)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// Dir returns a command that runs in the given dir. The command can be passed to one of the other
// funcs in this package to evaluate it and optionally get the output as a string. Dir is useful to
// construct one-liner command calls, because setting the dir is commonly needed and not settable
// with exec.Command directly.
func Dir(dir, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd
}

// MakeWorkDir creates a unique path inside the given root dir to use as a workspace. The name
// starts with the local time in a sortable format to help with browsing multiple workspaces. This
// function allows a command to run multiple times in sequence without overwriting or deleting the
// old data, for diagnostic purposes. This function uses os.MkdirAll to ensure the root dir exists.
func MakeWorkDir(rootDir string) (string, error) {
	pathDate := time.Now().Format("2006-01-02_15-04-05")
	if err := os.MkdirAll(rootDir, os.ModePerm); err != nil {
		return "", err
	}
	return os.MkdirTemp(rootDir, fmt.Sprintf("%s_*", pathDate))
}

// case5 字符串操作相关
// 剪切前缀
func CutPrefix(s, prefix string) (after string, found bool) {
	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):], true
	}
	return s, false
}

// 剪切后缀
func CutSuffix(s, suffix string) (before string, found bool) {
	if strings.HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)], true
	}
	return s, false
}

// 但在sep的最后一次出现时进行剪切，而不是在第一次出现时
func CutLast(s, sep string) (before, after string, found bool) {
	if i := strings.LastIndex(s, sep); i != -1 {
		return s[:i], s[i+len(sep):], true
	}
	return "", s, false
}

// 从指定文件中读取一个JSON值。支持BOM。
func ReadJSONFile(path string, i interface{}) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to open JSON file %v for reading: %w", path, err)
	}
	defer func() {
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()
	// ReadJSONFile 函数通过 unicode.BOMOverride 处理 BOM，从而避免 BOM 对 JSON 数据解码造成影响，确保能够正确读取包含 BOM 的 JSON 文件。
	content := transform.NewReader(f, unicode.BOMOverride(transform.Nop))
	d := json.NewDecoder(content)
	if err := d.Decode(i); err != nil {
		return fmt.Errorf("unable to decode JSON file %v: %w", path, err)
	}
	return nil
}

// BOM
// BOM 是位于文本文件开头的一组特殊字符，用于标识文件所使用的字符编码以及字节序。在 Unicode 编码中，不同的编码方式
//（如 UTF - 8、UTF - 16BE、UTF - 16LE 等）可能会使用不同的 BOM 来表明自己的身份。
//
// 具体编码对应的 BOM
//UTF - 8：BOM 是 EF BB BF（十六进制表示）。
//UTF - 16BE：BOM 是 FE FF。
//UTF - 16LE：BOM 是 FF FE。
//
//代码中 BOM 的作用
//在代码里，unicode.BOMOverride(transform.Nop) 这个函数用于处理可能存在于文件开头的 BOM。具体来说：
//unicode.BOMOverride：这是 Go 语言标准库 unicode/utf8 包中的一个函数，它会检查输入流的开头是否存在 BOM。如果存在，就会自动跳过 BOM，确保后续的解码操作不会受到 BOM 的干扰；如果不存在，就会正常处理数据。
//transform.Nop：它是一个无操作的转换器，在这里作为 unicode.BOMOverride 的参数，表示不进行额外的转换操作，仅处理 BOM。

// 将一个指定值作为缩进JSON写入文件，并带有尾随换行符。
func WriteJSONFile(path string, i interface{}) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to open JSON file %v for writing: %w", path, err)
	}
	defer func() {
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()

	d := json.NewEncoder(f)
	d.SetIndent("", "  ")
	if err := d.Encode(i); err != nil {
		return fmt.Errorf("unable to encode model into JSON file %v: %w", path, err)
	}
	return nil
}
