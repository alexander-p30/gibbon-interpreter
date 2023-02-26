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

func TestIntegerLiteral(t *testing.T) {
	assert := assert.New(t)

	input := []byte(`3;`)
	l := lexer.NewLexer(bytes.NewReader(input), "input")
	parser := NewParser(l)
	program := parser.ParseProgram()
	ensureNoErrors(t, parser)

	stmts := program.Statements

	assert.Len(stmts, 1)

	statement := program.Statements[0]
	integerExpression, ok := statement.(*ast.ExpressionStatement)
	if !assert.True(ok, "statement not of type *ast.ExpressionStatement") {
		assert.FailNow("")
	}

	testIntegerLiteral(t, integerExpression.Expression, 3)
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		input                []byte
		expectedOperator     string
		expectedIntegerValue int64
	}{
		{[]byte(`-5`), "-", 5},
		{[]byte(`!2`), "!", 2},
		{[]byte(`-7;`), "-", 7},
		{[]byte(`+99182346;`), "+", 99182346},
	}

	for _, test := range tests {
		assert := assert.New(t)

		l := lexer.NewLexer(bytes.NewReader(test.input), "input")
		parser := NewParser(l)
		program := parser.ParseProgram()

		ensureNoErrors(t, parser)

		assert.Len(program.Statements, 1)

		statement := program.Statements[0]
		stmt, ok := statement.(*ast.ExpressionStatement)
		if !assert.True(ok, "statement not of type *ast.ExpressionStatement") {
			assert.FailNow("")
		}

		prefixExpression, ok := stmt.Expression.(*ast.PrefixExpression)
		if !assert.True(ok, "expression not of type *ast.PrefixExpression") {
			assert.FailNow("")
		}

		assert.Equal(prefixExpression.Operator, test.expectedOperator)

		testIntegerLiteral(t, prefixExpression.Right, test.expectedIntegerValue)
	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input              []byte
		expectedOperator   string
		expectedLeftValue  int64
		expectedRightValue int64
	}{
		{[]byte(`3 - 2`), "-", 3, 2},
		{[]byte(`7 + 9`), "+", 7, 9},
		{[]byte(`1 * 16`), "*", 1, 16},
		{[]byte(`16 / 8`), "/", 16, 8},
		{[]byte(`30 > 60`), ">", 30, 60},
		{[]byte(`90 < 80`), "<", 90, 80},
		{[]byte(`1 == 0`), "==", 1, 0},
		{[]byte(`0 >= 20`), ">=", 0, 20},
		{[]byte(`4 <= 10`), "<=", 4, 10},
		{[]byte(`5 != 29`), "!=", 5, 29},
	}

	for _, test := range tests {
		assert := assert.New(t)

		l := lexer.NewLexer(bytes.NewReader(test.input), "input")
		parser := NewParser(l)
		program := parser.ParseProgram()

		ensureNoErrors(t, parser)

		assert.Len(program.Statements, 1)

		statement := program.Statements[0]
		stmt, ok := statement.(*ast.ExpressionStatement)
		if !assert.True(ok, "statement not of type *ast.ExpressionStatement") {
			assert.FailNow("")
		}

		testInfixExpression(
			t,
			stmt.Expression,
			test.expectedLeftValue,
			test.expectedOperator,
			test.expectedRightValue,
		)
	}
}

func TestExpressionOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input          []byte
		expectedString string
	}{
		{[]byte(`3 - 2`), "(3 - 2)"},
		{[]byte(`a + b + c`), "((a + b) + c)"},
		{[]byte(`a + b * c`), "(a + (b * c))"},
		{[]byte(`a * b / c + d`), "(((a * b) / c) + d)"},
		{[]byte(`!a <= b`), "((!a) <= b)"},
		{[]byte(`3 + 4 * 5 == 3 * 1 + 4 * 5`), "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{[]byte(`!a != !b`), "((!a) != (!b))"},
		{[]byte(`!a >= !b`), "((!a) >= (!b))"},
		{[]byte(`!a > !b < c`), "(((!a) > (!b)) < c)"},
	}

	for _, test := range tests {
		assert := assert.New(t)

		l := lexer.NewLexer(bytes.NewReader(test.input), "input")
		parser := NewParser(l)
		program := parser.ParseProgram()

		ensureNoErrors(t, parser)

		assert.Len(program.Statements, 1)

		statement := program.Statements[0]
		stmt, ok := statement.(*ast.ExpressionStatement)
		if !assert.True(ok, "statement not of type *ast.ExpressionStatement") {
			assert.FailNow("")
		}

		infixExpression, ok := stmt.Expression.(*ast.InfixExpression)
		if !assert.Truef(ok, "expression not of type *ast.InfixExpression %+v", stmt) {
			assert.FailNow("")
		}

		if !assert.Equal(infixExpression.String(), test.expectedString) {
			assert.FailNow("")
		}
	}
}

func TestReturnStatements(t *testing.T) {
	assert := assert.New(t)

	input := []byte(`
return 5;
return 10;
return 993322;
`)
	l := lexer.NewLexer(bytes.NewReader(input), "input")
	parser := NewParser(l)
	program := parser.ParseProgram()
	ensureNoErrors(t, parser)

	stmts := program.Statements

	assert.Len(stmts, 3, "Program statements do not contain 3 statemens.")

	for _, stmt := range stmts {
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !assert.Truef(ok, "stmt not *ast.ReturnStatement. got=%T", stmt) {
			continue
		}

		assert.Equal(returnStmt.TokenLiteral(), "return", "returnStmt.TokenLiteral not 'return'")
	}
}

// ------HELPERS------

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

func testIntegerLiteral(t *testing.T, exp ast.Expression, expectedValue int64) bool {
	assert := assert.New(t)

	integer, ok := exp.(*ast.IntegerLiteral)
	if !assert.True(ok, "statement not of type *ast.IntegerLiteral") {
		assert.FailNow("")
	}

	assert.Equal(integer.Value, expectedValue)
	return assert.Equal(integer.TokenLiteral(), fmt.Sprintf("%d", expectedValue))
}

func testLetStatement(t *testing.T, stmt ast.Statement, expectedIdentifier string) bool {
	assert := assert.New(t)

	assert.Equal(stmt.TokenLiteral(), "let", "s.TokenLiteral not 'let'")

	letStmt, ok := stmt.(*ast.LetStatement)
	// Prevent breaking on next lines
	if !assert.Truef(ok, "s not *ast.LetStatement. got: %T", stmt) {
		assert.FailNow("")
	}

	assert.Equalf(letStmt.Name.Value, expectedIdentifier, "letStmt.Name.Value not '%s'", expectedIdentifier)

	assert.Equalf(letStmt.Name.TokenLiteral(), expectedIdentifier, "letStmt.Name.TokenLiteral() not '%s'", expectedIdentifier)

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, expectedIdentifier string) bool {
	assert := assert.New(t)

	identifier, ok := exp.(*ast.Identifier)

	if !assert.Truef(ok, "Expression not of type *ast.Identifier, instead it is %T", identifier) {
		assert.FailNow("")
	}

	if !assert.Equal(expectedIdentifier, identifier.Value) {
		assert.FailNow("")
	}

	if !assert.Equal(expectedIdentifier, identifier.TokenLiteral()) {
		assert.FailNow("")
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}

	t.Errorf("expression type not valid, got=%T", exp)
	t.FailNow()
	return false
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) bool {
	infixExpression, ok := exp.(*ast.InfixExpression)

	assert := assert.New(t)

	if !assert.Truef(ok, "Received expression not of type *ast.InfixExpression, got=%T", exp) {
		t.FailNow()
	}

	if !testLiteralExpression(t, infixExpression.Left, left) {
		t.FailNow()
	}

	if !assert.Equal(operator, infixExpression.Operator) {
		t.FailNow()
	}

	if !testLiteralExpression(t, infixExpression.Right, right) {
		t.FailNow()
	}

	return true
}
