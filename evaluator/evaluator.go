package evaluator

import (
	"fmt"
	"gomonkey/ast"
	"gomonkey/object"
)

func Eval(node ast.Node) object.Object {
	fmt.Printf("Eval: %[1]T %[1]v\n", node)

	switch n := node.(type) {
	// 複数の文
	case *ast.Program:
		return evalStatements(n.Statements)

	// 単一の文
	case *ast.ExpressionStatement:
		return Eval(n.Expression)

	// 式
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}

	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	fmt.Printf("evalStatements: %[1]T %[1]v\n", stmts)

	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}
