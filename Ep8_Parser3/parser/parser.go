package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
	"strconv"
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
	//list of errors to return after parsing
	errors []string
	//hashmap of infix and prefix operators
	prefixParseFuncs map[token.TokenType]prefixParse
	infixParseFuncs  map[token.TokenType]infixParse
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	parser.prefixParseFuncs = make(map[token.TokenType]prefixParse)
	parser.addPrefixToken(token.ID, parser.parseIdentifier)
	parser.addPrefixToken(token.DIGIT, parser.parseIntegerLiteral)
	parser.addPrefixToken(token.MINUS, parser.parsePrefixExpression)
	parser.addPrefixToken(token.EXCLAM, parser.parsePrefixExpression)

	parser.getToken()
	parser.getToken()

	return parser
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	il := &ast.IntegerLiteral{Token: parser.currentToken}

	value, err := strconv.ParseInt(parser.currentToken.Lexeme, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", parser.currentToken.Lexeme)
		parser.errors = append(parser.errors, msg)
		return nil
	}

	il.Value = value
	return il
}

// Get the next token from our lexer
func (parser *Parser) getToken() {
	parser.currentToken = parser.nextToken
	parser.nextToken = parser.lexer.GetToken()
}

type (
	prefixParse func() ast.Expression
	infixParse  func(ast.Expression) ast.Expression
)

func (parser *Parser) addPrefixToken(tokenType token.TokenType, fnc prefixParse) {
	parser.prefixParseFuncs[tokenType] = fnc
}

func (parser *Parser) addInfixToken(tokenType token.TokenType, fnc infixParse) {
	parser.infixParseFuncs[tokenType] = fnc
}

// PROG = STMT_LIST
// STMT_LIST = STMT | STMT STMT_LIST
// above i define a statement list with a recursive nature, here we can just parse it iteratively
func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !parser.currentTokenIs(token.EOF) {
		stmt := parser.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		parser.getToken()
	}

	return program
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.currentToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: parser.currentToken}

	stmt.Expression = parser.parseExpression(NONE)

	if parser.nextTokenIs(token.SEMICOLON) {
		parser.getToken()
	}

	return stmt
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefixFunc := parser.prefixParseFuncs[parser.currentToken.Type]

	if prefixFunc == nil {
		return nil
	}

	leftExpr := prefixFunc()

	return leftExpr
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	pe := &ast.PrefixExpression{
		Token: parser.currentToken,
		Op:    parser.currentToken.Lexeme,
	}

	parser.getToken()

	pe.Value = parser.parseExpression(PREFIX)

	return pe
}

const (
	_ int = iota
	NONE
	EQUALS        // ==
	LESSERGREATER // < or >
	SUM           // +
	MULT          // *
	PREFIX        // !<condition>
	CALL          //function call
)

// LET_STMT = LET ID EQ EXPRESSION
func (parser *Parser) parseLetStatement() *ast.LetStatement {
	// LET
	stmt := &ast.LetStatement{Token: parser.currentToken}
	fmt.Printf("%+v\n", parser.currentToken)
	// ID
	if !parser.expect(token.ID) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}

	// EQ
	if !parser.expect(token.ASSIGN) {
		return nil
	}

	// EXPR(TODO)
	for !parser.currentTokenIs(token.SEMICOLON) {
		parser.getToken()
	}

	// Return our fresh LetStatement
	return stmt
}

// RETURN EXPR
func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	//RETURN
	stmt := &ast.ReturnStatement{Token: parser.currentToken}

	//EXPR
	//Skip the EXPR (TODO)
	parser.getToken()
	for !parser.currentTokenIs(token.SEMICOLON) {
		parser.getToken()
	}

	return stmt
}

// We expect the next token to be of type {tokenType}, if it is, we will 'consume' it, otherwise, return false and handle error
func (parser *Parser) expect(tokenType token.TokenType) bool {
	if parser.nextTokenIs(tokenType) {
		parser.getToken()
		return true
	} else {
		msg := fmt.Sprintf("expected next token to be %s, recieved %s instead", tokenType, parser.nextToken.Type)
		parser.errors = append(parser.errors, msg)
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

// Return the errors
func (parser *Parser) Errors() []string {
	return parser.errors
}
