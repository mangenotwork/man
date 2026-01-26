package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// -------------------------- 通用定义 --------------------------
// Position 原始代码位置（行/列从1开始）
type Position struct {
	Line int // 行号
	Col  int // 列号
}

func (p Position) String() string {
	return fmt.Sprintf("line %d, col %d", p.Line, p.Col)
}

// ParseError 解析错误
type ParseError struct {
	Pos   Position // 错误位置
	Msg   string   // 错误信息
	Level string   // 错误等级：lexer/parser/semantic
}

func (e ParseError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Level, e.Pos, e.Msg)
}

// TokenType Token类型
type TokenType int

const (
	TokenUnknown TokenType = iota
	TokenKeyword           // let/if/elseif/else/for/print
	TokenIdent             // 标识符
	TokenString            // 字符串字面量
	TokenNumber            // 数字字面量
	TokenBool              // 布尔字面量
	TokenAssign            // =
	TokenOp                // >/</>=/<=/!=/++/--
	TokenSemi              // ;
	TokenLBrace            // {
	TokenRBrace            // }
	TokenEOF               // 结束符
)

// Token 词法单元
type Token struct {
	Type  TokenType
	Value string
	Pos   Position
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q)", tokenTypeToString(t.Type), t.Value)
}

func tokenTypeToString(t TokenType) string {
	switch t {
	case TokenUnknown:
		return "Unknown"
	case TokenKeyword:
		return "Keyword"
	case TokenIdent:
		return "Ident"
	case TokenString:
		return "String"
	case TokenNumber:
		return "Number"
	case TokenBool:
		return "Bool"
	case TokenAssign:
		return "Assign"
	case TokenOp:
		return "Op"
	case TokenSemi:
		return "Semi"
	case TokenLBrace:
		return "LBrace"
	case TokenRBrace:
		return "RBrace"
	case TokenEOF:
		return "EOF"
	default:
		return "Unnamed"
	}
}

// -------------------------- 词法分析器 --------------------------
var (
	keywords     = map[string]struct{}{"let": {}, "if": {}, "elseif": {}, "else": {}, "for": {}, "print": {}}
	boolLiterals = map[string]struct{}{"true": {}, "false": {}}
	ops          = map[string]struct{}{">": {}, "<": {}, "!=": {}, ">=": {}, "<=": {}, "++": {}, "--": {}}
	identRegex   = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
)

// Lexer 词法分析器
type Lexer struct {
	input  string
	lines  []string     // 按行存储输入，用于精准定位
	pos    Position     // 当前位置
	rdPos  int          // 全局字符索引
	errors []ParseError // 词法错误
}

func NewLexer(input string) *Lexer {
	lines := strings.Split(input, "\n")
	// 确保每行都有内容（处理空行）
	for i, line := range lines {
		if line == "" {
			lines[i] = " "
		}
	}
	return &Lexer{
		input: input,
		lines: lines,
		pos:   Position{Line: 1, Col: 1},
	}
}

// NextToken 获取下一个Token
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.rdPos >= len(l.input) {
		return Token{Type: TokenEOF, Pos: l.pos}
	}

	c := l.input[l.rdPos]
	currentPos := l.pos

	// 字符串字面量
	if c == '"' {
		return l.scanString()
	}

	// 数字
	if unicode.IsDigit(rune(c)) {
		return l.scanNumber()
	}

	// 字母/下划线（标识符/关键字/布尔值）
	if unicode.IsLetter(rune(c)) || c == '_' {
		return l.scanIdentOrKeyword()
	}

	// 运算符/分隔符
	switch c {
	case '=':
		l.advance()
		return Token{Type: TokenAssign, Value: "=", Pos: currentPos}
	case ';':
		l.advance()
		return Token{Type: TokenSemi, Value: ";", Pos: currentPos}
	case '{':
		l.advance()
		return Token{Type: TokenLBrace, Value: "{", Pos: currentPos}
	case '}':
		l.advance()
		return Token{Type: TokenRBrace, Value: "}", Pos: currentPos}
	case '>':
		l.advance()
		if l.rdPos < len(l.input) && l.input[l.rdPos] == '=' {
			l.advance()
			return Token{Type: TokenOp, Value: ">=", Pos: currentPos}
		}
		return Token{Type: TokenOp, Value: ">", Pos: currentPos}
	case '<':
		l.advance()
		if l.rdPos < len(l.input) && l.input[l.rdPos] == '=' {
			l.advance()
			return Token{Type: TokenOp, Value: "<=", Pos: currentPos}
		}
		return Token{Type: TokenOp, Value: "<", Pos: currentPos}
	case '!':
		l.advance()
		if l.rdPos < len(l.input) && l.input[l.rdPos] == '=' {
			l.advance()
			return Token{Type: TokenOp, Value: "!=", Pos: currentPos}
		}
		l.addError(currentPos, "无效的运算符: '!' (仅支持 '!=')")
		return Token{Type: TokenUnknown, Value: "!", Pos: currentPos}
	case '+':
		l.advance()
		if l.rdPos < len(l.input) && l.input[l.rdPos] == '+' {
			l.advance()
			return Token{Type: TokenOp, Value: "++", Pos: currentPos}
		}
		l.addError(currentPos, "无效的运算符: '+' (仅支持 '++')")
		return Token{Type: TokenUnknown, Value: "+", Pos: currentPos}
	case '-':
		l.advance()
		if l.rdPos < len(l.input) && l.input[l.rdPos] == '-' {
			l.advance()
			return Token{Type: TokenOp, Value: "--", Pos: currentPos}
		}
		l.addError(currentPos, "无效的运算符: '-' (仅支持 '--')")
		return Token{Type: TokenUnknown, Value: "-", Pos: currentPos}
	default:
		l.addError(currentPos, fmt.Sprintf("未知字符: '%c'", c))
		l.advance()
		return Token{Type: TokenUnknown, Value: string(c), Pos: currentPos}
	}
}

