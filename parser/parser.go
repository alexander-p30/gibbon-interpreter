package parser

import (
	"fmt"
	"gibbon/ast"
	"gibbon/lexer"
	"gibbon/token"
	"strconv"
)

type prefixParserFn func() ast.Expression
type infixParserFn func(left ast.Expression) ast.Expression

type Parser struct {
	lexer                  *lexer.Lexer
	currentToken           token.Token
	peekToken              token.Token
	errors                 []Error
	infixExpressionParser  map[token.TokenType]infixParserFn
	prefixExpressionParser map[token.TokenType]prefixParserFn
}

type Error struct {
	message  string
	location token.TokenLocation
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: []Error{}}
	// initializing expression parsers
	p.initializeInfixParsers()
	p.initializePrefixParsers()
	// Setting both peek token and current token
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) initializeInfixParsers() {
	p.infixExpressionParser = make(map[token.TokenType]infixParserFn)
	p.registerInfixParser(token.EQUAL, p.parseInfixOperator)
	p.registerInfixParser(token.DIFFERENT, p.parseInfixOperator)
	p.registerInfixParser(token.LT, p.parseInfixOperator)
	p.registerInfixParser(token.GT, p.parseInfixOperator)
	p.registerInfixParser(token.LTE, p.parseInfixOperator)
	p.registerInfixParser(token.GTE, p.parseInfixOperator)
	p.registerInfixParser(token.PLUS, p.parseInfixOperator)
	p.registerInfixParser(token.MINUS, p.parseInfixOperator)
	p.registerInfixParser(token.ASTERISK, p.parseInfixOperator)
	p.registerInfixParser(token.SLASH, p.parseInfixOperator)
}

func (p *Parser) initializePrefixParsers() {
	p.prefixExpressionParser = make(map[token.TokenType]prefixParserFn)
	p.registerPrefixParser(token.IDENT, p.parseIdentifier)
	p.registerPrefixParser(token.INT, p.parseIntegerLiteral)
	p.registerPrefixParser(token.BANG, p.parsePrefixOperator)
	p.registerPrefixParser(token.MINUS, p.parsePrefixOperator)
	p.registerPrefixParser(token.PLUS, p.parsePrefixOperator)
}

func (p *Parser) registerPrefixParser(t token.TokenType, parser prefixParserFn) {
	p.prefixExpressionParser[t] = parser
}

func (p *Parser) registerInfixParser(t token.TokenType, parser infixParserFn) {
	p.infixExpressionParser[t] = parser
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)

	if err != nil {
		msg := fmt.Sprintf("Could not parse '%q' as integer literal", p.currentToken.Literal)
		p.errors = append(p.errors, Error{message: msg, location: p.currentToken.Location})
		return nil
	}

	return &ast.IntegerLiteral{Token: p.currentToken, Value: value}
}

func (p *Parser) parsePrefixOperator() ast.Expression {
	operatorToken := p.currentToken
	p.nextToken()
	right := p.parseExpression(PREFIX)
	return &ast.PrefixExpression{Token: operatorToken, Operator: operatorToken.Literal, Right: right}
}

func (p *Parser) parseInfixOperator(left ast.Expression) ast.Expression {
	operatorToken := p.currentToken
	p.nextToken()
	right := p.parseExpression(p.getTokenPrecedence(operatorToken))
	return &ast.InfixExpression{Token: operatorToken, Operator: operatorToken.Literal, Left: left, Right: right}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) Errors() []Error {
	return p.errors
}

func (p *Parser) peekError(expected token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", expected, p.peekToken.Type)
	p.errors = append(p.errors, Error{message: msg, location: p.peekToken.Location})
}

func (p *Parser) noPrefixParserFnError(t token.TokenType) {
	msg := fmt.Sprintf("token type %q has no registered PREFIX parser functions", t)
	p.errors = append(p.errors, Error{message: msg, location: p.currentToken.Location})
}

func (p *Parser) noInfixParserFnError(t token.TokenType) {
	msg := fmt.Sprintf("token type %q has no registered INFIX parser functions", t)
	p.errors = append(p.errors, Error{message: msg, location: p.currentToken.Location})
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token:      p.currentToken,
		Expression: p.parseExpression(LOWEST),
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL // method calls
)

var precedences = map[token.TokenType]int{
	token.EQUAL:     EQUALS,
	token.DIFFERENT: EQUALS,
	token.LTE:       LESSGREATER,
	token.GTE:       LESSGREATER,
	token.LT:        LESSGREATER,
	token.GT:        LESSGREATER,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.ASTERISK:  PRODUCT,
	token.SLASH:     PRODUCT,
}

func (p *Parser) getTokenPrecedence(t token.Token) int {
	if precedence, ok := precedences[t.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixParserFn := p.prefixExpressionParser[p.currentToken.Type]

	if prefixParserFn == nil {
		p.noPrefixParserFnError(p.currentToken.Type)
		return nil
	}

	leftExp := prefixParserFn()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.getTokenPrecedence(p.peekToken) {
		infixParserFn := p.infixExpressionParser[p.peekToken.Type]

		if infixParserFn == nil {
			p.noInfixParserFnError(p.peekToken.Type)
			return nil
		}

		p.nextToken()
		leftExp = infixParserFn(leftExp)
	}

	return leftExp
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}
