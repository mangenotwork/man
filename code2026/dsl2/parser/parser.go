package parser

/*

语法分析

*/

import (
	"dsl2/ast"
	"dsl2/lexer"
	"fmt"
	"log"
	"strconv"
)

type Parser struct {
	lexer  *lexer.Lexer
	curTok lexer.Token
	errors []string
	depth  int // 递归深度计数器
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: make([]string, 0),
		depth:  0,
	}
	p.nextToken() // 读取第一个token
	return p
}

// 最大递归深度
const maxParseDepth = 1000

func (p *Parser) checkDepth() bool {
	if p.depth > maxParseDepth {
		p.errors = append(p.errors, "解析递归深度超过限制")
		return false
	}
	return true
}

func (p *Parser) enter() {
	p.depth++
}

func (p *Parser) leave() {
	p.depth--
}

func (p *Parser) nextToken() {
	p.curTok = p.lexer.NextToken()
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curTok.Type == t
}

func (p *Parser) expect(t lexer.TokenType, msg ...string) bool {
	if !p.checkDepth() {
		return false
	}

	if p.curTokenIs(t) {
		p.nextToken()
		return true
	}

	// 构建错误消息
	errorMsg := fmt.Sprintf("第%d行第%d列: 期望 %s, 得到 %s",
		p.curTok.Line, p.curTok.Column, t, p.curTok.Type)

	if len(msg) > 0 {
		errorMsg += " (" + msg[0] + ")"
	}

	p.errors = append(p.errors, errorMsg)
	return false
}

func (p *Parser) Errors() []string {
	return p.errors
}

// Cleanup 清理解析器资源
func (p *Parser) Cleanup() {
	p.errors = nil
	p.lexer = nil
}

func (p *Parser) ParseProgram() *ast.Program {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	program := &ast.Program{
		StartPos: ast.Position{Line: 1, Column: 1},
	}

	for !p.curTokenIs(lexer.TokenEOF) {
		log.Println("p.curTok == ", p.curTok)
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		} else {
			// 如果解析失败，尝试恢复
			p.recoverFromError()
		}
	}

	return program
}

// recoverFromError 从解析错误中恢复
func (p *Parser) recoverFromError() {
	// 跳过当前token，直到找到语句结束标记
	for !p.curTokenIs(lexer.TokenEOF) &&
		!p.curTokenIs(lexer.TokenSemicolon) &&
		!p.curTokenIs(lexer.TokenRBrace) &&
		!p.curTokenIs(lexer.TokenLBrace) {
		p.nextToken()
	}

	if p.curTokenIs(lexer.TokenSemicolon) {
		p.nextToken()
	}
}

func (p *Parser) parseStatement() ast.Statement {
	if !p.checkDepth() {
		return nil
	}

	log.Println("parseStatement = ", p.curTok)

	p.enter()
	defer p.leave()

	switch p.curTok.Type {
	case lexer.TokenVar:
		return p.parseVarStatement()
	case lexer.TokenIf:
		return p.parseIfStatement()
	case lexer.TokenWhile:
		return p.parseWhileStatement()
	case lexer.TokenReturn:
		return p.parseReturnStatement()
	case lexer.TokenLBrace:
		return p.parseBlockStatement()
	case lexer.TokenChrome:
		return p.parseChromeStatement()
	case lexer.TokenBreak:
		return p.parseBreakStatement()
	case lexer.TokenContinue:
		return p.parseContinueStatement()
	case lexer.TokenFor:
		return p.parseForStatement()
	default:
		return p.parseSimpleStatement()
	}
}