// scanString 扫描字符串字面量
func (l *Lexer) scanString() Token {
	startPos := l.pos
	l.advance() // 跳过开头的"

	var sb strings.Builder
	for l.rdPos < len(l.input) {
		c := l.input[l.rdPos]
		if c == '"' {
			l.advance() // 跳过结尾的"
			return Token{Type: TokenString, Value: sb.String(), Pos: startPos}
		}
		if c == '\n' { // 换行符，字符串未闭合
			break
		}
		sb.WriteByte(c)
		l.advance()
	}

	// 字符串未闭合
	l.addError(startPos, "字符串字面量未闭合")
	return Token{Type: TokenString, Value: sb.String(), Pos: startPos}
}

// scanNumber 扫描数字
func (l *Lexer) scanNumber() Token {
	startPos := l.pos
	var sb strings.Builder
	hasDot := false

	for l.rdPos < len(l.input) {
		c := l.input[l.rdPos]
		if unicode.IsDigit(rune(c)) {
			sb.WriteByte(c)
			l.advance()
		} else if c == '.' && !hasDot {
			hasDot = true
			sb.WriteByte(c)
			l.advance()
		} else {
			break
		}
	}

	return Token{Type: TokenNumber, Value: sb.String(), Pos: startPos}
}

// scanIdentOrKeyword 扫描标识符/关键字/布尔值
func (l *Lexer) scanIdentOrKeyword() Token {
	startPos := l.pos
	var sb strings.Builder

	for l.rdPos < len(l.input) {
		c := l.input[l.rdPos]
		if unicode.IsLetter(rune(c)) || unicode.IsDigit(rune(c)) || c == '_' {
			sb.WriteByte(c)
			l.advance()
		} else {
			break
		}
	}

	val := sb.String()

	// 检查是否是关键字
	if _, ok := keywords[val]; ok {
		return Token{Type: TokenKeyword, Value: val, Pos: startPos}
	}

	// 检查是否是布尔值
	if _, ok := boolLiterals[val]; ok {
		return Token{Type: TokenBool, Value: val, Pos: startPos}
	}

	// 检查是否是合法标识符
	if !identRegex.MatchString(val) {
		l.addError(startPos, fmt.Sprintf("无效的标识符: '%s'", val))
		return Token{Type: TokenUnknown, Value: val, Pos: startPos}
	}

	return Token{Type: TokenIdent, Value: val, Pos: startPos}
}

// skipWhitespace 跳过空白符和注释
func (l *Lexer) skipWhitespace() {
	for l.rdPos < len(l.input) {
		c := l.input[l.rdPos]
		if c == ' ' || c == '\t' || c == '\r' {
			l.pos.Col++
			l.rdPos++
		} else if c == '\n' {
			l.pos.Line++
			l.pos.Col = 1
			l.rdPos++
		} else if c == '#' {
			// 跳过注释（到行尾）
			for l.rdPos < len(l.input) && l.input[l.rdPos] != '\n' {
				l.rdPos++
				l.pos.Col++
			}
		} else {
			break
		}
	}
}

// advance 前进一个字符，更新位置
func (l *Lexer) advance() {
	if l.rdPos >= len(l.input) {
		return
	}
	if l.input[l.rdPos] == '\n' {
		l.pos.Line++
		l.pos.Col = 1
	} else {
		l.pos.Col++
	}
	l.rdPos++
}

