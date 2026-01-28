package interpreter

/*

解释器

*/

import (
	"dsl2/ast"
	"fmt"
	"log"
	"os"
	"reflect"
)

// Value 值接口
type Value interface{}

// DictType 字典类型
type DictType map[Value]Value

// Function 函数定义
type Function func(args []Value) (Value, error)

// Context 执行上下文
type Context struct {
	parent      *Context
	variables   map[string]Value
	functions   map[string]Function
	returnVal   *Value
	hasReturn   bool
	hasBreak    bool
	hasContinue bool
}

// NewContext 创建上下文
func NewContext(parent *Context) *Context {
	return &Context{
		parent:    parent,
		variables: make(map[string]Value),
		functions: make(map[string]Function),
	}
}

// SetVar 设置变量
func (c *Context) SetVar(name string, value Value) {
	c.variables[name] = value
}

// GetVar 获取变量
func (c *Context) GetVar(name string) (Value, bool) {
	val, ok := c.variables[name]
	if !ok && c.parent != nil {
		return c.parent.GetVar(name)
	}
	return val, ok
}

// SetFunc 设置函数
func (c *Context) SetFunc(name string, fn Function) {
	c.functions[name] = fn
}

// GetFunc 获取函数
func (c *Context) GetFunc(name string) (Function, bool) {
	fn, ok := c.functions[name]
	if !ok && c.parent != nil {
		return c.parent.GetFunc(name)
	}
	return fn, ok
}

// Interpreter 解释器
type Interpreter struct {
	global *Context
	errors []error
}

// NewInterpreter 创建解释器
func NewInterpreter() *Interpreter {
	interp := &Interpreter{
		global: NewContext(nil),
		errors: []error{},
	}

	// 注册内置函数
	interp.registerBuiltins()

	return interp
}

// Global 返回全局上下文
func (i *Interpreter) Global() *Context {
	return i.global
}

// Interpret 执行AST
func (i *Interpreter) Interpret(program *ast.Program) (Value, error) {
	log.Println("执行AST ....")
	for n, stmt := range program.Statements {
		log.Println(n+1, " - Interpret ==> ", stmt)
		_ = i.evaluateStmt(stmt, i.global, n+1) // 忽略返回值

		if i.global.hasReturn {
			return *i.global.returnVal, nil
		}
	}

	return nil, nil
}

func (i *Interpreter) evaluateStmt(stmt ast.Statement, ctx *Context, hang int) Value {

	log.Println("evaluateStmt ==> ", stmt)

	switch s := stmt.(type) {
	case *ast.VarDecl:
		return i.evaluateVarDecl(s, ctx, hang)
	case *ast.AssignStmt:
		return i.evaluateAssignStmt(s, ctx, hang)
	case *ast.ExpressionStmt:
		return i.evaluateExpr(s.Expr, ctx, hang)
	case *ast.BlockStmt:
		return i.evaluateBlockStmt(s, ctx, hang)
	case *ast.IfStmt:
		return i.evaluateIfStmt(s, ctx, hang)
	case *ast.WhileStmt:
		return i.evaluateWhileStmt(s, ctx, hang)
	case *ast.ReturnStmt:
		return i.evaluateReturnStmt(s, ctx, hang)
	case *ast.ChromeStmt:
		return i.evaluateChromeStmt(s, ctx, hang)
	case *ast.BreakStmt:
		ctx.hasBreak = true
		return nil
	case *ast.ContinueStmt:
		ctx.hasContinue = true
		return nil
	case *ast.ForStmt:
		return i.evaluateForStmt(s, ctx, hang)
	case *ast.IndexAssignStmt: // 添加下标赋值处理
		return i.evaluateIndexAssignStmt(s, ctx, hang)
	default:
		log.Println("[Crash]len:", hang, " | ", fmt.Errorf("不支持的语句类型: %T", stmt))
		os.Exit(0)
	}
	return nil
}

