package lexer

import (
	"bytes"
	"gibbon/token"
	"testing"
)

func TestNextTokenWithBaseTokens(t *testing.T) {
	input := bytes.NewReader([]byte(`=(+){},;-!*/< >`))

	tests := []struct {
		expectedType              token.TokenType
		expectedLiteral           string
		expectedLocationLine      uint
		expectedLocationFirstChar uint
	}{
		{token.ASSIGN, "=", 0, 0},
		{token.LPAREN, "(", 0, 1},
		{token.PLUS, "+", 0, 2},
		{token.RPAREN, ")", 0, 3},
		{token.LBRACE, "{", 0, 4},
		{token.RBRACE, "}", 0, 5},
		{token.COMMA, ",", 0, 6},
		{token.SEMICOLON, ";", 0, 7},
		{token.MINUS, "-", 0, 8},
		{token.BANG, "!", 0, 9},
		{token.ASTERISK, "*", 0, 10},
		{token.SLASH, "/", 0, 11},
		{token.LT, "<", 0, 12},
		{token.GT, ">", 0, 14},
		{token.EOF, "", 0, 15},
	}

	l := NewLexer(input, "filename")

	for idx, test := range tests {
		token := l.NextToken()

		if token.Type != test.expectedType {
			t.Fatalf("Test [%d]\n\tExpected: %+v\n\tGot: %+v", idx, test.expectedType, token.Type)
		}

		if token.Literal != test.expectedLiteral {
			t.Fatalf("Test [%d]\n\tExpected: %+v\n\tGot: %+v", idx, test.expectedLiteral, token.Literal)
		}

		if token.Location.Line != test.expectedLocationLine {
			t.Fatalf("Test [%d]\n\tExpected line: %d\n\tGot: %d", idx, test.expectedLocationLine, token.Location.Line)
		}

		if token.Location.FirstCharIndex != test.expectedLocationFirstChar {
			t.Fatalf("Test [%d]\n\tExpected char index: %d\n\tGot: %d", idx, test.expectedLocationFirstChar, token.Location.FirstCharIndex)
		}
	}
}

func TestNextTokenWithCode(t *testing.T) {
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
  `))

	tests := []struct {
		expectedType              token.TokenType
		expectedLiteral           string
		expectedLocationLine      uint
		expectedLocationFirstChar uint
	}{
		{token.LET, "let", 0, 0},
		{token.IDENT, "five", 0, 4},
		{token.ASSIGN, "=", 0, 9},
		{token.INT, "5", 0, 11},
		{token.SEMICOLON, ";", 0, 12},
		{token.LET, "let", 1, 2},
		{token.IDENT, "ten", 1, 6},
		{token.ASSIGN, "=", 1, 10},
		{token.INT, "10", 1, 12},
		{token.SEMICOLON, ";", 1, 14},
		{token.LET, "let", 3, 2},
		{token.IDENT, "add", 3, 6},
		{token.ASSIGN, "=", 3, 10},
		{token.FUNCTION, "fn", 3, 12},
		{token.LPAREN, "(", 3, 14},
		{token.IDENT, "x", 3, 15},
		{token.COMMA, ",", 3, 16},
		{token.IDENT, "y", 3, 18},
		{token.RPAREN, ")", 3, 19},
		{token.LBRACE, "{", 3, 21},
		{token.IDENT, "x", 4, 4},
		{token.PLUS, "+", 4, 6},
		{token.IDENT, "y", 4, 8},
		{token.SEMICOLON, ";", 4, 9},
		{token.RBRACE, "}", 5, 2},
		{token.SEMICOLON, ";", 5, 3},
		{token.IF, "if", 7, 2},
		{token.LPAREN, "(", 7, 4},
		{token.BANG, "!", 7, 5},
		{token.TRUE, "true", 7, 6},
		{token.RPAREN, ")", 7, 10},
		{token.LBRACE, "{", 7, 12},
		{token.RETURN, "return", 8, 4},
		{token.FALSE, "false", 8, 11},
		{token.SEMICOLON, ";", 8, 16},
		{token.RBRACE, "}", 9, 2},
		{token.INT, "5", 11, 2},
		{token.EQUAL, "==", 11, 4},
		{token.INT, "10", 11, 7},
		{token.GTE, ">=", 12, 2},
		{token.LTE, "<=", 12, 5},
		{token.ILLEGAL, "<>", 12, 8},
		{token.DIFFERENT, "!=", 12, 11},
		{token.LET, "let", 14, 2},
		{token.IDENT, "result", 14, 6},
		{token.ASSIGN, "=", 14, 13},
		{token.IDENT, "add", 14, 15},
		{token.LPAREN, "(", 14, 18},
		{token.IDENT, "five", 14, 19},
		{token.COMMA, ",", 14, 23},
		{token.IDENT, "ten", 14, 25},
		{token.RPAREN, ")", 14, 28},
		{token.SEMICOLON, ";", 14, 29},
		{token.EOF, "", 15, 2},
	}

	l := NewLexer(input, "filename")

	for idx, test := range tests {
		token := l.NextToken()

		if token.Type != test.expectedType {
			t.Fatalf("Test [%d]\n\tWrong token type!\n\tExpected: %+v\n\tGot: %+v", idx, test, token)
		}

		if token.Literal != test.expectedLiteral {
			t.Fatalf("Test [%d]\n\tWrong token literal!\n\tExpected: %+v\n\tGot: %+v", idx, test, token)
		}

		if token.Location.Line != test.expectedLocationLine {
			t.Fatalf("Test [%d]\n\tExpected line: %d\n\tGot: %d", idx, test.expectedLocationLine, token.Location.Line)
		}

		if token.Location.FirstCharIndex != test.expectedLocationFirstChar {
			t.Fatalf("Test [%d]\n\tExpected char index: %d\n\tGot: %d", idx, test.expectedLocationFirstChar, token.Location.FirstCharIndex)
		}
	}
}

func TestEOFDetection(t *testing.T) {
	input := bytes.NewReader([]byte("some characters\nhere"))
	NewLexer(input, "filename")

}
