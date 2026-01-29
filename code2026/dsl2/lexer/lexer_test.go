package lexer

import (
	"testing"
)

func TestNextTokenBasic(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenAssign, "="},
		{TokenPlus, "+"},
		{TokenLParen, "("},
		{TokenRParen, ")"},
		{TokenLBrace, "{"},
		{TokenRBrace, "}"},
		{TokenComma, ","},
		{TokenSemicolon, ";"},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("测试[%d] - 类型错误。期望=%q, 得到=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("测试[%d] - 字面量错误。期望=%q, 得到=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextTokenOperators(t *testing.T) {
	input := `!-*/%<><=>===!=&&||`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenNot, "!"},
		{TokenMinus, "-"},
		{TokenAsterisk, "*"},
		{TokenSlash, "/"},
		{TokenMod, "%"},
		{TokenLT, "<"},
		{TokenGT, ">"},
		{TokenLE, "<="},
		{TokenGE, ">="},
		{TokenEQ, "=="},
		{TokenNE, "!="},
		{TokenAnd, "&&"},
		{TokenOr, "||"},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("测试[%d] - 类型错误。期望=%q, 得到=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("测试[%d] - 字面量错误。期望=%q, 得到=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextTokenKeywords(t *testing.T) {
	input := `var if else while return true false chrome break continue for`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenVar, "var"},
		{TokenIf, "if"},
		{TokenElse, "else"},
		{TokenWhile, "while"},
		{TokenReturn, "return"},
		{TokenTrue, "true"},
		{TokenFalse, "false"},
		{TokenChrome, "chrome"},
		{TokenBreak, "break"},
		{TokenContinue, "continue"},
		{TokenFor, "for"},
		{TokenEOF, ""},
		// todo ... 添加更多测试
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("测试[%d] - 类型错误。期望=%q, 得到=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("测试[%d] - 字面量错误。期望=%q, 得到=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextTokenIdentifiers(t *testing.T) {
	input := `foo bar123 _test var123 ifelse`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenIdent, "foo"},
		{TokenIdent, "bar123"},
		{TokenIdent, "_test"},
		{TokenIdent, "var123"}, // 以关键字开头但不是关键字
		{TokenIdent, "ifelse"}, // 包含关键字但不是关键字
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("测试[%d] - 类型错误。期望=%q, 得到=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("测试[%d] - 字面量错误。期望=%q, 得到=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextTokenIntegers(t *testing.T) {
	input := `123 456 0 999`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenInt, "123"},
		{TokenInt, "456"},
		{TokenInt, "0"},
		{TokenInt, "999"},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("测试[%d] - 类型错误。期望=%q, 得到=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("测试[%d] - 字面量错误。期望=%q, 得到=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextTokenStrings(t *testing.T) {
	input := `"hello" "world" "escaped \"quote\"" "line\nbreak"`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenString, "hello"},
		{TokenString, "world"},
		{TokenString, "escaped \"quote\""},
		{TokenString, "line\nbreak"},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("测试[%d] - 类型错误。期望=%q, 得到=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("测试[%d] - 字面量错误。期望=%q, 得到=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextTokenListsAndDicts(t *testing.T) {
	input := `[] [1, 2, 3] {} {"a": 1} .`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenLBracket, "["},
		{TokenRBracket, "]"},
		{TokenLBracket, "["},
		{TokenInt, "1"},
		{TokenComma, ","},
		{TokenInt, "2"},
		{TokenComma, ","},
		{TokenInt, "3"},
		{TokenRBracket, "]"},
		{TokenLBrace, "{"},
		{TokenRBrace, "}"},
		{TokenLBrace, "{"},
		{TokenString, "a"},
		{TokenColon, ":"},
		{TokenInt, "1"},
		{TokenRBrace, "}"},
		{TokenDot, "."}, // 链式调用操作符
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("测试[%d] - 类型错误。期望=%q, 得到=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("测试[%d] - 字面量错误。期望=%q, 得到=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextTokenComplexString(t *testing.T) {
	// 测试原始问题：包含转义引号和等号的字符串
	input := `print("[\"a\"] = ", 1)`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenIdent, "print"},
		{TokenLParen, "("},
		{TokenString, "[\"a\"] = "},
		{TokenComma, ","},
		{TokenInt, "1"},
		{TokenRParen, ")"},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("测试[%d] - 类型错误。期望=%q, 得到=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("测试[%d] - 字面量错误。期望=%q, 得到=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextTokenChainCall(t *testing.T) {
	input := `print("aa").print("bb").upper()`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenIdent, "print"},
		{TokenLParen, "("},
		{TokenString, "aa"},
		{TokenRParen, ")"},
		{TokenDot, "."},
		{TokenIdent, "print"},
		{TokenLParen, "("},
		{TokenString, "bb"},
		{TokenRParen, ")"},
		{TokenDot, "."},
		{TokenIdent, "upper"},
		{TokenLParen, "("},
		{TokenRParen, ")"},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("测试[%d] - 类型错误。期望=%q, 得到=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("测试[%d] - 字面量错误。期望=%q, 得到=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextTokenComments(t *testing.T) {
	input := `// 这是注释
var x = 10
# 这也是注释
print(x)`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenVar, "var"},
		{TokenIdent, "x"},
		{TokenAssign, "="},
		{TokenInt, "10"},
		{TokenIdent, "print"},
		{TokenLParen, "("},
		{TokenIdent, "x"},
		{TokenRParen, ")"},
		{TokenEOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("测试[%d] - 类型错误。期望=%q, 得到=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("测试[%d] - 字面量错误。期望=%q, 得到=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
