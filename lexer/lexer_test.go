package lexer

import (
	"gibbon/token"
	"testing"
)

func TestNextTokenWithBaseTokens(t *testing.T) {
	input := `=+(){},;-!*/< >`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.MINUS, "-"},
		{token.BANG, "!"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for idx, test := range tests {
		token := l.NextToken()

		if token.Type != test.expectedType {
			t.Fatalf("Test [%d]\n\tExpected: %q\n\tGot: %q", idx, test.expectedType, token.Type)
		}

		if token.Literal != test.expectedLiteral {
			t.Fatalf("Test [%d]\n\tExpected: %q\n\tGot: %q", idx, test.expectedLiteral, token.Literal)
		}
	}
}

func TestNextTokenWithCode(t *testing.T) {
	input := `let five = 5;
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
  `

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.BANG, "!"},
		{token.TRUE, "true"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "5"},
		{token.EQUAL, "=="},
		{token.INT, "10"},
		{token.GTE, ">="},
		{token.LTE, "<="},
		{token.ILLEGAL, "<>"},
		{token.DIFFERENT, "!="},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
	}

	l := NewLexer(input)

	for idx, test := range tests {
		token := l.NextToken()

		if token.Type != test.expectedType {
			t.Fatalf("Test [%d]\n\tWrong token type!\n\tExpected: %+v\n\tGot: %+v", idx, test, token)
		}

		if token.Literal != test.expectedLiteral {
			t.Fatalf("Test [%d]\n\tWrong token literal!\n\tExpected: %+v\n\tGot: %+v", idx, test, token)
		}
	}
}
