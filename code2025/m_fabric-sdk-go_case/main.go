package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

func main() {
	case1()

	case2()

	case5()

	case7()
}

// case1 new(T) 返回的是一个指向类型 T 零值的指针，而 *new(T) 会对这个指针进行解引用，从而得到类型 T 的零值实例。
// *new(T) 等同于 T{}（对于结构体类型）或者 [n]T{}（对于数组类型）。不过，在某些情况下，*new(T) 可能会让代码看起来更清晰，
// 特别是当类型名称比较长或者复杂的时候。
func case1() {
	type Person struct {
		Name string
		Age  int
	}
	p := *new(Person)
	fmt.Println(p)

	// 创建一个 *int 类型的零值实例
	ptr := *new(*int)
	fmt.Println(ptr)

}

/*
将http.Request返回值序列到 admin.AssignDomainWorkspacesByCapacitiesRequest
如例子的写法就很清晰明了

body, err := UnmarshalRequestAsJSON[admin.AssignDomainWorkspacesByCapacitiesRequest](req)
		if err != nil {
			return nil, err
		}

*/

func UnmarshalRequestAsJSON[T any](req *http.Request) (T, error) {
	tt := *new(T)
	if req.Body == nil {
		return tt, nil
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return tt, err
	}
	req.Body.Close()
	if err = json.Unmarshal(body, &tt); err != nil {
		return tt, err
	}
	return tt, nil
}

// case2 判断切片值的泛型优雅写法
func contains[T comparable](s []T, v T) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}
func case2() {
	if !contains([]int{1, 2, 3}, 4) {
		fmt.Println("false")
	}
}

// case3 func NopCloser(r io.Reader) io.ReadCloser
// 它接收一个 io.Reader 类型的参数 r，然后返回一个实现了 io.ReadCloser 接口的对象。这个返回的对象会将 Read 方法的调用转发给传入的 r，
// 而 Close 方法不会执行任何操作，只是简单地返回 nil
func SetResponseBody(resp *http.Response, body []byte, contentType string) *http.Response {
	if l := int64(len(body)); l > 0 {
		resp.Header.Set("Content-Type", contentType)
		resp.ContentLength = l
		resp.Body = io.NopCloser(bytes.NewReader(body))
	}
	return resp
}

// case4 返回指针的函数
// Ptr returns a pointer to the provided value.
func Ptr[T any](v T) *T {
	return &v
}

// SliceOfPtrs returns a slice of *T from the specified values.
func SliceOfPtrs[T any](vv ...T) []*T {
	slc := make([]*T, len(vv))
	for i := range vv {
		slc[i] = Ptr(vv[i])
	}
	return slc
}

// case5 reflect.ValueOf(body).IsZero() 借助反射机制来判定 body 这个变量是否为其类型的零值
// 使用场景
// 数据验证：在处理用户输入或者从外部数据源获取的数据时，你可能需要检查某个字段是否为空。借助反射，你可以编写通用的验证函数来处理不同类型的结构体。
// 默认值处理：当某个变量为零值时，你可以为其设置默认值。
// 注意事项
// 性能开销：反射操作通常会带来一定的性能开销，因为它需要在运行时进行类型检查和值的访问。所以，在性能敏感的场景中要谨慎使用。
// 类型安全：反射操作可能会绕过 Go 语言的类型系统，因此在使用时要确保类型的正确性，避免出现运行时错误。
func case5() {
	// 定义不同类型的变量
	var num int
	str := ""
	var flag bool
	type Person struct {
		Name string
		Age  int
	}
	var p Person

	// 使用 reflect.ValueOf().IsZero() 判断是否为零值
	fmt.Printf("num is zero: %v\n", reflect.ValueOf(num).IsZero())
	fmt.Printf("str is zero: %v\n", reflect.ValueOf(str).IsZero())
	fmt.Printf("flag is zero: %v\n", reflect.ValueOf(flag).IsZero())
	fmt.Printf("p is zero: %v\n", reflect.ValueOf(p).IsZero())
}

// case6 http.Request <-req.Context().Done() 是一个用于处理 HTTP 请求上下文取消信号的操作
// 使用场景
// 1. 处理客户端提前断开连接

// case7 url.IsAbs() 是 Go 语言标准库 net/url 包中的一个函数，用于判断给定的 URL 字符串是否为绝对 URL
// 绝对 URL 是指包含完整的协议（如 http://、https://）、主机名以及可选的端口号、路径等信息的 URL。例如，
// https://www.example.com/path/to/page 就是一个绝对 URL，而相对 URL 通常是相对于当前页面或服务器的路
// 径，如 ../path/to/page 或 page.html
func case7() {
	absURL := "https://www.example.com"
	relURL := "page.html"
	u1, _ := url.Parse(absURL)
	fmt.Println(u1.IsAbs()) // 输出: true
	u2, _ := url.Parse(relURL)
	fmt.Println(u2.IsAbs()) // 输出: false
}