func (p *Parser) parseSimpleStatement() ast.Statement {
	if !p.checkDepth() {
		return nil
	}
	log.Println("parseSimpleStatement = ", p.curTok)
	p.enter()
	defer p.leave()

	// 保存当前位置
	line := p.curTok.Line
	column := p.curTok.Column

	// 尝试解析表达式
	expr := p.parseExpression()
	if expr == nil {
		return nil
	}

	// 检查是否是赋值
	if p.curTokenIs(lexer.TokenAssign) {
		if ident, ok := expr.(*ast.Identifier); ok {
			p.nextToken() // 跳过 =
			right := p.parseExpression()
			if right == nil {
				return nil
			}

			stmt := &ast.AssignStmt{
				StartPos: ast.Position{
					Line:   line,
					Column: column,
				},
				Left: ident,
				Expr: right,
			}

			// 期望分号
			//p.expect(lexer.TokenSemicolon, "赋值语句后")
			if p.curTokenIs(lexer.TokenSemicolon) {
				p.nextToken() // 跳过 ;
			}
			return stmt
		} else {
			p.errors = append(p.errors,
				fmt.Sprintf("第%d行第%d列: 赋值目标必须是标识符", line, column))
			return nil
		}
	}

	// 如果不是赋值，就是表达式语句
	stmt := &ast.ExpressionStmt{
		StartPos: ast.Position{
			Line:   line,
			Column: column,
		},
		Expr: expr,
	}

	// 期望分号
	//p.expect(lexer.TokenSemicolon, "表达式语句后")

	if p.curTokenIs(lexer.TokenSemicolon) {
		p.nextToken() // 跳过 ;
	}

	return stmt
}

func (p *Parser) parseVarStatement() *ast.VarDecl {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	stmt := &ast.VarDecl{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
	}

	p.expect(lexer.TokenVar, "变量声明") // 跳过 var

	if !p.curTokenIs(lexer.TokenIdent) {
		p.errors = append(p.errors,
			fmt.Sprintf("第%d行第%d列: 期望标识符", p.curTok.Line, p.curTok.Column))
		return nil
	}

	stmt.Name = &ast.Identifier{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
		Name: p.curTok.Literal,
	}

	p.nextToken() // 跳过标识符

	// 可选类型注解
	if p.curTokenIs(lexer.TokenColon) {
		p.nextToken()
		if p.curTokenIs(lexer.TokenIdent) {
			stmt.Type = p.curTok.Literal
			p.nextToken()
		} else {
			p.errors = append(p.errors,
				fmt.Sprintf("第%d行第%d列: 期望类型标识符", p.curTok.Line, p.curTok.Column))
		}
	} else {
		stmt.Type = "auto"
	}

	// 可选初始化
	if p.curTokenIs(lexer.TokenAssign) {
		p.nextToken()
		stmt.Expr = p.parseExpression()
	}

	//p.expect(lexer.TokenSemicolon, "变量声明后")
	if p.curTokenIs(lexer.TokenSemicolon) {
		p.nextToken() // 跳过 ;
	}
	return stmt
}

func (p *Parser) parseChromeStatement() *ast.ChromeStmt {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()
	log.Println("parseChromeStatement = ", p.depth)

	p.expect(lexer.TokenChrome, "变量声明")

	var args []ast.Expression
	for i := 0; i <= p.depth; i++ {
		log.Println("[debug]", p.curTok)

		args = append(args, &ast.String{
			StartPos: ast.Position{
				Line:   p.curTok.Line,
				Column: p.curTok.Column,
			},
			Value: p.curTok.Literal,
		})
		p.nextToken()
	}

	log.Println("args  = ", args)

	stmt := &ast.ChromeStmt{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
		Args: args,
	}

	return stmt
}

