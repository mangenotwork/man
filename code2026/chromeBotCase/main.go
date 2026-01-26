package main

/*

通过语法定义关键字和结构来实现一个简单的解析器，用于解析和执行自定义脚本语言。

当前例子let解析器实现

当前例子还存在很多问题
但是 if 和 for 存在问题，错误无法正确识别

*/

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// LineRun 存储每行可执行的关键字和参数
type LineRun struct {
	LineNum   int        // 行号
	Keyword   string     // 关键字（let/if/elseif/else/for/print）
	Args      []string   // 参数
	Condition string     // if/for 的条件表达式
	Body      []*LineRun // 代码块内容（if/for 后的 {} 内）
}

// Lexer 解析器结构体
type Lexer struct {
	input              string     // 所有代码
	line               int        // 当前读取的行索引（从0开始）
	originalLineNumber int        // 原始代码行号（用户可见的行号）
	ch                 []string   // 当前行按空格分割后的切片
	variable           sync.Map   // 全局变量（key:变量名, value:变量值+类型）
	err                []string   // 错误信息
	codeLines          []string   // 按行分割后的完整代码（保留空行，用于行号对应）
	cleanedCodeLines   []string   // 清理后的代码行（无空行）
	lineMapping        []int      // 清理后的行索引对应原始行号的映射
	blockDepth         int        // 代码块嵌套深度（处理 {} 包裹的内容）
	inIfBlock          bool       // 是否在if代码块中（用于校验elseif/else）
	lastIfLine         int        // 上一个if的行号（用于关联elseif/else）
	currentBlock       []*LineRun // 当前代码块的内容
}

// 关键字映射（合法的关键字）
var keywordMap = map[string]struct{}{
	"let":    {}, // 定义变量
	"if":     {}, // 条件判断
	"else":   {}, // 条件判断
	"elseif": {}, // 条件判断
	"for":    {}, // 循环
	"print":  {}, // 输出打印
}

// 布尔值字面量
var boolLiterals = map[string]struct{}{
	"true":  {},
	"false": {},
}

// 变量类型枚举
const (
	VarTypeString = "string"
	VarTypeNumber = "number"
	VarTypeBool   = "bool"
)