func (i *Interpreter) evaluateExpr(expr ast.Expression, ctx *Context, hang int) Value {
	log.Println("evaluateExpr ==> ", expr)
	switch e := expr.(type) {
	case *ast.Integer:
		log.Println("evaluateExpr ast.Integer ==> ", e.Value)
		return e.Value
	case *ast.String:
		log.Println("evaluateExpr ast.String ==> ", e.Value)
		return e.Value
	case *ast.Boolean:
		log.Println("evaluateExpr ast.Boolean ==> ", e.Value)
		return e.Value
	case *ast.Identifier:
		log.Println("evaluateExpr ast.Identifier ==> ", e)
		if val, ok := ctx.GetVar(e.Name); ok {
			return val
		}
		log.Println("[Crash]len:", hang, " | ", fmt.Errorf("未定义的变量: %s", e.Name))
		os.Exit(0)
	case *ast.BinaryExpr:
		log.Println("evaluateExpr ast.BinaryExpr ==> ", e)
		return i.evaluateBinaryExpr(e, ctx, hang)
	case *ast.UnaryExpr:
		log.Println("evaluateExpr ast.UnaryExpr ==> ", e)
		return i.evaluateUnaryExpr(e, ctx, hang)

	case *ast.CallExpr:
		log.Println("evaluateExpr ast.CallExpr ==> ", e)
		return i.evaluateCallExpr(e, ctx, hang)

	case *ast.List: // 添加列表字面量求值
		log.Println("evaluateExpr ast.List ==> ", e)
		return i.evaluateList(e, ctx, hang)

	case *ast.IndexExpr: // 添加下标表达式求值
		log.Println("evaluateExpr ast.IndexExpr ==> ", e)
		return i.evaluateIndexExpr(e, ctx, hang)

	case *ast.Dict: // 添加字典字面量求值
		log.Println("evaluateExpr ast.Dict ==> ", e)
		return i.evaluateDict(e, ctx, hang)

	default:
		log.Println("[Crash]len:", hang, " | ", fmt.Errorf("不支持的表达式类型: %T", expr))
		os.Exit(0)
	}
	return nil
}

func (i *Interpreter) evaluateVarDecl(decl *ast.VarDecl, ctx *Context, hang int) Value {
	var value Value

	if decl.Expr != nil {
		value = i.evaluateExpr(decl.Expr, ctx, hang)
	} else {
		// 默认值
		switch decl.Type {
		case "int":
			value = int64(0)
		case "string":
			value = ""
		case "bool":
			value = false
		default:
			value = nil
		}
	}
	log.Println("evaluateVarDecl ==> ", decl.Name.Name, ":", value)
	ctx.SetVar(decl.Name.Name, value)
	return value
}

func (i *Interpreter) evaluateAssignStmt(stmt *ast.AssignStmt, ctx *Context, hang int) Value {
	value := i.evaluateExpr(stmt.Expr, ctx, hang)
	log.Println("evaluateAssignStmt ==> ", stmt.Left.Name, ":", value)
	ctx.SetVar(stmt.Left.Name, value)
	return value
}

func (i *Interpreter) evaluateBinaryExpr(expr *ast.BinaryExpr, ctx *Context, hang int) Value {
	left := i.evaluateExpr(expr.Left, ctx, hang)
	right := i.evaluateExpr(expr.Right, ctx, hang)

	// 类型检查和转换
	switch expr.Op {
	case "+":
		return i.add(left, right)
	case "-":
		return i.sub(left, right)
	case "*":
		return i.mul(left, right)
	case "/":
		return i.div(left, right)
	case "%":
		return i.mod(left, right)
	case "==":
		return i.equal(left, right)
	case "!=":
		return !i.equal(left, right)
	case "<":
		return i.less(left, right)
	case "<=":
		return i.less(left, right) || i.equal(left, right)
	case ">":
		return i.greater(left, right)
	case ">=":
		return i.greater(left, right) || i.equal(left, right)
	case "&&":
		return i.bool(left) && i.bool(right)
	case "||":
		return i.bool(left) || i.bool(right)
	default:
		i.errors = append(i.errors, fmt.Errorf("不支持的操作符: %s", expr.Op))
		return nil
	}
}

