package ast

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

// Program ノードは、Statementでもないし、Expressionでもないですね。
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return "" // ルートノードしかない場合ってこと(文が一切ない"プログラム"のとき)
	}
}
