package parser

import (
	"dsl2/ast"
	"dsl2/lexer"
	"fmt"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"var x = 5;", "x", int64(5)},
		{"var y = true;", "y", true},
		{"var foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements 不包含 1 条语句。得到=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.VarDecl).Expr
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements 不包含 3 条语句。得到=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStmt)
		if !ok {
			t.Errorf("stmt 不是 *ast.ReturnStatement。得到=%T", stmt)
			continue
		}
		if returnStmt.String() != "return" && returnStmt.Expr == nil {
			t.Errorf("returnStmt.Expr 是 nil")
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program 没有足够的语句。得到=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Statements[0] 不是 ast.ExpressionStmt。得到=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expr.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp 不是 *ast.Identifier。得到=%T", stmt.Expr)
	}
	if ident.Name != "foobar" {
		t.Errorf("ident.Name 不是 %s。得到=%s", "foobar", ident.Name)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program 没有足够的语句。得到=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Statements[0] 不是 ast.ExpressionStmt。得到=%T",
			program.Statements[0])
	}

	literal, ok := stmt.Expr.(*ast.Integer)
	if !ok {
		t.Fatalf("exp 不是 *ast.IntegerLiteral。得到=%T", stmt.Expr)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value 不是 %d。得到=%d", 5, literal.Value)
	}
	if literal.String() != "5" {
		t.Errorf("literal.String 不是 %s。得到=%s", "5", literal.String())
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program 没有足够的语句。得到=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("program.Statements[0] 不是 ast.ExpressionStmt。得到=%T",
				program.Statements[0])
		}

		boolExpr, ok := stmt.Expr.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp 不是 *ast.Boolean。得到=%T", stmt.Expr)
		}
		if boolExpr.Value != tt.expected {
			t.Errorf("boolExpr.Value 不是 %v。得到=%v", tt.expected, boolExpr.Value)
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStmt)
	literal, ok := stmt.Expr.(*ast.String)
	if !ok {
		t.Fatalf("exp 不是 *ast.StringLiteral。得到=%T", stmt.Expr)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value 不是 %q。得到=%q", "hello world", literal.Value)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", int64(5)},
		{"-15;", "-", int64(15)},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements 不包含 %d 条语句。得到=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("program.Statements[0] 不是 ast.ExpressionStmt。得到=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expr.(*ast.UnaryExpr)
		if !ok {
			t.Fatalf("stmt 不是 ast.UnaryExpr。得到=%T", stmt.Expr)
		}
		if exp.Op != tt.operator {
			t.Fatalf("exp.Operator 不是 '%s'。得到=%s",
				tt.operator, exp.Op)
		}
		if !testLiteralExpression(t, exp.Expr, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", int64(5), "+", int64(5)},
		{"5 - 5;", int64(5), "-", int64(5)},
		{"5 * 5;", int64(5), "*", int64(5)},
		{"5 / 5;", int64(5), "/", int64(5)},
		{"5 > 5;", int64(5), ">", int64(5)},
		{"5 < 5;", int64(5), "<", int64(5)},
		{"5 == 5;", int64(5), "==", int64(5)},
		{"5 != 5;", int64(5), "!=", int64(5)},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements 不包含 %d 条语句。得到=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("program.Statements[0] 不是 ast.ExpressionStmt。得到=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expr, tt.leftValue,
			tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("期望=%q, 得到=%q", tt.expected, actual)
		}
	}
}

func TestIfStatement(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body 不包含 %d 条语句。得到=%d\n",
			1, len(program.Statements))
	}

	// IfStmt 是语句，不是表达式
	stmt, ok := program.Statements[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("program.Statements[0] 不是 ast.IfStmt。得到=%T",
			program.Statements[0])
	}

	if !testInfixExpression(t, stmt.Condition, "x", "<", "y") {
		return
	}

	if len(stmt.Then.Stmts) != 1 {
		t.Errorf("then 块不是 1 条语句。得到=%d\n",
			len(stmt.Then.Stmts))
	}

	thenStmt, ok := stmt.Then.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("Statements[0] 不是 ast.ExpressionStmt。得到=%T",
			stmt.Then.Stmts[0])
	}

	if !testIdentifier(t, thenStmt.Expr, "x") {
		return
	}

	if stmt.Else != nil {
		t.Errorf("stmt.Else.Statements 不是 nil。得到=%+v", stmt.Else)
	}
}

