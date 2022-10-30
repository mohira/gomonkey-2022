package evaluator

import (
	"gomonkey/ast"
	"gomonkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
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

	case *ast.InfixExpression:
		left := Eval(n.Left)
		right := Eval(n.Right)

		return evalInfixExpression(n.Operator, left, right)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(n.Value)

	}

	return nil
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerObj && right.Type() == object.IntegerObj:
		return evalIntegerInfixExpression(operator, left, right)
	default:
		return NULL
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	// Boolean
	case "<":
		return &object.Boolean{Value: leftValue < rightValue}
	case ">":
		return &object.Boolean{Value: leftValue > rightValue}
	case "==":
		return &object.Boolean{Value: leftValue == rightValue}
	case "!=":
		return &object.Boolean{Value: leftValue != rightValue}
	default:
		return NULL
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.IntegerObj {
		return NULL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}

	// 書き換えると壊れるよ！
	// a = 1
	// b = -a
	//intObj := right.(*object.Integer)
	//intObj.Value = -intObj.Value
	//return intObj
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
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}
