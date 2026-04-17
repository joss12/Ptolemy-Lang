package tests

import (
	"testing"

	"github.com/joss12/ptolemy-lang/src"
)

func TestNextToken(t *testing.T) {
	input := `let x = 5 + 10

	cst PI = 3.14
	fn add(a, b){
		return a + b
	}

	//This is a comment
	let result = add(x, 20)

	if (result > 10){
		print("big")
	}else{
		print("small")
	}

	let arr = [1, 2, 3]
	let name = "BornToShine"
	`
	tests := []struct {
		expectedType    src.TokenType
		expectedLiteral string
	}{
		{src.LET, "let"},
		{src.IDENTIFIER, "x"},
		{src.ASSIGN, "="},
		{src.NUMBER, "5"},
		{src.PLUS, "+"},
		{src.NUMBER, "10"},

		{src.CST, "cst"},
		{src.IDENTIFIER, "PI"},
		{src.ASSIGN, "="},
		{src.NUMBER, "3.14"},

		{src.FN, "fn"},
		{src.IDENTIFIER, "add"},
		{src.LPARENT, "("},
		{src.IDENTIFIER, "a"},
		{src.COMMA, ","},
		{src.IDENTIFIER, "b"},
		{src.RPARENT, ")"},
		{src.LBRACE, "{"},
		{src.RETURN, "return"},
		{src.IDENTIFIER, "a"},
		{src.PLUS, "+"},
		{src.IDENTIFIER, "b"},
		{src.RBRACE, "}"},

		{src.LET, "let"},
		{src.IDENTIFIER, "result"},
		{src.ASSIGN, "="},
		{src.IDENTIFIER, "add"},
		{src.LPARENT, "("},
		{src.IDENTIFIER, "x"},
		{src.COMMA, ","},
		{src.NUMBER, "20"},
		{src.RPARENT, ")"},

		{src.IF, "if"},
		{src.LPARENT, "("},
		{src.IDENTIFIER, "result"},
		{src.GT, ">"},
		{src.NUMBER, "10"},
		{src.RPARENT, ")"},
		{src.LBRACE, "{"},
		{src.IDENTIFIER, "print"},
		{src.LPARENT, "("},
		{src.STRING, "big"},
		{src.RPARENT, ")"},
		{src.RBRACE, "}"},
		{src.ELSE, "else"},
		{src.LBRACE, "{"},
		{src.IDENTIFIER, "print"},
		{src.LPARENT, "("},
		{src.STRING, "small"},
		{src.RPARENT, ")"},
		{src.RBRACE, "}"},

		{src.LET, "let"},
		{src.IDENTIFIER, "arr"},
		{src.ASSIGN, "="},
		{src.LBRACKET, "["},
		{src.NUMBER, "1"},
		{src.COMMA, ","},
		{src.NUMBER, "2"},
		{src.COMMA, ","},
		{src.NUMBER, "3"},
		{src.RBRACKET, "]"},

		{src.LET, "let"},
		{src.IDENTIFIER, "name"},
		{src.ASSIGN, "="},
		{src.STRING, "BornToShine"},

		{src.EOF, ""},
	}

	l := src.NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestOperators(t *testing.T) {
	input := `== != < > <= >= && ||`

	tests := []struct {
		expectedType    src.TokenType
		expectedLiteral string
	}{
		{src.EQ, "=="},
		{src.NOT_EQ, "!="},
		{src.LT, "<"},
		{src.GT, ">"},
		{src.LTE, "<="},
		{src.GTE, ">="},
		{src.AND, "&&"},
		{src.OR, "||"},
		{src.EOF, ""},
	}

	l := src.NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestComments(t *testing.T) {
	input := `let x = 5 //this is a comment
	let y = 10
	`

	tests := []struct {
		expectedType    src.TokenType
		expectedLiteral string
	}{
		{src.LET, "let"},
		{src.IDENTIFIER, "x"},
		{src.ASSIGN, "="},
		{src.NUMBER, "5"},
		{src.LET, "let"},
		{src.IDENTIFIER, "y"},
		{src.ASSIGN, "="},
		{src.NUMBER, "10"},
		{src.EOF, ""}}

	l := src.NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
	}
}
