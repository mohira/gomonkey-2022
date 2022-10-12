package ast

import "gomonkey/token"

type Node interface {
	TokenLiteral() string
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

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier // Identifierノード(token.IDENTではない！)
	Value Expression
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) statementNode() {
	panic("implement me")
}

type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) expressionNode() {
	panic("implement me")
}

type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {
	panic("implement me")
}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// ExpressionStatement は 「1つの式」だけで構成される「文」
// x = 5; は let文だけど
// x + 5; は 式文
type ExpressionStatement struct {
	Token      token.Token // 式の最初のトークン
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {
	panic("implement me")
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.TokenLiteral()
}