func TestIfElseStatement(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body 不包含 %d 条语句。得到=%d\n",
			1, len(program.Statements))
	}

	// IfStmt 是语句，不是表达式
	stmt, ok := program.Statements[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("program.Statements[0] 不是 ast.IfStmt。得到=%T",
			program.Statements[0])
	}

	if !testInfixExpression(t, stmt.Condition, "x", "<", "y") {
		return
	}

	if len(stmt.Then.Stmts) != 1 {
		t.Errorf("then 块不是 1 条语句。得到=%d\n",
			len(stmt.Then.Stmts))
	}

	thenStmt, ok := stmt.Then.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("Statements[0] 不是 ast.ExpressionStmt。得到=%T",
			stmt.Then.Stmts[0])
	}

	if !testIdentifier(t, thenStmt.Expr, "x") {
		return
	}

	elseBlock, ok := stmt.Else.(*ast.BlockStmt)
	if !ok {
		t.Fatalf("stmt.Else 不是 ast.BlockStmt。得到=%T", stmt.Else)
	}

	if len(elseBlock.Stmts) != 1 {
		t.Errorf("else 块不是 1 条语句。得到=%d\n",
			len(elseBlock.Stmts))
	}

	elseStmt, ok := elseBlock.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("elseBlock.Statements[0] 不是 ast.ExpressionStmt。得到=%T",
			elseBlock.Stmts[0])
	}

	if !testIdentifier(t, elseStmt.Expr, "y") {
		return
	}
}