func (p *Parser) parseIfStatement() *ast.IfStmt {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	stmt := &ast.IfStmt{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
	}

	p.expect(lexer.TokenIf, "if语句") // 跳过 if

	// 不需要括号，直接解析条件
	stmt.Condition = p.parseExpression()
	if stmt.Condition == nil {
		return nil
	}

	// 解析 then 块
	stmt.Then = p.parseBlockStatement()
	if stmt.Then == nil {
		return nil
	}

	// 解析可选的 else
	if p.curTokenIs(lexer.TokenElse) {
		p.nextToken()

		if p.curTokenIs(lexer.TokenIf) {
			stmt.Else = p.parseIfStatement()
		} else {
			stmt.Else = p.parseBlockStatement()
		}
	}

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStmt {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	stmt := &ast.WhileStmt{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
	}

	p.expect(lexer.TokenWhile, "while语句") // 跳过 while

	// 不需要括号，直接解析条件
	stmt.Condition = p.parseExpression()
	if stmt.Condition == nil {
		return nil
	}

	// 解析循环体
	stmt.Body = p.parseBlockStatement()
	if stmt.Body == nil {
		return nil
	}

	return stmt
}

func (p *Parser) parseBreakStatement() *ast.BreakStmt {
	if !p.checkDepth() {
		return nil
	}

	log.Println("parseBreakStatement = ", p.curTok)

	p.enter()
	defer p.leave()

	stmt := &ast.BreakStmt{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
	}

	// 跳过 break
	p.nextToken()

	// 期待分号
	//if !p.curTokenIs(lexer.TokenSemicolon) {
	//	p.errors = append(p.errors, "break语句后缺少分号")
	//	return nil
	//}
	if p.curTokenIs(lexer.TokenSemicolon) {
		p.nextToken() // 跳过 ;
	}

	return stmt
}

func (p *Parser) parseContinueStatement() *ast.ContinueStmt {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	stmt := &ast.ContinueStmt{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
	}

	// 跳过 continue
	p.nextToken()

	// 期待分号
	//if !p.curTokenIs(lexer.TokenSemicolon) {
	//	p.errors = append(p.errors, "continue语句后缺少分号")
	//	return nil
	//}
	if p.curTokenIs(lexer.TokenSemicolon) {
		p.nextToken() // 跳过 ;
	}

	return stmt
}

// 单个循环能正常，但是嵌套循环存在问题
func (p *Parser) parseForStatement_old() *ast.ForStmt {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	stmt := &ast.ForStmt{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
	}

	// 跳过 for
	p.expect(lexer.TokenFor, "for语句")

	// 先检查最简单的两种情况
	// 情况1: for { ... } 无限循环
	if p.curTokenIs(lexer.TokenLBrace) {
		stmt.Body = p.parseBlockStatement()
		return stmt
	}

	// 情况2: for condition { ... } (类似 while)
	// 保存当前位置以便回溯
	// 由于我们没有保存位置，我们需要一种不同的方法
	// 我们可以先尝试解析表达式，然后检查它后面是什么

	// 首先尝试解析一个语句（可能是初始化）
	initStmt := p.parseInitStatement()

	if initStmt == nil {
		// 如果没有初始化语句，检查是否是分号
		if p.curTokenIs(lexer.TokenSemicolon) {
			// 是分号，说明是标准 for 循环但没有初始化
			stmt.Init = nil
			p.nextToken() // 跳过第一个分号
		} else {
			// 不是分号，尝试解析条件表达式
			// 这可能是 for condition { ... } 形式
			cond := p.parseExpression()
			if cond == nil {
				p.errors = append(p.errors, "for语句解析错误")
				return nil
			}

			// 检查下一个 token
			if p.curTokenIs(lexer.TokenLBrace) {
				// 是 {，说明是 for condition { ... } 形式
				stmt.Cond = cond
				stmt.Body = p.parseBlockStatement()
				return stmt
			} else if p.curTokenIs(lexer.TokenSemicolon) {
				// 是分号，说明是标准 for 循环但没有初始化
				stmt.Cond = cond
				p.nextToken() // 跳过分号
				// 继续解析标准 for 循环
			} else {
				p.errors = append(p.errors, "for语句语法错误")
				return nil
			}
		}
	} else {
		// 有初始化语句
		stmt.Init = initStmt

		// 检查下一个 token
		if p.curTokenIs(lexer.TokenSemicolon) {
			// 是分号，标准 for 循环
			p.nextToken() // 跳过分号
		} else if p.curTokenIs(lexer.TokenLBrace) {
			// 是 {，说明是 for init { ... } 形式
			// 这里 init 可能是条件表达式
			// 我们需要重新判断
			// 更安全的方式是：如果 init 是表达式语句，且后面是 {，则应该作为条件
			if exprStmt, ok := initStmt.(*ast.ExpressionStmt); ok {
				// 这其实是条件表达式
				stmt.Init = nil
				stmt.Cond = exprStmt.Expr
				stmt.Body = p.parseBlockStatement()
				return stmt
			}
		} else {
			p.errors = append(p.errors, "for语句期望分号或左大括号")
			return nil
		}
	}

	// 解析初始化部分
	// 可能有三种情况：
	// 1. 有分号：for (; i < 10; i++)
	// 2. 有初始化语句：for (var i = 0; i < 10; i++)
	// 3. 没有初始化：for (i < 10) { ... } (类似 while)
	if p.curTokenIs(lexer.TokenSemicolon) {
		// 没有初始化语句
		stmt.Init = nil
		p.nextToken() // 跳过第一个分号
	} else {
		// 解析初始化语句
		stmt.Init = p.parseInitStatement()
		if stmt.Init == nil && !p.curTokenIs(lexer.TokenSemicolon) {
			p.errors = append(p.errors, "for语句初始化部分解析错误")
			return nil
		}

		// 期待分号
		//if !p.expect(lexer.TokenSemicolon, "分号") {
		//	return nil
		//}
		if p.curTokenIs(lexer.TokenSemicolon) {
			p.nextToken() // 跳过 ;
		}
	}

	// 解析条件部分
	if p.curTokenIs(lexer.TokenSemicolon) {
		// 没有条件表达式，相当于 true
		stmt.Cond = nil
		p.nextToken() // 跳过第二个分号
	} else {
		// 解析条件表达式
		stmt.Cond = p.parseExpression()

		// 期待分号
		//if !p.expect(lexer.TokenSemicolon, "分号") {
		//	return nil
		//}
		if p.curTokenIs(lexer.TokenSemicolon) {
			p.nextToken() // 跳过 ;
		}
	}

	// 解析后置部分
	if p.curTokenIs(lexer.TokenLBrace) {
		// 没有后置语句
		stmt.Post = nil
	} else {
		// 解析后置语句
		stmt.Post = p.parsePostStatement()
		if stmt.Post == nil && !p.curTokenIs(lexer.TokenLBrace) {
			p.errors = append(p.errors, "for语句后置部分解析错误")
			return nil
		}
	}

	// 解析循环体
	stmt.Body = p.parseBlockStatement()
	if stmt.Body == nil {
		return nil
	}

	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStmt {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	stmt := &ast.ForStmt{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
	}

	// 跳过 for
	if !p.expect(lexer.TokenFor, "for") {
		return nil
	}

	// 情况1: for { ... } 无限循环
	if p.curTokenIs(lexer.TokenLBrace) {
		stmt.Body = p.parseBlockStatement()
		return stmt
	}

	// 先尝试解析第一个部分
	// 它可能是：
	// 1. var 声明 (for var i = 0; ...)
	// 2. 表达式 (for i < 5; ... 或 for i = 0; ...)
	// 3. 分号 (for ; ...)

	// 检查是否是 var 声明
	if p.curTokenIs(lexer.TokenVar) {
		// 情况: for var i = 0; i < 5; i++
		stmt.Init = p.parseInitStatement() // 解析 var 声明

		//// 必须跟分号
		//if !p.curTokenIs(lexer.TokenSemicolon) {
		//	p.errors = append(p.errors, "for语句期望分号")
		//	return nil
		//}
		//p.nextToken() // 跳过分号

		if p.curTokenIs(lexer.TokenSemicolon) {
			p.nextToken() // 跳过 ;
		}

		// 继续解析剩余部分
		return p.parseForRemaining(stmt)
	} else if p.curTokenIs(lexer.TokenSemicolon) {
		// 情况: for ; i < 5; i++
		p.nextToken() // 跳过分号
		return p.parseForRemaining(stmt)
	} else {
		// 尝试解析为表达式
		// 这可能是：
		// 1. 条件表达式: for i < 5 { ... }
		// 2. 赋值表达式: for i = 0; i < 5; i++

		// 保存当前位置以便回溯
		// 由于不能直接回溯，我们需要先解析表达式，然后判断
		expr := p.parseExpression()
		if expr == nil {
			p.errors = append(p.errors, "for语句解析错误")
			return nil
		}

		// 检查下一个 token
		if p.curTokenIs(lexer.TokenLBrace) {
			// 表达式后面是 {，说明是 for condition { ... } 形式
			stmt.Cond = expr
			stmt.Body = p.parseBlockStatement()
			return stmt
		} else if p.curTokenIs(lexer.TokenSemicolon) {
			// 表达式后面是分号，说明是标准 for 循环
			stmt.Init = &ast.ExpressionStmt{Expr: expr}
			p.nextToken() // 跳过分号
			return p.parseForRemaining(stmt)
		} else {
			p.errors = append(p.errors, fmt.Sprintf("for语句语法错误，期望 '{' 或 ';'，得到: %s", p.curTok.Literal))
			return nil
		}
	}
}

// 解析 for 循环的剩余部分（条件和后置）
func (p *Parser) parseForRemaining(stmt *ast.ForStmt) *ast.ForStmt {
	// 解析条件部分
	if p.curTokenIs(lexer.TokenSemicolon) {
		// 没有条件
		p.nextToken() // 跳过分号
	} else {
		stmt.Cond = p.parseExpression()
		if stmt.Cond == nil {
			p.errors = append(p.errors, "for语句条件解析错误")
			return nil
		}

		//if !p.expect(lexer.TokenSemicolon, "分号") {
		//	return nil
		//}
		if p.curTokenIs(lexer.TokenSemicolon) {
			p.nextToken() // 跳过 ;
		}
	}

	// 解析后置部分
	if !p.curTokenIs(lexer.TokenLBrace) {
		stmt.Post = p.parsePostStatement()
		if stmt.Post == nil {
			p.errors = append(p.errors, "for语句后置部分解析错误")
			return nil
		}
	}

	// 解析循环体
	stmt.Body = p.parseBlockStatement()
	if stmt.Body == nil {
		return nil
	}

	return stmt
}

// 解析 for 循环的初始化语句
func (p *Parser) parseInitStatement() ast.Statement {
	// 如果直接是分号，表示没有初始化语句
	if p.curTokenIs(lexer.TokenSemicolon) {
		return nil
	}

	// 可能是 var 声明
	if p.curTokenIs(lexer.TokenVar) {
		return p.parseVarStatement()
	}

	// 保存当前 token 以便回溯
	saveTok := p.curTok
	// 尝试解析为表达式语句
	expr := p.parseExpression()
	if expr == nil {
		// 解析失败，恢复保存的 token
		p.curTok = saveTok
		return nil
	}

	// 检查是否是赋值表达式
	switch expr.(type) {
	case *ast.Identifier:
		// 只是一个标识符，可能是赋值语句的一部分
		// 检查下一个 token
		if p.curTokenIs(lexer.TokenAssign) {
			// 恢复并解析为赋值语句
			p.curTok = saveTok
			return p.parseSimpleStatement() // .parseAssignStmt()
		}
		// 否则，只是一个表达式语句
		return &ast.ExpressionStmt{Expr: expr}
	default:
		// 其他表达式
		return &ast.ExpressionStmt{Expr: expr}
	}
}

// 解析 for 循环的后置语句
func (p *Parser) parsePostStatement() ast.Statement {
	if p.curTokenIs(lexer.TokenLBrace) || p.curTokenIs(lexer.TokenRParen) {
		return nil
	}

	// 解析表达式语句
	return p.parseSimpleStatement()
}

func (p *Parser) parseReturnStatement() *ast.ReturnStmt {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	stmt := &ast.ReturnStmt{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
	}

	p.expect(lexer.TokenReturn, "return语句") // 跳过 return

	// 解析返回值（可选）
	if !p.curTokenIs(lexer.TokenSemicolon) {
		stmt.Expr = p.parseExpression()
	}

	//p.expect(lexer.TokenSemicolon, "return语句后")
	if p.curTokenIs(lexer.TokenSemicolon) {
		p.nextToken() // 跳过 ;
	}
	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStmt {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	block := &ast.BlockStmt{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
	}

	p.expect(lexer.TokenLBrace, "代码块开始") // 跳过 {

	for !p.curTokenIs(lexer.TokenRBrace) && !p.curTokenIs(lexer.TokenEOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Stmts = append(block.Stmts, stmt)
		} else {
			// 解析失败，尝试恢复
			p.recoverFromError()
		}
	}

	p.expect(lexer.TokenRBrace, "代码块结束") // 跳过 }
	return block
}

// 表达式解析
func (p *Parser) parseExpression() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	log.Println("parseExpression = ", p.curTok)

	p.enter()

	defer p.leave()

	return p.parseLogicalOr()
}

func (p *Parser) parseLogicalOr() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	log.Println("parseLogicalOr = ", p.curTok)

	expr := p.parseLogicalAnd()
	if expr == nil {
		return nil
	}

	for p.curTokenIs(lexer.TokenOr) {

		op := p.curTok.Literal
		p.nextToken()
		right := p.parseLogicalAnd()
		if right == nil {
			return expr
		}

		expr = &ast.BinaryExpr{
			StartPos: ast.Position{
				Line:   p.curTok.Line,
				Column: p.curTok.Column,
			},
			Left:  expr,
			Op:    op,
			Right: right,
		}
	}

	return expr
}

