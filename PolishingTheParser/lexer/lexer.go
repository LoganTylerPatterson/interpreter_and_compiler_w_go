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
			tok = token.Token{Type: token.NEQ, Lexeme: string(ch) + string(l.ch)}
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
	case '[':
		tok = createToken(token.LBRACK, l.ch)
	case ']':
		tok = createToken(token.RBRACK, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Lexeme = l.readString()
	case 0:
		tok.Lexeme = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Lexeme = l.readId()
			tok.Type = token.DetermineTokenType(tok.Lexeme)
			return tok
		} else if isDigit(l.ch) {
			tok.Lexeme = l.readNumber()
			tok.Type = token.DIGIT
			return tok
		} else {
			tok = createToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readString() string {
	position := l.position

	for {
		l.readChar()
		if l.ch != '"' || l.ch != 0 {
			break
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) readId() string {
	pos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readNumber() string {
	pos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) eatWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func createToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Lexeme: string(ch)}
}
