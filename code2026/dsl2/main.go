// main.go
package main

import (
	"bufio"
	"dsl2/builtins"
	"dsl2/interpreter"
	"dsl2/lexer"
	"dsl2/logger"
	"dsl2/parser"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	logger.IsDebug = true

	//fmt.Printf("欢迎使用 漫语言 v%s\n", VERSION)
	//fmt.Printf("欢迎使用 漫语言 v%s\n", VERSION)
	//fmt.Printf("欢迎使用 漫语言 v%s\n", VERSION)
	//fmt.Printf("欢迎使用 漫语言 v%s\n", VERSION)
	//fmt.Println("输入代码并按回车执行，按Ctrl+Z(Windows)退出")
	//fmt.Println("使用 'exit' 或 'quit' 命令退出程序")
	//fmt.Println("===================================================================")
	//
	//runREPL()

	if len(os.Args) < 2 {
		// 如果没有参数，使用示例脚本

		//log.Println("所有示例测试")
		//runExample()
		//log.Println("所有示例测试完成 .......")

		runExample2()
		return
	}
	//
	//filename := os.Args[1]
	//source, err := ioutil.ReadFile(filename)
	//if err != nil {
	//	fmt.Printf("无法读取文件 %s: %v\n", filename, err)
	//	return
	//}
	//
	//runScript(string(source))
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

	log.Println("执行程序")

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
#case1 : chrome语句
chrome init prot=123 proty="127.0.0.1"

var c = 2

#case 2 : 注释
// chrome init prot=123
# 1231 chrome init prot=123 proty="127.0.0.1"

#case 3 : for 循环

for var i = 0; i < 5; i = i + 1 {
    print("i =", i);
}

for var i = 0; i < 5; i = i + 1 {
	if i == 2 {
		print("i == 2  continue ");
		continue
	}
 	print("i =", i);
}

for var i = 0; i < 5; i = i + 1 {
	if i == 2 {
		print("i == 2  break ");
		break
	}
 	print("i =", i);
}

var i = 0;
for i < 5 {
    print("i =", i);
    i = i + 1;
}

var i = 0;
for {
    if i >= 5 {
        break;
    }
    print("i =", i);
    i = i + 1;
}

for var i = 0; i < 3; i = i + 1 {
    for var j = 0; j < 3; j = j + 1 {
        if i == 1 && j == 1 {
            print("跳过 i=1,j=1");
            continue;
        }
        print("i =", i, "j =", j);
    }
}

#case 4 : 变量
var a1 = 1
print("a1 = ", a1)
var a2 = "aaa"
print("a2 = ", a2)
var a3 = true
print("a3 = ", a3)
`

	fmt.Println("执行示例脚本:")
	fmt.Println("======================================")
	runScript(script)
}

/*

case1 : chrome语句
chrome init prot=123 proty="127.0.0.1"

case 2 : 注释
// chrome init prot=123
# 1231 chrome init prot=123 proty="127.0.0.1"

case 3 : for 循环

for var i = 0; i < 5; i = i + 1 {
    print("i =", i);
}

for var i = 0; i < 5; i = i + 1 {
	if i == 2 {
		print("i == 2  continue ");
		continue
	}
 	print("i =", i);
}

for var i = 0; i < 5; i = i + 1 {
	if i == 2 {
		print("i == 2  break ");
		break
	}
 	print("i =", i);
}

var i = 0;
for i < 5 {
    print("i =", i);
    i = i + 1;
}

var i = 0;
for {
    if i >= 5 {
        break;
    }
    print("i =", i);
    i = i + 1;
}

for var i = 0; i < 3; i = i + 1 {
    for var j = 0; j < 3; j = j + 1 {
        if i == 1 && j == 1 {
            print("跳过 i=1,j=1");
            continue;
        }
        print("i =", i, "j =", j);
    }
}

case 4 : 变量
var a1 = 1
print("a1 = ", a1)
var a2 = "aaa"
print("a2 = ", a2)
var a3 = true
print("a3 = ", a3)
var a4 = [1,2,3]
print("a4 = ", a4)
var a5 = ["aa", "bb", "cc"]
print("a5 ", a5)
for var i = 0; i < 5; i = i + 1 {
	print("item ", a5[i])
}

case 5 列表 list
var list1 = [1, 2, 3]
var list2 = ["a", "b", "c"]
print("list1 = ", list1)
print("list2 = ", list2)
var first = list1[0]  // 1
var second = list2[1] // "b"
print("first = ", first)
print("second = ", second)

var length = len(list1)
print("length = ", length)

var i = 0
while i < len(list1) {
    print("list i = ", i, " 值: ", list1[i])
    i = i + 1
}

var combined = list1 + list2  // [1, 2, 3, "a", "b", "c"]
print("combined = ", combined)

if [1, 2] == [1, 2] {
    print("列表相等")
}

var matrix = [[1, 2], [3, 4]]
print("matrix = ", matrix)
var element = matrix[0][1]
print("element = ", element)


case 6 字典:

var dict = {"a": 1, "b": 2, "c": 3}
print("dict = ", dict)
print("dict a = ", dict["a"])

dict["d"] = 4
dict["a"] = 10
print("下标赋值 : ", dict["a"])  // 输出: 10

var length = len(dict)  // 4
print("length : ", length)


var keys = keys(dict)
var i = 0
while i < len(keys) {
    var key = keys[i]
    print(key, ":", dict[key])
    i = i + 1
}

var values = values(dict)
var j = 0
while j < len(values) {
    print("循环遍历值 = ", values[j])
    j = j + 1
}


var pairs = items(dict)
var k = 0
while k < len(pairs) {
    var pair = pairs[k]
    print("键:", pair[0], "值:", pair[1])
    k = k + 1
}

if has_key(dict, "a") {
    print("字典包含键 'a'")
}


delete(dict, "b")
print(len(dict))  // 输出: 3
print(dict["b"])


var dict1 = {"x": 1, "y": 2}
var dict2 = {"y": 3, "z": 4}  // y 会被覆盖
var merged = dict1 + dict2
print("字典合并 = ", merged["y"])  // 输出: 3


var d1 = {"a": 1, "b": 2}
var d2 = {"b": 2, "a": 1}
if d1 == d2 {
    print("字典相等（键顺序无关）")
}


var complex = {
    "users": ["Alice", "Bob", "Charlie"],
    "scores": {"Alice": 95, "Bob": 87},
    "active": true
}
print("嵌套结构 = ", complex["users"][0])    // 输出: Alice
print("嵌套结构 = ", complex["scores"]["Alice"])  // 输出: 95

var s = "test";
var s1 = upper(s).repeat(2)
print(s1)


bug:
var s = "test"
s.print()

*/

func runExample2() {
	// 一个简单的示例脚本
	script := `
chrome init prot=123 proty="127.0.0.1"
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

//  ==============================  像 python 那样的命令行

// 常量定义
const (
	PROMPT      = ">>> "
	PROMPT_CONT = "... "
	VERSION     = "0.1.0"
	EXIT_MSG    = "再见！"
)

// REPL主循环
func runREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	interp := interpreter.NewInterpreter()

	// 设置输入提示
	fmt.Print(PROMPT)

	var inputLines []string
	braceCount := 0
	parenCount := 0
	bracketCount := 0

	for scanner.Scan() {
		line := scanner.Text()

		// 检查退出命令
		if shouldExit(line) {
			fmt.Println(EXIT_MSG)
			return
		}

		// 统计括号数量以确定是否继续输入
		braceCount += countChars(line, '{', '}')
		parenCount += countChars(line, '(', ')')
		bracketCount += countChars(line, '[', ']')

		inputLines = append(inputLines, line)

		// 如果所有括号都匹配，执行代码
		if braceCount == 0 && parenCount == 0 && bracketCount == 0 {
			// 拼接所有行
			fullInput := strings.Join(inputLines, "\n")

			// 执行代码
			executeCode(fullInput, interp)

			// 重置状态
			inputLines = nil
			braceCount = 0
			parenCount = 0
			bracketCount = 0

			fmt.Print(PROMPT)
		} else {
			// 需要继续输入
			fmt.Print(PROMPT_CONT)
		}
	}

	// 处理可能的剩余输入
	if len(inputLines) > 0 {
		fullInput := strings.Join(inputLines, "\n")
		executeCode(fullInput, interp)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取输入错误: %v\n", err)
	}
}

