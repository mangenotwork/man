package ast

import (
	"strings"
	"testing"
)

func TestNodePositions(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected Position
	}{
		{
			name: "Identifier position",
			node: &Identifier{
				StartPos: Position{Line: 1, Column: 1},
				Name:     "x",
			},
			expected: Position{Line: 1, Column: 1},
		},
		{
			name: "Integer position",
			node: &Integer{
				StartPos: Position{Line: 2, Column: 3},
				Value:    42,
			},
			expected: Position{Line: 2, Column: 3},
		},
		{
			name: "String position",
			node: &String{
				StartPos: Position{Line: 3, Column: 5},
				Value:    "hello",
			},
			expected: Position{Line: 3, Column: 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := tt.node.Pos()
			if pos.Line != tt.expected.Line || pos.Column != tt.expected.Column {
				t.Errorf("期望位置 Line=%d, Column=%d, 得到 Line=%d, Column=%d",
					tt.expected.Line, tt.expected.Column, pos.Line, pos.Column)
			}
		})
	}
}

func TestNodeString(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected string
	}{
		{
			name: "Identifier string",
			node: &Identifier{
				Name: "myVar",
			},
			expected: "myVar",
		},
		{
			name: "Integer string",
			node: &Integer{
				Value: 123,
			},
			expected: "123",
		},
		{
			name: "String string",
			node: &String{
				Value: "test",
			},
			expected: "\"test\"",
		},
		{
			name: "Boolean true",
			node: &Boolean{
				Value: true,
			},
			expected: "true",
		},
		{
			name: "Boolean false",
			node: &Boolean{
				Value: false,
			},
			expected: "false",
		},
		{
			name: "Binary expression",
			node: &BinaryExpr{
				Left:  &Identifier{Name: "x"},
				Op:    "+",
				Right: &Identifier{Name: "y"},
			},
			expected: "(x + y)",
		},
		{
			name: "Unary expression",
			node: &UnaryExpr{
				Op:   "!",
				Expr: &Identifier{Name: "flag"},
			},
			expected: "(!flag)",
		},
		{
			name: "Call expression",
			node: &CallExpr{
				Function: &Identifier{Name: "print"},
				Args: []Expression{
					&String{Value: "hello"},
					&Integer{Value: 42},
				},
			},
			expected: "print(\"hello\", 42)",
		},
		{
			name: "List literal",
			node: &List{
				Elements: []Expression{
					&Integer{Value: 1},
					&Integer{Value: 2},
					&Integer{Value: 3},
				},
			},
			expected: "[1, 2, 3]",
		},
		{
			name: "Dict literal",
			node: &Dict{
				Pairs: map[Expression]Expression{
					&String{Value: "a"}: &Integer{Value: 1},
					&String{Value: "b"}: &Integer{Value: 2},
				},
			},
			// 注意：map的遍历顺序是不确定的
			expected: "",
		},
		{
			name: "Chain call expression",
			node: &ChainCallExpr{
				Calls: []*CallExpr{
					{
						Function: &Identifier{Name: "print"},
						Args: []Expression{
							&String{Value: "hello"},
						},
					},
					{
						Function: &Identifier{Name: "upper"},
						Args:     []Expression{},
					},
				},
			},
			expected: "print(\"hello\").upper()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.node.String()

			if tt.name == "Dict literal" {
				// 对于字典，只检查基本格式
				if result[0] != '{' || result[len(result)-1] != '}' {
					t.Errorf("字典字符串格式错误: %s", result)
				}
				// 检查是否包含预期的键值对
				if !strings.Contains(result, "\"a\": 1") && !strings.Contains(result, "\"b\": 2") {
					t.Errorf("字典字符串不包含预期的键值对: %s", result)
				}
				if !(strings.Contains(result, "\"a\": 1") || strings.Contains(result, "\"b\": 2")) {
					t.Errorf("字典字符串格式错误: %s", result)
				}
			} else if result != tt.expected {
				t.Errorf("期望 %q, 得到 %q", tt.expected, result)
			}
		})
	}
}

func TestStatementString(t *testing.T) {
	tests := []struct {
		name     string
		stmt     Statement
		expected string
	}{
		{
			name: "Var declaration with init",
			stmt: &VarDecl{
				Name: &Identifier{Name: "x"},
				Type: "int",
				Expr: &Integer{Value: 10},
			},
			expected: "var x int = 10",
		},
		{
			name: "Var declaration without init",
			stmt: &VarDecl{
				Name: &Identifier{Name: "y"},
				Type: "auto",
			},
			expected: "var y auto",
		},
		{
			name: "Assign statement",
			stmt: &AssignStmt{
				Left: &Identifier{Name: "x"},
				Expr: &Integer{Value: 20},
			},
			expected: "x = 20",
		},
		{
			name: "If statement without else",
			stmt: &IfStmt{
				Condition: &Identifier{Name: "flag"},
				Then: &BlockStmt{
					Stmts: []Statement{
						&ExpressionStmt{
							Expr: &CallExpr{
								Function: &Identifier{Name: "print"},
								Args: []Expression{
									&String{Value: "true"},
								},
							},
						},
					},
				},
			},
			expected: "if flag {\n  print(\"true\")\n}",
		},
		{
			name: "While statement",
			stmt: &WhileStmt{
				Condition: &Identifier{Name: "i"},
				Body: &BlockStmt{
					Stmts: []Statement{
						&ExpressionStmt{
							Expr: &CallExpr{
								Function: &Identifier{Name: "print"},
								Args: []Expression{
									&String{Value: "looping"},
								},
							},
						},
					},
				},
			},
			expected: "while i {\n  print(\"looping\")\n}",
		},
		{
			name: "Return statement with value",
			stmt: &ReturnStmt{
				Expr: &Integer{Value: 42},
			},
			expected: "return 42",
		},
		{
			name:     "Return statement without value",
			stmt:     &ReturnStmt{},
			expected: "return",
		},
		{
			name: "For statement",
			stmt: &ForStmt{
				Init: &VarDecl{
					Name: &Identifier{Name: "i"},
					Type: "int",
					Expr: &Integer{Value: 0},
				},
				Cond: &BinaryExpr{
					Left:  &Identifier{Name: "i"},
					Op:    "<",
					Right: &Integer{Value: 5},
				},
				Post: &AssignStmt{
					Left: &Identifier{Name: "i"},
					Expr: &BinaryExpr{
						Left:  &Identifier{Name: "i"},
						Op:    "+",
						Right: &Integer{Value: 1},
					},
				},
				Body: &BlockStmt{
					Stmts: []Statement{
						&ExpressionStmt{
							Expr: &CallExpr{
								Function: &Identifier{Name: "print"},
								Args: []Expression{
									&Identifier{Name: "i"},
								},
							},
						},
					},
				},
			},
			expected: "for var i int = 0; (i < 5); i = (i + 1) {\n  print(i)\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.stmt.String()
			if result != tt.expected {
				t.Errorf("期望 %q, 得到 %q", tt.expected, result)
			}
		})
	}
}
