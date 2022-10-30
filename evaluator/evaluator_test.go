package evaluator_test

import (
	"gomonkey/evaluator"
	"gomonkey/lexer"
	"gomonkey/object"
	"gomonkey/parser"
	"testing"
)

func TestIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		testIntegerObject(t, evaluated, tt.expected)
	}

}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	integerObj, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("obj is not *object.Integer. got=%T", obj)
		return false
	}

	if integerObj.Value != expected {
		t.Errorf("integerObj.Value not %d, got %d", expected, integerObj.Value)
		return false
	}

	return true
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return evaluator.Eval(program)
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	booleanObj, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("obj is not *object.Boolean. got=%[1]T, (%+[1]v)", obj)
		return false
	}

	if booleanObj.Value != expected {
		t.Errorf("booleanObj.Value is not %t, got %t", expected, booleanObj.Value)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},

		// 特殊だぞ！
		{"!5", false},

		// 2連続のやつ
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},

		// Integerならとにかくfalseになる仕様。!0でも!1でもとにかくfalse
		{"!1", false},
		{"!0", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		testBooleanObject(t, evaluated, tt.expected)
	}
}
