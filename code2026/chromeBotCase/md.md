===================  第一版本 草案

实现解析器

规则

一行语句为一条执行指令，第一个单词必须是<执行指令关键词>，例如 input， 随后空格后输入指令对应参数，例如 定位的Xpath,
参数可以使用内置方法函数(变量),方法名是大小开头， 例如  click Search("百度一下")

<执行指令关键词> <执行变量>  例如 input //*[@id="chat-textarea"]
内置方法函数(变量),方法名是大小开头,    例如  click Search("百度一下")
也可以将主变量（一个）放在最前面链式接多个方法，如果方法就一个输入并且是主变量者不需加括号，否则加括号内输入变量  "百度一下".Search.Scope("button")
脚本注释用 #
to x 将前面的执行输出到变量x，一般使用场景在采集页面html的时候

例子1 打开百度输入github点击百度一下，然后采集页面的列表

```  baidu_case.cas
chrome init  # 重新初始化chrome，电脑新打开chrome
table https://www.baidu.com  # 在第一个table页输入百度的链接
input //*[@id="chat-textarea"] "github"  # 在Xpath定位出输入 github
click Search("百度一下")  # 在页面搜索"百度一下"文本定位位置，搜索到了点击
let result = ""
collect //div[@class="result"] to result.DelHtml.Save("./result.txt")  # 提取页面定位的Xpath保存到x,x执行删除HTMl标签然后保存到./result.txt
```

执行指令关键词

- chrome  操作chrome浏览器进程,  open 打开新chrome， proxy 设置代理, close 关闭chrome ....
- table   操作table页面,  add 添加页面， close 关闭页面 ....
- input   浏览器上执行输入操作，后面第一个参数跟Xpath定位,或内置函数，第二个参数是输入内容
- click   浏览器上的点击操作，后面第一个参数跟Xpath定位,或内置函数
- collect  采集页面，后面第一个参数跟Xpath定位,或内置函数
- stop   停止，后面第一个参数是停止时间 1s  100ms 1m 1d ....
- let 定义变量  类型有字符串,数值类型,布尔类型
- if true, false, 支持运算 >,<,=,!= 条件判断
- for 循环 i=0, i<10, i++  {} 内是循环体
- break 退出循环
- continue 跳过本次循环
- to  将页面的html赋予给变量
- check 检查浏览器上的Xpath, 后面第一个参数跟Xpath定位，返回true或false
- fn 定义函数 存在return值
- print 在控制台输出
- list [] 数组类型
- exit 



例子2
``` baidu_case2.cas 打开百度输入github点击百度一下，然后点击页面的列表然后再返回再点击
chrome init 

# 定义进入github搜索页
fn inGithub() {
	table https://www.baidu.com  # 在第一个table页输入百度的链接
    input //*[@id="chat-textarea"] "github"  # 在Xpath定位出输入 github
    click Search("百度一下")  # 在页面搜索"百度一下"文本定位位置，搜索到了点击
}

# 定义列表上点击子项
fn clickItem(i) bool {
	click //div[@class="result"][%d] i # 点击每一个
	stop 5s  # 休息5s
	let a = check //title
	table CallBack # 执行退回上一个页面 
	return a
}

# 执行
inGithub()
for i=1; i<10; i++ {
   if clickItem(i) == false {
       print("第%d个没有title",i)
   }
}

```

例子3 打开豆包问问题
```
chrome init 
table "https://www.doubao.com/chat/"
list [
"介绍一下golang",
"给安装教程",
"写一段能运行的代码"
]
for i=0; i<list.length; i++ {
    input //textarea[@data-testid='chat_input_input'] list[i]
    click //*[@id='flow-end-msg-send']
    let x = false
    for i=0; i<10; i++ {
        if check "//div[contains(@class, 'send-btn-wrapper') and (contains(@class, '!hidden'))]" { # 检查是否还在回复
            break
        }
        if i==9 { # 检查超时了，这种情况一般是网络不稳定或ai出现一直循环回复
            print("超时")
            x = true
        }
    }
    stop 1s  # 休息1s
    if x == true {
       print("中间等待回复的时候超时了，结束本次循环")
       break
    }
}
chrome close # 关闭当前chrome的进程
```

例子4  最终形态，输入文本让AI能理解后转化为 例子3 
```
打开浏览器到"https://www.doubao.com/chat/"
list [
"介绍一下golang",
"给安装教程",
"写一段能运行的代码"
]
将list循环执行操作
输入 //textarea[@data-testid='chat_input_input'] list[i] 
点击 //*[@id='flow-end-msg-send']
循环等待检查10次 "//div[contains(@class, 'send-btn-wrapper') and (contains(@class, '!hidden'))]" 如何是true则跳出循环，附加如果是第9次都不为true则超时
遇到超时，跳出list循环
关闭浏览器
```


```python
var keywordMap = map[string]struct{}{
	"let":     struct{}{},
	"chrome":  struct{}{},
	"table":   struct{}{},
	"input":   struct{}{},
	"click":   struct{}{},
	"collect": struct{}{},
	"stop":    struct{}{},
	"if":      struct{}{},
	"for":     struct{}{},
	"to":      struct{}{},
	"check":   struct{}{},
}

```


test case 
```python
let sName = "百度一下"
let 
let sName
let sName = 
let sName = "
le
let chat-textarea-Xpath = "//*[@id="chat-textarea"]"
let inputTxt = "github"
chrome open
table "https://www.baidu.com"
input chat-textarea-Xpath inputTxt
click Search(sName)
collect //div[@class="result"] to result.DelHtml.Save("./result.txt")
```


===================  第二版本 草案

基于 AST 实现，内置封装了很多方法，语法要简单，不能太像编程语言，没有过多的语法，完全过程化编程，专注于流程，表达清楚为目标

1. 支持变量以及算数运算
2. 支持逻辑判断
3. 支持循环
4. 内置关键词
5. 变量值如果是字符类型必须双引号
6. 注释为 #
7. 支持变量类型 数值类型，字符串，数组，bool类型
8. 支持很多内置函数方法
9. 支持函数链式调用
10. 

例子1 
```
let url = "https://www.baidu.com" 
chrome init # 重新初始化chrome，电脑新打开chrome
let ciList = ["github", "ai", "golang"] # 定义一个数组
for i=0;i<Len(ciList);i++ { # 循环这个数组
    table 1 open url # 在第一个table页输入百度的链接
    input "//*[@id="chat-textarea"]" ciList[i]   # 在Xpath定位出输入 数据
    click SearchXpath("百度一下")   # 在页面搜索"百度一下"文本定位位置，搜索到了点击
    if CheckXpath("//div[@class="result"]") { # 判断是否存在这个Xpath
        log("找到数据并保存到 ./result.txt")
        collect "//div[@class="result"]" to x.DelHtml().Save(Str("./result_%s.txt",ciList[i])) # 提取页面定位的Xpath保存到x,x执行删除HTMl标签然后保存到./result.txt
    } else {
        log(Str("%s没有找到数据", ciList[i]))
    }
    stop 1s # 停止1s
}
chrome close # 关闭浏览器
```



