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
	switch n := node.(type) {
	// 複数の文
	case *ast.Program:
		return evalProgram(n)

	case *ast.BlockStatement:
		return evalBlockStatement(n)

	// 単一の文
	case *ast.ExpressionStatement:
		return Eval(n.Expression)

	case *ast.ReturnStatement:
		val := Eval(n.ReturnValue)
		return &object.ReturnValue{Value: val}
	// 式
	case *ast.IfExpression:
		return evalIfExpression(n)

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

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

func evalIfExpression(n *ast.IfExpression) object.Object {
	condition := Eval(n.Condition)

	if isTruthy(condition) {
		return Eval(n.Consequence)
	} else if n.Alternative != nil {
		return Eval(n.Alternative)
	} else {
		// 条件式false かつ elseブロックがないときってこと
		return NULL
	}
}

func isTruthy(condition object.Object) bool {
	// switchよりこっちの方が、Definitionな感じなので良いと思う！
	return condition != NULL && condition != FALSE
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerObj && right.Type() == object.IntegerObj:
		// MEMO: 整数オペランド同士の==演算とかはここでで処理されている
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() == object.BooleanObj && right.Type() == object.BooleanObj:
		// どっちもBOOLEAN
		// unknown operator: BOOLEAN + BOOLEAN
		msg := fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type())
		return &object.Error{Message: msg}

	default:
		//type mismatch: INTEGER + BOOLEAN
		msg := fmt.Sprintf("type mismatch: %s %s %s", left.Type(), operator, right.Type())
		return &object.Error{Message: msg}
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
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
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

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt)

		// アンラップしないので、型情報だけでいい
		// が、nilの場合に.Type()のアクセスをするとpanicになるので、
		// 短絡評価を使っている感じだと思う
		if result != nil && result.Type() == object.ReturnValueObj {
			return result // 返すけど、アンラップはしません！
		}
	}

	return result
}
