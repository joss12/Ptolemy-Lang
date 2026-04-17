package src

import (
	"fmt"
	"strconv"
)

const (
	_ int = iota
	PREC_LOWEST
	PREC_ASSIGN
	PREC_OR
	PREC_AND
	PREC_EQUALS
	PREC_LESSGREATER
	PREC_SUM
	PREC_PRODUCT
	PREC_PREFIX
	PREC_CALL
	PREC_INDEX
)

// precedences maps token types to their precedence level
var precedences = map[TokenType]int{
	ASSIGN:   PREC_ASSIGN,
	OR:       PREC_OR,
	AND:      PREC_AND,
	EQ:       PREC_EQUALS,
	NOT_EQ:   PREC_EQUALS,
	LT:       PREC_LESSGREATER,
	GT:       PREC_LESSGREATER,
	LTE:      PREC_LESSGREATER,
	GTE:      PREC_LESSGREATER,
	PLUS:     PREC_SUM,
	MINUS:    PREC_SUM,
	SLASH:    PREC_PRODUCT,
	ASTERISK: PREC_PRODUCT,
	MODULO:   PREC_PRODUCT,
	LPARENT:  PREC_CALL,
	LBRACKET: PREC_INDEX,
}

// Parser takes tokens from lexer and builds AST
type Parser struct {
	lexer  *Lexer
	errors []string

	curToken  Token
	peekToken Token

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

// Function types for Pratt parsing
type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

// NewParser creates a new Parser
func NewParser(l *Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}

	// Initialize prefix parse functions
	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(NUMBER, p.parseNumberLiteral)
	p.registerPrefix(STRING, p.parseStringLiteral)
	p.registerPrefix(TRUE, p.parseBooleanLiteral)
	p.registerPrefix(FALSE, p.parseBooleanLiteral)
	p.registerPrefix(NULL, p.parseNullLiteral)
	p.registerPrefix(BANG, p.parsePrefixExpression)
	p.registerPrefix(MINUS, p.parsePrefixExpression)
	p.registerPrefix(LPARENT, p.parseGroupedExpression)
	p.registerPrefix(LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(FN, p.parseFunctionLiteral)
	p.registerPrefix(IF, p.parseIfExpression)

	// Initialize infix parse functions
	p.infixParseFns = make(map[TokenType]infixParseFn)
	p.registerInfix(ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(PLUS, p.parseInfixExpression)
	p.registerInfix(MINUS, p.parseInfixExpression)
	p.registerInfix(SLASH, p.parseInfixExpression)
	p.registerInfix(ASTERISK, p.parseInfixExpression)
	p.registerInfix(MODULO, p.parseInfixExpression)
	p.registerInfix(EQ, p.parseInfixExpression)
	p.registerInfix(NOT_EQ, p.parseInfixExpression)
	p.registerInfix(LT, p.parseInfixExpression)
	p.registerInfix(GT, p.parseInfixExpression)
	p.registerInfix(LTE, p.parseInfixExpression)
	p.registerInfix(GTE, p.parseInfixExpression)
	p.registerInfix(AND, p.parseInfixExpression)
	p.registerInfix(OR, p.parseInfixExpression)
	p.registerInfix(LPARENT, p.parseCallExpression)
	p.registerInfix(LBRACKET, p.parseIndexExpression)

	// Read two tokens so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// Errors returns the parser's errors
func (p *Parser) Errors() []string {
	return p.errors
}

// nextToken advances both curToken and peekToken
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// curTokenIs checks if current token is of given type
func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs checks if peek token is of given type
func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek advances if peek token matches expected type
func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// peekPrecedence returns the precedence of peek token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return PREC_LOWEST
}

// curPrecedence returns the precedence of current token
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return PREC_LOWEST
}

// peekError adds an error for unexpected peek token
func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead at line %d",
		t, p.peekToken.Type, p.peekToken.Line)
	p.errors = append(p.errors, msg)
}

// noPrefixParseFnError adds an error when no prefix parse function is found
func (p *Parser) noPrefixParseFnError(t TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found at line %d",
		t, p.curToken.Line)
	p.errors = append(p.errors, msg)
}

// registerPrefix registers a prefix parse function
func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix registers an infix parse function
func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// ParseProgram parses the entire program
func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for !p.curTokenIs(EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement parses a single statement
func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case LET:
		return p.parseLetStatement()
	case CST:
		return p.parseConstStatement()
	case RETURN:
		return p.parseReturnStatement()
	case WHILE:
		return p.parseWhileStatement()
	case FOR:
		return p.parseForStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement parses: let x = 5
func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Token: p.curToken}

	if !p.expectPeek(IDENTIFIER) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(PREC_LOWEST)

	// Optional semicolon
	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseConstStatement parses: cst PI = 3.14