func (i *Interpreter) evaluateUnaryExpr(expr *ast.UnaryExpr, ctx *Context, hang int) Value {
	right := i.evaluateExpr(expr.Expr, ctx, hang)

	switch expr.Op {
	case "-":
		switch v := right.(type) {
		case int64:
			return -v
		case float64:
			return -v
		default:
			i.errors = append(i.errors, fmt.Errorf("不支持的操作: -%T", right))
			return nil
		}
	case "!":
		return !i.bool(right)
	default:
		i.errors = append(i.errors, fmt.Errorf("不支持的操作符: %s", expr.Op))
		return nil
	}
}

func (i *Interpreter) evaluateCallExpr(expr *ast.CallExpr, ctx *Context, hang int) Value {
	fn, ok := ctx.GetFunc(expr.Function.Name)
	if !ok {
		i.errors = append(i.errors, fmt.Errorf("未定义的函数: %s", expr.Function.Name))
		return nil
	}

	args := make([]Value, len(expr.Args))
	for idx, arg := range expr.Args {
		args[idx] = i.evaluateExpr(arg, ctx, hang)
	}

	result, err := fn(args)
	if err != nil {
		i.errors = append(i.errors, fmt.Errorf("函数调用错误 %s: %v", expr.Function.Name, err))
		return nil
	}

	return result
}

func (i *Interpreter) evaluateChromeStmt(expr *ast.ChromeStmt, ctx *Context, hang int) Value {
	log.Println("evaluateChromeStmt args = ", expr.Args)
	fn, ok := ctx.GetFunc("chrome")
	if !ok {
		i.errors = append(i.errors, fmt.Errorf("未定义Chrome"))
		return nil
	}
	args := make([]Value, len(expr.Args))
	for idx, arg := range expr.Args {
		args[idx] = i.evaluateExpr(arg, ctx, hang)
	}

	result, err := fn(args)
	if err != nil {
		i.errors = append(i.errors, fmt.Errorf("Chrome调用错误: %v", err))
		return nil
	}

	return result
}

func (i *Interpreter) evaluateBlockStmt(block *ast.BlockStmt, ctx *Context, hang int) Value {
	// 创建一个新的作用域
	newCtx := NewContext(ctx)
	log.Println("evaluateBlockStmt ==> ", block.Stmts)

	for _, stmt := range block.Stmts {
		log.Println("evaluateBlockStmt is stmt item ==> ", stmt)

		switch stmt.(type) {
		case *ast.BreakStmt:
			ctx.hasBreak = true
			return nil
		case *ast.ContinueStmt:
			ctx.hasContinue = true
			return nil
		default:
			_ = i.evaluateStmt(stmt, newCtx, hang)

			if newCtx.hasReturn || newCtx.hasBreak || newCtx.hasContinue {
				ctx.hasReturn = true
				ctx.returnVal = newCtx.returnVal
				return *ctx.returnVal
			}
		}
	}

	return nil
}

func (i *Interpreter) evaluateIfStmt(stmt *ast.IfStmt, ctx *Context, hang int) Value {
	condition := i.evaluateExpr(stmt.Condition, ctx, hang)
	log.Println("evaluateBlockStmt ==> ", stmt)
	if i.bool(condition) {
		return i.evaluateBlockStmt(stmt.Then, ctx, hang)
	} else if stmt.Else != nil {
		switch e := stmt.Else.(type) {
		case *ast.BlockStmt:
			return i.evaluateBlockStmt(e, ctx, hang)
		case *ast.IfStmt:
			return i.evaluateIfStmt(e, ctx, hang)
		}
	}

	return nil
}

func (i *Interpreter) evaluateWhileStmt(stmt *ast.WhileStmt, ctx *Context, hang int) Value {
	log.Println("evaluateWhileStmt ==> ", stmt)
	for {
		// 检查循环条件
		condition := i.evaluateExpr(stmt.Condition, ctx, hang)
		if !i.bool(condition) {
			break
		}

		// 执行循环体
		loopCtx := NewContext(ctx)

		// 执行循环体
		for _, stmtItem := range stmt.Body.Stmts {
			_ = i.evaluateStmt(stmtItem, loopCtx, hang)

			// 检查是否需要提前退出
			if loopCtx.hasReturn || loopCtx.hasBreak || loopCtx.hasContinue {
				//log.Println("检查到要提前退出")
				break
			}
		}

		// 处理控制流
		if loopCtx.hasReturn {
			ctx.hasReturn = true
			ctx.returnVal = loopCtx.returnVal
			return *ctx.returnVal
		}

		if loopCtx.hasBreak {
			break
		}

		// 如果父作用域中本来没有这个变量，但现在有了，也要设置
		for k, v := range loopCtx.variables {
			ctx.SetVar(k, v) // 直接设置，覆盖原有的值
		}

		if loopCtx.hasContinue {
			continue
		}

	}

	return nil
}