// 检查是否应该退出
func shouldExit(line string) bool {
	trimmed := strings.TrimSpace(strings.ToLower(line))
	return trimmed == "exit" || trimmed == "quit" || trimmed == ":q"
}

// 统计括号数量
func countChars(line string, openChar, closeChar byte) int {
	count := 0
	for i := 0; i < len(line); i++ {
		if line[i] == openChar {
			count++
		} else if line[i] == closeChar {
			count--
		}
	}
	return count
}

// 执行代码并输出结果
func executeCode(input string, interp *interpreter.Interpreter) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	// 处理特殊命令
	if handleSpecialCommands(input) {
		return
	}

	// 记录执行前的状态
	// 这里可以记录一些状态，如果需要的话

	// 解析和执行
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		printErrors("解析错误", p.Errors())
		return
	}

	result, err := interp.Interpret(program)
	if err != nil {
		fmt.Printf("执行错误: %v\n", err)
		return
	}

	// 关键：只输出包含return或print的结果
	// 或者只输出表达式的值
	shouldOutput := false

	// 检查是否应该输出
	if result != nil {
		trimmedInput := strings.TrimSpace(input)

		// 如果输入以return开头，应该输出
		if strings.HasPrefix(trimmedInput, "return ") {
			shouldOutput = true
		} else if isExpression(trimmedInput) { // 如果输入是表达式（不以语句关键字开头），应该输出
			shouldOutput = true
		} else if strings.Contains(trimmedInput, "print(") { // 如果输入包含print，已经在print函数中输出了
			shouldOutput = false
		}
	}

	if shouldOutput {
		printResult(result)
	}
}

