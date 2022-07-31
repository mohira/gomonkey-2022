package ast

import (
	"bytes"
	"gomonkey/token"
)

type Node interface {
	TokenLiteral() string
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

// Program ノードは、Statementでもないし、Expressionでもないですね。
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer // 今なら strings.Builder かも

	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return "" // ルートノードしかない場合ってこと(文が一切ない"プログラム"のとき)
	}
}

// LetStatement は 当然、Nodeだし、Statementですね。Expressionではないですね。
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier // <identifier>
	Value Expression  // <expression>
}

func (ls *LetStatement) String() string {
	// ex: "let x = 5;" みたいな文字列が手に入るよ
	var out bytes.Buffer

	out.WriteString(ls.Token.Literal + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) statementNode() {
	panic("implement me")
}

// Identifier は Node です。さらに、monkeyの仕様では、 Identifier は Expression なのです！ もちろん Statement ではありません。
type Identifier struct {
	Token token.Token // token.IDENT トークン
	Value string
}

func (i *Identifier) String() string {
	return i.Value
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

func (rs *ReturnStatement) String() string {
	// ex: "return 5;" みたいな文字列が手に入るよ

	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) statementNode() {
	panic("implement me")
}

type ExpressionStatement struct {
	Token      token.Token // <式>の最初のトークンを持つらしい。謎。
	Expression Expression
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) statementNode() {
	panic("implement me")
}

type IntegerLiteral struct {
	Token token.Token // token.INT
	Value int64
}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) expressionNode() {
	panic("implement me")
}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token // 演算を意味する前置トークン: 「!」とか「-」とか
	Operator string
	Right    Expression // <operator><expression> だから、位置関係的に「右」
}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

func (pe *PrefixExpression) expressionNode() {
	panic("implement me")
}

type InfixExpression struct {
	Token    token.Token // どの演算かを表す情報
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String() + " " + ie.Operator + " " + ie.Right.String())
	out.WriteString(")")

	return out.String()
}

func (ie *InfixExpression) expressionNode() {
	panic("implement me")
}

type Boolean struct {
	Token token.Token // token.TRUE | token.FALSE
	Value bool
}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) expressionNode() {
	panic("implement me")
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

type IfExpression struct {
	// GoはIf<文>だからちょっと違うね。
	Token       token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

func (ie *IfExpression) expressionNode() {
	panic("implement me")
}

// BlockStatement は単数形だけど、複数のStatementを持っているっていう構造
// https://pkg.go.dev/go/ast#BlockStmt と同じ
type BlockStatement struct {
	Token      token.Token // { トークン が入るらしい
	Statements []Statement
}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

func (bs *BlockStatement) statementNode() {
	panic("implement me")
}