// addError 添加词法错误
func (l *Lexer) addError(pos Position, msg string) {
	l.errors = append(l.errors, ParseError{
		Pos:   pos,
		Msg:   msg,
		Level: "lexer",
	})
}

// Errors 获取词法错误
func (l *Lexer) Errors() []ParseError {
	return l.errors
}

// -------------------------- AST定义 --------------------------
type ASTNode interface {
	Position() Position
}

// Program 根节点
type Program struct {
	Stmts []ASTNode
}

func (p *Program) Position() Position {
	if len(p.Stmts) == 0 {
		return Position{Line: 1, Col: 1}
	}
	return p.Stmts[0].Position()
}

// LetStmt let语句: let ident = expr
type LetStmt struct {
	IdentPos Position
	Name     string
	Value    Expr
}

func (l *LetStmt) Position() Position {
	return l.IdentPos
}

// PrintStmt print语句: print expr
type PrintStmt struct {
	StmtPos Position
	Arg     Expr
}

func (p *PrintStmt) Position() Position {
	return p.StmtPos
}

// IfStmt if语句: if expr { ... } elseif expr { ... } else { ... }
type IfStmt struct {
	IfPos   Position
	Cond    Expr
	Body    *BlockStmt
	ElseIfs []*IfStmt
	Else    *BlockStmt
}

func (i *IfStmt) Position() Position {
	return i.IfPos
}

// ForStmt for语句: for expr; expr; expr { ... }
type ForStmt struct {
	ForPos Position
	Init   Expr
	Cond   Expr
	Post   Expr
	Body   *BlockStmt
}

func (f *ForStmt) Position() Position {
	return f.ForPos
}

// BlockStmt 块语句: { ... }
type BlockStmt struct {
	LBracePos Position
	Stmts     []ASTNode
}

func (b *BlockStmt) Position() Position {
	return b.LBracePos
}

// Expr 表达式接口
type Expr interface {
	ASTNode
	exprNode()
}

// LiteralExpr 字面量表达式
type LiteralExpr struct {
	ExprPos Position
	Type    string // string/number/bool
	Value   interface{}
}

func (l *LiteralExpr) Position() Position { return l.ExprPos }
func (l *LiteralExpr) exprNode()          {}

// IdentExpr 标识符表达式
type IdentExpr struct {
	ExprPos Position
	Name    string
}

func (i *IdentExpr) Position() Position { return i.ExprPos }
func (i *IdentExpr) exprNode()          {}

// AssignExpr 赋值表达式: ident = expr
type AssignExpr struct {
	ExprPos Position
	Left    *IdentExpr
	Op      string
	Right   Expr
}

func (a *AssignExpr) Position() Position { return a.ExprPos }
func (a *AssignExpr) exprNode()          {}

// OpExpr 运算表达式: expr op expr 或 expr++/expr--
type OpExpr struct {
	ExprPos Position
	Left    Expr
	Op      string
	Right   Expr
}

func (o *OpExpr) Position() Position { return o.ExprPos }
func (o *OpExpr) exprNode()          {}

// -------------------------- 语法分析器 --------------------------
type Parser struct {
	lexer   *Lexer
	curTok  Token
	peekTok Token
	errors  []ParseError
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken()
	p.nextToken()
	return p
}

// Parse 解析整个程序
func (p *Parser) Parse() *Program {
	program := &Program{}

	for p.curTok.Type != TokenEOF {
		stmt := p.parseStmt()
		if stmt != nil {
			program.Stmts = append(program.Stmts, stmt)
		}
		// 错误恢复：如果当前token不是语句开头，继续前进
		for p.curTok.Type != TokenEOF && !p.isStmtStart(p.curTok) {
			p.nextToken()
		}
	}

	return program
}

// parseStmt 解析单个语句
func (p *Parser) parseStmt() ASTNode {
	if p.curTok.Type == TokenEOF {
		return nil
	}

	switch p.curTok.Type {
	case TokenKeyword:
		switch p.curTok.Value {
		case "let":
			return p.parseLetStmt()
		case "print":
			return p.parsePrintStmt()
		case "if":
			return p.parseIfStmt()
		case "for":
			return p.parseForStmt()
		case "elseif", "else":
			// 单独出现的elseif/else是语法错误
			p.addError(p.curTok.Pos, fmt.Sprintf("关键字 '%s' 不能单独出现", p.curTok.Value))
			p.nextToken()
			return nil
		default:
			p.addError(p.curTok.Pos, fmt.Sprintf("不支持的关键字: '%s'", p.curTok.Value))
			p.nextToken()
			return nil
		}
	case TokenLBrace:
		return p.parseBlockStmt()
	default:
		p.addError(p.curTok.Pos, fmt.Sprintf("无效的语句开头: %s", p.curTok))
		p.nextToken()
		return nil
	}
}

