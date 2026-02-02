package parser

/*

语法分析

*/

import (
	"dsl2/ast"
	"dsl2/lexer"
	"dsl2/logger"
	"fmt"
	"strconv"
	"strings"
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
		logger.Debug("p.curTok == ", p.curTok)
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

	logger.Debug("parseStatement = ", p.curTok)

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
	logger.Debug("parseSimpleStatement = ", p.curTok)
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
		// 检查左边是否是标识符或下标表达式
		switch left := expr.(type) {
		case *ast.Identifier:
			// 普通变量赋值
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
				Left: left,
				Expr: right,
			}

			// 跳过分号
			if p.curTokenIs(lexer.TokenSemicolon) {
				p.nextToken() // 跳过 ;
			}
			return stmt

		case *ast.IndexExpr:
			// 字典或列表元素赋值
			p.nextToken() // 跳过 =
			right := p.parseExpression()
			if right == nil {
				return nil
			}

			stmt := &ast.IndexAssignStmt{
				StartPos: ast.Position{
					Line:   line,
					Column: column,
				},
				Target: left,
				Expr:   right,
			}

			// 跳过分号
			if p.curTokenIs(lexer.TokenSemicolon) {
				p.nextToken() // 跳过 ;
			}
			return stmt

		default:
			p.errors = append(p.errors,
				fmt.Sprintf("第%d行第%d列: 赋值目标必须是标识符或下标表达式", line, column))
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

	// 跳过分号
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

	if p.curTokenIs(lexer.TokenSemicolon) {
		p.nextToken() // 跳过 ;
	}
	return stmt
}

// Chrome 语法解析 chrome arg1 arg2=123 ...
func (p *Parser) parseChromeStatement() *ast.ChromeStmt {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()
	logger.Debug("parseChromeStatement = ", p.depth)

	p.expect(lexer.TokenChrome, "变量声明")

	var args []ast.Expression
	for i := 0; i < p.depth; i++ {
		logger.Debug("[debug]", p.curTok)

		args = append(args, &ast.String{
			StartPos: ast.Position{
				Line:   p.curTok.Line,
				Column: p.curTok.Column,
			},
			Value: p.curTok.Literal,
		})
		p.nextToken()
	}

	logger.Debug("args  = ", args)

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

	logger.Debug("parseBreakStatement = ", p.curTok)

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

	// 跳过分号
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

	// 跳过期待分号
	if p.curTokenIs(lexer.TokenSemicolon) {
		p.nextToken() // 跳过 ;
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

	// 检查是否是 var 声明
	if p.curTokenIs(lexer.TokenVar) {
		// 情况: for var i = 0; i < 5; i++
		stmt.Init = p.parseInitStatement() // 解析 var 声明
		// 跳过分号
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

	logger.Debug("parseExpression: 开始解析表达式，当前token: %v", p.curTok)

	p.enter()
	defer p.leave()

	expr := p.parseLogicalOr()
	if expr == nil {
		// 解析失败，尝试跳过当前token
		logger.Debug("parseExpression: 解析失败，跳过token: %v", p.curTok)
		p.nextToken()
		return nil
	}

	logger.Debug("parseExpression: 解析结果: %T(%v)", expr, expr)
	return expr
}

func (p *Parser) parseLogicalOr() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	logger.Debug("parseLogicalOr = ", p.curTok)

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
	logger.Debug("parseLogicalAnd = ", p.curTok)
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
	logger.Debug("parseEquality = ", p.curTok)
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
	logger.Debug("parseComparison = ", p.curTok)
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
	logger.Debug("parseTerm =  ", p.curTok)
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
	logger.Debug("parseFactor =  ", p.curTok)
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

	logger.Debug("parseUnary =  ", p.curTok)

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

	// 修改这里，让 parseUnary 调用 parseCallOrIndex
	return p.parseCallOrIndex()
}

func (p *Parser) parsePrimary() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	logger.Debug("parsePrimary =  ", p.curTok)

	switch p.curTok.Type {
	case lexer.TokenIdent:
		return p.parseIdentifierOrCall()
	case lexer.TokenInt:
		// 检查是否是浮点数（包含小数点）
		if strings.Contains(p.curTok.Literal, ".") {
			return p.parseFloat()
		}
		return p.parseInteger()
	case lexer.TokenFloat:
		return p.parseFloat()
	case lexer.TokenString:
		return p.parseString()
	case lexer.TokenTrue, lexer.TokenFalse:
		return p.parseBoolean()
	case lexer.TokenLParen:
		return p.parseGroupedExpression()
	case lexer.TokenLBracket: // 添加列表字面量解析
		return p.parseListLiteral()
	case lexer.TokenLBrace: // 字典字面量
		return p.parseDictLiteral()
	default:

		errStr := fmt.Sprintf("第%d行第%d列: 期望表达式，得到: %s",
			p.curTok.Line, p.curTok.Column, p.curTok.Type)
		logger.Debug(errStr)
		p.errors = append(p.errors, errStr)
		return nil
	}
}

