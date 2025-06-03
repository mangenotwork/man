package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
)

func main() {

	case1()

	case4()

	case7()
}

// 检查url的可靠性
func case1() {
	var check = func(input string) (string, error) {
		value, err := url.ParseRequestURI(input)
		if err != nil || value.Host == "" { // 如果解析失败，或者没有host，就返回错误
			urlErr := &url.Error{}
			if errors.As(err, &urlErr) {
				err = urlErr.Unwrap()
			} else {
				err = errors.New("invalid URI for request")
			}

			return "", err
		}

		return value.String(), nil
	}

	result, err := check("http://www.baidu.com")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("check url: %s", result)
}

// 使用 switch .(type) 进行类型判断
/*
**类型选择（type switch）**语法，用于判断一个接口值的具体动态类型。它通常用在 interface{} 类型的变量上，用来实现类似于其他语言中的“运行时类型判断”功能

- 注意：switch .(type) 只能用在 interface{} 类型上，不能用于具体类型的变量。

- 和普通类型断言的区别
方式	是否支持多个类型判断	是否简洁	是否适用于多个类型分支
类型断言 v.(T)	❌ 否	✅ 是	❌ 不适合多分支
类型选择 switch .(type)	✅ 是	✅ 简洁	✅ 非常适合多分支判断


- 常见用途
接收任意类型参数并根据不同类型做不同处理；
解析 JSON、YAML 等结构化数据时判断字段类型；
实现插件系统、泛型逻辑等灵活场景；

*/
func case2() {
	var t interface{} = 1

	switch val := t.(type) {
	case int:
		fmt.Println("val is int. ", val)
	case string:
		fmt.Println("val is string. ", val)
	}

}

// 接口的静态检查
/*

使用空白标识符 _ 结合接口类型声明是一种常见的技巧，用于确保某个类型实现了特定的接口。这种做法通常被称为“接口实现的静态检查”，它能在编译期就
验证你的类型是否正确地实现了所需接口，而不是等到运行时才发现问题。

- 优点
编译期检查：这种方式能够在编译阶段就发现问题，而不需要等待运行时出现错误。
清晰明确：直接在代码中表明了类型的意图和约束条件，增加了代码的可读性和可维护性。
无副作用：由于使用的是 nil，不会实际调用任何方法或构造任何对象，因此没有性能损耗或者意外行为的风险

*/

type MyType struct{}

// 实现 io.Writer 接口的 Write 方法
func (m *MyType) Write(p []byte) (n int, err error) {
	fmt.Println("Writing bytes:", string(p))
	return len(p), nil
}

// 静态检查 MyType 是否实现了 io.Writer 接口
var _ io.Writer = (*MyType)(nil)

// case4 json相关

func IsJSON(content string) bool {
	return json.Valid([]byte(content))
}

func jsonTransform(content *string, prettyJSON bool) error { //revive:disable-line:flag-parameter
	var result any
	if err := json.Unmarshal([]byte(*content), &result); err != nil {
		return err
	}

	var resultRaw []byte
	var err error

	if prettyJSON {
		resultRaw, err = json.MarshalIndent(result, "", "  ")
	} else {
		resultRaw, err = json.Marshal(result)
	}

	if err != nil {
		return err
	}

	*content = strings.TrimSpace(string(resultRaw))

	return nil
}

func case4() {
	a := `{"name":"test","age":18}`
	log.Println(IsJSON(a))
	log.Println(jsonTransform(&a, true))
	log.Println(a)
	log.Println(jsonTransform(&a, false))
	log.Println(a)
}

// base64编解码和压缩解压相关

func Base64Decode(content *string) error {
	if content == nil {
		return nil
	}

	payload, err := base64ToByte(*content)
	if err != nil {
		return err
	}

	*content = strings.TrimSpace(string(payload))

	return nil
}

func Base64Encode[T string | []byte](content T) (string, error) {
	var payload string

	switch v := any(content).(type) {
	case string:
		payload = byteToBase64([]byte(v))
	case []byte:
		payload = byteToBase64(v)
	}

	return payload, nil
}

func Base64GzipEncode(content *string) error {
	if content == nil {
		return nil
	}

	var bu bytes.Buffer
	gz := gzip.NewWriter(&bu)

	if _, err := gz.Write([]byte(*content)); err != nil {
		return err
	}

	if err := gz.Flush(); err != nil {
		return err
	}

	if err := gz.Close(); err != nil {
		return err
	}

	*content = byteToBase64(bu.Bytes())

	return nil
}

func Base64GzipDecode(content *string) error {
	if content == nil {
		return nil
	}

	data, err := base64ToByte(*content)
	if err != nil {
		return nil
	}

	rd, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil
	}
	defer rd.Close()

	dataDecomp, err := io.ReadAll(rd)
	if err != nil {
		return nil
	}

	*content = string(dataDecomp)

	return nil
}

func base64ToByte(content string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func byteToBase64(content []byte) string {
	return base64.StdEncoding.EncodeToString(content)
}

// slices.Values 函数是在 Go 1.23 版本中才引
/*
// 输出为字符串的切片
func ConvertEnumsToStringSlices[T any](values []T) []string { //revive:disable-line:flag-parameter
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = fmt.Sprintf("%v", v)
	}
	return result
}

func SortMapStringByKeys[T any](m map[string]T) map[string]T {
	result := make(map[string]T, len(m))
	for _, k := range slices.Sorted(maps.Keys(m)) {
		result[k] = m[k]
	}

	return result
}

func case6() {
	log.Println(ConvertEnumsToStringSlices([]int{1, 2, 3}))
	log.Println(ConvertEnumsToStringSlices([]string{"a", "2", "c"}))
	log.Println(ConvertEnumsToStringSlices([]bool{true, false, false}))

	log.Println(SortMapStringByKeys(map[string]int{"a": 1, "b": 2}))
}
*/

// sha256 的泛型实现

func Sha256[T string | []byte](content T) string {
	var hash [32]byte

	switch v := any(content).(type) {
	case string:
		hash = sha256.Sum256([]byte(v))
	case []byte:
		hash = sha256.Sum256(v)
	}

	return hex.EncodeToString(hash[:])
}

func case7() {
	log.Println(Sha256("hello world"))
	log.Println(Sha256([]byte("hello world")))
}