// parseLetStmt 解析let语句
func (p *Parser) parseLetStmt() *LetStmt {
	letPos := p.curTok.Pos
	p.nextToken() // 跳过let

	// 检查变量名
	if p.curTok.Type != TokenIdent {
		p.addError(letPos, "let关键字后必须跟合法的变量名")
		// 错误恢复：跳过直到下一个语句
		p.skipUntilNextStmt()
		return nil
	}

	ident := &IdentExpr{
		ExprPos: p.curTok.Pos,
		Name:    p.curTok.Value,
	}
	p.nextToken() // 跳过敏识符

	// 检查赋值符
	if p.curTok.Type != TokenAssign {
		p.addError(p.curTok.Pos, "变量名后必须跟赋值符 '='")
		// 错误恢复：跳过直到下一个语句
		p.skipUntilNextStmt()
		return nil
	}
	p.nextToken() // 跳过=

	// 解析赋值表达式
	value := p.parseExpr()
	if value == nil {
		p.addError(p.curTok.Pos, "赋值符后必须跟合法的表达式")
		// 错误恢复：跳过直到下一个语句
		p.skipUntilNextStmt()
		return nil
	}

	return &LetStmt{
		IdentPos: letPos,
		Name:     ident.Name,
		Value:    value,
	}
}

// parsePrintStmt 解析print语句
func (p *Parser) parsePrintStmt() *PrintStmt {
	printPos := p.curTok.Pos
	p.nextToken() // 跳过print

	// 解析打印参数
	arg := p.parseExpr()
	if arg == nil {
		p.addError(p.curTok.Pos, "print关键字后必须跟合法的表达式")
		p.skipUntilNextStmt()
		return nil
	}

	return &PrintStmt{
		StmtPos: printPos,
		Arg:     arg,
	}
}

// parseIfStmt 解析if语句
func (p *Parser) parseIfStmt() *IfStmt {
	ifPos := p.curTok.Pos
	p.nextToken() // 跳过if

	// 解析条件表达式
	cond := p.parseExpr()
	if cond == nil {
		p.addError(p.curTok.Pos, "if关键字后必须跟合法的条件表达式")
		// 错误恢复：跳过直到块语句
		p.skipUntilBlock()
		return nil
	}

	// 解析if体
	body := p.parseBlockStmt()
	if body == nil {
		p.addError(p.curTok.Pos, "条件表达式后必须跟块语句 '{'")
		// 错误恢复：跳过直到elseif/else/下一个语句
		p.skipUntilElseOrEOF()
		return nil
	}

	ifStmt := &IfStmt{
		IfPos:   ifPos,
		Cond:    cond,
		Body:    body,
		ElseIfs: []*IfStmt{},
	}

	// 解析elseif分支
	for p.curTok.Type == TokenKeyword && p.curTok.Value == "elseif" {
		elseifPos := p.curTok.Pos
		p.nextToken() // 跳过elseif

		// 解析elseif条件
		elseifCond := p.parseExpr()
		if elseifCond == nil {
			p.addError(elseifPos, "elseif关键字后必须跟合法的条件表达式")
			// 错误恢复：跳过直到块语句
			p.skipUntilBlock()
			continue
		}

		// 解析elseif体
		elseifBody := p.parseBlockStmt()
		if elseifBody == nil {
			p.addError(p.curTok.Pos, "elseif条件表达式后必须跟块语句 '{'")
			// 错误恢复：跳过直到elseif/else/下一个语句
			p.skipUntilElseOrEOF()
			continue
		}

		ifStmt.ElseIfs = append(ifStmt.ElseIfs, &IfStmt{
			IfPos: elseifPos,
			Cond:  elseifCond,
			Body:  elseifBody,
		})
	}

	// 解析else分支
	if p.curTok.Type == TokenKeyword && p.curTok.Value == "else" {
		p.nextToken() // 跳过else

		// 解析else体
		elseBody := p.parseBlockStmt()
		if elseBody == nil {
			p.addError(p.curTok.Pos, "else关键字后必须跟块语句 '{'")
			return ifStmt
		}

		ifStmt.Else = elseBody
	}

	return ifStmt
}