func (p *Parser) parseLogicalAnd() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()
	log.Println("parseLogicalAnd = ", p.curTok)
	expr := p.parseEquality()
	if expr == nil {
		return nil
	}

	for p.curTokenIs(lexer.TokenAnd) {

		op := p.curTok.Literal
		p.nextToken()
		right := p.parseEquality()
		if right == nil {
			return expr
		}

		expr = &ast.BinaryExpr{
			StartPos: ast.Position{
				Line:   p.curTok.Line,
				Column: p.curTok.Column,
			},
			Left:  expr,
			Op:    op,
			Right: right,
		}
	}

	return expr
}

func (p *Parser) parseEquality() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()
	log.Println("parseEquality = ", p.curTok)
	expr := p.parseComparison()
	if expr == nil {
		return nil
	}

	for p.curTokenIs(lexer.TokenEQ) || p.curTokenIs(lexer.TokenNE) {

		op := p.curTok.Literal
		p.nextToken()
		right := p.parseComparison()
		if right == nil {
			return expr
		}

		expr = &ast.BinaryExpr{
			StartPos: ast.Position{
				Line:   p.curTok.Line,
				Column: p.curTok.Column,
			},
			Left:  expr,
			Op:    op,
			Right: right,
		}
	}

	return expr
}

func (p *Parser) parseComparison() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()
	log.Println("parseComparison = ", p.curTok)
	expr := p.parseTerm()
	if expr == nil {
		return nil
	}

	for p.curTokenIs(lexer.TokenLT) || p.curTokenIs(lexer.TokenLE) ||
		p.curTokenIs(lexer.TokenGT) || p.curTokenIs(lexer.TokenGE) {

		op := p.curTok.Literal
		p.nextToken()
		right := p.parseTerm()
		if right == nil {
			return expr
		}

		expr = &ast.BinaryExpr{
			StartPos: ast.Position{
				Line:   p.curTok.Line,
				Column: p.curTok.Column,
			},
			Left:  expr,
			Op:    op,
			Right: right,
		}
	}

	return expr
}

