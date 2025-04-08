package main

func main() {

}

// case1  go:generate 是什么
/*
在 Go 语言中，//go:generate 是什么？

//go:generate 是 Go 语言中的一种代码生成指令，允许开发者在代码中嵌入命令，用于在构建前自动生成或修改代码或配置文件。通过 go generate 命令触发，属于 Go 工具链的一部分，常用于简化重复性代码的管理。

核心概念
指令格式
在 Go 文件中以注释形式声明：

复制
//go:generate <command> <arguments>
示例：

复制
//go:generate stringer -type=Pill
触发方式
在终端运行命令：

复制
go generate ./...
工具会扫描当前目录及子目录下的 .go 文件，执行所有 //go:generate 指令。

主要用途
自动生成代码

为枚举类型生成 String() 方法（使用 stringer 工具）。

生成协议缓冲区代码（如 protoc-gen-go）。

生成 Mock 对象（如 mockgen）。

资源嵌入
将静态文件（如 HTML、CSS）转换为 Go 代码（例如通过 go-bindata）。

配置管理
自动生成版本号、编译时间等元信息。

典型示例
1. 为枚举生成 String() 方法

go
复制
//go:generate stringer -type=Color
type Color int

const (
    Red Color = iota
    Blue
    Green
)
运行 go generate 后生成 color_string.go，包含 Color.String() 方法。

2. 生成 gRPC 代码

go
复制
//go:generate protoc --go_out=. --go-grpc_out=. ./proto/hello.proto
运行 go generate 调用 protoc 生成 gRPC 代码。

关键注意事项
手动触发
go generate 不会在 go build 或 go run 时自动执行，需显式调用。

工具依赖
指令中的命令（如 stringer、protoc）需提前安装到系统环境变量。

跨平台兼容性
命令需兼容当前操作系统（如路径差异）。

依赖顺序
多个 //go:generate 指令按文件中的声明顺序执行。

常见工具
stringer：为枚举生成 String() 方法。

mockgen：生成接口的 Mock 实现。

protoc-gen-go：生成 Protocol Buffers 代码。

go-bindata：将文件嵌入 Go 代码。

swag：生成 Swagger 文档。

最佳实践
集中管理生成逻辑
在项目中添加 gen.go 文件统一管理指令。

版本控制生成代码
将生成的代码提交到仓库，避免依赖他人环境。

结合 go:build 标签
控制生成代码的条件（如仅针对特定平台）。

总结
//go:generate 是 Go 生态中实现代码自动化的核心工具，适用于以下场景：

减少模板代码维护（如枚举、协议、Mocks）。

嵌入外部资源（配置文件、静态文件）。

确保生成代码与源码同步。

合理使用可显著提升开发效率和代码可维护性。
*/

// CancellationReason defines the possible reasons a recognition result might be canceled.
type CancellationReason int

const (
	// Error indicates that an error occurred during speech recognition.
	Error CancellationReason = 1

	// EndOfStream indicates that the end of the audio stream was reached.
	EndOfStream CancellationReason = 2

	// CancelledByUser indicates that request was cancelled by the user.
	// Added in version 1.17.0
	CancelledByUser CancellationReason = 3
)

//go:generate stringer -type=CancellationReason -output=cancellation_reason_string.go

// case2  在 Go 语言中，unsafe.Pointer 是一种特殊的指针类型，允许直接操作内存地址并绕过 Go 的类型安全机制。
//它是 unsafe 包的核心功能之一，主要用于实现底层内存操作（如与 C 语言交互或优化性能）
/*
1. 核心特性
通用指针：
unsafe.Pointer 可以指向任意类型的值（类似于 C 中的 void*），能与其他指针类型（如 *int）或 uintptr 相互转换。

绕过类型系统：
允许直接读写内存，不受 Go 类型安全的约束（需谨慎使用，可能导致未定义行为）。

无垃圾回收保护：
指向的内存可能被 Go 的垃圾回收器回收，需确保指针引用的对象始终存活。

2. 使用场景
场景 1：类型转换
将一种类型的指针转换为另一种类型的指针（例如将 *float64 转为 *int64）：

go
复制
var f float64 = 3.1415
// 将 *float64 转为 *int64
ptr := (*int64)(unsafe.Pointer(&f))
fmt.Println(*ptr) // 输出底层内存的二进制表示对应的整数
场景 2：访问结构体未导出字段
通过计算内存偏移量访问结构体的私有字段：

go
复制
type Secret struct {
    public  int
    private int // 未导出字段
}

s := Secret{public: 42, private: 100}
// 获取结构体地址的 unsafe.Pointer
p := unsafe.Pointer(&s)
// 计算 private 字段的偏移量（需保证结构体布局不变）
privatePtr := (*int)(unsafe.Pointer(uintptr(p) + unsafe.Offsetof(s.public)))
fmt.Println(*privatePtr) // 输出 100（假设内存对齐无填充）
场景 3：与 C 语言交互
在 CGO 中将 Go 指针传递给 C 函数：

go
复制

#include <stdio.h>
void print(void *p) { printf("%d\n", *(int*)p); }

import "C"
import "unsafe"

func main() {
	x := 42
	C.print(unsafe.Pointer(&x)) // 将 Go 指针转为 C 的 void*
}
场景 4：高性能内存操作
避免复制数据，直接操作内存（如解析二进制流）：

go
复制
data := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
// 将字节切片转为字符串（零拷贝）
str := *(*string)(unsafe.Pointer(&data))
fmt.Println(str) // 输出 "Hello"


3. 关键函数与规则
unsafe.Pointer 转换规则：
Go 定义了以下合法转换路径：

复制
*T <-> unsafe.Pointer <-> uintptr <-> C 指针
直接转换 *T 到 uintptr 是非法的。

相关函数：

unsafe.Sizeof(x)：获取变量 x 的内存大小。

unsafe.Offsetof(s.f)：获取结构体字段 f 的偏移量。

unsafe.Alignof(x)：获取变量 x 的内存对齐要求。


4. 风险与注意事项
内存安全：
错误的指针操作可能导致程序崩溃（如访问已释放内存）。

不可移植性：
内存布局依赖编译器和平台（如结构体字段偏移量可能变化）。

绕过类型系统：
可能破坏 Go 的内存安全保证（如修改只读内存）。



5. 最佳实践
仅在必要时使用：优先使用标准库的安全方法（如 binary.Read）。

限制作用范围：将 unsafe 代码封装在独立模块中，并添加详细注释。

避免悬垂指针：确保 unsafe.Pointer 引用的对象不会被垃圾回收。

测试充分：针对不同平台和编译器版本验证内存操作。

*/

