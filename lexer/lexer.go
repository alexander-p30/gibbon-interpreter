package lexer

import (
	"gibbon/token"
)

const EOL_CHAR = 0

type Lexer struct {
	input               string // input being parsed
	currentCharPosition int    // current position on input (current char position)
	nextCharPosition    int    // current reading position (next char to be parsed position)
	currentChar         byte   // current char in examination
}

// TODO receive file name and io.Reader to allow for file position storage
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.nextCharPosition >= len(l.input) {
		l.currentChar = EOL_CHAR
	} else {
		l.currentChar = l.input[l.nextCharPosition]
	}

	l.currentCharPosition = l.nextCharPosition
	l.nextCharPosition++
}

func (l *Lexer) NextToken() token.Token {
	var nextToken token.Token

	l.skipWhitespace()
	switch l.currentChar {
	// Operators
	case '=':
		nextToken = newToken(token.ASSIGN, l.currentChar)
	case '+':
		nextToken = newToken(token.PLUS, l.currentChar)

		// Delimiters
	case ',':
		nextToken = newToken(token.COMMA, l.currentChar)
	case ';':
		nextToken = newToken(token.SEMICOLON, l.currentChar)
	case '(':
		nextToken = newToken(token.LPAREN, l.currentChar)
	case ')':
		nextToken = newToken(token.RPAREN, l.currentChar)
	case '{':
		nextToken = newToken(token.LBRACE, l.currentChar)
	case '}':
		nextToken = newToken(token.RBRACE, l.currentChar)

		// Special
	case EOL_CHAR:
		nextToken.Literal = ""
		nextToken.Type = token.EOF
	default:
		if isCharValidInIdentifier(l.currentChar) {
			nextToken.Literal = l.readIdentifier()
			nextToken.Type = token.GetIdentTokenType(nextToken.Literal)
			return nextToken
		} else if isDigit(l.currentChar) {
			nextToken.Type = token.INT
			nextToken.Literal = l.readNumber()
			return nextToken
		} else {
			newToken(token.ILLEGAL, l.currentChar)
		}
	}

	l.readChar()

	return nextToken
}

func (l *Lexer) readIdentifier() string {
	firstIdentifierCharPosition := l.currentCharPosition

	for isCharValidInIdentifier(l.currentChar) {
		l.readChar()
	}

	return l.input[firstIdentifierCharPosition:l.currentCharPosition]
}

func isCharValidInIdentifier(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func (l *Lexer) readNumber() string {
	firstIdentifierCharPosition := l.currentCharPosition

	for isDigit(l.currentChar) {
		l.readChar()
	}

	return l.input[firstIdentifierCharPosition:l.currentCharPosition]
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.currentChar == ' ' || l.currentChar == '\t' || l.currentChar == '\n' || l.currentChar == '\r' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, tokenChar byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(tokenChar)}
}
