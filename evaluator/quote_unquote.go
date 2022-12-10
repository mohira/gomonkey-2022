package evaluator

import (
	"fmt"
	"gomonkey/ast"
	"gomonkey/object"
	"gomonkey/token"
)

func quote(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquoteCalls(node, env)

	return &object.Quote{Node: node}
}

func evalUnquoteCalls(quoted ast.Node, env *object.Environment) ast.Node {
	// to俺: ast.Modifyをまず呼びだしているからな！
	// 第2引数の 関数 はその後やで！
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}

		// CallExpressionが確定しているのでエラーチェックいらないでしょ？ と思うのでそうしてます。
		callExpr := node.(*ast.CallExpression)

		if len(callExpr.Arguments) != 1 {
			// 他のとこでチェックかもしれんが、とりあえずここでチェックしておく。やや些末なので別にいいでしょって思ってる。
			return node
		}

		evaluated := Eval(callExpr.Arguments[0], env)

		return convertObjectToASTNode(evaluated)

	})
}

func isUnquoteCall(node ast.Node) bool {
	callExpr, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}

	return callExpr.Function.TokenLiteral() == "unquote"
}

func convertObjectToASTNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    "INT",
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{
			Token: t,
			Value: obj.Value,
		}
	case *object.Boolean:
		var t token.Token

		if obj.Value {
			t = token.Token{Type: token.TRUE, Literal: "true"}
		} else {
			t = token.Token{Type: token.FALSE, Literal: "false"}
		}

		return &ast.Boolean{Token: t, Value: obj.Value}

	default:
		return nil
	}
}
