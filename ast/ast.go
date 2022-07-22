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

func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStatement) statementNode() {
	panic("implement me")
}

type ExpressionStatement struct {
	Token      token.Token // <式>の最初のトークンを持つらしい。謎。
	Expression Expression
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) statementNode() {
	panic("implement me")
}