// 校验变量名是否合法（字母开头，仅含字母/数字/下划线）
var validVarNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func main() {
	casetxt := `
let sName = "百度一下"
let 
let sName
let sName = 
let sName = "
le
let a =1
let b =2
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

	// 初始化解析器
	lexer := NewLexer(casetxt)

	// 存储解析后的可执行行
	lineRuns := make([]*LineRun, 0)

	// 逐行解析代码（使用清理后的行，但行号对应原始代码）
	for lexer.line < len(lexer.cleanedCodeLines) {
		currentLine := lexer.cleanedCodeLines[lexer.line]
		// 获取原始行号（用户可见的行号）
		lexer.originalLineNumber = lexer.lineMapping[lexer.line]
		lexer.line++

		// 按空格分割当前行
		lexer.ch = splitLineToParts(currentLine)
		if len(lexer.ch) == 0 {
			continue
		}

		// 处理包含 } 和 elseif/else 的复合行
		trimmedLine := strings.TrimSpace(currentLine)
		hasClosingBrace := false
		remainingLine := trimmedLine

		// 场景：} 后跟 elseif/else（如 "} elseif {"、"} else {"）
		if strings.HasPrefix(trimmedLine, "}") {
			// 分割 } 和后续内容
			parts := strings.SplitN(trimmedLine, "}", 2)
			bracePart := strings.TrimSpace(parts[0])
			remainingLine = strings.TrimSpace(parts[1])

			// 仅处理单个 } 的情况
			if bracePart == "}" {
				lexer.blockDepth--
				// 不再报多余的 } 错误，优先处理后续的 elseif/else
				hasClosingBrace = true
			}

			// 如果还有剩余内容，重新处理
			if remainingLine != "" {
				currentLine = remainingLine
				lexer.ch = splitLineToParts(currentLine)
				if len(lexer.ch) == 0 {
					continue
				}
			} else {
				continue
			}
		}

		// 场景：纯 } 行
		if trimmedLine == "}" {
			lexer.blockDepth--
			// 仅当深度为负时才报错（真正的多余 }）
			if lexer.blockDepth < 0 {
				lexer.addError(fmt.Sprintf("多余的代码块结束符 }"))
				lexer.blockDepth = 0 // 重置，避免后续错误累积
			}
			continue
		}

		// 关键字校验
		firstToken := lexer.ch[0]
		if _, isKeyword := keywordMap[firstToken]; !isKeyword {
			lexer.addError(fmt.Sprintf("无效的关键字: %s", firstToken))
			continue
		}

		// 按关键字分支解析
		switch firstToken {
		case "let":
			run, err := lexer.parseLet()
			if err != nil {
				lexer.addError(err.Error())
				continue
			}
			if run != nil {
				if lexer.blockDepth > 0 {
					lexer.currentBlock = append(lexer.currentBlock, run)
				} else {
					lineRuns = append(lineRuns, run)
				}
			}

		case "if":
			run, err := lexer.parseIf(currentLine)
			if err != nil {
				lexer.addError(err.Error())
				continue
			}
			if run != nil {
				lexer.inIfBlock = true
				lexer.lastIfLine = lexer.originalLineNumber
				lineRuns = append(lineRuns, run)
			}

		case "elseif":
			// 校验elseif是否在合法位置（如果是从} elseif来的，已经处理过blockDepth）
			if !lexer.inIfBlock && !hasClosingBrace {
				lexer.addError("elseif 必须跟在if或另一个elseif之后")
				continue
			}
			run, err := lexer.parseElseIf(currentLine)
			if err != nil {
				lexer.addError(err.Error())
				continue
			}
			if run != nil {
				// 如果是从} elseif来的，需要增加blockDepth
				if hasClosingBrace {
					lexer.blockDepth++
				}
				lineRuns = append(lineRuns, run)
			}

		case "else":
			// 校验else是否在合法位置
			if !lexer.inIfBlock && !hasClosingBrace {
				lexer.addError("else 必须跟在if或elseif之后")
				continue
			}
			run, err := lexer.parseElse(currentLine)
			if err != nil {
				lexer.addError(err.Error())
				continue
			}
			if run != nil {
				// 如果是从} else来的，需要增加blockDepth
				if hasClosingBrace {
					lexer.blockDepth++
				}
				lexer.inIfBlock = false
				lineRuns = append(lineRuns, run)
			}

		case "for":
			run, err := lexer.parseFor(currentLine)
			if err != nil {
				lexer.addError(err.Error())
				continue
			}
			if run != nil {
				lineRuns = append(lineRuns, run)
			}

		case "print":
			run := &LineRun{
				LineNum: lexer.originalLineNumber,
				Keyword: "print",
				Args:    lexer.ch[1:],
			}
			if lexer.blockDepth > 0 {
				lexer.currentBlock = append(lexer.currentBlock, run)
			} else {
				lineRuns = append(lineRuns, run)
			}

		default:
			lexer.addError(fmt.Sprintf("暂未实现的关键字解析: %s", firstToken))
		}
	}

	// 检查是否有未闭合的代码块
	if lexer.blockDepth > 0 {
		lexer.addError(fmt.Sprintf("代码块未闭合，剩余未结束的代码块深度：%d", lexer.blockDepth))
	}

	// 输出解析错误
	if len(lexer.err) > 0 {
		log.Println("===== 解析错误 =====")
		for _, e := range lexer.err {
			log.Println(e)
		}
	} else {
		log.Println("===== 解析成功 =====")
	}

	// 输出解析结果
	log.Println("\n===== 解析结果 =====")
	for _, run := range lineRuns {
		log.Printf("行%d | 关键字: %s | 参数: %v | 条件: %s",
			run.LineNum, run.Keyword, run.Args, run.Condition)
	}

	// 输出已定义的变量
	log.Println("\n===== 已定义变量 =====")
	lexer.variable.Range(func(key, value interface{}) bool {
		varInfo := value.(map[string]interface{})
		log.Printf("变量名: %s | 值: %v | 类型: %s",
			key, varInfo["value"], varInfo["type"])
		return true
	})
}

// NewLexer 创建并初始化解析器
func NewLexer(input string) *Lexer {
	// 按行分割原始代码（保留所有行，包括空行）
	originalLines := strings.Split(input, "\n")
	cleanedLines := make([]string, 0)
	lineMapping := make([]int, 0) // 清理后的行索引 -> 原始行号

	// 构建清理后的行和行号映射（原始行号从1开始）
	for idx, line := range originalLines {
		originalLineNum := idx + 1
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleanedLines = append(cleanedLines, line)
			lineMapping = append(lineMapping, originalLineNum)
		}
	}

	return &Lexer{
		input:              input,
		codeLines:          originalLines,
		cleanedCodeLines:   cleanedLines,
		lineMapping:        lineMapping,
		variable:           sync.Map{},
		err:                make([]string, 0),
		blockDepth:         0,
		inIfBlock:          false,
		lastIfLine:         0,
		currentBlock:       make([]*LineRun, 0),
		line:               0,
		originalLineNumber: 0,
	}
}

// parseLet 解析 let 关键字（变量定义）
func (l *Lexer) parseLet() (*LineRun, error) {
	// 1. 基础语法校验：let 后必须有内容
	if len(l.ch) < 2 {
		return nil, errors.New("let 关键字后缺少变量定义（正确格式：let 变量名=值）")
	}

	// 2. 拼接变量定义部分
	varDefStr := strings.Join(l.ch[1:], "")
	equalCount := strings.Count(varDefStr, "=")

	// 3. 校验等号数量
	if equalCount != 1 {
		return nil, fmt.Errorf("变量定义必须且仅能有一个等号（当前：%d个），正确格式：let 变量名=值", equalCount)
	}

	// 4. 分割变量名和变量值
	sn := strings.SplitN(varDefStr, "=", 2)
	varName := strings.TrimSpace(sn[0])
	varValueStr := strings.TrimSpace(sn[1])

	// 5. 校验变量名合法性
	if !validVarNameRegex.MatchString(varName) {
		return nil, fmt.Errorf("变量名 %s 不合法（必须以字母/下划线开头，仅含字母/数字/下划线）", varName)
	}

	// 6. 校验变量值不能为空
	if varValueStr == "" {
		return nil, errors.New("变量值不能为空（正确格式：let 变量名=值）")
	}

	// 7. 识别变量类型并校验
	var varType string
	var varValue interface{}
	var err error

	// 7.1 字符串类型
	if len(varValueStr) >= 2 && varValueStr[0] == '"' {
		if varValueStr[len(varValueStr)-1] != '"' {
			return nil, fmt.Errorf("字符串值 %s 缺少闭合的双引号", varValueStr)
		}
		varType = VarTypeString
		varValue = varValueStr[1 : len(varValueStr)-1]

		// 7.2 布尔类型
	} else if _, isBool := boolLiterals[varValueStr]; isBool {
		varType = VarTypeBool
		varValue, err = strconv.ParseBool(varValueStr)
		if err != nil {
			return nil, fmt.Errorf("布尔值解析失败：%s", err.Error())
		}

		// 7.3 数值类型
	} else if isNumber(varValueStr) {
		varType = VarTypeNumber
		if intVal, errInt := strconv.Atoi(varValueStr); errInt == nil {
			varValue = intVal
		} else if floatVal, errFloat := strconv.ParseFloat(varValueStr, 64); errFloat == nil {
			varValue = floatVal
		} else {
			return nil, fmt.Errorf("数值类型解析失败：%s", varValueStr)
		}

		// 7.4 未知类型
	} else {
		return nil, fmt.Errorf("变量值 %s 类型不支持（仅支持字符串/数值/布尔类型）", varValueStr)
	}

	// 8. 存储变量
	l.variable.Store(varName, map[string]interface{}{
		"value": varValue,
		"type":  varType,
	})

	// 9. 返回解析结果
	return &LineRun{
		LineNum: l.originalLineNumber,
		Keyword: "let",
		Args:    []string{varName, fmt.Sprintf("%v", varValue), varType},
	}, nil
}

// parseIf 解析 if 关键字
func (l *Lexer) parseIf(line string) (*LineRun, error) {
	condition, err := extractCondition(line, "if")
	if err != nil {
		return nil, err
	}

	// 校验条件表达式类型
	if err := validateConditionType(l, condition); err != nil {
		return nil, fmt.Errorf("if 条件类型不合法：%s", err.Error())
	}

	l.blockDepth++

	return &LineRun{
		LineNum:   l.originalLineNumber,
		Keyword:   "if",
		Condition: condition,
	}, nil
}

// parseElseIf 解析 elseif 关键字
func (l *Lexer) parseElseIf(line string) (*LineRun, error) {
	condition, err := extractCondition(line, "elseif")
	if err != nil {
		return nil, err
	}

	// 校验elseif条件不能为空
	if condition == "" {
		return nil, errors.New("elseif 条件表达式不能为空")
	}

	// 校验条件表达式类型
	if err := validateConditionType(l, condition); err != nil {
		return nil, fmt.Errorf("elseif 条件类型不合法：%s", err.Error())
	}

	// 如果不是从} elseif来的，才增加blockDepth
	if l.blockDepth == 0 && !l.inIfBlock {
		l.blockDepth++
	}

	return &LineRun{
		LineNum:   l.originalLineNumber,
		Keyword:   "elseif",
		Condition: condition,
	}, nil
}

// parseElse 解析 else 关键字
func (l *Lexer) parseElse(line string) (*LineRun, error) {
	trimmedLine := strings.TrimSpace(line)
	if !strings.HasSuffix(trimmedLine, "{") || strings.TrimSpace(strings.TrimSuffix(trimmedLine, "{")) != "else" {
		return nil, errors.New("else 语法错误（正确格式：else {）")
	}

	// 如果不是从} else来的，才增加blockDepth
	if l.blockDepth == 0 && !l.inIfBlock {
		l.blockDepth++
	}

	return &LineRun{
		LineNum: l.originalLineNumber,
		Keyword: "else",
	}, nil
}

// parseFor 解析 for 关键字
func (l *Lexer) parseFor(line string) (*LineRun, error) {
	condition, err := extractCondition(line, "for")
	if err != nil {
		return nil, err
	}

	parts := strings.Split(condition, ";")
	if len(parts) != 3 {
		return nil, errors.New("for 循环表达式必须是三段式（正确格式：for 初始化; 条件; 增量 {）")
	}

	// 基础校验每一部分
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, fmt.Errorf("for 循环第%d部分不能为空（初始化; 条件; 增量）", i+1)
		}
	}

	l.blockDepth++

	return &LineRun{
		LineNum:   l.originalLineNumber,
		Keyword:   "for",
		Condition: condition,
	}, nil
}

// validateConditionType 校验条件表达式类型
func validateConditionType(l *Lexer, condition string) error {
	// 1. 布尔字面量直接通过
	if _, isBool := boolLiterals[condition]; isBool {
		return nil
	}

	// 2. 检查变量类型
	condition = strings.TrimSpace(condition)
	varNames := extractVariablesFromCondition(condition)

	for _, varName := range varNames {
		if isOperator(varName) {
			continue
		}

		val, ok := l.variable.Load(varName)
		if !ok {
			return fmt.Errorf("变量 %s 未定义", varName)
		}

		varInfo := val.(map[string]interface{})
		if varInfo["type"] != VarTypeBool {
			return fmt.Errorf("变量 %s 类型为 %s，条件表达式要求布尔类型", varName, varInfo["type"])
		}
	}

	// 3. 运算符校验
	validOps := []string{">", "<", "=", "!=", ">=", "<="}
	hasValidOp := false
	for _, op := range validOps {
		if strings.Contains(condition, op) {
			hasValidOp = true
			break
		}
	}

	_, isBoolLiteral := boolLiterals[condition]
	if !hasValidOp && !isBoolLiteral {
		return errors.New("条件必须包含合法运算符（> < = != >= <=）或布尔值（true/false）")
	}

	return nil
}

// extractVariablesFromCondition 提取变量名
func extractVariablesFromCondition(condition string) []string {
	ops := []string{">", "<", "=", "!=", ">=", "<=", " ", "\t"}
	cleaned := condition
	for _, op := range ops {
		cleaned = strings.ReplaceAll(cleaned, op, "|")
	}

	parts := strings.Split(cleaned, "|")
	var vars []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" && validVarNameRegex.MatchString(part) {
			vars = append(vars, part)
		}
	}
	return vars
}

// isOperator 判断是否是运算符
func isOperator(s string) bool {
	ops := map[string]struct{}{
		">":  {},
		"<":  {},
		"=":  {},
		"!=": {},
		">=": {},
		"<=": {},
		"++": {},
		"--": {},
	}
	_, ok := ops[s]
	return ok
}

// addError 添加错误信息（使用正确的原始行号）
func (l *Lexer) addError(msg string) {
	l.err = append(l.err, fmt.Sprintf("Err: line %d, position 0, %s", l.originalLineNumber, msg))
}

// splitLineToParts 分割行
func splitLineToParts(line string) []string {
	return strings.Fields(line)
}

// isNumber 判断是否是数值
func isNumber(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return true
	}
	return false
}

// extractCondition 提取条件表达式
func extractCondition(line, keyword string) (string, error) {
	trimmedLine := strings.TrimSpace(line)
	conditionPart := strings.TrimPrefix(trimmedLine, keyword)
	conditionPart = strings.TrimSpace(conditionPart)

	if !strings.HasSuffix(conditionPart, "{") {
		return "", fmt.Errorf("%s 后必须跟条件表达式和 {（正确格式：%s 条件 {）", keyword, keyword)
	}

	condition := strings.TrimSpace(strings.TrimSuffix(conditionPart, "{"))
	if condition == "" && keyword != "elseif" {
		return "", fmt.Errorf("%s 条件表达式不能为空", keyword)
	}

	return condition, nil
}
