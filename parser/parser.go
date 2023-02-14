package parser

import (
	"gibbon/lexer"
	"gibbon/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []Error
}

type Error struct {
	message  string
	location token.TokenLocation
}