func (p *Parser) parseConstStatement() Statement {
	stmt := &ConstStatement{Token: p.curToken}

	if !p.expectPeek(IDENTIFIER) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(PREC_LOWEST)

	// Optional semicolon
	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement parses: return x + 5
func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(PREC_LOWEST)

	// Optional semicolon
	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpressionStatement parses expressions used as statements
func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(PREC_LOWEST)

	// Optional semicolon
	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseWhileStatement parses: while (condition) { body }
func (p *Parser) parseWhileStatement() *WhileStatement {
	stmt := &WhileStatement{Token: p.curToken}

	if !p.expectPeek(LPARENT) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(PREC_LOWEST)

	if !p.expectPeek(RPARENT) {
		return nil
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseForStatement parses: for (init; condition; increment) { body }
func (p *Parser) parseForStatement() *ForStatement {
	stmt := &ForStatement{Token: p.curToken}

	if !p.expectPeek(LPARENT) {
		return nil
	}

	// Parse init (optional)
	// FIX: peek ahead instead of consuming - avoids semicolon conflict
	p.nextToken()
	if !p.curTokenIs(SEMICOLON) {
		stmt.Init = p.parseStatement()
		// parseStatement may consume the semicolon already, check if we need to advance
		if p.peekTokenIs(SEMICOLON) {
			p.nextToken()
		}
	}

	// Parse condition (optional)
	p.nextToken()
	if !p.curTokenIs(SEMICOLON) {
		stmt.Condition = p.parseExpression(PREC_LOWEST)
	}

	if !p.expectPeek(SEMICOLON) {
		return nil
	}

	// Parse increment (optional)
	p.nextToken()
	if !p.curTokenIs(RPARENT) {
		stmt.Increment = p.parseExpression(PREC_LOWEST)
	}

	if !p.expectPeek(RPARENT) {
		return nil
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseBlockStatement parses: { stmt1; stmt2; }
func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.curToken}
	block.Statements = []Statement{}

	p.nextToken()

	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseExpression is the main expression parser using Pratt parsing
func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

// parseIdentifier parses an identifier
func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parseNumberLiteral parses a number
func (p *Parser) parseNumberLiteral() Expression {
	lit := &NumberLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as number", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

// parseStringLiteral parses a string
func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parseBooleanLiteral parses true or false
func (p *Parser) parseBooleanLiteral() Expression {
	return &BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(TRUE)}
}

// parseNullLiteral parses null
func (p *Parser) parseNullLiteral() Expression {
	return &NullLiteral{Token: p.curToken}
}

// parsePrefixExpression parses: -5 or !true
func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREC_PREFIX)

	return expression
}

// parseInfixExpression parses: 5 + 3, x > y, etc.
func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// parseAssignmentExpression parses: x = 10
func (p *Parser) parseAssignmentExpression(left Expression) Expression {
	expr := &AssignmentExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()
	expr.Value = p.parseExpression(PREC_LOWEST)

	return expr
}

// parseGroupedExpression parses: (5 + 3)
func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()

	exp := p.parseExpression(PREC_LOWEST)

	if !p.expectPeek(RPARENT) {
		return nil
	}

	return exp
}

// parseArrayLiteral parses: [1, 2, 3]
func (p *Parser) parseArrayLiteral() Expression {
	array := &ArrayLiteral{Token: p.curToken}

	array.Elements = p.parseExpressionList(RBRACKET)

	return array
}

// parseIndexExpression parses: arr[0]
func (p *Parser) parseIndexExpression(left Expression) Expression {
	exp := &IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(PREC_LOWEST)

	if !p.expectPeek(RBRACKET) {
		return nil
	}

	return exp
}

// parseCallExpression parses: add(5, 3)
func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(RPARENT)
	return exp
}

// parseExpressionList parses comma-separated expressions
func (p *Parser) parseExpressionList(end TokenType) []Expression {
	list := []Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(PREC_LOWEST))

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(PREC_LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// parseFunctionLiteral parses: fn(x, y) { return x + y }
func (p *Parser) parseFunctionLiteral() Expression {
	lit := &FunctionLiteral{Token: p.curToken}

	// Handle optional function name: fn fibonacci(n) { }
	if p.peekTokenIs(IDENTIFIER) {
		p.nextToken()
		lit.Name = p.curToken.Literal
	}

	if !p.expectPeek(LPARENT) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	// If named, store in environment so it can call itself (recursion)
	return lit
}

// parseFunctionParameters parses: (x, y, z)
func (p *Parser) parseFunctionParameters() []*Identifier {
	identifiers := []*Identifier{}

	if p.peekTokenIs(RPARENT) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(RPARENT) {
		return nil
	}

	return identifiers
}

// parseIfExpression parses: if (condition) { consequence } else { alternative }
func (p *Parser) parseIfExpression() Expression {
	expression := &IfStatement{Token: p.curToken}

	if !p.expectPeek(LPARENT) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(PREC_LOWEST)

	if !p.expectPeek(RPARENT) {
		return nil
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	// Handle else and else if
	if p.peekTokenIs(ELSE) {
		p.nextToken()

		// Check if it's "else if" instead of just "else"
		if p.peekTokenIs(IF) {
			p.nextToken()

			// Parse the if statement recursively
			elseIfStmt := p.parseIfExpression()

			// Wrap it in a block statement
			expression.Alternative = &BlockStatement{
				Token: p.curToken,
				Statements: []Statement{
					&ExpressionStatement{
						Token:      p.curToken,
						Expression: elseIfStmt,
					},
				},
			}
		} else {
			// Regular else block
			if !p.expectPeek(LBRACE) {
				return nil
			}

			expression.Alternative = p.parseBlockStatement()
		}
	}

	return expression
}
