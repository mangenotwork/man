package ast

/*

AST

*/

import (
	"fmt"
	"strings"
)

// Node AST节点接口
type Node interface {
	Pos() Position
	String() string
}

// Position 源代码位置
type Position struct {
	Line   int
	Column int
}

// Program 程序节点
type Program struct {
	StartPos   Position
	Statements []Statement
}

func (p *Program) Pos() Position  { return p.StartPos }
func (p *Program) String() string { return "Program" }

// Statement 语句接口
type Statement interface {
	Node
	stmtNode()
}

type ExpressionStmt struct {
	StartPos Position
	Expr     Expression
}

func (e *ExpressionStmt) Pos() Position  { return e.StartPos }
func (e *ExpressionStmt) String() string { return e.Expr.String() }
func (e *ExpressionStmt) stmtNode()      {}

// Expression 表达式接口
type Expression interface {
	Node
	exprNode()
}

// Identifier 基础节点 标识符
type Identifier struct {
	StartPos Position
	Name     string
}

func (i *Identifier) Pos() Position  { return i.StartPos }
func (i *Identifier) String() string { return i.Name }
func (i *Identifier) exprNode()      {}

// Integer 整数字面量
type Integer struct {
	StartPos Position
	Value    int64
}

func (i *Integer) Pos() Position  { return i.StartPos }
func (i *Integer) String() string { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) exprNode()      {}

// String 字符串字面量
type String struct {
	StartPos Position
	Value    string
}

func (s *String) Pos() Position  { return s.StartPos }
func (s *String) String() string { return fmt.Sprintf("\"%s\"", s.Value) }
func (s *String) exprNode()      {}

// Boolean 布尔值
type Boolean struct {
	StartPos Position
	Value    bool
}

func (b *Boolean) Pos() Position  { return b.StartPos }
func (b *Boolean) String() string { return fmt.Sprintf("%v", b.Value) }
func (b *Boolean) exprNode()      {}

// BinaryExpr 表达式节点
type BinaryExpr struct {
	StartPos Position
	Left     Expression
	Op       string
	Right    Expression
}

func (b *BinaryExpr) Pos() Position { return b.StartPos }
func (b *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left, b.Op, b.Right)
}
func (b *BinaryExpr) exprNode() {}

// UnaryExpr 一元表达式
type UnaryExpr struct {
	StartPos Position
	Op       string
	Expr     Expression
}

func (u *UnaryExpr) Pos() Position { return u.StartPos }
func (u *UnaryExpr) String() string {
	return fmt.Sprintf("(%s%s)", u.Op, u.Expr)
}
func (u *UnaryExpr) exprNode() {}

// CallExpr 函数调用
type CallExpr struct {
	StartPos Position
	Function *Identifier
	Args     []Expression
}

func (c *CallExpr) Pos() Position { return c.StartPos }
func (c *CallExpr) String() string {
	args := make([]string, len(c.Args))
	for i, arg := range c.Args {
		args[i] = arg.String()
	}
	return fmt.Sprintf("%s(%s)", c.Function.Name, strings.Join(args, ", "))
}
func (c *CallExpr) exprNode() {}

// AssignStmt 赋值语句
type AssignStmt struct {
	StartPos Position
	Left     *Identifier
	Expr     Expression
}

func (a *AssignStmt) Pos() Position { return a.StartPos }
func (a *AssignStmt) String() string {
	return fmt.Sprintf("%s = %s", a.Left.Name, a.Expr.String())
}
func (a *AssignStmt) stmtNode() {}

// VarDecl 变量声明
type VarDecl struct {
	StartPos Position
	Name     *Identifier
	Type     string // 类型，如 "int", "string", "auto"
	Expr     Expression
}

func (v *VarDecl) Pos() Position { return v.StartPos }
func (v *VarDecl) String() string {
	if v.Expr != nil {
		return fmt.Sprintf("var %s %s = %s", v.Name.Name, v.Type, v.Expr.String())
	}
	return fmt.Sprintf("var %s %s", v.Name.Name, v.Type)
}
func (v *VarDecl) stmtNode() {}

// BlockStmt 块语句
type BlockStmt struct {
	StartPos Position
	Stmts    []Statement
}

func (b *BlockStmt) Pos() Position { return b.StartPos }
func (b *BlockStmt) String() string {
	stmts := make([]string, len(b.Stmts))
	for i, stmt := range b.Stmts {
		stmts[i] = stmt.String()
	}
	return fmt.Sprintf("{\n  %s\n}", strings.Join(stmts, "\n  "))
}
func (b *BlockStmt) stmtNode() {}

