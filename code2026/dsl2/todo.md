- [ok] 新增 chrome 关键词
- [ok] 将注释 // 改为 #
- [ok] 新增 break  continue  
- [ok] 新增 for 循环
- [ok] 新增数据类型 列表  
- [ok] 新增数据类型 字典 也就是哈希表
- [ok] 解决bug  bug1 :  print("dict["a"] = ", dict["a"])    报错，语法没过 , 错误定位是 "dict[\"a\"] = " 被当成Token在处理
- [ok] 新增函数链式调用和链式中的管道传递 print("aa").print("bb")  // aa\nbb
- [ok] 新增测试，需要能断言，跑所有的测试脚本
- 扩展链式调用，变量可以跟后面的函数
- 新增语法错误信息提示
- 新增 table 例如 table 1 open url # 在第一个table页输入百度的链接
- 新增 input 例如 input "//*[@id="chat-textarea"]" ciList[i]   # 在Xpath定位出输入 数据
- 新增 click 例如 click SearchXpath("百度一下")   # 在页面搜索"百度一下"文本定位位置，搜索到了点击


# 测试
1. ast/ast_test.go- AST节点测试
2. lexer/lexer_test.go- 词法分析器测试
3. parser/parser_test.go- 语法分析器测试
4. interpreter/interpreter_test.go- 解释器测试
5. integration_test.go- 集成测试


还需要多测试脚本语法，挖掘其他的语法特性
