package src

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

const (
	ILLEGAL TokenType = "ILLEGAL" //Unknown token
	EOF     TokenType = "EOF"     // End Of file

	//identifiers & literals
	IDENTIFIER TokenType = "IDENTIFIER" // variable names, and function names
	NUMBER     TokenType = "NUMBER"     //123, 3.14
	STRING     TokenType = "STRING"     //hello

	//Operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	MODULO   TokenType = "%"

	BANG TokenType = "!"

	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="
	LT     TokenType = "<"
	GT     TokenType = ">"
	LTE    TokenType = "<="
	GTE    TokenType = ">="

	AND TokenType = "&&"
	OR  TokenType = "||"

	// Deliminaters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"

	LPARENT  TokenType = "("
	RPARENT  TokenType = ")"
	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"

	// keywords
	LET    TokenType = "LET"
	CST    TokenType = "CST"
	FN     TokenType = "FN"
	RETURN TokenType = "RETURN"
	IF     TokenType = "IF"
	ELSE   TokenType = "ELSE"
	WHILE  TokenType = "WHILE"
	FOR    TokenType = "FOR"
	TRUE   TokenType = "TRUE"
	FALSE  TokenType = "FALSE"
	NULL   TokenType = "NULL"
)

// keybwords maps keyword strings too thier token types
var keywords = map[string]TokenType{
	"let":    LET,
	"cst":    CST,
	"fn":     FN,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
	"for":    FOR,
	"true":   TRUE,
	"false":  FALSE,
	"null":   NULL,
}

// function LookUpIndentifier checks if identifier is a keyword
// If it is, return the keyword's TokenType, otherwise return IDENTIFIER
func LookupIdentifier(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENTIFIER
}

// Function NewToken creates a new token with line and column information
func NewToken(tokenType TokenType, literal string, line, column int) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
		Line:    line,
		Column:  column,
	}
}