// IfStmt 控制流
type IfStmt struct {
	StartPos  Position
	Condition Expression
	Then      *BlockStmt
	Else      Statement // 可以是 nil, *BlockStmt, 或 *IfStmt
}

func (i *IfStmt) Pos() Position { return i.StartPos }
func (i *IfStmt) String() string {
	elseStr := ""
	if i.Else != nil {
		elseStr = fmt.Sprintf(" else %s", i.Else.String())
	}
	return fmt.Sprintf("if %s %s%s", i.Condition.String(), i.Then.String(), elseStr)
}
func (i *IfStmt) stmtNode() {}

// WhileStmt While循环
type WhileStmt struct {
	StartPos  Position
	Condition Expression
	Body      *BlockStmt
}

func (w *WhileStmt) Pos() Position { return w.StartPos }
func (w *WhileStmt) String() string {
	return fmt.Sprintf("while %s %s", w.Condition.String(), w.Body.String())
}
func (w *WhileStmt) stmtNode() {}

// BreakStmt break 语句
type BreakStmt struct {
	StartPos Position
}

func (b *BreakStmt) Pos() Position  { return b.StartPos }
func (b *BreakStmt) String() string { return "break" }
func (b *BreakStmt) stmtNode()      {}

// ContinueStmt continue 语句
type ContinueStmt struct {
	StartPos Position
}

func (c *ContinueStmt) Pos() Position  { return c.StartPos }
func (c *ContinueStmt) String() string { return "continue" }
func (c *ContinueStmt) stmtNode()      {}

// ReturnStmt 返回语句
type ReturnStmt struct {
	StartPos Position
	Expr     Expression
}

func (r *ReturnStmt) Pos() Position { return r.StartPos }
func (r *ReturnStmt) String() string {
	if r.Expr != nil {
		return fmt.Sprintf("return %s", r.Expr.String())
	}
	return "return"
}
func (r *ReturnStmt) stmtNode() {}

// ForStmt for循环
type ForStmt struct {
	StartPos Position
	Init     Statement  // 初始化语句，可以为 nil
	Cond     Expression // 条件表达式，可以为 nil
	Post     Statement  // 后置语句，可以为 nil
	Body     *BlockStmt
}

func (f *ForStmt) Pos() Position { return f.StartPos }
func (f *ForStmt) String() string {
	initStr := ""
	if f.Init != nil {
		initStr = f.Init.String()
	}

	condStr := ""
	if f.Cond != nil {
		condStr = f.Cond.String()
	}

	postStr := ""
	if f.Post != nil {
		postStr = f.Post.String()
	}

	return fmt.Sprintf("for %s; %s; %s %s", initStr, condStr, postStr, f.Body.String())
}
func (f *ForStmt) stmtNode() {}

// ChromeStmt chrome 关键字，操作chrome
type ChromeStmt struct {
	StartPos Position
	Args     []Expression
}

func (c *ChromeStmt) Pos() Position { return c.StartPos }
func (c *ChromeStmt) String() string {
	args := make([]string, len(c.Args))
	for i, arg := range c.Args {
		args[i] = arg.String()
	}
	return fmt.Sprintf("Chrome %s ", strings.Join(args, " "))
}
func (c *ChromeStmt) stmtNode() {}

// List 列表字面量
type List struct {
	StartPos Position
	Elements []Expression
}

func (l *List) Pos() Position { return l.StartPos }
func (l *List) String() string {
	elements := make([]string, len(l.Elements))
	for i, elem := range l.Elements {
		elements[i] = elem.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}
func (l *List) exprNode() {}

// IndexExpr 下标表达式
type IndexExpr struct {
	StartPos Position
	Left     Expression
	Index    Expression
}

func (i *IndexExpr) Pos() Position { return i.StartPos }
func (i *IndexExpr) String() string {
	return fmt.Sprintf("%s[%s]", i.Left.String(), i.Index.String())
}
func (i *IndexExpr) exprNode() {}

// Dict 字典字面量
type Dict struct {
	StartPos Position
	Pairs    map[Expression]Expression // 键值对
}

func (d *Dict) Pos() Position { return d.StartPos }
func (d *Dict) String() string {
	pairs := make([]string, 0, len(d.Pairs))
	for key, value := range d.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", key.String(), value.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}
func (d *Dict) exprNode() {}

// IndexAssignStmt 下标赋值语句
type IndexAssignStmt struct {
	StartPos Position
	Target   *IndexExpr
	Expr     Expression
}

func (a *IndexAssignStmt) Pos() Position { return a.StartPos }
func (a *IndexAssignStmt) String() string {
	return fmt.Sprintf("%s = %s", a.Target.String(), a.Expr.String())
}
func (a *IndexAssignStmt) stmtNode() {}