// 函数解析
func (p *Parser) parseIdentifierOrCall() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	logger.Debug("parseIdentifierOrCall: 当前token: %v", p.curTok)

	// 保存标识符位置
	line := p.curTok.Line
	column := p.curTok.Column
	name := p.curTok.Literal

	// 创建标识符
	ident := &ast.Identifier{
		StartPos: ast.Position{
			Line:   line,
			Column: column,
		},
		Name: name,
	}

	// 跳过标识符
	p.nextToken()

	logger.Debug("parseIdentifierOrCall: 返回标识符: %s", name)
	return ident

}

func (p *Parser) parseInteger() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	// 检查是否是浮点数
	if strings.Contains(p.curTok.Literal, ".") {
		// 这实际上是浮点数，调用 parseFloat
		return p.parseFloat()
	}

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

func (p *Parser) parseFloat() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	// 将字符串转换为float64
	value, err := strconv.ParseFloat(p.curTok.Literal, 64)
	if err != nil {
		p.errors = append(p.errors,
			fmt.Sprintf("第%d行第%d列: 无法解析浮点数: %s",
				p.curTok.Line, p.curTok.Column, p.curTok.Literal))
		return nil
	}

	expr := &ast.Float{
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

	expr := &ast.Boolean{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
		Value: p.curTok.Type == lexer.TokenTrue,
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

func (p *Parser) parseListLiteral() *ast.List {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	list := &ast.List{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
	}

	p.nextToken() // 跳过 [

	// 解析列表元素
	if !p.curTokenIs(lexer.TokenRBracket) {
		for {
			// 解析元素
			element := p.parseExpression()
			if element != nil {
				list.Elements = append(list.Elements, element)
			} else {
				// 元素解析失败
				p.errors = append(p.errors, "列表元素解析失败")
				return nil
			}

			// 检查是否有更多元素
			if p.curTokenIs(lexer.TokenComma) {
				p.nextToken() // 跳过逗号
				continue
			}

			break
		}
	}

	// 期望右括号
	if !p.expect(lexer.TokenRBracket, "列表字面量后") {
		return nil
	}

	return list
}

func (p *Parser) parseCallOrIndex() ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	logger.Debug("parseCallOrIndex: 开始")

	// 先解析基本表达式
	expr := p.parsePrimary()
	if expr == nil {
		return nil
	}

	logger.Debug("parseCallOrIndex: 基本表达式: %T(%v)", expr, expr)

	// 处理后续的调用、下标或链式调用
	return p.parsePostfix(expr)
}

func (p *Parser) parsePostfix(expr ast.Expression) ast.Expression {
	if !p.checkDepth() {
		return expr
	}

	p.enter()
	defer p.leave()

	logger.Debug("parsePostfix: 开始，expr=%T(%v), curTok=%v", expr, expr, p.curTok)

	// 循环处理多个下标、调用或链式调用
	for {
		switch p.curTok.Type {
		case lexer.TokenLParen:
			// 函数调用
			expr = p.parseCall(expr)
			logger.Debug("parsePostfix: 解析调用后，expr=%T(%v)", expr, expr)
		case lexer.TokenLBracket:
			// 下标表达式
			expr = p.parseIndex(expr)
		case lexer.TokenDot:
			// 检查是否是数字字面量的一部分
			switch expr.(type) {
			case *ast.Integer, *ast.Float:
				// 数字后面不能跟点号进行链式调用
				return expr
			default:
				// 这是链式调用的开始
				// 回退当前token，让parseChainCall重新处理点号
				// 但实际上我们已经在点号位置，所以直接解析链式调用
				chain := p.parseChainCall(expr)
				if chain != nil {
					expr = chain
					logger.Debug("parsePostfix: 解析链式调用后，expr=%T(%v)", expr, expr)

					// 检查是否还有后续的点号
					// 链式调用可能已经处理了多个点号，所以这里要检查
					continue
				} else {
					// 链式调用解析失败
					return expr
				}
			}
		default:
			// 既不是调用也不是下标也不是链式，返回当前表达式
			logger.Debug("parsePostfix: 返回表达式: %T(%v)", expr, expr)
			return expr
		}
	}
}

func (p *Parser) parseCall(left ast.Expression) ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	logger.Debug("parseCall: 开始解析调用，left=%T(%v)", left, left)

	// 检查左边是否是标识符
	ident, ok := left.(*ast.Identifier)
	if !ok {
		// 如果不是标识符，可能是其他表达式
		// 例如，在链式调用中，left 可能是其他类型的表达式
		logger.Debug("parseCall: 左边不是标识符: %T", left)

		// 对于非标识符，我们可能无法直接调用
		// 这里返回一个错误或尝试处理
		p.errors = append(p.errors, "函数调用必须是标识符")
		return nil
	}

	logger.Debug("parseCall: 函数名: %s", ident.Name)
	p.nextToken() // 跳过 (

	var args []ast.Expression

	// 解析参数列表
	if !p.curTokenIs(lexer.TokenRParen) {
		for {
			// 解析参数
			arg := p.parseExpression()
			logger.Debug("parseCall: 解析参数: %v", arg)
			if arg != nil {
				args = append(args, arg)
			} else {
				// 参数解析失败
				break
			}

			// 检查是否有更多参数
			if p.curTokenIs(lexer.TokenComma) {
				p.nextToken() // 跳过逗号
				continue
			}

			break
		}
	}

	// 期望右括号
	if !p.expect(lexer.TokenRParen, "函数调用参数列表后") {
		return nil
	}

	return &ast.CallExpr{
		StartPos: ast.Position{
			Line:   ident.StartPos.Line,
			Column: ident.StartPos.Column,
		},
		Function: ident,
		Args:     args,
	}
}