// parseForStmt 解析for语句
func (p *Parser) parseForStmt() *ForStmt {
	forPos := p.curTok.Pos
	p.nextToken() // 跳过for

	// 解析初始化表达式
	init := p.parseExpr()
	if init == nil {
		p.addError(forPos, "for关键字后必须跟初始化表达式")
		p.skipUntilSemiOrBrace()
		return nil
	}

	// 检查第一个分号
	if p.curTok.Type != TokenSemi {
		p.addError(p.curTok.Pos, "初始化表达式后必须跟分号 ';'")
		p.skipUntilSemiOrBrace()
		return nil
	}
	p.nextToken() // 跳过;

	// 解析条件表达式
	cond := p.parseExpr()
	if cond == nil {
		p.addError(p.curTok.Pos, "第一个分号后必须跟条件表达式")
		p.skipUntilSemiOrBrace()
		return nil
	}

	// 检查第二个分号
	if p.curTok.Type != TokenSemi {
		p.addError(p.curTok.Pos, "条件表达式后必须跟分号 ';'")
		p.skipUntilSemiOrBrace()
		return nil
	}
	p.nextToken() // 跳过;

	// 解析增量表达式
	post := p.parseExpr()
	if post == nil {
		p.addError(p.curTok.Pos, "第二个分号后必须跟增量表达式")
		p.skipUntilBlock()
		return nil
	}

	// 解析for体
	body := p.parseBlockStmt()
	if body == nil {
		p.addError(p.curTok.Pos, "增量表达式后必须跟块语句 '{'")
		return nil
	}

	return &ForStmt{
		ForPos: forPos,
		Init:   init,
		Cond:   cond,
		Post:   post,
		Body:   body,
	}
}

// parseBlockStmt 解析块语句
func (p *Parser) parseBlockStmt() *BlockStmt {
	if p.curTok.Type != TokenLBrace {
		return nil
	}

	lbPos := p.curTok.Pos
	p.nextToken() // 跳过{

	block := &BlockStmt{
		LBracePos: lbPos,
		Stmts:     []ASTNode{},
	}

	// 解析块内语句
	for p.curTok.Type != TokenRBrace && p.curTok.Type != TokenEOF {
		stmt := p.parseStmt()
		if stmt != nil {
			block.Stmts = append(block.Stmts, stmt)
		}
	}

	// 检查右大括号
	if p.curTok.Type != TokenRBrace {
		p.addError(lbPos, "块语句未找到闭合的 '}'")
		return block
	}

	p.nextToken() // 跳过}
	return block
}

// parseExpr 解析表达式
func (p *Parser) parseExpr() Expr {
	return p.parseOpExpr()
}

// parseOpExpr 解析运算表达式
func (p *Parser) parseOpExpr() Expr {
	expr := p.parsePrimaryExpr()
	if expr == nil {
		return nil
	}

	// 处理二元运算符
	for p.curTok.Type == TokenOp && p.isBinaryOp(p.curTok.Value) {
		opPos := p.curTok.Pos
		op := p.curTok.Value
		p.nextToken()

		right := p.parsePrimaryExpr()
		if right == nil {
			p.addError(opPos, fmt.Sprintf("运算符 '%s' 后必须跟合法的表达式", op))
			return expr
		}

		expr = &OpExpr{
			ExprPos: opPos,
			Left:    expr,
			Op:      op,
			Right:   right,
		}
	}

	// 处理一元后缀运算符 (++/--)
	if p.curTok.Type == TokenOp && (p.curTok.Value == "++" || p.curTok.Value == "--") {
		opPos := p.curTok.Pos
		op := p.curTok.Value
		p.nextToken()

		expr = &OpExpr{
			ExprPos: opPos,
			Left:    expr,
			Op:      op,
			Right:   nil,
		}
	}

	// 处理赋值运算符
	if p.curTok.Type == TokenAssign {
		opPos := p.curTok.Pos
		op := p.curTok.Value
		p.nextToken()

		// 赋值左值必须是标识符
		ident, ok := expr.(*IdentExpr)
		if !ok {
			p.addError(opPos, "赋值运算符左边必须是标识符")
			return expr
		}

		right := p.parsePrimaryExpr()
		if right == nil {
			p.addError(opPos, "赋值运算符后必须跟合法的表达式")
			return expr
		}

		expr = &AssignExpr{
			ExprPos: opPos,
			Left:    ident,
			Op:      op,
			Right:   right,
		}
	}

	return expr
}

