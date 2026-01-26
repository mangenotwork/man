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
- to  将页面的html赋予给变量
- check 检查浏览器上的Xpath, 后面第一个参数跟Xpath定位，返回true或false


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


