package main

import (
	"bytes"
	"container/list"
	_ "embed"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io/ioutil"
	"log"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type CompanyVideoItem struct {
	SnapshotKey string `json:"snapshot_key"`
	SnapshotUrl string `json:"snapshot_url"`
	VideoKey    string `json:"video_key"`
	VideoName   string `json:"video_name"`
	VideoUrl    string `json:"video_url"`
}

func main() {
	case4()

	case5()

	case6()

	case7()

	case8()

	case9()
}

// case1
// google.golang.org/protobuf/types/known/timestamppb 是 Go 语言中用于处理 Google Protocol Buffers（protobuf）中 Timestamp 类型的库。
//在 protobuf 中，Timestamp 是一种用于表示时间戳的标准类型。这个库提供了相关的函数和方法，用于在 Go 语言中方便地对 Timestamp 类型进行操作，
//比如将 Go 标准库中的 time.Time 类型转换为 Timestamp 类型，以及将 Timestamp 类型转换回 time.Time 类型等，使得在使用 protobuf 进行数
//据序列化和反序列化时，能够更方便地处理时间相关的数据。
// ex：
/*
func WithStartTime(startTime time.Time) NewOrchestrationOptions {
	return func(req *protos.CreateInstanceRequest) error {
		req.ScheduledStartTimestamp = timestamppb.New(startTime)
		return nil
	}
}
*/

// case2
// google.golang.org/protobuf/types/known/wrapperspb 是 Go 语言中用于处理 Google Protocol Buffers（protobuf）中包装类型的库。
//在 protobuf 中，基本数据类型（如 int、string、bool 等）有时需要表示 “不存在” 或 “未设置” 的状态，包装类型就用于解决这个问题。该库提供
//了一系列的包装类型，如 Int32Value、StringValue、BoolValue 等，每个包装类型都包含一个对应的基本数据类型字段以及相关的方法，方便在 protobuf
//消息中对这些可空的基本数据类型进行处理。例如，当一个 protobuf 消息中的某个字段可能不被设置时，可以使用相应的包装类型来表示，这样在编码和解码过
//程中能够更准确地处理数据的存在性和值的情况。
// ex:
/*
func WithInput(input any) NewOrchestrationOptions {
	return func(req *protos.CreateInstanceRequest) error {
		bytes, err := json.Marshal(input)
		if err != nil {
			return err
		}
		req.Input = wrapperspb.String(string(bytes))
		return nil
	}
}
*/

// case3
// //go:embed 是 Go 语言中的一个编译指令，用于将文件或目录嵌入到 Go 程序的二进制文件中。它是 Go 1.16 版本引入的一项特性。
// 使用 //go:embed 指令可以方便地将一些静态资源，如 HTML 文件、JavaScript 文件、配置文件、图片等，与 Go 程序打包在一起，避免在运行时依赖外部文件。
// 需要 import _ "embed"
//
//go:embed hello.txt
var data []byte

func case3() {
	content, err := ioutil.ReadAll(bytes.NewReader(data))
	if err != nil {
		fmt.Println("Error reading embedded file:", err)
		return
	}
	fmt.Println(string(content))
}

/*
其它例子，初始化数据
...
//go:embed schema.sql
var schema string
...
// Initialize database
if _, err := be.db.Exec(context.Background(), schema); err != nil {
    panic(fmt.Errorf("failed to initialize the database: %w", err))
}
...
*/

// case4 log.Flags() 函数能让你获取标准日志记录器当前的输出标志，结合 log.SetFlags() 函数，你可以灵活地控制日志的输出格式。
func case4() {
	// 获取当前的日志标志
	flags := log.Flags()
	log.Printf("当前的日志标志值为: %d\n", flags)

	// 打印不同标志组合下的日志
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("这是包含日期和时间的日志")

	log.SetFlags(log.Lshortfile)
	log.Println("这是包含简短文件名和行号的日志")

	// 恢复默认标志
	log.SetFlags(flags)
	log.Println("恢复到原来的日志标志设置")
}

// case5
// wrapperspb 包提供的这些包装类型在需要表示可空的基本数据类型时非常有用，特别是在使用 Protocol Buffers 进行数据序列化和反序列化时。
func case5() {
	// 创建并使用 BoolValue
	boolValue := wrapperspb.Bool(true)
	fmt.Printf("BoolValue: %v\n", boolValue.Value)

	// 创建并使用 StringValue
	stringValue := wrapperspb.String("Hello, World!")
	fmt.Printf("StringValue: %v\n", stringValue.Value)

	// 创建并使用 Int32Value
	int32Value := wrapperspb.Int32(42)
	fmt.Printf("Int32Value: %v\n", int32Value.Value)

	// 创建并使用 DoubleValue
	doubleValue := wrapperspb.Double(3.1415926)
	fmt.Printf("DoubleValue: %v\n", doubleValue.Value)
}

// case6
// 通过 google.golang.org/protobuf/types/known/timestamppb 库，可以方便地在 Protobuf 消息中进行时间戳的处理和转换，
// 确保时间数据在不同系统和语言之间的准确表示和传输。
func case6() {
	// 获取当前时间的Timestamp
	now := timestamppb.Now()
	fmt.Printf("当前时间的 Timestamp: %v\n", now)

	// 将time.Time转换为Timestamp
	t := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
	ts := timestamppb.New(t)
	fmt.Printf("转换后的 Timestamp: %v\n", ts)

	// 将Timestamp转换为time.Time
	t2 := ts.AsTime()
	fmt.Printf("转换后的 time.Time: %v\n", t2)

	fmt.Println(ts.String())
	fmt.Println(ts.GetSeconds())
	fmt.Println(ts.GetNanos())
}

func case7() {
	fmt.Println(GetTaskFunctionName("123"))
	fmt.Println(GetTaskFunctionName(case6))
}

func GetTaskFunctionName(f any) string {
	if name, ok := f.(string); ok {
		return name
	} else {
		// this gets the full module path (github.com/org/module/package.function)
		name = runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		log.Println(name)
		startIndex := strings.LastIndexByte(name, '.') // 取最后一个.的位置，通过+1就到了函数名
		log.Println(startIndex)
		if startIndex > 0 {
			name = name[startIndex+1:]
		}
		return name
	}
}

// case8
// json.MarshalIndent 将 p 转换为 JSON 格式的字节切片，并使用 4 个空格进行缩进。最后将字节切片转换为字符串并打印输出，
// 得到的 JSON 数据具有良好的缩进格式，便于阅读。
func case8() {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	p := Person{
		Name: "John",
		Age:  30,
	}
	jsonData1, err := json.Marshal(p)
	fmt.Println(string(jsonData1))
	jsonData, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))
}

// case9
// container/list 库提供了一种灵活的数据结构，适用于需要频繁插入和删除元素的场景，例如实现队列、栈等数据结构。
func case9() {
	l := list.New()
	fmt.Printf("新创建的链表: %v\n", l)

	// 在尾部插入元素
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)

	// 在头部插入元素
	l.PushFront(0)

	// 在 e1 之前插入元素
	l.InsertBefore(1.5, e1)

	// 在 e2 之后插入元素
	l.InsertAfter(2.5, e2)

	fmt.Printf("插入元素后的链表: ")
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Printf("%v ", e.Value)
	}
	fmt.Println()

	// 删除元素 e1
	removed := l.Remove(e1)
	fmt.Printf("删除的元素: %v\n", removed)

	fmt.Printf("删除元素后的链表: ")
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Printf("%v ", e.Value)
	}
	fmt.Println()

	fmt.Printf("遍历链表: ")
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Printf("%v ", e.Value)
	}
	fmt.Println()
}