func (p *Parser) parseTerm() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()
	log.Println("parseTerm =  ", p.curTok)
	expr := p.parseFactor()
	if expr == nil {
		return nil
	}

	for p.curTokenIs(lexer.TokenPlus) || p.curTokenIs(lexer.TokenMinus) {

		op := p.curTok.Literal
		p.nextToken()
		right := p.parseFactor()
		if right == nil {
			return expr
		}

		expr = &ast.BinaryExpr{
			StartPos: ast.Position{
				Line:   p.curTok.Line,
				Column: p.curTok.Column,
			},
			Left:  expr,
			Op:    op,
			Right: right,
		}
	}

	return expr
}

func (p *Parser) parseFactor() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()
	log.Println("parseFactor =  ", p.curTok)
	expr := p.parseUnary()
	if expr == nil {
		return nil
	}

	for p.curTokenIs(lexer.TokenAsterisk) || p.curTokenIs(lexer.TokenSlash) || p.curTokenIs(lexer.TokenMod) {
		op := p.curTok.Literal
		p.nextToken()
		right := p.parseUnary()
		if right == nil {
			return expr
		}

		expr = &ast.BinaryExpr{
			StartPos: ast.Position{
				Line:   p.curTok.Line,
				Column: p.curTok.Column,
			},
			Left:  expr,
			Op:    op,
			Right: right,
		}
	}

	return expr
}