func (i *Interpreter) evaluateForStmt(stmt *ast.ForStmt, ctx *Context, hang int) Value {
	// 1. 执行初始化语句
	if stmt.Init != nil {
		_ = i.evaluateStmt(stmt.Init, ctx, hang)
		if ctx.hasReturn || ctx.hasBreak || ctx.hasContinue {
			return nil
		}
	}

	// 2. 循环主体
	for {
		// 2.1 检查循环条件
		if stmt.Cond != nil {
			condition := i.evaluateExpr(stmt.Cond, ctx, hang)
			if !i.bool(condition) {
				break
			}
		}
		// 如果没有条件表达式，相当于条件永远为 true

		// 2.2 执行循环体
		for _, stmtItem := range stmt.Body.Stmts {
			_ = i.evaluateStmt(stmtItem, ctx, hang)

			// 检查控制流
			if ctx.hasReturn {
				return *ctx.returnVal
			}
			if ctx.hasBreak {
				ctx.hasBreak = false
				return nil
			}
			if ctx.hasContinue {
				ctx.hasContinue = false
				break
			}
		}

		// 如果是 continue，直接执行后置语句
		if ctx.hasContinue {
			ctx.hasContinue = false
		}

		// 2.3 执行后置语句
		if stmt.Post != nil {
			_ = i.evaluateStmt(stmt.Post, ctx, hang)
			if ctx.hasReturn || ctx.hasBreak || ctx.hasContinue {
				return nil
			}
		}

		// 如果遇到 break，已经在上面的检查中返回
	}

	return nil
}

func (i *Interpreter) evaluateReturnStmt(stmt *ast.ReturnStmt, ctx *Context, hang int) Value {
	var value Value
	if stmt.Expr != nil {
		value = i.evaluateExpr(stmt.Expr, ctx, hang)
	}
	log.Println("evaluateReturnStmt ==> ", stmt, value)
	ctx.hasReturn = true
	ctx.returnVal = &value
	return value
}

// 辅助函数
func (i *Interpreter) bool(v Value) bool {
	switch val := v.(type) {
	case bool:
		return val
	case int64:
		return val != 0
	case float64:
		return val != 0
	case string:
		return val != ""
	default:
		return v != nil
	}
}

func (i *Interpreter) add(left, right Value) Value {
	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return l + r
		case float64:
			return float64(l) + r
		case string:
			return fmt.Sprintf("%d%s", l, r)
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return l + float64(r)
		case float64:
			return l + r
		case string:
			return fmt.Sprintf("%f%s", l, r)
		}
	case string:
		return fmt.Sprintf("%s%v", l, right)

	case []Value: // 列表加法（连接）
		switch r := right.(type) {
		case []Value:
			// 连接两个列表
			result := make([]Value, len(l)+len(r))
			copy(result, l)
			copy(result[len(l):], r)
			return result
		default:
			// 将其他值添加到列表末尾
			result := make([]Value, len(l)+1)
			copy(result, l)
			result[len(l)] = r
			return result
		}

	case DictType: // 字典合并
		switch r := right.(type) {
		case DictType:
			// 合并两个字典
			result := make(DictType)
			// 先复制左边字典的所有元素
			for k, v := range l {
				result[k] = v
			}
			// 然后复制右边字典的所有元素（右边的会覆盖左边的）
			for k, v := range r {
				result[k] = v
			}
			return result
		}
	}

	i.errors = append(i.errors, fmt.Errorf("不支持的操作: %T + %T", left, right))
	return nil
}

