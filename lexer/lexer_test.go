package lexer

import (
	"bytes"
	"gibbon/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNextTokenWithBaseTokens(t *testing.T) {
	assert := assert.New(t)

	input := bytes.NewReader([]byte(`=(+){},;-!*/< >`))

	tests := []struct {
		expectedType              token.TokenType
		expectedLiteral           string
		expectedLocationLine      uint
		expectedLocationFirstChar uint
	}{
		{token.ASSIGN, "=", 1, 1},
		{token.LPAREN, "(", 1, 2},
		{token.PLUS, "+", 1, 3},
		{token.RPAREN, ")", 1, 4},
		{token.LBRACE, "{", 1, 5},
		{token.RBRACE, "}", 1, 6},
		{token.COMMA, ",", 1, 7},
		{token.SEMICOLON, ";", 1, 8},
		{token.MINUS, "-", 1, 9},
		{token.BANG, "!", 1, 10},
		{token.ASTERISK, "*", 1, 11},
		{token.SLASH, "/", 1, 12},
		{token.LT, "<", 1, 13},
		{token.GT, ">", 1, 15},
		{token.EOF, "", 1, 16},
	}

	l := NewLexer(input, "filename")

	for i, test := range tests {
		token := l.NextToken()

		if !assert.Equal(test.expectedType, token.Type) {
			assert.FailNowf("", "Failed on test line %d", i)
		}
		if !assert.Equal(test.expectedLiteral, token.Literal) {
			assert.FailNowf("", "Failed on test line %d", i)
		}
		if !assert.Equal(test.expectedLocationLine, token.Location.Line) {
			assert.FailNowf("", "Failed on test line %d", i)
		}
		if !assert.Equal(test.expectedLocationFirstChar, token.Location.FirstCharIndex) {
			assert.FailNowf("", "Failed on test line %d", i)
		}
	}
}

func TestNextTokenWithCode(t *testing.T) {
	assert := assert.New(t)

	input := bytes.NewReader([]byte(`let five = 5;
  let ten = 10;

  let add = fn(x, y) { 
    x + y;
  };

  if(!true) {
    return false;
  }

  5 == 10
  >= <= <> !=

  let result = add(five, ten);
  ''`))

	tests := []struct {
		expectedType              token.TokenType
		expectedLiteral           string
		expectedLocationLine      uint
		expectedLocationFirstChar uint
	}{
		{token.LET, "let", 1, 1},
		{token.IDENT, "five", 1, 5},
		{token.ASSIGN, "=", 1, 10},
		{token.INT, "5", 1, 12},
		{token.SEMICOLON, ";", 1, 13},
		{token.LET, "let", 2, 3},
		{token.IDENT, "ten", 2, 7},
		{token.ASSIGN, "=", 2, 11},
		{token.INT, "10", 2, 13},
		{token.SEMICOLON, ";", 2, 15},
		{token.LET, "let", 4, 3},
		{token.IDENT, "add", 4, 7},
		{token.ASSIGN, "=", 4, 11},
		{token.FUNCTION, "fn", 4, 13},
		{token.LPAREN, "(", 4, 15},
		{token.IDENT, "x", 4, 16},
		{token.COMMA, ",", 4, 17},
		{token.IDENT, "y", 4, 19},
		{token.RPAREN, ")", 4, 20},
		{token.LBRACE, "{", 4, 22},
		{token.IDENT, "x", 5, 5},
		{token.PLUS, "+", 5, 7},
		{token.IDENT, "y", 5, 9},
		{token.SEMICOLON, ";", 5, 10},
		{token.RBRACE, "}", 6, 3},
		{token.SEMICOLON, ";", 6, 4},
		{token.IF, "if", 8, 3},
		{token.LPAREN, "(", 8, 5},
		{token.BANG, "!", 8, 6},
		{token.TRUE, "true", 8, 7},
		{token.RPAREN, ")", 8, 11},
		{token.LBRACE, "{", 8, 13},
		{token.RETURN, "return", 9, 5},
		{token.FALSE, "false", 9, 12},
		{token.SEMICOLON, ";", 9, 17},
		{token.RBRACE, "}", 10, 3},
		{token.INT, "5", 12, 3},
		{token.EQUAL, "==", 12, 5},
		{token.INT, "10", 12, 8},
		{token.GTE, ">=", 13, 3},
		{token.LTE, "<=", 13, 6},
		{token.ILLEGAL, "<>", 13, 9},
		{token.DIFFERENT, "!=", 13, 12},
		{token.LET, "let", 15, 3},
		{token.IDENT, "result", 15, 7},
		{token.ASSIGN, "=", 15, 14},
		{token.IDENT, "add", 15, 16},
		{token.LPAREN, "(", 15, 19},
		{token.IDENT, "five", 15, 20},
		{token.COMMA, ",", 15, 24},
		{token.IDENT, "ten", 15, 26},
		{token.RPAREN, ")", 15, 29},
		{token.SEMICOLON, ";", 15, 30},
		{token.ILLEGAL, "'", 16, 3},
		{token.ILLEGAL, "'", 16, 4},
		{token.EOF, "", 16, 5},
	}

	l := NewLexer(input, "filename")

	for i, test := range tests {
		token := l.NextToken()

		if !assert.Equal(test.expectedType, token.Type) {
			assert.FailNowf("", "Failed on test line %d", i)
		}
		if !assert.Equal(test.expectedLiteral, token.Literal) {
			assert.FailNowf("", "Failed on test line %d", i)
		}
		if !assert.Equal(test.expectedLocationLine, token.Location.Line) {
			assert.FailNowf("", "Failed on test line %d", i)
		}
		if !assert.Equal(test.expectedLocationFirstChar, token.Location.FirstCharIndex) {
			assert.FailNowf("", "Failed on test line %d", i)
		}
	}
}

func TestEOFDetection(t *testing.T) {
	input := bytes.NewReader([]byte("some characters\nhere"))
	NewLexer(input, "filename")

}