func (p *Parser) parseUnary() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	log.Println("parseUnary =  ", p.curTok)

	if p.curTokenIs(lexer.TokenNot) || p.curTokenIs(lexer.TokenMinus) {

		op := p.curTok.Literal
		p.nextToken()
		expr := p.parseUnary()
		if expr == nil {
			return nil
		}

		return &ast.UnaryExpr{
			StartPos: ast.Position{
				Line:   p.curTok.Line,
				Column: p.curTok.Column,
			},
			Op:   op,
			Expr: expr,
		}
	}

	return p.parsePrimary()
}

func (p *Parser) parsePrimary() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	log.Println("parsePrimary =  ", p.curTok)

	switch p.curTok.Type {
	case lexer.TokenIdent:
		return p.parseIdentifierOrCall()
	case lexer.TokenInt:
		return p.parseInteger()
	case lexer.TokenString:
		return p.parseString()
	case lexer.TokenTrue, lexer.TokenFalse:
		return p.parseBoolean()
	case lexer.TokenLParen:
		return p.parseGroupedExpression()
	default:

		errStr := fmt.Sprintf("第%d行第%d列: 期望表达式，得到: %s",
			p.curTok.Line, p.curTok.Column, p.curTok.Type)
		log.Println(errStr)
		p.errors = append(p.errors, errStr)
		return nil
	}
}

