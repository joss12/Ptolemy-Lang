package tests

import (
	"github.com/joss12/ptolemy-lang/src"
	"testing"
)

func TestTokenTypes(t *testing.T) {
	tests := []struct {
		input    string
		expected src.TokenType
	}{
		{"let", src.LET},
		{"fn", src.FN},
		{"myVar", src.IDENTIFIER},
		{"return", src.RETURN},
	}

	for _, tt := range tests {
		result := src.LookupIdentifier(tt.input)
		if result != tt.expected {
			t.Errorf("LookupIdentifier(%q) = %q, expected %q",
				tt.input, result, tt.expected)
		}
	}
}

func TestNewToken(t *testing.T) {
	token := src.NewToken(src.NUMBER, "42", 1, 5)

	if token.Type != src.NUMBER {
		t.Errorf("Expected type NUMBER, got %q", token.Type)
	}
	if token.Literal != "42" {
		t.Errorf("Expected literal '4', got %q", token.Literal)
	}
	if token.Line != 1 {
		t.Errorf("Expected line 1, got %d", token.Line)
	}
	if token.Column != 5 {
		t.Errorf("Expected column 5, got %d", token.Column)
	}
}
