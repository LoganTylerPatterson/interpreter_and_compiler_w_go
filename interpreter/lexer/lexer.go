package lexer

import (
	"interpreter/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
		l.position = l.readPosition
		return
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) GetToken() token.Token {
	var tok token.Token

	l.eatWhitespace()

	switch l.ch {
	case '{':
		tok = createToken(token.LBRACE, l.ch)
	case '}':
		tok = createToken(token.RBRACE, l.ch)
	case '(':
		tok = createToken(token.LPAREN, l.ch)
	case ')':
		tok = createToken(token.RPAREN, l.ch)
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Lexeme: string(ch) + string(l.ch)}
		} else {
			tok = createToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Lexeme: string(ch) + string(l.ch)}
		} else {
			tok = createToken(token.EXCLAM, l.ch)
		}
	case '-':
		tok = createToken(token.MINUS, l.ch)
	case '+':
		tok = createToken(token.PLUS, l.ch)
	case '*':
		tok = createToken(token.MULT, l.ch)
	case '/':
		tok = createToken(token.DIV, l.ch)
	case '<':
		tok = createToken(token.LT, l.ch)
	case '>':
		tok = createToken(token.GT, l.ch)
	case ';':
		tok = createToken(token.SEMICOLON, l.ch)
	case ',':
		tok = createToken(token.COMMA, l.ch)
	case 0:
		tok.Lexeme = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Lexeme = l.readId()
			tok.Type = token.DetermineTokenType(tok.Lexeme)
			return tok
		} else {
			tok = createToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) eatWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readId() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func createToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Lexeme: string(ch)}
}