func (p *Parser) parseIndex(left ast.Expression) ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	logger.Debug("parseIndex: 开始解析下标，当前token:", p.curTok)

	// 跳过 [
	p.nextToken()
	logger.Debug("parseIndex: 跳过[后，当前token:", p.curTok)

	// 解析索引表达式
	index := p.parseExpression()
	if index == nil {
		p.errors = append(p.errors,
			fmt.Sprintf("第%d行第%d列: 下标表达式解析失败",
				p.curTok.Line, p.curTok.Column))

		// 尝试恢复：跳过直到遇到 ] 或文件结束
		for p.curTok.Type != lexer.TokenRBracket && p.curTok.Type != lexer.TokenEOF {
			p.nextToken()
		}

		// 跳过 ]
		if p.curTokenIs(lexer.TokenRBracket) {
			p.nextToken()
		}

		return nil
	}

	logger.Debug("parseIndex: 索引表达式解析成功:", index)

	// 期望 ]
	if !p.curTokenIs(lexer.TokenRBracket) {
		p.errors = append(p.errors,
			fmt.Sprintf("第%d行第%d列: 下标表达式后期望]，得到 %s",
				p.curTok.Line, p.curTok.Column, p.curTok.Type))

		// 尝试恢复：跳过直到遇到 ] 或文件结束
		for p.curTok.Type != lexer.TokenRBracket && p.curTok.Type != lexer.TokenEOF {
			p.nextToken()
		}

		// 跳过 ]
		if p.curTokenIs(lexer.TokenRBracket) {
			p.nextToken()
		}

		return nil
	}

	// 跳过 ]
	p.nextToken()

	return &ast.IndexExpr{
		StartPos: left.Pos(),
		Left:     left,
		Index:    index,
	}
}

