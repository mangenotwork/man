package main

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
)

func main() {
	case1()

	case3()

	case5()
}

// 使用 google/uuid 库可以来做 session id
func case1() {
	fmt.Println(uuid.New().String())
}

// case2  用到字符串的地方最好是定义常量
/*
const (
	// header keys
	headerKeyAccept              = "Accept"
	headerKeyAuthorization       = "Authorization"
	headerKeyContentType         = "Content-Type"
	HeaderKeyContinuationToken   = "X-MS-ContinuationToken"
	headerKeyFedAuthRedirect     = "X-TFS-FedAuthRedirect"
	headerKeyForceMsaPassThrough = "X-VSS-ForceMsaPassThrough"
	headerKeySession             = "X-TFS-Session"
	headerUserAgent              = "User-Agent"

	// media types
	MediaTypeTextPlain       = "text/plain"
	MediaTypeApplicationJson = "application/json"
)

*/

// case3
// bytes.TrimPrefix 是 bytes 包中的一个函数，其用途是从字节切片的开头移除指定的前缀。如果字节切片以指定的前缀起始，该函数会返回去掉前缀后的字节切片；若不以指定前缀起始，则直接返回原字节切片。
func case3() {
	// 示例字节切片
	s := []byte("Hello, World!")
	prefix := []byte("Hello, ")

	// 使用 TrimPrefix 移除前缀
	result := bytes.TrimPrefix(s, prefix)

	fmt.Printf("原字节切片: %s\n", string(s))
	fmt.Printf("移除前缀后的字节切片: %s\n", string(result))

	// 测试不以指定前缀开头的情况
	nonMatchingPrefix := []byte("Hi, ")
	result2 := bytes.TrimPrefix(s, nonMatchingPrefix)
	fmt.Printf("使用不匹配的前缀移除后的字节切片: %s\n", string(result2))
}

// case4
/*
golang 结构体  type WrappedError struct {
	ExceptionId      *string   }  与  type WrappedError struct {
	ExceptionId      string}  有什么区别


内存使用
ExceptionId *string：指针类型除了要存储字符串本身占用的内存，还需要额外的内存来存储指针变量（通常是 8 字节，取决于系统架构）。并且，由于指针指向的字符串可能存储在堆上，这会增加内存管理的复杂度。
ExceptionId string：直接存储字符串值，内存使用相对简单，仅为字符串本身占用的内存。


适用场景
ExceptionId *string：适用于需要在多个地方共享同一个字符串，或者字符串可能为空，并且希望通过 nil 来明确表示这种空状态的情况。
ExceptionId string：适用于每个结构体实例需要独立存储字符串值，且不需要用 nil 来表示特殊状态的情况。
*/

// case5
/*
结构体内 `json:"$id,omitempty"`

omitempty 的含义
omitempty 是一个可选的修饰符，其作用是当结构体字段的值为零值（对于数值类型是 0，对于字符串类型是 ""，对于指针类型是 nil 等）时，在 JSON 序列化过程中忽略该字段，不会将其包含在生成的 JSON 数据里。

*/

// case5 账号状态
type AccountStatus string

type accountStatusValuesType struct {
	None     AccountStatus
	Enabled  AccountStatus
	Disabled AccountStatus
	Deleted  AccountStatus
	Moved    AccountStatus
}

var AccountStatusValues = accountStatusValuesType{
	None: "none",
	// This hosting account is active and assigned to a customer.
	Enabled: "enabled",
	// This hosting account is disabled.
	Disabled: "disabled",
	// This account is part of deletion batch and scheduled for deletion.
	Deleted: "deleted",
	// This account is not mastered locally and has physically moved.
	Moved: "moved",
}

func case5() {
	fmt.Println(AccountStatusValues.Enabled)
}

// case6 对于http参数要用 url.Values{} 来接收和处理

// case7 Impl 后缀清晰地表明它是接口的具体实现
// 定义一个接口
type Animal interface {
	MakeSound()
}

// 实现类型，使用 Impl 后缀
type DogImpl struct{}

func (d DogImpl) MakeSound() {
	println("汪汪汪")
}

// case8  可以用 if resp != nil && (resp.StatusCode < 200 || resp.StatusCode >= 300)  来判断请求是否错误
/*
resp, err := client.client.Do(request)
	if resp != nil && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		err = client.UnwrapError(resp)
	}
	return resp, err
*/

// case8  .github workflows是什么
/*
GitHub Actions的workflows是一种自动化工具，用于在GitHub仓库中执行预定义的任务和命令。‌

基本概念和用途
GitHub Actions的配置文件称为workflow文件，这些文件存放在代码仓库的.github/workflows目录下。每个workflow文件采用YAML格式，
文件名可以任意取，但后缀名必须为.yml或.yaml。一个仓库可以配置多个workflow文件，每个文件定义了一组任务和命令，用于自动化软件开发周
期中的各种操作，如测试、构建、部署等‌。

触发条件和配置方式
Workflow的触发条件通常由on字段定义，可以指定特定的事件，如代码推送（push）、拉取请求（pull_request）等。此外，还可以设置分支和
标签的过滤条件。例如，on: push: branches: - master表示只有在master分支发生push事件时才会触发workflow‌。

Workflow的主体是jobs字段，表示要执行的一项或多项任务。每个任务由job_id标识，可以定义多个步骤（steps），每个步骤可以执行特定的命
令或调用其他Actions。例如，一个测试任务的配置可能包括检查代码、设置环境、运行测试等步骤‌。

实际应用场景
GitHub Actions的workflows在持续集成/持续部署（CI/CD）场景中非常有用。它们可以用于运行测试用例、构建项目、打包应用以及部署到不
同的环境。例如，可以在每次代码提交后自动运行测试，确保代码质量；在代码合并到主分支时自动构建和部署应用到生产环境‌。

通过合理配置workflows，开发者可以大大提高开发效率，确保代码质量，并简化部署流程。

*/
