package ast

import "interpreter/token"

type Node interface {
	TokenLexeme() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLexeme() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLexeme()
	} else {
		return ""
	}
}

//	LET STATEMENTS //
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) TokenLexeme() string { return ls.Token.Lexeme }
func (ls *LetStatement) statementNode()      {}

//	IDs //
type Identifier struct {
	Token token.Token
	Value string
}

func (id *Identifier) TokenLexeme() string { return id.Token.Lexeme }
func (id *Identifier) expressionNode()     {}

/** 	RETURN STATEMENTS		**/
type ReturnStatement struct {
	// RETURN token
	Token token.Token
	// Value from expression to be returned
	Value Expression
}

func (rt *ReturnStatement) TokenLexeme() string { return rt.Token.Lexeme }
func (rt *ReturnStatement) statementNode()      {}