func (p *Parser) parseDictLiteral() *ast.Dict {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	dict := &ast.Dict{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
		Pairs: make(map[ast.Expression]ast.Expression),
	}

	// 记录当前位置以便调试
	currentLine := p.curTok.Line
	currentColumn := p.curTok.Column

	// 跳过 {
	if !p.expect(lexer.TokenLBrace, "字典开始") {
		return nil
	}

	// 如果立即遇到 }，返回空字典
	if p.curTokenIs(lexer.TokenRBrace) {
		p.nextToken() // 跳过 }
		return dict
	}

	// 解析字典键值对
	for {
		// 调试：记录当前位置
		logger.Debug("解析字典键，当前位置: 行=%d, 列=%d, token=%v, literal=%q",
			p.curTok.Line, p.curTok.Column, p.curTok.Type, p.curTok.Literal)

		// 解析键
		key := p.parseExpression()
		if key == nil {
			p.errors = append(p.errors,
				fmt.Sprintf("第%d行第%d列: 字典键解析失败",
					p.curTok.Line, p.curTok.Column))
			return nil
		}

		// 检查键类型：不允许布尔值作为键
		switch key.String() {
		case "true:", "false:":
			p.errors = append(p.errors,
				fmt.Sprintf("第%d行第%d列: 布尔值不能作为字典键",
					p.curTok.Line, p.curTok.Column))
			return nil
		}

		logger.Debug("键解析成功: %v", key)

		// 期望冒号
		if !p.expect(lexer.TokenColon, "字典键后期望冒号") {
			p.errors = append(p.errors,
				fmt.Sprintf("第%d行第%d列: 字典键后期望冒号，得到 %s",
					currentLine, currentColumn, p.curTok.Type))
			return nil
		}

		// 解析值
		value := p.parseExpression()
		if value == nil {
			p.errors = append(p.errors,
				fmt.Sprintf("第%d行第%d列: 字典值解析失败",
					p.curTok.Line, p.curTok.Column))
			return nil
		}

		// 将键值对添加到字典
		dict.Pairs[key] = value

		// 检查是否有更多键值对
		if !p.curTokenIs(lexer.TokenComma) {
			break
		}

		// 跳过逗号
		p.nextToken()

		// 如果逗号后立即遇到 }，这是语法错误
		if p.curTokenIs(lexer.TokenRBrace) {
			p.errors = append(p.errors,
				fmt.Sprintf("第%d行第%d列: 字典中多余的逗号",
					p.curTok.Line, p.curTok.Column))
			return nil
		}
	}

	// 期望右花括号
	if !p.expect(lexer.TokenRBrace, "字典字面量后") {
		return nil
	}

	return dict
}

