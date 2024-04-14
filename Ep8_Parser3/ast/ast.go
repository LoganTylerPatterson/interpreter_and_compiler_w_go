package ast

import (
	"bytes"
	"interpreter/token"
)

type Node interface {
	TokenLexeme() string
	String() string
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

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// LET STATEMENTS //
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) TokenLexeme() string { return ls.Token.Lexeme }
func (ls *LetStatement) statementNode()      {}
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLexeme() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// IDs //
type Identifier struct {
	Token token.Token
	Value string
}

func (id *Identifier) TokenLexeme() string { return id.Token.Lexeme }
func (id *Identifier) expressionNode()     {}
func (id *Identifier) String() string {
	return id.Value
}

/** 	RETURN STATEMENTS		**/
type ReturnStatement struct {
	// RETURN token
	Token token.Token
	// Value from expression to be returned
	Value Expression
}

func (rt *ReturnStatement) TokenLexeme() string { return rt.Token.Lexeme }
func (rt *ReturnStatement) statementNode()      {}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLexeme() + " ")

	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

/**  	EXPRESSION STATEMENTS 	**/
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) TokenLexeme() string { return es.Token.Lexeme }
func (es *ExpressionStatement) statementNode()      {}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

/** INTEGER LITERAL EXPRESSIONS **/
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) TokenLexeme() string { return il.Token.Lexeme }
func (il *IntegerLiteral) expressionNode()     {}
func (il *IntegerLiteral) String() string      { return il.Token.Lexeme }

/** PREFIX EXPRESSIONS **/
type PrefixExpression struct {
	Token token.Token
	Op    string
	Value Expression
}

func (pe *PrefixExpression) expressionNode()     {}
func (pe *PrefixExpression) TokenLexeme() string { return pe.Token.Lexeme }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Op)
	out.WriteString(pe.Value.String())
	out.WriteString(")")

	return out.String()
}
