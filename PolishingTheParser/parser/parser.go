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
	INDEX
)

var precedences = map[token.TokenType]int{
	token.EQ:     EQUALS,
	token.NEQ:    EQUALS,
	token.LT:     LESSERGREATER,
	token.GT:     LESSERGREATER,
	token.PLUS:   SUM,
	token.MINUS:  SUM,
	token.DIV:    MULT,
	token.MULT:   MULT,
	token.LPAREN: CALL,
	token.LBRACK: INDEX,
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
	parser.addPrefixToken(token.LPAREN, parser.parseGroupedExpression)
	parser.addPrefixToken(token.TRUE, parser.parseBoolean)
	parser.addPrefixToken(token.FALSE, parser.parseBoolean)
	parser.addPrefixToken(token.IF, parser.parseIfExpression)
	parser.addPrefixToken(token.FUNC, parser.parseFunctionLiteral)
	parser.addPrefixToken(token.LBRACK, parser.parseArrayLiteral)

	parser.infixParseFuncs = make(map[token.TokenType]infixParse)
	parser.addInfixToken(token.PLUS, parser.parseInfixExpression)
	parser.addInfixToken(token.MINUS, parser.parseInfixExpression)
	parser.addInfixToken(token.DIV, parser.parseInfixExpression)
	parser.addInfixToken(token.MULT, parser.parseInfixExpression)
	parser.addInfixToken(token.EQ, parser.parseInfixExpression)
	parser.addInfixToken(token.NEQ, parser.parseInfixExpression)
	parser.addInfixToken(token.LT, parser.parseInfixExpression)
	parser.addInfixToken(token.GT, parser.parseInfixExpression)
	parser.addInfixToken(token.LPAREN, parser.parseCallExpression)
	parser.addInfixToken(token.LBRACK, parser.parseIndexExpression)

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

func (parser *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: parser.currentToken}

	array.Items = parser.parseExpressionList(token.RBRACK)

	return array
}

func (parser *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: parser.currentToken, Left: left}

	parser.getToken()

	exp.Index = parser.parseExpression(NONE)

	if !parser.expect(token.RBRACK) {
		return nil
	}

	return exp
}

func (parser *Parser) parseExpressionList(endToken token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if parser.nextTokenIs(endToken) {
		parser.getToken()
		return list
	}

	parser.getToken()
	list = append(list, parser.parseExpression(NONE))

	for parser.nextTokenIs(token.COMMA) {
		parser.getToken()
		parser.getToken()
		list = append(list, parser.parseExpression(NONE))
	}

	if !parser.expect(endToken) {
		return nil
	}

	return list
}

func (parser *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: parser.currentToken, Function: function}
	exp.Arguments = parser.parseExpressionList(token.RPAREN)
	return exp
}

func (parser *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: parser.currentToken}

	if !parser.expect(token.LPAREN) {
		return nil
	}

	parser.getToken()
	expression.Condition = parser.parseExpression(NONE)

	if !parser.expect(token.RPAREN) {
		return nil
	}

	if !parser.expect(token.LBRACE) {
		return nil
	}

	expression.Consequence = parser.parseStatementBlock()

	if parser.nextTokenIs(token.ELSE) {
		parser.getToken()
		if !parser.expect(token.LBRACE) {
			return nil
		}
		expression.Alternative = parser.parseStatementBlock()
	}
	return expression
}

func (parser *Parser) parseStatementBlock() *ast.StatementBlock {
	block := &ast.StatementBlock{Token: parser.currentToken}
	block.Statements = []ast.Statement{}
	parser.getToken()

	for !parser.currentTokenIs(token.RBRACE) && !parser.currentTokenIs(token.EOF) {
		stmt := parser.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		parser.getToken()
	}
	return block
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

func (parser *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanExpression{Token: parser.currentToken, Value: parser.currentTokenIs(token.TRUE)}
}

func (parser *Parser) parseGroupedExpression() ast.Expression {
	parser.getToken()

	exp := parser.parseExpression(NONE)

	if !parser.expect(token.RPAREN) {
		return nil
	}

	return exp
}

func (parser *Parser) parseFunctionLiteral() ast.Expression {
	fl := &ast.FunctionLiteral{
		Token: parser.currentToken,
	}

	if !parser.expect(token.LPAREN) {
		return nil
	}

	fl.Parameters = parser.parseFunctionParameters()

	if !parser.expect(token.LBRACE) {
		return nil
	}

	fl.Body = parser.parseStatementBlock()
	return fl
}

func (parser *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if parser.nextTokenIs(token.RPAREN) {
		parser.getToken()
		return identifiers
	}

	parser.getToken()
	ident := &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}
	identifiers = append(identifiers, ident)

	for parser.nextTokenIs(token.COMMA) {
		parser.getToken()
		parser.getToken()
		ident := &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}
		identifiers = append(identifiers, ident)
	}

	if !parser.expect(token.RPAREN) {
		return nil
	}

	return identifiers
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
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: parser.currentToken}

	if !parser.expect(token.ID) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}

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