// 判断是否是表达式
func isExpression(input string) bool {
	trimmed := strings.TrimSpace(input)

	// 空行不是表达式
	if trimmed == "" {
		return false
	}

	// 检查是否是语句
	statements := []string{
		"var ", "if ", "while ", "for ",
		"func ", "break ", "continue ",
		"print(",
	}

	for _, stmt := range statements {
		if strings.HasPrefix(trimmed, stmt) {
			return false
		}
	}

	// 检查是否包含赋值（但不是比较）
	if strings.Contains(trimmed, "=") && !strings.Contains(trimmed, "==") {
		// 包含单个等号，可能是赋值语句
		return false
	}

	// 其他情况可能是表达式
	return true
}

// 处理特殊命令
func handleSpecialCommands(input string) bool {
	trimmed := strings.TrimSpace(input)

	switch strings.ToLower(trimmed) {
	case "help", "?":
		printHelp()
		return true
	case "version", "ver":
		fmt.Printf("DSL2 解释器 v%s\n", VERSION)
		return true
	case "clear", "cls":
		clearScreen()
		return true
	case "env", "variables":
		// 这里可以添加查看环境变量的功能
		fmt.Println("环境变量功能待实现")
		return true
	}

	return false
}

// 打印帮助信息
func printHelp() {
	helpText := `
可用命令:
  help, ?      - 显示此帮助信息
  version, ver - 显示版本信息
  clear, cls   - 清屏
  exit, quit   - 退出程序
  env          - 查看环境变量

语法示例:
  >>> var x = 10
  >>> x + 5
  15
  >>> if (x > 5) { print("x大于5") }
  x大于5
  >>> for var i = 0; i < 3; i = i + 1 { print(i) }
  0
  1
  2
  >>> {"a": 1, "b": 2}
  {a: 1, b: 2}
`
	fmt.Println(helpText)
}

// 清屏
func clearScreen() {
	// 简单的清屏：打印多个空行
	for i := 0; i < 50; i++ {
		fmt.Println()
	}
}

// 打印错误信息
func printErrors(prefix string, errors []string) {
	fmt.Printf("%s:\n", prefix)
	for _, err := range errors {
		fmt.Printf("  %s\n", err)
	}
}

// 格式化并打印结果
func printResult(result interpreter.Value) {
	switch v := result.(type) {
	case int64:
		fmt.Println(v)
	case float64:
		fmt.Println(v)
	case string:
		fmt.Printf("%q\n", v)
	case bool:
		fmt.Println(v)
	case []interpreter.Value:
		// 列表
		fmt.Print("[")
		for i, item := range v {
			if i > 0 {
				fmt.Print(", ")
			}
			printValue(item)
		}
		fmt.Println("]")
	case interpreter.DictType:
		// 字典
		fmt.Print("{")
		first := true
		for key, value := range v {
			if !first {
				fmt.Print(", ")
			}
			first = false
			fmt.Printf("%v: ", key)
			printValue(value)
		}
		fmt.Println("}")
	case nil:
		// 不输出nil
	default:
		fmt.Printf("%v\n", v)
	}
}

// 打印单个值
func printValue(value interpreter.Value) {
	switch v := value.(type) {
	case string:
		fmt.Printf("%q", v)
	default:
		fmt.Printf("%v", v)
	}
}
