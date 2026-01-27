// parser/parser.go
package parser

import (
	"fmt"
	"log"
	"my-dsl/ast"
	"my-dsl/lexer"
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
