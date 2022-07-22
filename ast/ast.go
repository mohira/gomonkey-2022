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

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) statementNode() {
	panic("implement me")
}