func (p *Parser) parseChainCall(left ast.Expression) ast.Expression {
	if !p.checkDepth() {
		return nil
	}

	p.enter()
	defer p.leave()

	logger.Debug("========= parseChainCall 开始 =========")
	logger.Debug("左边表达式: %T(%v)", left, left)
	logger.Debug("当前token: %v", p.curTok)

	// 创建一个链式调用节点
	chain := &ast.ChainCallExpr{
		StartPos: left.Pos(),
		Calls:    make([]*ast.CallExpr, 0),
	}

	// 处理左边表达式
	if existingChain, ok := left.(*ast.ChainCallExpr); ok {
		// 如果左边已经是链式调用，复用它的调用列表
		chain.Calls = existingChain.Calls
		logger.Debug("parseChainCall: 复用已有链式调用，包含 %d 个调用", len(chain.Calls))
	} else {
		// 左边不是链式调用
		// 创建第一个调用（可能是虚拟调用）
		firstCall := p.createFirstChainCall(left)
		if firstCall == nil {
			logger.Debug("parseChainCall: 无法创建第一个调用")
			return nil
		}
		chain.Calls = append(chain.Calls, firstCall)
		logger.Debug("parseChainCall: 添加第一个调用: %s", firstCall.Function.Name)
	}

	// 当前token应该是点号（调用者应该已经确保）
	if !p.curTokenIs(lexer.TokenDot) {
		logger.Debug("parseChainCall: 当前token不是点号: %v", p.curTok)
		// 但可能我们只有一个调用
		if len(chain.Calls) == 1 {
			// 如果是虚拟调用，返回原始表达式
			if chain.Calls[0].Function.Name == "_value" {
				return chain.Calls[0].Args[0]
			}
			return chain.Calls[0]
		}
	}

	// 解析链式调用的每个部分
	for p.curTokenIs(lexer.TokenDot) {
		logger.Debug("parseChainCall: 处理点号")

		// 跳过当前点号
		p.nextToken()
		logger.Debug("parseChainCall: 跳过了点号，当前token: %v", p.curTok)

		// 解析方法名
		if !p.curTokenIs(lexer.TokenIdent) {
			p.errors = append(p.errors, "期望方法名")
			return nil
		}

		methodName := p.curTok.Literal
		logger.Debug("parseChainCall: 方法名: %s", methodName)

		// 创建方法标识符
		methodIdent := &ast.Identifier{
			StartPos: ast.Position{
				Line:   p.curTok.Line,
				Column: p.curTok.Column,
			},
			Name: methodName,
		}

		// 跳过方法名
		p.nextToken()

		// 检查是否有括号
		if !p.curTokenIs(lexer.TokenLParen) {
			p.errors = append(p.errors, "期望 '(' 开始方法调用")
			return nil
		}

		// 解析方法调用
		methodCall := p.parseCall(methodIdent)
		if methodCall == nil {
			logger.Debug("parseChainCall: 方法调用解析失败")
			return nil
		}

		// 添加到链中
		if callExpr, ok := methodCall.(*ast.CallExpr); ok {
			chain.Calls = append(chain.Calls, callExpr)
			logger.Debug("parseChainCall: 添加方法调用: %s", callExpr.Function.Name)
		} else {
			p.errors = append(p.errors, "链式调用中的元素必须是函数调用")
			return nil
		}

		logger.Debug("parseChainCall: 当前链有 %d 个调用", len(chain.Calls))
	}

	// 如果只有一个调用，直接返回它
	if len(chain.Calls) == 1 {
		logger.Debug("parseChainCall: 只有一个调用，直接返回")
		// 如果是虚拟调用，返回原始表达式
		if chain.Calls[0].Function.Name == "_value" {
			return chain.Calls[0].Args[0]
		}
		return chain.Calls[0]
	}

	logger.Debug("parseChainCall: 返回链式调用，包含 %d 个调用", len(chain.Calls))
	return chain
}

// 创建链式调用的第一个调用
func (p *Parser) createFirstChainCall(expr ast.Expression) *ast.CallExpr {
	logger.Debug("createFirstChainCall: expr=%T(%v)", expr, expr)

	switch v := expr.(type) {
	case *ast.CallExpr:
		// 已经是函数调用
		return v
	case *ast.Identifier:
		// 标识符（变量），创建虚拟调用
		return &ast.CallExpr{
			StartPos: v.Pos(),
			Function: &ast.Identifier{
				StartPos: v.Pos(),
				Name:     "_value",
			},
			Args: []ast.Expression{v},
		}
	case *ast.String, *ast.Integer, *ast.Float, *ast.Boolean:
		// 字面量，创建虚拟调用
		return &ast.CallExpr{
			StartPos: v.Pos(),
			Function: &ast.Identifier{
				StartPos: v.Pos(),
				Name:     "_value",
			},
			Args: []ast.Expression{v},
		}
	default:
		// 其他表达式，也创建虚拟调用
		return &ast.CallExpr{
			StartPos: v.Pos(),
			Function: &ast.Identifier{
				StartPos: v.Pos(),
				Name:     "_value",
			},
			Args: []ast.Expression{v},
		}
	}
}

// 解析方法调用
func (p *Parser) parseMethodCall(methodName string) ast.Expression {
	// 期望当前位置是左括号
	if !p.curTokenIs(lexer.TokenLParen) {
		p.errors = append(p.errors, "期望 '(' 开始方法调用")
		return nil
	}

	// 创建方法名标识符
	methodIdent := &ast.Identifier{
		StartPos: ast.Position{
			Line:   p.curTok.Line,
			Column: p.curTok.Column,
		},
		Name: methodName,
	}

	// 解析调用
	return p.parseCall(methodIdent)
}

// 将表达式转换为调用表达式
func (p *Parser) convertToCallExpr(expr ast.Expression) *ast.CallExpr {
	// 对于简单情况，如果表达式是标识符，但还没有被处理
	// 这里需要根据你的语法进行调整
	return nil
}
