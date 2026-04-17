package src

// lexer reads source code and produces tokens
type Lexer struct {
	input        string
	position     int  // current position input (points to current char)
	readPosition int  // current reading position in input (after curent char)
	ch           byte // current char under examination
	line         int  // current line number
	column       int  //current column number
}

// NewLexer creates a new lexer instance
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar() // Initialize by reading the first character
	return l
}

// readChar reads the next character and advances position in the input
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for "NUL" - signals end of input
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++

	//Track line numbers
	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
}

// peekChar looks at the next character without advancing position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	// Save current position for token creation
	line := l.line
	column := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = NewToken(EQ, string(ch)+string(l.ch), line, column)
		} else {
			tok = NewToken(ASSIGN, string(l.ch), line, column)
		}
	case '+':
		tok = NewToken(PLUS, string(l.ch), line, column)
	case '-':
		tok = NewToken(MINUS, string(l.ch), line, column)
	case '*':
		tok = NewToken(ASTERISK, string(l.ch), line, column)
	case '/':
		//Check for comments
		if l.peekChar() == '/' {
			l.skipComment()
			return l.NextToken() // Skip comment and get next real token
		}
		tok = NewToken(SLASH, string(l.ch), line, column)
	case '%':
		tok = NewToken(MODULO, string(l.ch), line, column)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = NewToken(NOT_EQ, string(ch)+string(l.ch), line, column)
		} else {
			tok = NewToken(BANG, string(l.ch), line, column)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = NewToken(LTE, string(ch)+string(l.ch), line, column)
		} else {
			tok = NewToken(LT, string(l.ch), line, column)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = NewToken(GTE, string(ch)+string(l.ch), line, column)
		} else {
			tok = NewToken(GT, string(l.ch), line, column)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = NewToken(AND, string(ch)+string(l.ch), line, column)
		} else {
			tok = NewToken(ILLEGAL, string(l.ch), line, column)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = NewToken(OR, string(ch)+string(l.ch), line, column)
		} else {
			tok = NewToken(ILLEGAL, string(l.ch), line, column)
		}
	case ',':
		tok = NewToken(COMMA, string(l.ch), line, column)
	case ';':
		tok = NewToken(SEMICOLON, string(l.ch), line, column)
	case '(':
		tok = NewToken(LPARENT, string(l.ch), line, column)
	case ')':
		tok = NewToken(RPARENT, string(l.ch), line, column)
	case '{':
		tok = NewToken(LBRACE, string(l.ch), line, column)
	case '}':
		tok = NewToken(RBRACE, string(l.ch), line, column)
	case '[':
		tok = NewToken(LBRACKET, string(l.ch), line, column)
	case ']':
		tok = NewToken(RBRACKET, string(l.ch), line, column)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		tok.Line = line
		tok.Column = column
	case 0:
		tok = NewToken(EOF, "", line, column)
	default:
		if isLetter(l.ch) {
			tok.Line = line
			tok.Column = column
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdentifier(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Line = line
			tok.Column = column
			tok.Literal = l.readNumber()
			tok.Type = NUMBER
			return tok //Early return because readNumber already advanced
		} else {
			tok = NewToken(ILLEGAL, string(l.ch), line, column)
		}
	}
	l.readChar()
	return tok
}

// skipWhitespace skips whitespace characters (space, tab, newline, carriage return)
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// skipComment skips single-line comments (// ...)
func (l *Lexer) skipComment() {
	// Skip the tow slashes
	l.readChar()
	l.readChar()

	// Read until end of line or end of file
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

// readIdentifier reads an identifier (variable name, function name, keyword)
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number (interger or float)
func (l *Lexer) readNumber() string {
	position := l.position

	//Read digits
	for isDigit(l.ch) {
		l.readChar()
	}

	//Check for decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar() // consume '.'

		// Read fractional part
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

// readString reads a string literal (everything between quotes)
func (l *Lexer) readString() string {
	position := l.position + 1 // Skip openning quotes

	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// isLetter checks if a characters is a letter or underscore
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
