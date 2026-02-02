package lexer

/*

词法分析

*/

import (
	"dsl2/logger"
	"strings"
	"unicode"
)

// TokenType 令牌类型
type TokenType int

const (
	TokenIllegal TokenType = iota
	TokenEOF

	// 标识符和字面量
	TokenIdent
	TokenInt
	TokenFloat
	TokenString
	TokenBool

	// 运算符
	TokenAssign   // =
	TokenPlus     // +
	TokenMinus    // -
	TokenAsterisk // *
	TokenSlash    // /
	TokenMod      // %
	TokenEQ       // ==
	TokenNE       // !=
	TokenLT       // <
	TokenLE       // <=
	TokenGT       // >
	TokenGE       // >=
	TokenAnd      // &&
	TokenOr       // ||
	TokenNot      // !

	// 分隔符
	TokenComma     // ,
	TokenSemicolon // ;
	TokenColon     // :
	TokenLParen    // (
	TokenRParen    // )
	TokenLBrace    // {
	TokenRBrace    // }
	TokenLBracket  // [
	TokenRBracket  // ]
	TokenDot       // .

	// 关键字
	TokenVar
	TokenIf
	TokenElse
	TokenWhile
	TokenReturn
	TokenTrue
	TokenFalse
	TokenChrome
	TokenBreak
	TokenContinue
	TokenFor
)

var tokenTypeStrings = map[TokenType]string{
	TokenIllegal:   "ILLEGAL",
	TokenEOF:       "EOF",
	TokenIdent:     "IDENT",
	TokenInt:       "INT",
	TokenFloat:     "FLOAT",
	TokenString:    "STRING",
	TokenBool:      "BOOL",
	TokenAssign:    "=",
	TokenPlus:      "+",
	TokenMinus:     "-",
	TokenAsterisk:  "*",
	TokenSlash:     "/",
	TokenMod:       "%",
	TokenEQ:        "==",
	TokenNE:        "!=",
	TokenLT:        "<",
	TokenLE:        "<=",
	TokenGT:        ">",
	TokenGE:        ">=",
	TokenAnd:       "&&",
	TokenOr:        "||",
	TokenNot:       "!",
	TokenComma:     ",",
	TokenSemicolon: ";",
	TokenColon:     ":",
	TokenLParen:    "(",
	TokenRParen:    ")",
	TokenLBrace:    "{",
	TokenRBrace:    "}",
	TokenLBracket:  "[",
	TokenRBracket:  "]",
	TokenDot:       ".",
	TokenVar:       "var",
	TokenIf:        "if",
	TokenElse:      "else",
	TokenWhile:     "while",
	TokenReturn:    "return",
	TokenTrue:      "true",
	TokenFalse:     "false",
	TokenChrome:    "chrome",
	TokenBreak:     "break",
	TokenContinue:  "continue",
	TokenFor:       "for",
}

func (t TokenType) String() string {
	if s, ok := tokenTypeStrings[t]; ok {
		return s
	}
	return "UNKNOWN"
}

// Token 令牌
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Lexer 词法分析器
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune
	line         int
	column       int
}

// New 创建词法分析器
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = rune(l.input[l.readPosition])
	}

	l.position = l.readPosition
	l.readPosition++
	l.column++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
	logger.Debug("readChar: 后: ch=%v, position=%d, readPosition=%d",
		l.ch, l.position, l.readPosition)
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return rune(l.input[l.readPosition])
}