// parsePrimaryExpr 解析基础表达式
func (p *Parser) parsePrimaryExpr() Expr {
	switch p.curTok.Type {
	case TokenIdent:
		expr := &IdentExpr{
			ExprPos: p.curTok.Pos,
			Name:    p.curTok.Value,
		}
		p.nextToken()
		return expr
	case TokenString:
		expr := &LiteralExpr{
			ExprPos: p.curTok.Pos,
			Type:    "string",
			Value:   p.curTok.Value,
		}
		p.nextToken()
		return expr
	case TokenNumber:
		var val interface{}
		var typ string

		if intVal, err := strconv.Atoi(p.curTok.Value); err == nil {
			val = intVal
			typ = "number"
		} else if floatVal, err := strconv.ParseFloat(p.curTok.Value, 64); err == nil {
			val = floatVal
			typ = "number"
		} else {
			p.addError(p.curTok.Pos, fmt.Sprintf("无效的数字格式: '%s'", p.curTok.Value))
			p.nextToken()
			return nil
		}

		expr := &LiteralExpr{
			ExprPos: p.curTok.Pos,
			Type:    typ,
			Value:   val,
		}
		p.nextToken()
		return expr
	case TokenBool:
		val, _ := strconv.ParseBool(p.curTok.Value)
		expr := &LiteralExpr{
			ExprPos: p.curTok.Pos,
			Type:    "bool",
			Value:   val,
		}
		p.nextToken()
		return expr
	default:
		return nil
	}
}

// 辅助函数
func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.lexer.NextToken()
}

func (p *Parser) addError(pos Position, msg string) {
	p.errors = append(p.errors, ParseError{
		Pos:   pos,
		Msg:   msg,
		Level: "parser",
	})
}

// isStmtStart 判断token是否是语句开头
func (p *Parser) isStmtStart(tok Token) bool {
	if tok.Type == TokenKeyword {
		switch tok.Value {
		case "let", "print", "if", "for":
			return true
		}
	}
	return tok.Type == TokenLBrace
}

// skipUntilNextStmt 跳过直到下一个语句开头
func (p *Parser) skipUntilNextStmt() {
	for p.curTok.Type != TokenEOF && !p.isStmtStart(p.curTok) {
		p.nextToken()
	}
}

// skipUntilBlock 跳过直到块语句
func (p *Parser) skipUntilBlock() {
	for p.curTok.Type != TokenEOF && p.curTok.Type != TokenLBrace {
		p.nextToken()
	}
}

// skipUntilElseOrEOF 跳过直到elseif/else/EOF
func (p *Parser) skipUntilElseOrEOF() {
	for p.curTok.Type != TokenEOF {
		if p.curTok.Type == TokenKeyword && (p.curTok.Value == "elseif" || p.curTok.Value == "else") {
			break
		}
		if p.isStmtStart(p.curTok) {
			break
		}
		p.nextToken()
	}
}

// skipUntilSemiOrBrace 跳过直到分号/块语句
func (p *Parser) skipUntilSemiOrBrace() {
	for p.curTok.Type != TokenEOF && p.curTok.Type != TokenSemi && p.curTok.Type != TokenLBrace {
		p.nextToken()
	}
}

// isBinaryOp 判断是否是二元运算符
func (p *Parser) isBinaryOp(op string) bool {
	switch op {
	case ">", "<", ">=", "<=", "!=":
		return true
	default:
		return false
	}
}

// Errors 获取语法错误
func (p *Parser) Errors() []ParseError {
	return p.errors
}

// -------------------------- 语义分析器 --------------------------
type VarInfo struct {
	Type  string
	Value interface{}
	Pos   Position
}

type Checker struct {
	prog   *Program
	scope  map[string]VarInfo
	errors []ParseError
}

func NewChecker(prog *Program) *Checker {
	if prog == nil {
		prog = &Program{}
	}
	return &Checker{
		prog:   prog,
		scope:  make(map[string]VarInfo),
		errors: []ParseError{},
	}
}

// Check 执行语义检查
func (c *Checker) Check() {
	for _, stmt := range c.prog.Stmts {
		if stmt == nil {
			continue
		}
		c.checkStmt(stmt)
	}
}

// checkStmt 检查单个语句
func (c *Checker) checkStmt(stmt ASTNode) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *LetStmt:
		c.checkLetStmt(s)
	case *PrintStmt:
		c.checkPrintStmt(s)
	case *IfStmt:
		c.checkIfStmt(s)
	case *ForStmt:
		c.checkForStmt(s)
	case *BlockStmt:
		c.checkBlockStmt(s)
	default:
		c.addError(stmt.Position(), fmt.Sprintf("不支持的语句类型: %T", s))
	}
}

