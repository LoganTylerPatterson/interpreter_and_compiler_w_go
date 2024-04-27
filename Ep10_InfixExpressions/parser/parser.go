package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
	"strconv"
)

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

var precedences = map[token.TokenType]int{
	token.EQ:    EQUALS,
	token.NEQ:   EQUALS,
	token.LT:    LESSERGREATER,
	token.GT:    LESSERGREATER,
	token.PLUS:  SUM,
	token.MINUS: SUM,
	token.DIV:   MULT,
	token.MULT:  MULT,
}

type (
	prefixParse func() ast.Expression
	infixParse  func(ast.Expression) ast.Expression
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

	parser.infixParseFuncs = make(map[token.TokenType]infixParse)
	parser.addInfixToken(token.PLUS, parser.parseInfixExpression)
	parser.addInfixToken(token.MINUS, parser.parseInfixExpression)
	parser.addInfixToken(token.DIV, parser.parseInfixExpression)
	parser.addInfixToken(token.MULT, parser.parseInfixExpression)
	parser.addInfixToken(token.EQ, parser.parseInfixExpression)
	parser.addInfixToken(token.NEQ, parser.parseInfixExpression)
	parser.addInfixToken(token.LT, parser.parseInfixExpression)
	parser.addInfixToken(token.GT, parser.parseInfixExpression)

	//set our current token and peek token
	parser.getToken()
	parser.getToken()

	return parser
}

// Get the next token from our lexer
func (parser *Parser) getToken() {
	parser.currentToken = parser.nextToken
	parser.nextToken = parser.lexer.GetToken()
}

// We expect the next token to be of type {tokenType}, if it is, we will 'consume' it, otherwise, return false and handle error
func (parser *Parser) expect(tokenType token.TokenType) bool {
	if parser.nextTokenIs(tokenType) {
		parser.getToken()
		return true
	} else {
		parser.peekError(tokenType)
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

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf(
		"expected next token to be %s, got %s instead",
		t,
		parser.currentToken.Type,
	)
	parser.errors = append(parser.errors, msg)
}

func (parser *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf(
		"no prefix parse function for %s found",
		t,
	)
	parser.errors = append(parser.errors, msg)
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
		fmt.Print("parsing expression")
		return parser.parseExpressionStatement()
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
	if !parser.expect(token.ASSIGN) {
		return nil
	}

	parser.getToken()

	stmt.Value = parser.parseExpression(NONE)

	if parser.nextTokenIs(token.SEMICOLON) {
		parser.getToken()
	}

	return stmt
}

// RETURN EXPR
func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	//RETURN
	stmt := &ast.ReturnStatement{Token: parser.currentToken}

	parser.getToken()

	stmt.Value = parser.parseExpression(NONE)

	if parser.nextTokenIs(token.SEMICOLON) {
		parser.getToken()
	}

	return stmt
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
	prefix := parser.prefixParseFuncs[parser.currentToken.Type]
	if prefix == nil {
		parser.noPrefixParseFnError(parser.currentToken.Type)
		return nil
	}
	leftExpr := prefix()

	for !parser.nextTokenIs(token.SEMICOLON) && precedence < parser.peekPrecedence() {
		infix := parser.infixParseFuncs[parser.nextToken.Type]
		if infix == nil {
			return leftExpr
		}

		parser.getToken()

		leftExpr = infix(leftExpr)
	}
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

func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token: parser.currentToken,
		Op:    parser.currentToken.Lexeme,
		Left:  left,
	}

	precedence := parser.curPrecedence()
	parser.getToken()
	expr.Right = parser.parseExpression(precedence)

	return expr
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.nextToken.Type]; ok {
		return p
	}
	return NONE
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return NONE
}

func (parser *Parser) addPrefixToken(tokenType token.TokenType, fnc prefixParse) {
	parser.prefixParseFuncs[tokenType] = fnc
}

func (parser *Parser) addInfixToken(tokenType token.TokenType, fnc infixParse) {
	parser.infixParseFuncs[tokenType] = fnc
}
