// example.dsl
// 你的DSL脚本示例

// 变量声明
var count = 10;
var name = "World";
var is_active = true;

// 函数调用
print("Hello, " + name + "!");
print("Count:", count);

// 条件语句
if count > 5 {
    print("Count 大于 5");
} else {
    print("Count 小于等于 5");
}

// 循环
var i = 0;
while i < 3 {
    print("循环次数:", i);
    i = i + 1;
}

// 调用你的Go工具函数
var users = db_query("SELECT * FROM users WHERE active = true");
log_info("查询到用户:", len(str(users)));

for user in users {
    var processed = process_data(user);
    log_info("处理用户:", processed);
}

// 字符串处理
var text = "  Hello, DSL!  ";
print("原始:", text);
print("大写:", upper(text));
print("小写:", lower(text));
print("修剪:", trim(text));

// 类型转换
var num_str = "123";
var num = int(num_str);
print("转换后的数字:", num + 100);

// 返回结果
return "脚本执行完成";