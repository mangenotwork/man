// main.js

// 导入模块：
// - 命名导出需要用 {} 包裹，名称必须与导出时一致
// - 默认导出可以自定义名称（这里用 divideFunc）
import divideFunc, { add, multiply, subtract, pi } from './mathModule.js';

// 使用导入的功能
console.log("加法：", add(2, 3)); // 输出：加法：5
console.log("乘法：", multiply(4, 5)); // 输出：乘法：20
console.log("减法：", subtract(10, 4)); // 输出：减法：6
console.log("圆周率：", pi); // 输出：圆周率：3.14159
console.log("除法：", divideFunc(8, 2)); // 输出：除法：4