func TestWhileStatement(t *testing.T) {
	input := `while (x < 10) { x = x + 1 }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body 不包含 %d 条语句。得到=%d\n",
			1, len(program.Statements))
	}

	// WhileStmt 是语句
	stmt, ok := program.Statements[0].(*ast.WhileStmt)
	if !ok {
		t.Fatalf("program.Statements[0] 不是 ast.WhileStmt。得到=%T",
			program.Statements[0])
	}

	if !testInfixExpression(t, stmt.Condition, "x", "<", int64(10)) {
		return
	}

	if len(stmt.Body.Stmts) != 1 {
		t.Errorf("while 块不是 1 条语句。得到=%d\n",
			len(stmt.Body.Stmts))
	}
}

func TestListLiteral(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Statements[0] 不是 ast.ExpressionStmt。得到=%T",
			program.Statements[0])
	}

	list, ok := stmt.Expr.(*ast.List)
	if !ok {
		t.Fatalf("exp 不是 *ast.List。得到=%T", stmt.Expr)
	}

	if len(list.Elements) != 3 {
		t.Fatalf("len(list.Elements) 不是 3。得到=%d", len(list.Elements))
	}

	testIntegerLiteral(t, list.Elements[0], 1)
	testInfixExpression(t, list.Elements[1], 2, "*", 2)
	testInfixExpression(t, list.Elements[2], 3, "+", 3)
}

func TestDictLiteral(t *testing.T) {
	// 字典字面量作为变量初始值
	input := `var myDict = {"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements 长度错误。期望 1，得到=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.VarDecl)
	if !ok {
		t.Fatalf("第一条语句不是变量声明。得到=%T", program.Statements[0])
	}

	dict, ok := stmt.Expr.(*ast.Dict)
	if !ok {
		t.Fatalf("初始化表达式不是 *ast.Dict。得到=%T", stmt.Expr)
	}

	if len(dict.Pairs) != 3 {
		t.Fatalf("dict.Pairs 长度错误。期望 3，得到=%d", len(dict.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range dict.Pairs {
		literal, ok := key.(*ast.String)
		if !ok {
			t.Errorf("key 不是 *ast.StringLiteral。得到=%T", key)
			continue
		}

		expectedValue := expected[literal.Value]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("exp 不是 *ast.ExpressionStmt。得到=%T", program.Statements[0])
	}

	indexExp, ok := stmt.Expr.(*ast.IndexExpr)
	if !ok {
		t.Fatalf("exp 不是 *ast.IndexExpression。得到=%T", stmt.Expr)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestChainCallExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`print("hello").upper()`,
			`print("hello").upper()`,
		},
		{
			`getData().process().print()`,
			`getData().process().print()`,
		},
		{
			`list().sort().filter(x > 0)`,
			`list().sort().filter((x > 0))`,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("期望=%q, 得到=%q", tt.expected, actual)
		}
	}
}

func TestComplexChainCall(t *testing.T) {
	input := `print("aa").print("bb")`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body 不包含 1 条语句。得到=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Statements[0] 不是 ast.ExpressionStmt。得到=%T",
			program.Statements[0])
	}

	chainCall, ok := stmt.Expr.(*ast.ChainCallExpr)
	if !ok {
		t.Fatalf("stmt.Expr 不是 *ast.ChainCallExpr。得到=%T", stmt.Expr)
	}

	if len(chainCall.Calls) != 2 {
		t.Fatalf("chainCall.Calls 长度不是 2。得到=%d", len(chainCall.Calls))
	}

	// 检查第一个调用
	firstCall := chainCall.Calls[0]
	if firstCall.Function.Name != "print" {
		t.Errorf("第一个调用函数名不是 'print'。得到=%s", firstCall.Function.Name)
	}
	if len(firstCall.Args) != 1 {
		t.Fatalf("第一个调用参数长度不是 1。得到=%d", len(firstCall.Args))
	}
	testStringLiteral(t, firstCall.Args[0], "aa")

	// 检查第二个调用
	secondCall := chainCall.Calls[1]
	if secondCall.Function.Name != "print" {
		t.Errorf("第二个调用函数名不是 'print'。得到=%s", secondCall.Function.Name)
	}
	if len(secondCall.Args) != 1 {
		t.Fatalf("第二个调用参数长度不是 1。得到=%d", len(secondCall.Args))
	}
	testStringLiteral(t, secondCall.Args[0], "bb")
}

// 辅助函数保持不变...
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("解析器有 %d 个错误", len(errors))
	for _, msg := range errors {
		t.Errorf("解析器错误: %s", msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.String()[:3] != "var" {
		t.Errorf("s.String() 不以 'var' 开头。得到=%q", s.String())
		return false
	}

	letStmt, ok := s.(*ast.VarDecl)
	if !ok {
		t.Errorf("s 不是 *ast.VarDecl。得到=%T", s)
		return false
	}

	if letStmt.Name.Name != name {
		t.Errorf("letStmt.Name.Name 不是 '%s'。得到=%s", name, letStmt.Name.Name)
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.Integer)
	if !ok {
		t.Errorf("il 不是 *ast.Integer。得到=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value 不是 %d。得到=%d", value, integ.Value)
		return false
	}

	expectedStr := fmt.Sprintf("%d", value)
	if integ.String() != expectedStr {
		t.Errorf("integ.String 不是 %s。得到=%s", expectedStr, integ.String())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp 不是 *ast.Identifier。得到=%T", exp)
		return false
	}

	if ident.Name != value {
		t.Errorf("ident.Name 不是 %s。得到=%s", value, ident.Name)
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("exp 类型无法处理。得到=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.BinaryExpr)
	if !ok {
		t.Errorf("exp 不是 *ast.BinaryExpr。得到=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Op != operator {
		t.Errorf("exp.Operator 不是 '%s'。得到=%q", operator, opExp.Op)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp 不是 *ast.Boolean。得到=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value 不是 %t。得到=%t", value, bo.Value)
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
	str, ok := exp.(*ast.String)
	if !ok {
		t.Errorf("exp 不是 *ast.String。得到=%T", exp)
		return false
	}

	if str.Value != value {
		t.Errorf("str.Value 不是 %q。得到=%q", value, str.Value)
		return false
	}

	return true
}