func (i *Interpreter) sub(left, right Value) Value {
	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return l - r
		case float64:
			return float64(l) - r
		default:
			i.errors = append(i.errors, fmt.Errorf("不支持的操作: int64 - %T", right))
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return l - float64(r)
		case float64:
			return l - r
		default:
			i.errors = append(i.errors, fmt.Errorf("不支持的操作: float64 - %T", right))
		}
	default:
		i.errors = append(i.errors, fmt.Errorf("不支持的操作: %T - %T", left, right))
	}
	return nil
}

func (i *Interpreter) mul(left, right Value) Value {
	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return l * r
		case float64:
			return float64(l) * r
		default:
			i.errors = append(i.errors, fmt.Errorf("不支持的操作: int64 * %T", right))
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return l * float64(r)
		case float64:
			return l * r
		default:
			i.errors = append(i.errors, fmt.Errorf("不支持的操作: float64 * %T", right))
		}
	default:
		i.errors = append(i.errors, fmt.Errorf("不支持的操作: %T * %T", left, right))
	}
	return nil
}

func (i *Interpreter) div(left, right Value) Value {
	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			if r == 0 {
				i.errors = append(i.errors, fmt.Errorf("除零错误"))
				return nil
			}
			return l / r
		case float64:
			if r == 0 {
				i.errors = append(i.errors, fmt.Errorf("除零错误"))
				return nil
			}
			return float64(l) / r
		default:
			i.errors = append(i.errors, fmt.Errorf("不支持的操作: int64 / %T", right))
		}
	case float64:
		switch r := right.(type) {
		case int64:
			if r == 0 {
				i.errors = append(i.errors, fmt.Errorf("除零错误"))
				return nil
			}
			return l / float64(r)
		case float64:
			if r == 0 {
				i.errors = append(i.errors, fmt.Errorf("除零错误"))
				return nil
			}
			return l / r
		default:
			i.errors = append(i.errors, fmt.Errorf("不支持的操作: float64 / %T", right))
		}
	default:
		i.errors = append(i.errors, fmt.Errorf("不支持的操作: %T / %T", left, right))
	}
	return nil
}

func (i *Interpreter) mod(left, right Value) Value {
	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			if r == 0 {
				i.errors = append(i.errors, fmt.Errorf("模零错误"))
				return nil
			}
			return l % r
		default:
			i.errors = append(i.errors, fmt.Errorf("不支持的操作: int64 %% %T", right))
		}
	default:
		i.errors = append(i.errors, fmt.Errorf("不支持的操作: %T %% %T", left, right))
	}
	return nil
}

func (i *Interpreter) equal(left, right Value) bool {
	// 如果是列表，需要深度比较
	if l, ok := left.([]Value); ok {
		if r, ok := right.([]Value); ok {
			if len(l) != len(r) {
				return false
			}
			for idx := 0; idx < len(l); idx++ {
				if !i.valuesEqual(l[idx], r[idx]) { // 使用辅助函数比较元素
					return false
				}
			}
			return true
		}
		return false
	}
	return reflect.DeepEqual(left, right)
}

// 添加辅助函数
func (i *Interpreter) valuesEqual(left, right Value) bool {
	// 递归处理嵌套列表和字典
	switch l := left.(type) {
	case []Value: // 列表
		if r, ok := right.([]Value); ok {
			if len(l) != len(r) {
				return false
			}
			for idx := 0; idx < len(l); idx++ {
				if !i.valuesEqual(l[idx], r[idx]) {
					return false
				}
			}
			return true
		}
		return false

	case DictType: // 字典
		if r, ok := right.(DictType); ok {
			if len(l) != len(r) {
				return false
			}
			for key, lval := range l {
				if rval, exists := r[key]; exists {
					if !i.valuesEqual(lval, rval) {
						return false
					}
				} else {
					return false
				}
			}
			return true
		}
		return false

	default:
		return reflect.DeepEqual(left, right)
	}
}
func (i *Interpreter) less(left, right Value) bool {
	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return l < r
		case float64:
			return float64(l) < r
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return l < float64(r)
		case float64:
			return l < r
		}
	case string:
		switch r := right.(type) {
		case string:
			return l < r
		}
	}

	i.errors = append(i.errors, fmt.Errorf("不支持的操作: %T < %T", left, right))
	return false
}

