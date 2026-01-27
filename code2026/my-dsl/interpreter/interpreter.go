package interpreter

/*

解释器

*/

import (
	"fmt"
	"my-dsl/ast"
	"reflect"
)

// Value 值接口
type Value interface{}

// Function 函数定义
type Function func(args []Value) (Value, error)

// Context 执行上下文
type Context struct {
	parent    *Context
	variables map[string]Value
	functions map[string]Function
	returnVal *Value
	hasReturn bool
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
	for _, stmt := range program.Statements {
		_ = i.evaluateStmt(stmt, i.global) // 忽略返回值

		if i.global.hasReturn {
			return *i.global.returnVal, nil
		}
	}

	return nil, nil
}

func (i *Interpreter) evaluateStmt(stmt ast.Statement, ctx *Context) Value {
	switch s := stmt.(type) {
	case *ast.VarDecl:
		return i.evaluateVarDecl(s, ctx)
	case *ast.AssignStmt:
		return i.evaluateAssignStmt(s, ctx)
	case *ast.ExpressionStmt:
		return i.evaluateExpr(s.Expr, ctx)
	case *ast.BlockStmt:
		return i.evaluateBlockStmt(s, ctx)
	case *ast.IfStmt:
		return i.evaluateIfStmt(s, ctx)
	case *ast.WhileStmt:
		return i.evaluateWhileStmt(s, ctx)
	case *ast.ReturnStmt:
		return i.evaluateReturnStmt(s, ctx)
	default:
		i.errors = append(i.errors, fmt.Errorf("不支持的语句类型: %T", stmt))
		return nil
	}
}

func (i *Interpreter) evaluateExpr(expr ast.Expression, ctx *Context) Value {
	switch e := expr.(type) {
	case *ast.Integer:
		return e.Value
	case *ast.String:
		return e.Value
	case *ast.Boolean:
		return e.Value
	case *ast.Identifier:
		if val, ok := ctx.GetVar(e.Name); ok {
			return val
		}
		i.errors = append(i.errors, fmt.Errorf("未定义的变量: %s", e.Name))
		return nil
	case *ast.BinaryExpr:
		return i.evaluateBinaryExpr(e, ctx)
	case *ast.UnaryExpr:
		return i.evaluateUnaryExpr(e, ctx)
	case *ast.CallExpr:
		return i.evaluateCallExpr(e, ctx)
	default:
		i.errors = append(i.errors, fmt.Errorf("不支持的表达式类型: %T", expr))
		return nil
	}
}

func (i *Interpreter) evaluateVarDecl(decl *ast.VarDecl, ctx *Context) Value {
	var value Value

	if decl.Expr != nil {
		value = i.evaluateExpr(decl.Expr, ctx)
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

	ctx.SetVar(decl.Name.Name, value)
	return value
}

func (i *Interpreter) evaluateAssignStmt(stmt *ast.AssignStmt, ctx *Context) Value {
	value := i.evaluateExpr(stmt.Expr, ctx)
	ctx.SetVar(stmt.Left.Name, value)
	return value
}

func (i *Interpreter) evaluateBinaryExpr(expr *ast.BinaryExpr, ctx *Context) Value {
	left := i.evaluateExpr(expr.Left, ctx)
	right := i.evaluateExpr(expr.Right, ctx)

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

func (i *Interpreter) evaluateUnaryExpr(expr *ast.UnaryExpr, ctx *Context) Value {
	right := i.evaluateExpr(expr.Expr, ctx)

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

func (i *Interpreter) evaluateCallExpr(expr *ast.CallExpr, ctx *Context) Value {
	fn, ok := ctx.GetFunc(expr.Function.Name)
	if !ok {
		i.errors = append(i.errors, fmt.Errorf("未定义的函数: %s", expr.Function.Name))
		return nil
	}

	args := make([]Value, len(expr.Args))
	for idx, arg := range expr.Args {
		args[idx] = i.evaluateExpr(arg, ctx)
	}

	result, err := fn(args)
	if err != nil {
		i.errors = append(i.errors, fmt.Errorf("函数调用错误 %s: %v", expr.Function.Name, err))
		return nil
	}

	return result
}

func (i *Interpreter) evaluateBlockStmt(block *ast.BlockStmt, ctx *Context) Value {
	// 创建一个新的作用域
	newCtx := NewContext(ctx)

	for _, stmt := range block.Stmts {
		_ = i.evaluateStmt(stmt, newCtx) // 忽略返回值

		if newCtx.hasReturn {
			ctx.hasReturn = true
			ctx.returnVal = newCtx.returnVal
			return *ctx.returnVal
		}
	}

	return nil
}

func (i *Interpreter) evaluateIfStmt(stmt *ast.IfStmt, ctx *Context) Value {
	condition := i.evaluateExpr(stmt.Condition, ctx)

	if i.bool(condition) {
		return i.evaluateBlockStmt(stmt.Then, ctx)
	} else if stmt.Else != nil {
		switch e := stmt.Else.(type) {
		case *ast.BlockStmt:
			return i.evaluateBlockStmt(e, ctx)
		case *ast.IfStmt:
			return i.evaluateIfStmt(e, ctx)
		}
	}

	return nil
}

func (i *Interpreter) evaluateWhileStmt(stmt *ast.WhileStmt, ctx *Context) Value {
	for {
		condition := i.evaluateExpr(stmt.Condition, ctx)
		if !i.bool(condition) {
			break
		}

		// 为循环体创建新上下文
		loopCtx := NewContext(ctx)

		// 执行循环体
		for _, stmt := range stmt.Body.Stmts {
			_ = i.evaluateStmt(stmt, loopCtx)

			if loopCtx.hasReturn {
				ctx.hasReturn = true
				ctx.returnVal = loopCtx.returnVal
				return *ctx.returnVal
			}
		}

		// 关键：从循环体作用域传播变量到父作用域
		for k, v := range loopCtx.variables {
			if _, ok := ctx.GetVar(k); ok {
				ctx.SetVar(k, v)
			}
		}
	}

	return nil
}

func (i *Interpreter) evaluateReturnStmt(stmt *ast.ReturnStmt, ctx *Context) Value {
	var value Value
	if stmt.Expr != nil {
		value = i.evaluateExpr(stmt.Expr, ctx)
	}

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
	return reflect.DeepEqual(left, right)
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
		default:
			return nil, fmt.Errorf("len() 不支持的类型: %T", args[0])
		}
	})
}