//case3
/*
在 Go 语言中，C.free 是通过 CGO（Go 与 C 语言交互机制）调用的 C 标准库函数 free()，用于释放由 C 代码分配的内存（如通过 C.malloc
或 C 函数返回的指针）。由于 Go 的垃圾回收机制不会管理 C 层分配的内存，必须显式调用 C.free 来手动释放，否则会导致内存泄漏。


#include <stdlib.h>

import "C"
import "unsafe"

func main() {
	// 创建 C 字符串（Go 字符串 -> C 的 char*）
	cStr := C.CString("Hello, CGO!")
	defer C.free(unsafe.Pointer(cStr)) // 必须手动释放

	// 调用 C 函数（假设某个 C 函数返回动态内存）
	cPtr := C.some_c_function()         // 返回 *C.char
	defer C.free(unsafe.Pointer(cPtr))   // 释放内存
}


最佳实践
立即绑定 defer
分配内存后立即用 defer 确保释放：

go
复制
cStr := C.CString("Safe")
defer C.free(unsafe.Pointer(cStr)) // 确保函数退出前释放
避免长期持有 C 指针
尽快将 C 数据转换为 Go 类型（如 C.GoString）并释放：

go
复制
cResult := C.get_data()
goData := C.GoString(cResult) // 复制数据到 Go 内存
C.free(unsafe.Pointer(cResult))
封装内存管理
将 C 指针和释放操作封装到结构体中：

go
复制
type SafeCString struct {
    ptr *C.char
}

func NewSafeCString(s string) *SafeCString {
    return &SafeCString{ptr: C.CString(s)}
}

func (s *SafeCString) Free() {
    C.free(unsafe.Pointer(s.ptr))
    s.ptr = nil // 防止重复释放
}

// 使用
s := NewSafeCString("Hello")
defer s.Free()
使用 runtime.SetFinalizer
为 C 指针关联析构函数（需谨慎，不推荐依赖）：

go
复制
cStr := C.CString("Auto-free")
ptr := unsafe.Pointer(cStr)
runtime.SetFinalizer(&ptr, func(_ *unsafe.Pointer) {
    C.free(ptr)
})

*/

//case4
/*

在 Go 语言中，uintptr 是一个无符号整数类型，其大小足以存储任意指针的地址值（即内存地址的整数值）。它通常与 unsafe 包结合使用，
用于底层内存操作（如指针运算或与 C 语言交互），但需谨慎使用，因为它会绕过 Go 的类型安全和垃圾回收机制。

核心特性
本质是整数
uintptr 是一个数值类型，不持有指针的语义，仅表示内存地址的整数值。

go
复制
var ptr uintptr = 0x7ffe12345678 // 直接表示内存地址
与指针的转换
可以通过 unsafe.Pointer 与其他指针类型（如 *int）相互转换：

go
复制
x := 42
ptr := uintptr(unsafe.Pointer(&x)) // 将 *int 转为 uintptr
不阻止垃圾回收
uintptr 仅保存地址的数值，不持有对象的引用。如果对象被垃圾回收，uintptr 可能指向无效内存（悬垂指针）。


主要用途
1. 指针运算
计算内存偏移量（例如访问结构体未导出字段）：

go
复制
type Data struct {
    a int
    b string // 未导出字段
}

d := Data{a: 10, b: "secret"}
// 获取结构体起始地址
basePtr := unsafe.Pointer(&d)
// 计算字段 b 的偏移量
offset := unsafe.Offsetof(d.b)
// 通过 uintptr 进行地址运算
bPtr := (*string)(unsafe.Pointer(uintptr(basePtr) + offset))
fmt.Println(*bPtr) // 输出 "secret"
2. 与 C 语言交互
将 Go 指针传递给 C 函数时，可能需要将指针转换为 uintptr：

go
复制

#include <stdio.h>
void print_addr(void *p) { printf("%p\n", p); }

import "C"
import "unsafe"

func main() {
	x := 42
	addr := uintptr(unsafe.Pointer(&x))
	C.print_addr(unsafe.Pointer(uintptr(addr))) // 转换为 C 的 void*
}
3. 反射和底层操作
在反射或自定义内存分配器中直接操作内存布局：

go
复制
// 获取切片底层数组的地址
slice := []int{1, 2, 3}
dataPtr := uintptr(unsafe.Pointer(&slice[0]))

*/

