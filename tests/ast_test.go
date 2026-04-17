package tests

import (
	"testing"

	"github.com/joss12/ptolemy-lang/src"
)

func TestStrict(t *testing.T) {
	program := &src.Program{
		Statements: []src.Statement{
			&src.LetStatement{
				Token: src.Token{Type: src.LET, Literal: "let"},
				Name: &src.Identifier{
					Token: src.Token{Type: src.IDENTIFIER, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &src.Identifier{
					Token: src.Token{Type: src.IDENTIFIER, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

func TestInfixExpression(t *testing.T) {
	expr := src.InfixExpression{
		Token: src.Token{Type: src.PLUS, Literal: "+"},
		Left: &src.NumberLiteral{
			Token: src.Token{Type: src.NUMBER, Literal: "5"},
			Value: 5,
		},
		Operator: "+",
		Right: &src.NumberLiteral{
			Token: src.Token{Type: src.NUMBER, Literal: "3"},
			Value: 3,
		},
	}
	if expr.String() != "(5 + 3)" {
		t.Errorf("expr.String() wrong. got=%q", expr.String())
	}
}

func TestArrayLiteral(t *testing.T) {
	arr := &src.ArrayLiteral{
		Token: src.Token{Type: src.LBRACE, Literal: "["},
		Elements: []src.Expression{
			&src.NumberLiteral{
				Token: src.Token{Type: src.NUMBER, Literal: "1"},
				Value: 1,
			},
			&src.NumberLiteral{
				Token: src.Token{Type: src.NUMBER, Literal: "2"},
				Value: 2,
			},
			&src.NumberLiteral{
				Token: src.Token{Type: src.NUMBER, Literal: "3"},
				Value: 3,
			},
		},
	}
	if arr.String() != "[1, 2, 3]" {
		t.Errorf("arr.String() wrong. got=%q", arr.String())
	}
}

func TestFunctionLiteral(t *testing.T) {
	fn := &src.FunctionLiteral{
		Token: src.Token{Type: src.FN, Literal: "fn"},
		Parameters: []*src.Identifier{
			{
				Token: src.Token{Type: src.IDENTIFIER, Literal: "x"},
				Value: "x",
			},
			{
				Token: src.Token{Type: src.IDENTIFIER, Literal: "y"},
				Value: "y",
			},
		},
		Body: &src.BlockStatement{
			Token: src.Token{Type: src.LBRACE, Literal: "{"},
			Statements: []src.Statement{
				&src.ReturnStatement{
					Token: src.Token{Type: src.RETURN, Literal: "return"},
					ReturnValue: &src.InfixExpression{
						Token: src.Token{Type: src.PLUS, Literal: "+"},
						Left: &src.Identifier{
							Token: src.Token{Type: src.IDENTIFIER, Literal: "x"},
							Value: "x",
						},
						Operator: "+",
						Right: &src.Identifier{
							Token: src.Token{Type: src.IDENTIFIER, Literal: "y"},
							Value: "y",
						},
					},
				},
			},
		},
	}
	expected := "fn(x, y) return (x + y);"
	if fn.String() != expected {
		t.Errorf("fn.String() wrong. go=%q, expected=%q", fn.String(), expected)
	}
}