func (p *Parser) parseIdentifierOrCall() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	// 保存标识符位置
	line := p.curTok.Line
	column := p.curTok.Column
	name := p.curTok.Literal
	log.Println("parseIdentifierOrCall = ", p.curTok.Type)
	p.nextToken()
	log.Println("parseIdentifierOrCall 2= ", p.curTok.Type)
	// 检查是否是函数调用
	if p.curTokenIs(lexer.TokenLParen) {
		p.nextToken() // 跳过 (

		var args []ast.Expression
		log.Println("parseIdentifierOrCall 3= ", p.curTok.Type)
		// 解析参数列表
		if !p.curTokenIs(lexer.TokenRParen) {
			for {
				log.Println("parseIdentifierOrCall 4= ", p.curTok.Type)
				// 解析参数
				arg := p.parseExpression()
				log.Println("parseIdentifierOrCall 5= ", arg, p.curTok.Type)
				if arg != nil {
					args = append(args, arg)
				} else {
					// 参数解析失败
					break
				}
				//log.Println(p.curTok.Type)
				// 检查是否有更多参数
				if p.curTokenIs(lexer.TokenComma) {
					p.nextToken() // 跳过逗号
					continue
				}

				break
			}
		}

		// 期望右括号
		log.Println(p.curTok.Type)
		if !p.expect(lexer.TokenRParen, "函数调用参数列表后") {
			return nil
		}

		return &ast.CallExpr{
			StartPos: ast.Position{
				Line:   line,
				Column: column,
			},
			Function: &ast.Identifier{
				StartPos: ast.Position{
					Line:   line,
					Column: column,
				},
				Name: name,
			},
			Args: args,
		}
	}

	// 不是函数调用，返回标识符
	return &ast.Identifier{
		StartPos: ast.Position{
			Line:   line,
			Column: column,
		},
		Name: name,
	}
}

func (p *Parser) parseInteger() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	value, err := strconv.ParseInt(p.curTok.Literal, 10, 64)
	if err != nil {
		p.errors = append(p.errors,
			fmt.Sprintf("第%d行第%d列: 无法解析整数: %s",
				p.curTok.Line, p.curTok.Column, p.curTok.Literal))
		return nil
	}

	expr := &ast.Integer{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
		Value: value,
	}

	p.nextToken()
	return expr
}

func (p *Parser) parseString() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	expr := &ast.String{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
		Value: p.curTok.Literal,
	}

	p.nextToken()
	return expr
}

func (p *Parser) parseBoolean() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	value := p.curTok.Type == lexer.TokenTrue
	expr := &ast.Boolean{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
		Value: value,
	}

	p.nextToken()
	return expr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	p.expect(lexer.TokenLParen, "分组表达式开始")
	expr := p.parseExpression()
	p.expect(lexer.TokenRParen, "分组表达式结束")
	return expr
}
