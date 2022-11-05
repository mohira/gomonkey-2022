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

var Environment map[string]object.Object

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

	case *ast.LetStatement:
		v := Eval(n.Value)

		if isError(v) {
			return v
		}

		// 値を登録する
		Environment[n.Name.Value] = v

		return nil // ???? とりあえずこれにしとく

	case *ast.ReturnStatement:
		val := Eval(n.ReturnValue)

		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	// 式
	case *ast.IfExpression:
		return evalIfExpression(n)

	case *ast.PrefixExpression: // !true, !5, !!false
		right := Eval(n.Right)

		if isError(right) {
			return right
		}

		return evalPrefixExpression(n.Operator, right)

	case *ast.InfixExpression:
		left := Eval(n.Left)
		if isError(left) {
			return left
		}

		right := Eval(n.Right)
		if isError(right) {
			return right
		}

		return evalInfixExpression(n.Operator, left, right)

	case *ast.Identifier:
		if _, ok := Environment[n.Value]; !ok {
			return newError("あとで")
		}

		return Environment[n.Value]

	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(n.Value)

	}

	return nil
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorObj
	}

	return false
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalIfExpression(n *ast.IfExpression) object.Object {
	condition := Eval(n.Condition)

	if isError(condition) {
		// ERRORオブジェクトは実は truthy だった！
		// truthy := 「NULLでない かつ falseでない」 なので！！！
		return condition
	}
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
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.IntegerObj {
		// `-`という単項演算子が許されるのは(決めの問題でもあるが)、ふつーは、数値だけなので、
		// 条件判定は、 INTEGERオブジェクトじゃないとき でよさげ！
		return newError("unknown operator: -%s", right.Type())
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

		if result != nil {
			rt := result.Type()
			if rt == object.ReturnValueObj || rt == object.ErrorObj {
				return result
			}
		}

		// アンラップしないので、型情報だけでいい
		// が、nilの場合に.Type()のアクセスをするとpanicになるので、
		// 短絡評価を使っている感じだと思う
		// このベタベタ実装もきらいじゃないよ！
		// if result != nil && result.Type() == object.ReturnValueObj {
		// 	   return result // 返すけど、アンラップはしません！
		// }
		//
		// if result != nil && result.Type() == object.ErrorObj {
		//  	return result
		// }

	}

	return result
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
