package ast

import (
	"gomonkey/token"
)

/*
NodeはTokenLiteral()が必要。ノードが関連付けられているトークンのリテラル値がわからないとダメ。

Nodeはいくつかの種類がある
- ただのNode
- Statementインタフェース実装しているNode
- Expressionインタフェースを実装しているNode

*/

// Q. Node と Token は何が違うの？
// Q. キーワード let は Node ですか？ それとも Token ですか？

type Node interface {
	TokenLiteral() string // そのノードが関連付けられているトークン の リテラル値 を返す。デバッグ専用。
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program プログラム は 複数の 文 から構成されるって感じ。
// Programノードは、ルートノードですよ
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

// LetStatement
// let は Token。
// let文 は Node
type LetStatement struct {
	// let <identifier> = <expression>;
	Token token.Token // token.LET っていうトークンが入るだけじゃんね。
	Name  *Identifier // 要は、左辺の<識別子>
	Value Expression  // こっちは、右辺の<式>
}

func (ls *LetStatement) TokenLiteral() string {
	// let文 は Node なので TokenLiteral() を実装しないとだめだよね。
	return ls.Token.Literal
}

func (ls *LetStatement) statementNoda() {
	// let文 は Statement でもあるので、 statementNode() を実装しないといけないよね
	panic("implement me")
}

type Identifier struct {
	Token token.Token // 何かしらのトークンでもあるよね
	Value string      // 識別子の"実際の値"とでも言えばいいかな。
}

func (i *Identifier) TokenLiteral() string {
	// <identifier> は Node だよね。
	return i.Token.Literal
}

func (i *Identifier) expressionNode() {
	// <Identifier> は <式> でもあるよね
}
