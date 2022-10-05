package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	// Special
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	BANG      = "!"
	ASTERISK  = "*"
	SLASH     = "/"
	LT        = "<"
	GT        = ">"
	LTE       = "<="
	GTE       = ">="
	EQUAL     = "=="
	DIFFERENT = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	RETURN   = "RETURN"
	LET      = "LET"
	IF       = "IF"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
)

var keywordTypes = map[string]TokenType{
	"fn":     FUNCTION,
	"return": RETURN,
	"let":    LET,
	"if":     IF,
	"true":   TRUE,
	"false":  FALSE,
}

var operators = [...]string{ASSIGN, BANG, LT, GT, LTE, GTE, EQUAL, DIFFERENT}

func GetOperatorTokenType(operator string) TokenType {
	for _, tokenType := range operators {
		if tokenType == operator {
			return TokenType(tokenType)
		}
	}

	return ILLEGAL
}

func GetIdentTokenType(identifier string) TokenType {
	if tokenType, ok := keywordTypes[identifier]; ok {
		return tokenType
	}

	return IDENT
}
