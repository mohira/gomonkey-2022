package ast_test

import (
	"gomonkey/ast"
	"reflect"
	"testing"
)

func TestModify(t *testing.T) {
	// 毎回新しいIntegerLiteralを返す「関数」なのをお忘れなく！
	one := func() ast.Expression { return &ast.IntegerLiteral{Value: 1} }
	two := func() ast.Expression { return &ast.IntegerLiteral{Value: 2} }

	turnOneIntoTwo := func(node ast.Node) ast.Node {
		integer, ok := node.(*ast.IntegerLiteral)
		if !ok {
			return node
		}

		if integer.Value != 1 {
			return node
		}

		integer.Value = 2
		return integer
	}

	tests := []struct {
		input    ast.Node
		expected ast.Node
	}{
		{one(), two()},
		{&ast.Program{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: one()}}},
			&ast.Program{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: two()}}},
		},
	}

	for _, tt := range tests {
		// ノード変更関数を渡す感じ
		modified := ast.Modify(tt.input, turnOneIntoTwo)

		equal := reflect.DeepEqual(modified, tt.expected)
		if !equal {
			t.Errorf("not equal. got=%#v, want=%#v", modified, tt.expected)
		}

	}

}