func (i *Interpreter) greater(left, right Value) bool {
	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return l > r
		case float64:
			return float64(l) > r
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return l > float64(r)
		case float64:
			return l > r
		}
	case string:
		switch r := right.(type) {
		case string:
			return l > r
		}
	}

	i.errors = append(i.errors, fmt.Errorf("不支持的操作: %T > %T", left, right))
	return false
}

// 注册内置函数
func (i *Interpreter) registerBuiltins() {
	// 打印函数
	i.global.SetFunc("print", func(args []Value) (Value, error) {
		for _, arg := range args {
			fmt.Print(arg, " ")
		}
		fmt.Println()
		return nil, nil
	})

	i.global.SetFunc("println", func(args []Value) (Value, error) {
		for _, arg := range args {
			fmt.Println(arg)
		}
		return nil, nil
	})

	// 类型转换
	i.global.SetFunc("int", func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("int() 需要一个参数")
		}

		switch v := args[0].(type) {
		case int64:
			return v, nil
		case float64:
			return int64(v), nil
		case string:
			var result int64
			_, err := fmt.Sscanf(v, "%d", &result)
			if err != nil {
				return nil, fmt.Errorf("无法转换字符串为int: %s", v)
			}
			return result, nil
		default:
			return nil, fmt.Errorf("无法转换为int: %T", args[0])
		}
	})

	i.global.SetFunc("str", func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("str() 需要一个参数")
		}
		return fmt.Sprintf("%v", args[0]), nil
	})

	// 数学函数
	i.global.SetFunc("len", func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("len() 需要一个参数")
		}

		switch v := args[0].(type) {
		case string:
			return int64(len(v)), nil
		case []Value: // 添加对列表的支持
			return int64(len(v)), nil
		case DictType: // 字典
			return int64(len(v)), nil
		default:
			return nil, fmt.Errorf("len() 不支持的类型: %T", args[0])
		}
	})

	// 字典的 keys
	i.global.SetFunc("keys", func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("keys() 需要一个参数")
		}

		dict, ok := args[0].(DictType)
		if !ok {
			return nil, fmt.Errorf("keys() 只支持字典，得到: %T", args[0])
		}

		// 获取所有键
		keys := make([]Value, 0, len(dict))
		for key := range dict {
			keys = append(keys, key)
		}
		return keys, nil
	})

	// 字典的 values
	i.global.SetFunc("values", func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("values() 需要一个参数")
		}

		dict, ok := args[0].(DictType)
		if !ok {
			return nil, fmt.Errorf("values() 只支持字典，得到: %T", args[0])
		}

		// 获取所有值
		values := make([]Value, 0, len(dict))
		for _, value := range dict {
			values = append(values, value)
		}
		return values, nil
	})

	// 字典的 items
	i.global.SetFunc("items", func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("items() 需要一个参数")
		}

		dict, ok := args[0].(DictType)
		if !ok {
			return nil, fmt.Errorf("items() 只支持字典，得到: %T", args[0])
		}

		// 获取所有键值对（每个键值对是一个包含两个元素的列表）
		items := make([]Value, 0, len(dict))
		for key, value := range dict {
			pair := []Value{key, value}
			items = append(items, pair)
		}
		return items, nil
	})

	// 字典的 has_key
	i.global.SetFunc("has_key", func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("has_key() 需要两个参数: dict, key")
		}

		dict, ok := args[0].(DictType)
		if !ok {
			return nil, fmt.Errorf("has_key() 第一个参数必须是字典，得到: %T", args[0])
		}

		// 检查键是否存在
		_, exists := dict[args[1]]
		return exists, nil
	})

	// 字典的 delete
	i.global.SetFunc("delete", func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("delete() 需要两个参数: dict, key")
		}

		dict, ok := args[0].(DictType)
		if !ok {
			return nil, fmt.Errorf("delete() 第一个参数必须是字典，得到: %T", args[0])
		}

		// 删除键
		delete(dict, args[1])
		return nil, nil
	})
}

