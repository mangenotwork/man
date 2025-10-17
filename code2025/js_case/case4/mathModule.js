// mathModule.js

// 命名导出：导出单个函数
export function add(a, b) {
    return a + b;
  }
  
  // 命名导出：导出另一个函数
  export function multiply(a, b) {
    return a * b;
  }
  
  // 先定义再批量导出
  function subtract(a, b) {
    return a - b;
  }
  
  const pi = 3.14159;
  
  // 批量命名导出（可以导出多个变量/函数）
  export { subtract, pi };
  
  // 默认导出：一个模块只能有一个默认导出
  export default function divide(a, b) {
    if (b === 0) throw new Error("除数不能为0");
    return a / b;
  }