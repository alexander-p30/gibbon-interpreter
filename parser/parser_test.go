package parser

import (
	"bytes"
	"fmt"
	"gibbon/ast"
	"gibbon/lexer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLetStatementParse(t *testing.T) {
	assert := assert.New(t)

	input := []byte(`
let a = 1;
let b = +47138471;
let something = -47194738292;
`)

	l := lexer.NewLexer(bytes.NewReader(input), "input")
	parser := NewParser(l)
	program := parser.ParseProgram()
	ensureNoErrors(t, parser)

	assert.NotNil(program, "ParseProgram() returned nil")

	assert.Len(program.Statements, 3, "program.Statements does not contain 3 statements")

	tests := []struct {
		expectedIdentifier string
	}{
		{"a"},
		{"b"},
		{"something"},
	}

	for i, test := range tests {
		if !testLetStatement(t, program.Statements[i], test.expectedIdentifier) {
			return
		}
	}
}

func TestLetStatementParseError(t *testing.T) {
	assert := assert.New(t)

	input := []byte(`
let a = 1;
let a = 1
let a = ;
let a 1;
let 1;
`)

	l := lexer.NewLexer(bytes.NewReader(input), "input")
	parser := NewParser(l)
	parser.ParseProgram()
	errors := parser.Errors()

	if !assert.Len(errors, 2) {
		for i, err := range errors {
			t.Errorf("Error [%d] %s:%d:%d: \"%s\"",
				i,
				parser.lexer.SourceFile(),
				err.location.Line,
				err.location.FirstCharIndex,
				err.message,
			)
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	assert := assert.New(t)

	input := []byte(`someIdentifier;`)
	l := lexer.NewLexer(bytes.NewReader(input), "input")
	parser := NewParser(l)
	program := parser.ParseProgram()
	ensureNoErrors(t, parser)

	stmts := program.Statements

	assert.Len(stmts, 1)

	statement := program.Statements[0]
	identifierExpression, ok := statement.(*ast.ExpressionStatement)
	if !assert.True(ok, "statement not of type *ast.ExpressionStatement") {
		assert.FailNow("")
	}

	identifier, ok := identifierExpression.Expression.(*ast.Identifier)
	if !assert.True(ok, "statement not of type *ast.Identifier") {
		assert.FailNow("")
	}

	assert.Equal(identifier.Value, "someIdentifier")
	assert.Equal(identifier.TokenLiteral(), "someIdentifier")
}
func testLetStatement(t *testing.T, stmt ast.Statement, expectedIdentifier string) bool {
	assert := assert.New(t)

	assert.Equal(stmt.TokenLiteral(), "let", "s.TokenLiteral not 'let'")

	letStmt, ok := stmt.(*ast.LetStatement)
	assert.Truef(ok, "s not *ast.LetStatement. got: %T", stmt)

	assert.Equalf(letStmt.Name.Value, expectedIdentifier, "letStmt.Name.Value not '%s'", expectedIdentifier)

	assert.Equalf(letStmt.Name.TokenLiteral(), expectedIdentifier, "letStmt.Name.TokenLiteral() not '%s'", expectedIdentifier)

	return true
}

func ensureNoErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))

	for i, err := range errors {
		t.Errorf(
			"Parser error [%d] at %s:%d:%d: \"%s\"",
			i,
			p.lexer.SourceFile(),
			err.location.Line,
			err.location.FirstCharIndex,
			err.message,
		)
	}

	t.FailNow()
}