// checkLetStmt 检查let语句
func (c *Checker) checkLetStmt(s *LetStmt) {
	if s == nil {
		return
	}

	// 检查变量是否已定义
	if _, exists := c.scope[s.Name]; exists {
		c.addError(s.Position(), fmt.Sprintf("变量 '%s' 已重复定义", s.Name))
		return
	}

	// 检查赋值表达式
	if s.Value == nil {
		c.addError(s.Position(), "变量赋值表达式不能为空")
		return
	}

	// 检查表达式类型
	valType, val := c.checkExpr(s.Value)
	if valType == "" {
		c.addError(s.Value.Position(), "变量赋值表达式类型不合法")
		return
	}

	// 注册变量
	c.scope[s.Name] = VarInfo{
		Type:  valType,
		Value: val,
		Pos:   s.Position(),
	}
}

// checkPrintStmt 检查print语句
func (c *Checker) checkPrintStmt(s *PrintStmt) {
	if s == nil {
		return
	}

	// 检查打印参数
	if s.Arg == nil {
		c.addError(s.Position(), "print语句参数不能为空")
		return
	}

	// 检查参数类型
	_, _ = c.checkExpr(s.Arg)
}

// checkIfStmt 检查if语句
func (c *Checker) checkIfStmt(s *IfStmt) {
	if s == nil {
		return
	}

	// 检查条件表达式
	if s.Cond == nil {
		c.addError(s.Position(), "if条件表达式不能为空")
		return
	}

	// 检查条件表达式类型（必须是布尔型）
	condType, _ := c.checkExpr(s.Cond)
	if condType != "bool" {
		c.addError(s.Cond.Position(), fmt.Sprintf("if条件表达式必须是布尔类型，当前为 '%s' 类型", condType))
	}

	// 检查if体
	if s.Body != nil {
		c.checkBlockStmt(s.Body)
	}

	// 检查elseif分支
	for _, elseif := range s.ElseIfs {
		if elseif == nil {
			continue
		}

		// 检查elseif条件
		if elseif.Cond == nil {
			c.addError(elseif.Position(), "elseif条件表达式不能为空")
			continue
		}

		// 检查elseif条件类型
		elseifType, _ := c.checkExpr(elseif.Cond)
		if elseifType != "bool" {
			c.addError(elseif.Cond.Position(), fmt.Sprintf("elseif条件表达式必须是布尔类型，当前为 '%s' 类型", elseifType))
		}

		// 检查elseif体
		if elseif.Body != nil {
			c.checkBlockStmt(elseif.Body)
		}
	}

	// 检查else分支
	if s.Else != nil {
		c.checkBlockStmt(s.Else)
	}
}

// checkForStmt 检查for语句
func (c *Checker) checkForStmt(s *ForStmt) {
	if s == nil {
		return
	}

	// 检查初始化表达式
	if s.Init != nil {
		_, _ = c.checkExpr(s.Init)
	} else {
		c.addError(s.Position(), "for初始化表达式不能为空")
	}

	// 检查条件表达式
	if s.Cond != nil {
		condType, _ := c.checkExpr(s.Cond)
		if condType != "bool" {
			c.addError(s.Cond.Position(), fmt.Sprintf("for条件表达式必须是布尔类型，当前为 '%s' 类型", condType))
		}
	} else {
		c.addError(s.Position(), "for条件表达式不能为空")
	}

	// 检查增量表达式
	if s.Post != nil {
		_, _ = c.checkExpr(s.Post)
	} else {
		c.addError(s.Position(), "for增量表达式不能为空")
	}

	// 检查for体
	if s.Body != nil {
		c.checkBlockStmt(s.Body)
	}
}

// checkBlockStmt 检查块语句
func (c *Checker) checkBlockStmt(s *BlockStmt) {
	if s == nil {
		return
	}

	for _, stmt := range s.Stmts {
		c.checkStmt(stmt)
	}
}

