package ast

type ModifyFunc func(Node) Node

func Modify(node Node, modifier ModifyFunc) Node {
	switch node := node.(type) {
	case *Program:
		for i, statement := range node.Statements {
			// MEMO: エラー処理は！？
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}

	case *InfixExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Right, _ = Modify(node.Right, modifier).(Expression)

	case *PrefixExpression:
		node.Right, _ = Modify(node.Right, modifier).(Expression)

	case *IndexExpression:
		node.Index, _ = Modify(node.Index, modifier).(Expression)
		node.Left, _ = Modify(node.Left, modifier).(Expression)
	case *ExpressionStatement:
		// MEMO: エラー処理は！？
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)
	}

	return modifier(node)
}