// NextToken 获取下一个令牌
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	tok.Line = l.line
	tok.Column = l.column

	logger.Debug("NextToken -> ", string(l.ch))

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			logger.Debug("TokenEQ")
			l.readChar()
			tok.Type = TokenEQ
			tok.Literal = "=="
		} else {
			logger.Debug("TokenEQ !! ")

			//// 检查前面是否有标识符
			//// 如果前面是标识符字符，则 = 应该属于标识符的一部分
			//if l.position > 0 && (isLetter(rune(l.input[l.position-1])) || isDigit(rune(l.input[l.position-1])) ||
			//	l.input[l.position-1] == '_' || l.input[l.position-1] == '"') {
			//
			//	logger.Debug("回退一个字符，然后读取整个标识符 ")
			//	// 回退一个字符，然后读取整个标识符
			//	l.position--
			//	l.readPosition--
			//	l.column--
			//	l.ch = rune(l.input[l.position])
			//
			//	// 读取整个标识符（包含 = 和后面的值）
			//	tok.Literal = l.readIdentifier()
			//	logger.Debug("tok.Literal = ", tok.Literal)
			//	tok.Type = l.lookupIdent(tok.Literal)
			//	return tok
			//} else {
			//	logger.Debug("TokenAssign")
			//	tok.Type = TokenAssign
			//	tok.Literal = string(l.ch)
			//}

			logger.Debug("TokenAssign")
			tok.Type = TokenAssign
			tok.Literal = string(l.ch)

			//l.readChar() // 关键：读取下一个字符
			//return tok
		}
	case '+':
		tok.Type = TokenPlus
		tok.Literal = string(l.ch)
	case '-':
		tok.Type = TokenMinus
		tok.Literal = string(l.ch)
	case '*':
		tok.Type = TokenAsterisk
		tok.Literal = string(l.ch)
	case '/':
		if l.peekChar() == '/' {
			l.skipComment()
			return l.NextToken()
		}
		tok.Type = TokenSlash
		tok.Literal = string(l.ch)
	case '#':
		l.skipComment()
		return l.NextToken()
	case '%':
		tok.Type = TokenMod
		tok.Literal = string(l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TokenNE
			tok.Literal = "!="
		} else {
			tok.Type = TokenNot
			tok.Literal = string(l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TokenLE
			tok.Literal = "<="
		} else {
			tok.Type = TokenLT
			tok.Literal = string(l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TokenGE
			tok.Literal = ">="
		} else {
			tok.Type = TokenGT
			tok.Literal = string(l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok.Type = TokenAnd
			tok.Literal = "&&"
		} else {
			tok.Type = TokenIllegal
			tok.Literal = string(l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok.Type = TokenOr
			tok.Literal = "||"
		} else {
			tok.Type = TokenIllegal
			tok.Literal = string(l.ch)
		}
	case ',':
		tok.Type = TokenComma
		tok.Literal = string(l.ch)
	case ';':
		tok.Type = TokenSemicolon
		tok.Literal = string(l.ch)
	case ':':
		tok.Type = TokenColon
		tok.Literal = string(l.ch)
	case '(':
		tok.Type = TokenLParen
		tok.Literal = string(l.ch)
	case ')':
		tok.Type = TokenRParen
		tok.Literal = string(l.ch)
	case '{':
		tok.Type = TokenLBrace
		tok.Literal = string(l.ch)
	case '}':
		tok.Type = TokenRBrace
		tok.Literal = string(l.ch)
	case '[':
		tok.Type = TokenLBracket
		tok.Literal = string(l.ch)
	case ']':
		tok.Type = TokenRBracket
		tok.Literal = string(l.ch)
	case '"':
		tok.Type = TokenString
		tok.Literal = l.readString()
		logger.Debug("读取到字符串: '%s'", tok.Literal)
	case '.':
		// 检查是否是浮点数的一部分
		if isDigit(l.ch) || (l.ch == '.' && isDigit(l.peekChar())) {
			// 读取数字，可能是整数或浮点数
			number := l.readNumber()
			// 检查是否包含小数点
			if strings.Contains(number, ".") {
				tok.Type = TokenFloat
				tok.Literal = number
			} else {
				tok.Type = TokenInt
				tok.Literal = number
			}
		} else {
			// 这是链式调用操作符
			tok.Type = TokenDot
			tok.Literal = string(l.ch)
		}
	case 0:
		tok.Type = TokenEOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = l.lookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = TokenInt
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok.Type = TokenIllegal
			tok.Literal = string(l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	l.readChar() // 跳过换行符
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for {
		// 允许的字符：字母、数字、下划线、等号、点、冒号、减号
		if isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
			l.readChar()

		} else if l.ch == '"' {
			// 处理带引号的部分
			l.readChar() // 读取引号
			// 读取引号内的内容
			for l.ch != '"' && l.ch != 0 {
				if l.ch == '\\' {
					l.readChar() // 跳过转义符
					if l.ch != 0 {
						l.readChar() // 读取转义字符
					}
				} else {
					l.readChar()
				}
			}
			// 跳过结束引号
			if l.ch == '"' {
				l.readChar()
			}
		} else {
			// 遇到其他字符，停止读取
			break
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	// 检查是否有小数点
	if l.ch == '.' {
		// 读取小数点
		l.readChar()
		// 读取小数部分
		for isDigit(l.ch) {
			l.readChar()
		}
		return l.input[position:l.position]
	}

	// 读取整数部分
	for isDigit(l.ch) {
		l.readChar()
	}

	// 检查是否有小数点
	if l.ch == '.' && isDigit(l.peekChar()) {
		// 读取小数点
		l.readChar()
		// 读取小数部分
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) readFloat() string {
	position := l.position

	// 如果以小数点开头，添加前导0
	if l.ch == '.' {
		l.readChar()
	}

	// 读取小数部分
	for isDigit(l.ch) {
		l.readChar()
	}

	number := l.input[position:l.position]
	// 如果以小数点开头，添加前导0
	if number[0] == '.' {
		number = "0" + number
	}

	return number
}

func (l *Lexer) readString() string {
	position := l.position + 1 // 跳过开始的引号

	for {
		l.readChar()

		if l.ch == 0 {
			// 文件结束
			break
		}

		if l.ch == '"' {
			break
		}

		if l.ch == '\\' {
			// 跳过转义字符和它后面的字符
			l.readChar()
			continue
		}

	}

	if l.ch == '"' {
		// 正常结束，有右引号
		str := l.input[position:l.position]
		return unescapeString(str)
	}

	// 字符串没有正确结束
	return l.input[position:l.position]
}

func unescapeString(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	for i := 0; i < len(s); {
		if s[i] == '\\' && i+1 < len(s) {
			i++ // 跳过反斜杠

			// 只处理这些特定的转义序列
			switch s[i] {
			case '\\':
				result.WriteByte('\\')
			case '"':
				result.WriteByte('"')
			case 'n':
				result.WriteByte('\n')
			case 't':
				result.WriteByte('\t')
			case 'r':
				result.WriteByte('\r')
			case '0':
				result.WriteByte('\x00')
			default:
				// 关键：对于 \t, \n 等，如果不是转义序列的一部分
				// 应该保持为 \t, \n
				// 所以写回反斜杠和字符
				result.WriteByte('\\')
				result.WriteByte(s[i])
			}

		} else {
			result.WriteByte(s[i])

		}
		i++
	}

	return result.String()
}

func (l *Lexer) lookupIdent(ident string) TokenType {
	switch ident {
	case "var":
		return TokenVar
	case "if":
		return TokenIf
	case "else":
		return TokenElse
	case "while":
		return TokenWhile
	case "return":
		return TokenReturn
	case "true":
		return TokenTrue
	case "false":
		return TokenFalse
	case "chrome":
		return TokenChrome
	case "break":
		return TokenBreak
	case "continue":
		return TokenContinue
	case "for":
		return TokenFor
	default:
		return TokenIdent
	}
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}
