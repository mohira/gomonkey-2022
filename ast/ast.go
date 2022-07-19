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
