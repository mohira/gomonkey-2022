package ast_test

import (
	"fmt"
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
		// フィールド名の幅が揃ったほうが見やすいのでaとbにした。
		a ast.Node
		b ast.Node
	}{
		// IntegerLiteral
		{one(), two()},

		// Programノード
		{
			&ast.Program{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: one()}}},
			&ast.Program{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: two()}}},
		},

		// 中置演算式
		{
			&ast.InfixExpression{Left: one(), Operator: "+", Right: two()},
			&ast.InfixExpression{Left: two(), Operator: "+", Right: two()},
		},
		{
			&ast.InfixExpression{Left: two(), Operator: "+", Right: one()},
			&ast.InfixExpression{Left: two(), Operator: "+", Right: two()},
		},

		// 前置演算式
		{
			&ast.PrefixExpression{Operator: "-", Right: one()},
			&ast.PrefixExpression{Operator: "-", Right: two()},
		},

		// 添字演算子式
		{
			// `1[1]` -> `2[2]` は評価エラーだけど、いまはノードの話です
			&ast.IndexExpression{Left: one(), Index: one()},
			&ast.IndexExpression{Left: two(), Index: two()},
		},

		// if式

		{
			// if (1) {1} else {1}
			&ast.IfExpression{
				Condition: one(),
				Consequence: &ast.BlockStatement{
					Statements: []ast.Statement{&ast.ExpressionStatement{Expression: one()}},
				},
				Alternative: &ast.BlockStatement{
					Statements: []ast.Statement{&ast.ExpressionStatement{Expression: one()}},
				},
			},
			// if (2) {2} else {2}
			&ast.IfExpression{
				Condition: two(),
				Consequence: &ast.BlockStatement{
					Statements: []ast.Statement{&ast.ExpressionStatement{Expression: two()}},
				},
				Alternative: &ast.BlockStatement{
					Statements: []ast.Statement{&ast.ExpressionStatement{Expression: two()}},
				},
			},
		},
		// elseがないif式はvalidだからテストする
		{
			// if (1) {1}
			&ast.IfExpression{
				Condition: one(),
				Consequence: &ast.BlockStatement{
					Statements: []ast.Statement{&ast.ExpressionStatement{Expression: one()}},
				},
			},
			// if (2) {2}
			&ast.IfExpression{
				Condition: two(),
				Consequence: &ast.BlockStatement{
					Statements: []ast.Statement{&ast.ExpressionStatement{Expression: two()}},
				},
			},
		},

		{
			&ast.ReturnStatement{ReturnValue: one()},
			&ast.ReturnStatement{ReturnValue: two()},
		},
		{
			&ast.LetStatement{Value: one()},
			&ast.LetStatement{Value: two()},
		},

		// 関数リテラル
		{
			// fn() { 1 }
			&ast.FunctionLiteral{
				Parameters: []*ast.Identifier{},
				Body: &ast.BlockStatement{
					Statements: []ast.Statement{
						&ast.ExpressionStatement{Expression: one()},
					},
				},
			},
			// fn() { 2 }
			&ast.FunctionLiteral{
				Parameters: []*ast.Identifier{},
				Body: &ast.BlockStatement{
					Statements: []ast.Statement{
						&ast.ExpressionStatement{Expression: two()},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		// ノード変更関数を渡す感じ
		t.Run(fmt.Sprintf("test#%d", i), func(t *testing.T) {

			modified := ast.Modify(tt.a, turnOneIntoTwo)

			equal := reflect.DeepEqual(modified, tt.b)
			if !equal {
				t.Errorf("not equal. got=%#v, want=%#v", modified, tt.b)
			}
		})

	}

}
