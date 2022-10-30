package evaluator

import (
	"fmt"
	"gomonkey/ast"
	"gomonkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
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
	case *ast.PrefixExpression: // !true, !5, !!false
		right := Eval(n.Right)

		return evalPrefixExpression(n.Operator, right)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(n.Value)

	}

	return nil
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	default:
		return NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func nativeBoolToBooleanObject(value bool) object.Object {
	if value {
		return TRUE
	} else {
		return FALSE
	}
}

func evalStatements(stmts []ast.Statement) object.Object {
	fmt.Printf("evalStatements: %[1]T %[1]v\n", stmts)

	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}
