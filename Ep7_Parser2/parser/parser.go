package parser

import (
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
)

// we need tokens provided by {lexer}
// currentToken represents the current token we are trying to parse
// nextToken represents the next token we will parse
// using the values together will help us determine what is coming next
// ie. (9) + 9. The ) lets us know the current expression is ending, and the + lets us know that we
// are dealing with an ADD expression
type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	nextToken    token.Token
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer}

	parser.getToken()
	parser.getToken()

	return parser
}

// Get the next token from our lexer
func (parser *Parser) getToken() {
	parser.currentToken = parser.nextToken
	parser.nextToken = parser.lexer.GetToken()
}

// PROG = STMT_LIST
// STMT_LIST = STMT | STMT STMT_LIST
// above i define a statement list with a recursive nature, here we can just parse it iteratively
func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	program.Statements = []ast.Statement{}

	for parser.currentToken.Type != token.EOF {
		stmt := parser.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		parser.getToken()
	}

	return program
}

// STMT = LET_STMT
func (parser *Parser) parseStatement() ast.Statement {
	switch parser.currentToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	default:
		return nil
	}
}

// LET_STMT = LET ID EQ EXPRESSION
func (parser *Parser) parseLetStatement() *ast.LetStatement {
	// LET
	stmt := &ast.LetStatement{Token: parser.currentToken}

	// ID
	if !parser.expect(token.ID) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}

	// EQ
	if !parser.expect(token.EQ) {
		return nil
	}

	// EXPR(TODO)
	for !parser.currentTokenIs(token.SEMICOLON) {
		parser.getToken()
	}

	// Return our fresh LetStatement
	return stmt
}

// We expect the next token to be of type {tokenType}, if it is, we will 'consume' it, otherwise, return false and handle error
func (parser *Parser) expect(tokenType token.TokenType) bool {
	if parser.nextTokenIs(tokenType) {
		parser.getToken()
		return true
	} else {
		return false
	}
}

// Check if the currentToken is of the provided {tokenTYpe}
func (parser *Parser) currentTokenIs(tokenType token.TokenType) bool {
	return parser.currentToken.Type == tokenType
}

// Check if the nextToken is of the provided {tokenType}
func (parser *Parser) nextTokenIs(tokenType token.TokenType) bool {
	return parser.nextToken.Type == tokenType
}
