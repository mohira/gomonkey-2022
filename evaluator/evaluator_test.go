package evaluator_test

import (
	"gomonkey/evaluator"
	"gomonkey/lexer"
	"gomonkey/object"
	"gomonkey/parser"
	"testing"
)

// 私は「-」前置演算子のために新しいテスト関数を書くのではなく、このテストを拡張することにした。
// それには2つ理由がある。
// 第一に、前置の「-」演算子がサポートするオペランドは整数だけだからだ。
// 第二に、このテスト関数は全ての整数演算を含むように成長させ、期待する振る舞いを明確で整理された書き方で1つの場所にまとめておくためだ。
func TestIntegerExpression(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected int64
	}{
		// ast.IntegerLiteralなやつ
		{"5", 5},
		{"10", 10},

		// ast.PrefixExpressionなやつ
		{"-5", -5},
		{"-10", -10},

		// もうちょいテストケースを筋肉質にできると思うけど？
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			evaluated := testEval(tt.input)

			testIntegerObject(t, evaluated, tt.expected)

		})
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
	t.Parallel()

	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},

		// 比較演算とか
		{"1 < 2", true},
		{"1 < 1", false},
		{"2 < 1", false},

		{"2 > 1", true}, // 追加してやったぞ！
		{"1 > 2", false},
		{"1 > 1", false},

		{"1 == 1", true},
		{"1 == 2", false},

		{"1 != 2", true},
		{"1 != 1", false},

		// == と != だけサポートしているよ！
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},

		{"true != true", false},
		{"false != false", false},
		{"true != false", true},

		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			evaluated := testEval(tt.input)
			testBooleanObject(t, evaluated, tt.expected)
		})
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
	t.Parallel()

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
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			evaluated := testEval(tt.input)

			testBooleanObject(t, evaluated, tt.expected)
		})
	}
}

func TestIfElseExpressions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected any
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},

		// consequence が複数のstatement
		{"if (true) { 10; 20; } else { 30 }", 20},
		{"if (true) { 10; 20; }", 20},

		{"if (false) { 10 } else { 20; 30; }", 30},

		// 現状では 識別子(*ast.Identifier)を評価する用になってないから、
		// conditionがnilになる。(monkeyのNULLではない！)
		// ホスト言語のnilは、monkeyのNULLでもなけれあFALSEでもないので、truthyになる
		{"if (a) { 10 }", 10},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			evaluated := testEval(tt.input)
			integer, ok := tt.expected.(int) // 最初から int64 じゃだめなの？ あとで試す
			if ok {
				testIntegerObject(t, evaluated, int64(integer))
			} else {
				testNullObject(t, evaluated)
			}
		})
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != evaluator.NULL {
		t.Errorf("obj is not NULL. got=%[1]T (%+[1]v)", obj)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},

		// return したら次の文は評価しないよね
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},

		// Returnオブジェクトという存在の意義
		// Returnオブジェクトという存在をつくることで
		// 1. 「私をReturnしてください」という情報をオブジェクト自身で表現できる！
		//   いいかえると、ネストレベルと言うか、外部の情報がなくても管理できる！
		//   RETURN(INT(10))みたいになってるので、オブジェクトからReturnしてほしいアピールがある感じ！
		//   木構造と再帰の相性のパワーも活かせる感じある！
		// 2. 逆に INT(10) という、結局の評価結果を扱うスタイルにすると
		// リターンするかどうかの情報をたさないと行けない！ ← ネストレベルみたいな概念が必要
		// それでもたぶんいけるんだけど、実装は複雑になると思う。
		{`
if (10 > 1) {
    if (10 > 1) {
        return 10;
    }
    
    return 1;
}`, 10},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			evaluated := testEval(tt.input)
			testIntegerObject(t, evaluated, tt.expected)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		// type mismatch: オペランド同士の型が一致していない
		{"3 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"true + 3;", "type mismatch: BOOLEAN + INTEGER"},

		// unknown operator: オペランド同士の型は一致しているが、演算子がおかしい
		{"true + false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},

		// 単項演算子
		{"-true", "unknown operator: -BOOLEAN"},

		// 実行時のエラーの後に文があるケースは、中断を実装する必要があるね？
		{"3; true + false; 4;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"3 + true; 4;", "type mismatch: INTEGER + BOOLEAN"},

		// 実験: % は ILLEGALなトークン(トークンとして認めてないのでParse時点で失敗する
		//		{"3 % 4;", "unknown operator: INTEGER % INTEGER"},

		// >> 1 + true
		// 💥 ERROR:type mismatch: INTEGER + BOOLEAN
		// >> - (1 + true)
		// 💥 ERROR:unknown operator: -ERROR
		{"- (1 + true)", "type mismatch: INTEGER + BOOLEAN"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)

			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("Errorオブジェクトじゃないよ！ got=%[1]T(%+[1]v)", evaluated)
				return
			}

			if errObj.Message != tt.expectedMessage {
				t.Errorf("エラーオブジェクトのMessageがちがうよ！ want=%s, got=%s", tt.expectedMessage, errObj.Message)
			}

		})
	}
}
