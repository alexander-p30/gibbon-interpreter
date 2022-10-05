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
	case '+':
		nextToken = newToken(token.PLUS, l.currentChar)
	case '-':
		nextToken = newToken(token.MINUS, l.currentChar)
	case '*':
		nextToken = newToken(token.ASTERISK, l.currentChar)
	case '/':
		nextToken = newToken(token.SLASH, l.currentChar)

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
		if isOperator(l.currentChar) {
			nextToken.Literal = l.readMultiCharToken(isOperator)
			nextToken.Type = token.GetOperatorTokenType(nextToken.Literal)
			return nextToken
		} else if isValidInIdentifier(l.currentChar) {
			nextToken.Literal = l.readMultiCharToken(isValidInIdentifier)
			nextToken.Type = token.GetIdentTokenType(nextToken.Literal)
			return nextToken
		} else if isDigit(l.currentChar) {
			nextToken.Type = token.INT
			nextToken.Literal = l.readMultiCharToken(isDigit)
			return nextToken
		} else {
			newToken(token.ILLEGAL, l.currentChar)
		}
	}

	l.readChar()

	return nextToken
}

func isValidInIdentifier(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.currentChar == ' ' || l.currentChar == '\t' || l.currentChar == '\n' || l.currentChar == '\r' {
		l.readChar()
	}
}

var multiCharOperatorsInitials = [...]byte{token.LTE[0], token.GTE[0], token.EQUAL[0], token.DIFFERENT[0]}

func isOperator(char byte) bool {
	for _, multiCharOperatorInitial := range multiCharOperatorsInitials {
		if char == multiCharOperatorInitial {
			return true
		}
	}

	return false
}

func (l *Lexer) readMultiCharToken(verifierFunc func(byte) bool) string {
	firstTokenCharPosition := l.currentCharPosition

	for verifierFunc(l.currentChar) {
		l.readChar()
	}

	return l.input[firstTokenCharPosition:l.currentCharPosition]
}

func newToken(tokenType token.TokenType, tokenChar byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(tokenChar)}
}
