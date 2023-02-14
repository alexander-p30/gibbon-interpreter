package lexer

import (
	"gibbon/token"
	"io"
)

const EOF_CHAR = 0

type bytePosition struct {
	byte uint
	line uint
}

type Lexer struct {
	input               io.ByteReader // input being parsed
	fileName            string        // name of the file being lexed
	currentChar         byte          // current char in examination
	currentCharPosition bytePosition  // current byte position on file being parsed
	nextCharPosition    bytePosition  // next byte position on file being parsed
	eofReached          bool          // indicates wether has been completely read
}

func NewLexer(input io.ByteReader, fileName string) *Lexer {
	l := &Lexer{input: input, fileName: fileName}
	l.readChar()
	return l
}

func (l *Lexer) SourceFile() string {
	return l.fileName
}

func (l *Lexer) readChar() {
	byte, err := l.input.ReadByte()

	l.currentChar = byte

	l.currentCharPosition = l.nextCharPosition

	if err == io.EOF {
		l.currentCharPosition = l.nextCharPosition
		return
	} else if byte == '\n' {
		l.nextCharPosition.line++
		l.nextCharPosition.byte = 0
	} else {
		l.nextCharPosition.byte++
	}
}

func (l *Lexer) NextToken() token.Token {
	var nextToken token.Token

	l.skipWhitespace()
	switch l.currentChar {
	// Operators
	case '+':
		nextToken = newToken(token.PLUS, l.currentChar, l.currentCharPosition)
	case '-':
		nextToken = newToken(token.MINUS, l.currentChar, l.currentCharPosition)
	case '*':
		nextToken = newToken(token.ASTERISK, l.currentChar, l.currentCharPosition)
	case '/':
		nextToken = newToken(token.SLASH, l.currentChar, l.currentCharPosition)

		// Delimiters
	case ',':
		nextToken = newToken(token.COMMA, l.currentChar, l.currentCharPosition)
	case ';':
		nextToken = newToken(token.SEMICOLON, l.currentChar, l.currentCharPosition)
	case '(':
		nextToken = newToken(token.LPAREN, l.currentChar, l.currentCharPosition)
	case ')':
		nextToken = newToken(token.RPAREN, l.currentChar, l.currentCharPosition)
	case '{':
		nextToken = newToken(token.LBRACE, l.currentChar, l.currentCharPosition)
	case '}':
		nextToken = newToken(token.RBRACE, l.currentChar, l.currentCharPosition)

		// Special
	case EOF_CHAR:
		nextToken.Location = token.TokenLocation{Line: l.currentCharPosition.line, FirstCharIndex: l.currentCharPosition.byte}
		nextToken.Literal = ""
		nextToken.Type = token.EOF
	default:
		if isOperator(l.currentChar) {
			nextToken.Location = token.TokenLocation{Line: l.currentCharPosition.line, FirstCharIndex: l.currentCharPosition.byte}
			nextToken.Literal = l.readMultiCharToken(isOperator)
			nextToken.Type = token.GetOperatorTokenType(nextToken.Literal)
			return nextToken
		} else if isValidInIdentifier(l.currentChar) {
			nextToken.Location = token.TokenLocation{Line: l.currentCharPosition.line, FirstCharIndex: l.currentCharPosition.byte}
			nextToken.Literal = l.readMultiCharToken(isValidInIdentifier)
			nextToken.Type = token.GetIdentTokenType(nextToken.Literal)
			return nextToken
		} else if isDigit(l.currentChar) {
			nextToken.Location = token.TokenLocation{Line: l.currentCharPosition.line, FirstCharIndex: l.currentCharPosition.byte}
			nextToken.Type = token.INT
			nextToken.Literal = l.readMultiCharToken(isDigit)
			return nextToken
		} else {
			nextToken = newToken(token.ILLEGAL, l.currentChar, l.currentCharPosition)
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
	readChars := []byte{}

	for verifierFunc(l.currentChar) {
		readChars = append(readChars, l.currentChar)
		l.readChar()
	}

	return string(readChars)
}

func newToken(tokenType token.TokenType, tokenChar byte, bytePosition bytePosition) token.Token {
	return token.Token{Type: tokenType, Literal: string(tokenChar), Location: token.TokenLocation{Line: bytePosition.line, FirstCharIndex: bytePosition.byte}}
}