// checkExpr 检查表达式
func (c *Checker) checkExpr(expr Expr) (string, interface{}) {
	if expr == nil {
		return "", nil
	}

	switch e := expr.(type) {
	case *LiteralExpr:
		return e.Type, e.Value

	case *IdentExpr:
		// 检查变量是否已定义
		if info, exists := c.scope[e.Name]; exists {
			return info.Type, info.Value
		}
		c.addError(e.Position(), fmt.Sprintf("未定义的变量: '%s'", e.Name))
		return "", nil

	case *AssignExpr:
		// 检查左值
		if e.Left == nil {
			c.addError(e.Position(), "赋值表达式左值不能为空")
			return "", nil
		}

		// 检查右值
		rightType, rightVal := c.checkExpr(e.Right)
		if rightType == "" {
			c.addError(e.Right.Position(), "赋值表达式右值类型不合法")
			return "", nil
		}

		// 更新变量类型和值
		c.scope[e.Left.Name] = VarInfo{
			Type:  rightType,
			Value: rightVal,
			Pos:   e.Position(),
		}

		return rightType, rightVal

	case *OpExpr:
		// 一元运算符 (++/--)
		if e.Right == nil && (e.Op == "++" || e.Op == "--") {
			leftType, leftVal := c.checkExpr(e.Left)
			if leftType != "number" {
				c.addError(e.Position(), fmt.Sprintf("运算符 '%s' 仅适用于数字类型，当前为 '%s' 类型", e.Op, leftType))
				return "", nil
			}

			// 模拟自增/自减
			switch v := leftVal.(type) {
			case int:
				if e.Op == "++" {
					leftVal = v + 1
				} else {
					leftVal = v - 1
				}
			case float64:
				if e.Op == "++" {
					leftVal = v + 1.0
				} else {
					leftVal = v - 1.0
				}
			}

			// 更新变量值（如果是标识符）
			if ident, ok := e.Left.(*IdentExpr); ok {
				if info, exists := c.scope[ident.Name]; exists {
					c.scope[ident.Name] = VarInfo{
						Type:  info.Type,
						Value: leftVal,
						Pos:   info.Pos,
					}
				}
			}

			return leftType, leftVal
		}

		// 二元运算符
		leftType, _ := c.checkExpr(e.Left)
		rightType, _ := c.checkExpr(e.Right)

		if leftType == "" || rightType == "" {
			return "", nil
		}

		// 检查类型匹配
		if leftType != rightType {
			c.addError(e.Position(), fmt.Sprintf("运算符 '%s' 左右操作数类型不匹配: '%s' vs '%s'",
				e.Op, leftType, rightType))
			return "", nil
		}

		// 比较运算符返回布尔类型
		return "bool", true

	default:
		c.addError(expr.Position(), fmt.Sprintf("不支持的表达式类型: %T", e))
		return "", nil
	}
}

// addError 添加语义错误
func (c *Checker) addError(pos Position, msg string) {
	c.errors = append(c.errors, ParseError{
		Pos:   pos,
		Msg:   msg,
		Level: "semantic",
	})
}

// Errors 获取语义错误
func (c *Checker) Errors() []ParseError {
	return c.errors
}

// GetScope 获取变量作用域
func (c *Checker) GetScope() map[string]VarInfo {
	return c.scope
}

// -------------------------- 主函数 --------------------------
func main() {
	// 测试用例
	testCode := `
let sName = "百度一下"
let 
let sName
let sName = 
let sName = "
le
let a = 1
let b = 2
if a > b {
	print sName
} elseif  {   # 这里应该会报错 没有条件语句
	print a
} else {
	print b
}

let c = true
if c {
	print c
}

if a { # 这里应该会报错，条件语句不是布尔值
	print a
}

for i=0; i<10; i++ {
	print i
}

for j=3 {  # 这里应该会报错，循环条件缺失
	print j
}
`

	// 1. 词法分析
	lexer := NewLexer(testCode)
	lexerErrors := lexer.Errors()
	if len(lexerErrors) > 0 {
		log.Println("===== 词法错误 =====")
		for _, err := range lexerErrors {
			log.Println(err)
		}
	}

	// 2. 语法分析
	parser := NewParser(lexer)
	program := parser.Parse()
	parserErrors := parser.Errors()
	if len(parserErrors) > 0 {
		log.Println("\n===== 语法错误 =====")
		for _, err := range parserErrors {
			log.Println(err)
		}
	}

	// 3. 语义分析
	checker := NewChecker(program)
	checker.Check()
	semanticErrors := checker.Errors()
	if len(semanticErrors) > 0 {
		log.Println("\n===== 语义错误 =====")
		for _, err := range semanticErrors {
			log.Println(err)
		}
	}

	// 4. 输出已定义的变量
	log.Println("\n===== 已定义的变量 =====")
	scope := checker.GetScope()
	for name, info := range scope {
		log.Printf("变量名: %s | 类型: %s | 值: %v | 定义位置: %s",
			name, info.Type, info.Value, info.Pos)
	}

	// 5. 总结
	totalErrors := len(lexerErrors) + len(parserErrors) + len(semanticErrors)
	log.Printf("\n===== 总结 =====")
	log.Printf("词法错误: %d 个", len(lexerErrors))
	log.Printf("语法错误: %d 个", len(parserErrors))
	log.Printf("语义错误: %d 个", len(semanticErrors))
	log.Printf("总错误数: %d 个", totalErrors)
}
