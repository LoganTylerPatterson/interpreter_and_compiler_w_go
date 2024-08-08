package token

type TokenType string

type Token struct {
	Type   TokenType
	Lexeme string
}

const (
	EOF       = "EOF"
	ILLEGAL   = "ILLEGAL"
	ID        = "ID"
	DIGIT     = "DIGIT"
	ASSIGN    = "="
	EQ        = "=="
	NEQ       = "!="
	PLUS      = "+"
	MINUS     = "-"
	MULT      = "*"
	DIV       = "/"
	LT        = "<"
	GT        = ">"
	COMMA     = ","
	EXCLAM    = "!"
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	FUNC      = "FUNC"
	LET       = "LET"
	IF        = "IF"
	ELSE      = "ELSE"
	RETURN    = "RETURN"
	TRUE      = "TRUE"
	FALSE     = "FALSE"
)

var keywords = map[string]TokenType{
	"func":   FUNC,
	"let":    LET,
	"false":  FALSE,
	"true":   TRUE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func DetermineTokenType(id string) TokenType {
	if tokenType, found := keywords[id]; found {
		return tokenType
	}
	return ID
}
