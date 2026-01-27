// main.go
package main

import (
	"fmt"
	"io/ioutil"
	"my-dsl/builtins"
	"my-dsl/interpreter"
	"my-dsl/lexer"
	"my-dsl/parser"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		// 如果没有参数，使用示例脚本
		runExample()
		return
	}

	filename := os.Args[1]
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("无法读取文件 %s: %v\n", filename, err)
		return
	}

	runScript(string(source))
}

func runScript(source string) {
	// 词法分析
	l := lexer.New(source)

	// 语法分析
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("解析错误:")
		for _, err := range p.Errors() {
			fmt.Println("  " + err)
		}
		return
	}

	// 创建解释器
	interp := interpreter.NewInterpreter()

	// 注册内置函数
	builtins.RegisterBuiltins(interp)

	// 注册你的Go工具函数
	registerCustomFunctions(interp)

	// 执行程序
	result, err := interp.Interpret(program)
	if err != nil {
		fmt.Printf("执行错误: %v\n", err)
		return
	}

	if result != nil {
		fmt.Printf("程序返回值: %v\n", result)
	}
}

func runExample() {
	// 一个简单的示例脚本
	script := `
print("aaaaa")
// 变量声明
var count = 10
var name = "World"
var is_active = true;

// 函数调用
print("Hello, " + name + "!")
print("Count:", count)

// 条件语句
if count > 5 {
    print("Count 大于 5");
} else {
    print("Count 小于等于 5")
}

if is_active {
	print("is_active");
}

// 算数运算
var a = 2 + 1;
print("a:", a);

// 循环
var i = 0;
while i < 3 {
    print("循环次数:", i);
    i = i + 1;
	break;
}

// 调用你的Go工具函数
log_info("i = ", i);

`

	fmt.Println("执行示例脚本:")
	fmt.Println("======================================")
	runScript(script)
}

// 注册你的Go工具函数
func registerCustomFunctions(interp *interpreter.Interpreter) {

	interp.Global().SetFunc("log_info", func(args []interpreter.Value) (interpreter.Value, error) {
		fmt.Print("[INFO] ")
		for _, arg := range args {
			fmt.Print(arg, " ")
		}
		fmt.Println()
		return nil, nil
	})

	interp.Global().SetFunc("log_error", func(args []interpreter.Value) (interpreter.Value, error) {
		fmt.Print("[ERROR] ")
		for _, arg := range args {
			fmt.Print(arg, " ")
		}
		fmt.Println()
		return nil, nil
	})
}
