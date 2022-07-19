package ast

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