func (i *Interpreter) evaluateList(list *ast.List, ctx *Context, hang int) Value {
	elements := make([]Value, len(list.Elements))
	for idx, element := range list.Elements {
		elements[idx] = i.evaluateExpr(element, ctx, hang)
	}
	return elements
}

func (i *Interpreter) evaluateIndexExpr(expr *ast.IndexExpr, ctx *Context, hang int) Value {
	// 求值左边的表达式（应该是列表或字典）
	left := i.evaluateExpr(expr.Left, ctx, hang)

	// 求值下标
	index := i.evaluateExpr(expr.Index, ctx, hang)

	// 检查左边是列表还是字典
	switch container := left.(type) {
	case []Value: // 列表
		// 检查下标是否是整数
		idx, ok := index.(int64)
		if !ok {
			i.errors = append(i.errors, fmt.Errorf("列表下标必须是整数，得到: %T", index))
			return nil
		}

		// 检查下标是否越界
		if idx < 0 || idx >= int64(len(container)) {
			i.errors = append(i.errors, fmt.Errorf("列表下标越界: 长度=%d, 下标=%d", len(container), idx))
			return nil
		}

		return container[idx]

	case DictType: // 字典
		// 检查键是否是可哈希的类型
		if !i.isHashable(index) {
			i.errors = append(i.errors, fmt.Errorf("字典键必须是可哈希的类型，得到: %T", index))
			return nil
		}

		// 查找键对应的值
		value, exists := container[index]
		if !exists {
			i.errors = append(i.errors, fmt.Errorf("字典中不存在键: %v", index))
			return nil
		}

		return value

	default:
		i.errors = append(i.errors, fmt.Errorf("下标操作只支持列表或字典，得到: %T", left))
		return nil
	}
}

func (i *Interpreter) evaluateDict(dict *ast.Dict, ctx *Context, hang int) Value {
	result := make(DictType)

	for keyExpr, valueExpr := range dict.Pairs {
		// 求值键
		key := i.evaluateExpr(keyExpr, ctx, hang)

		// 检查键的类型（在Go中，只有可比较的类型才能作为map的键）
		// 我们只支持基本类型作为键
		if !i.isHashable(key) {
			i.errors = append(i.errors, fmt.Errorf("字典键必须是可哈希的类型，得到: %T", key))
			return nil
		}

		// 求值值
		value := i.evaluateExpr(valueExpr, ctx, hang)

		// 添加到字典
		result[key] = value
	}

	return result
}

// 检查值是否可以作为字典的键
func (i *Interpreter) isHashable(value Value) bool {
	switch value.(type) {
	case int64, float64, string, bool:
		return true
	default:
		return false
	}
}

func (i *Interpreter) evaluateIndexAssignStmt(stmt *ast.IndexAssignStmt, ctx *Context, hang int) Value {
	// 求值右边的表达式
	value := i.evaluateExpr(stmt.Expr, ctx, hang)

	// 获取目标容器和键/下标
	target := stmt.Target.Left
	key := stmt.Target.Index

	// 求值目标容器
	container := i.evaluateExpr(target, ctx, hang)

	// 求值键/下标
	index := i.evaluateExpr(key, ctx, hang)

	// 检查容器类型并赋值
	switch c := container.(type) {
	case []Value: // 列表
		// 检查下标是否是整数
		idx, ok := index.(int64)
		if !ok {
			i.errors = append(i.errors, fmt.Errorf("列表下标必须是整数，得到: %T", index))
			return nil
		}

		// 检查下标是否越界
		if idx < 0 || idx >= int64(len(c)) {
			i.errors = append(i.errors, fmt.Errorf("列表下标越界: 长度=%d, 下标=%d", len(c), idx))
			return nil
		}

		// 赋值
		c[idx] = value

	case DictType: // 字典
		// 检查键是否是可哈希的类型
		if !i.isHashable(index) {
			i.errors = append(i.errors, fmt.Errorf("字典键必须是可哈希的类型，得到: %T", index))
			return nil
		}

		// 赋值（添加或修改）
		c[index] = value

	default:
		i.errors = append(i.errors, fmt.Errorf("下标赋值只支持列表或字典，得到: %T", container))
		return nil
	}

	return value
}