//case5
/*

在 Go 语言中，//export cgoAudioOutputCallWriteCallback 是一个 CGO 指令，用于将 Go 函数导出为 C 语言可调用的函数。
它通常用于实现 跨语言回调，例如在音频处理场景中，C 库需要调用 Go 编写的回调函数来处理音频数据。

核心作用
导出 Go 函数
通过 //export <函数名> 注释，告诉 CGO 将指定的 Go 函数编译为 C 可调用的形式（生成对应的 C 头文件）。

跨语言回调
允许 C 代码（如音频库）调用 Go 函数，实现事件驱动逻辑（如音频数据写入、实时处理）。

*/

// case5
/*
在编程里，“Wrapper”（直译为 “包装器”）是一种常见的编程概念和技术，它指的是一个代码结构（如类、函数、模块等），
其作用是封装或包装另一个代码组件，以提供额外的功能、简化使用方式、适配不同接口等。


当需要让一个已有的组件与另一个接口不兼容的系统协同工作时，可创建一个包装器来实现适配。
// 第三方库提供的接口
type ThirdPartyService struct{}

func (t *ThirdPartyService) DoSomethingSpecial() {
    fmt.Println("第三方服务执行特殊操作")
}

// 自定义接口
type MyInterface interface {
    DoWork()
}

// 包装器结构体，将第三方服务适配到自定义接口
type ServiceWrapper struct {
    thirdParty *ThirdPartyService
}

func (w *ServiceWrapper) DoWork() {
    w.thirdParty.DoSomethingSpecial()
}

func main() {
    thirdParty := &ThirdPartyService{}
    wrapper := &ServiceWrapper{thirdParty: thirdParty}

    var myInterface MyInterface = wrapper
    myInterface.DoWork()
}

在上述代码中，ServiceWrapper 结构体作为包装器，把 ThirdPartyService 适配到 MyInterface 接口，使得第三方服务能以自定义接口的方式被使用。

*/

// case6
/*

CLI项目聚合命令方法，设计每个对外提供方法入参一致那么就可以使用map定义map函数，学习的这个项目使用的这种方案:


type functionMap = map[string]func(string, string, string)

func printHelp(executableName string, samples functionMap) {
	fmt.Println("Input not valid")
	fmt.Println("Usage: ")
	fmt.Println(executableName, " <subscription> <region> <file> <sample>")
	fmt.Println("Where sample is of the format <scenario>:<sample>")
	fmt.Println("Available samples:")
	for id, _ := range samples {
		fmt.Println(" -- ", id)
	}
}

func main() {
	samples := functionMap{
		"speech_recognizer:RecognizeOnceFromWavFile":        recognizer.RecognizeOnceFromWavFile,
		"speech_recognizer:RecognizeOnceFromCompressedFile": recognizer.RecognizeOnceFromCompressedFile,
		"speech_recognizer:RecognizeOnceFromALAWFile":       recognizer.RecognizeOnceFromALAWFile,
		"speech_recognizer:ContinuousFromMicrophone":        recognizer.ContinuousFromMicrophone,
		"speech_recognizer:RecognizeContinuousUsingWrapper": recognizer.RecognizeContinuousUsingWrapper,
		"dialog_service_connector:ListenOnce":               dialog_service_connector.ListenOnce,
		"dialog_service_connector:KWS":                      dialog_service_connector.KWS,
		"dialog_service_connector:ListenOnceFromStream":     dialog_service_connector.ListenOnceFromStream,
		"speech_synthesizer:SynthesisToSpeaker":             synthesizer.SynthesisToSpeaker,
		"speech_synthesizer:SynthesisToAudioDataStream":     synthesizer.SynthesisToAudioDataStream,
		"speaker_recognizer:IndependentIdentification":      speaker_recognition.IndependentIdentification,
		"speaker_recognizer:IndependentVerification":        speaker_recognition.IndependentVerification,
	}
	args := os.Args[1:]
	if len(args) != 4 {
		printHelp(os.Args[0], samples)
		return
	}
	subscription := args[0]
	region := args[1]
	file := args[2]
	sample := args[3]
	sampleFunction := samples[sample]
	if sampleFunction == nil {
		printHelp(os.Args[0], samples)
		return
	}
	sampleFunction(subscription, region, file)
}

*/
