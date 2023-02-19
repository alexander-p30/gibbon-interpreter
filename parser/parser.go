package parser

import (
	"fmt"
	"gibbon/ast"
	"gibbon/lexer"
	"gibbon/token"
	"strconv"
)

type prefixParserFn func() ast.Expression
type infixParserFn func() ast.Expression

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
}

func (p *Parser) initializePrefixParsers() {
	p.prefixExpressionParser = make(map[token.TokenType]prefixParserFn)
	p.registerPrefixParser(token.IDENT, p.parseIdentifier)
	p.registerPrefixParser(token.INT, p.parseIntegerLiteral)
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

func (p *Parser) parseExpression(precedence int) ast.Expression {
	parserFn := p.prefixExpressionParser[p.currentToken.Type]

	if parserFn == nil {
		return nil
	}

	return parserFn()
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

