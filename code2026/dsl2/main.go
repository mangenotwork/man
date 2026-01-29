// main.go
package main

import (
	"dsl2/builtins"
	"dsl2/interpreter"
	"dsl2/lexer"
	"dsl2/parser"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		// 如果没有参数，使用示例脚本

		//log.Println("所有示例测试")
		//runExample()
		//log.Println("所有示例测试完成 .......")

		runExample2()
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

*/

func runExample2() {
	// 一个简单的示例脚本
	script := `
print("aa").print("bb")
print("hello").upper().print()
print("hi").repeat(3).print()
var x = "test"
print(x).upper().repeat(2).print()
var list = [1, 2, 3]
print(list).print(len(list))
`

	/*





	 */

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
