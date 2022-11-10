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
		t.Errorf("obj is not *object.Integer. got=%[1]T %[1]v", obj)
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
	env := object.NewEnvironment()

	return evaluator.Eval(program, env)
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
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
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
		{`"Hello" * 3`, "type mismatch: STRING * INTEGER"},

		// unknown operator: オペランド同士の型は一致しているが、演算子がおかしい
		{"true + false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{`"Hello" - "World"`, "unknown operator: STRING - STRING"},

		// 単項演算子
		{"-true", "unknown operator: -BOOLEAN"},

		// 実行時のエラーの後に文があるケースは、中断を実装する必要があるね？
		{"3; true + false; 4;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"3 + true; 4;", "type mismatch: INTEGER + BOOLEAN"},

		// 実験: % は ILLEGALなトークン(トークンとして認めてないのでParse時点で失敗する
		//		{"3 % 4;", "unknown operator: INTEGER % INTEGER"},

		// ERRORオブジェクト が演算に入っちゃって、エラーメッセージがちゃんとしなくなるのはダメだよ！ な、ケース。
		// オペランドにERRORオブジェクトが入るケース。
		// >> 1 + true
		// 💥 ERROR:type mismatch: INTEGER + BOOLEAN
		// >> - (1 + true)
		// 💥 ERROR:unknown operator: -ERROR
		{"- (1 + true)", "type mismatch: INTEGER + BOOLEAN"},
		{"(1 + true) + 2", "type mismatch: INTEGER + BOOLEAN"},
		{"1 + (true + 2)", "type mismatch: BOOLEAN + INTEGER"},

		// ERRORオブジェクトは実は truthy だった！
		// truthy := NULLでない かつ falseでない なので！！！
		{"if (1 + true) { return 2; }", "type mismatch: INTEGER + BOOLEAN"},

		// 未定義な識別子へのアクセス
		{"foobar; ", "identifier not found: foobar"},
		{"if (a) { 10; }", "identifier not found: a"},

		// 関数を呼び出したときに引数の過不足がある
		{"let add = fn(x, y) { return x + y;}; add(1, 2, 9)", "argument error: wrong number of arguments (given 3, expected 2)"},
		{"let add = fn(x, y) { return x + y;}; add(1)", "argument error: wrong number of arguments (given 1, expected 2)"},

		// ["foo" + 1]
		{`["foo" + 1]`, "type mismatch: STRING + INTEGER"},
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

func TestLetStatements(t *testing.T) {
	// let文の評価と、ちゃんと名前で値が保存されているかを確認するわけだよ
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},

		// 再代入ありなケース
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},

		// あとで: let文単体だった場合はどうするっけ?
		// {"let a = 5;", ?????}

		// あとで: 識別子未定義の場合のテストもどっかに書く

	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)

			testIntegerObject(t, evaluated, tt.expected)
		})
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)

	fnObj, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not *object.Function. got=%[1]T(%+[1]v)", evaluated)
	}

	if len(fnObj.Parameters) != 1 {
		t.Fatalf("パラメータ数が違うよ。got=%+v", fnObj.Parameters)
	}

	if fnObj.Parameters[0].String() != "x" {
		t.Fatalf("パラメータが 'x' じゃないよ！ got=%q", fnObj.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fnObj.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fnObj.Body.String())
	}

}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let double = fn(x) { x * 2 ;}; double(5);", 10},

		// 暗黙の戻り値
		{"let identity = fn(x) { x; }; identity(5);", 5},

		// return 文による値の返却
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; 10;}; identity(5);", 5},

		// 複数のパラメータ
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},

		// 関数は式です
		{"fn(x) { x; }(5)", 5},

		// 関数に渡す前の引数の評価
		// 関数の引数を前から順に評価しているので、
		// 第2引数の式を評価するときに、その第2引数が、また引数をもつような式だと、
		// xという名前で、評価済みの、第1引数を上書きしちゃうので、こわれる！
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"let add = fn(x, y) { x + y; }; add(5, add(90, 10));", 105},

		// こういう外側の環境が必要なやつは、またあとできっとやるでしょう
		//{"let a=1; fn(x){ a + x;} ", 5},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)
			testIntegerObject(t, evaluated, tt.expected)
		})
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
	fn(y) { x + y };
};

let addTwo = newAdder(2);

fn(x) {
	addTwo(3);
}(9999)
`
	evaluated := testEval(input)

	testIntegerObject(t, evaluated, 5)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello, world!"`

	evaluated := testEval(input)

	strObj, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("*object.Stringじゃないよ.got=%[1]T(%+[1]v)", evaluated)
	}

	if strObj.Value != "Hello, world!" {
		t.Errorf("want %s, got %s", "Hello, world!", strObj.Value)
	}
}

func TestStringExpression(t *testing.T) {
	tests := []struct {
		input string
		want  any
	}{
		{`"Hello, " + "world!"`, "Hello, world!"},
		{`"Hello" == "Hello"`, true},
		{`"Hello" == "World"`, false},

		{`"Hello" != "Hello"`, false},
		{`"Hello" != "World"`, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch tt.want.(type) {
			case string:
				strObj, ok := evaluated.(*object.String)
				if !ok {
					t.Fatalf("*object.Stringじゃないよ.got=%[1]T(%+[1]v)", evaluated)
				}

				if strObj.Value != tt.want {
					t.Errorf("want %s, got %s", "Hello, world!", strObj.Value)
				}
			case bool:
				boolObj, ok := evaluated.(*object.Boolean)
				if !ok {
					t.Fatalf("*object.Booleanじゃないよ.got=%[1]T(%+[1]v)", evaluated)
				}

				if boolObj.Value != tt.want {
					t.Errorf("want %t, got %t", tt.want, boolObj.Value)
				}

			}

		})
	}

}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},

		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "argument error: wrong number of arguments (given 2, expected 1)"},
		{`len()`, "argument error: wrong number of arguments (given 0, expected 1)"},

		// TODO: 組み込み関数は事前宣言されているだけで、あとから上書きありにする？ どうする？
		// {`let len = fn(x) {return 1;}; len("hoge")`, 1},

		// lenで配列の要素数もわかるよ
		{"len([1, 2, 3])", 3},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case int:
				testIntegerObject(t, evaluated, int64(expected))
			case string:
				errObj, ok := evaluated.(*object.Error)
				if !ok {
					t.Fatalf("評価結果がERRORオブジェクトになってないよ. got=%[1]T(%+[1]v)", evaluated)
				}

				if errObj.Message != expected {
					t.Errorf("おかしいよ。expected=%s, got=%s", expected, errObj.Message)
				}
			}

		})
	}

}

func TestArrayLiteral(t *testing.T) {
	input := "[1, 2 * 3, 4 + 5]"

	evaluated := testEval(input)

	array, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("*object.Array じゃないよ。got=%T", evaluated)
	}

	if len(array.Elements) != 3 {
		t.Errorf("要素数が 3 じゃないよ。got=%d", len(array.Elements))
	}

	// [1, 2 * 3, 4 + 5] -> [1, 6, 9]
	testIntegerObject(t, array.Elements[0], 1)
	testIntegerObject(t, array.Elements[1], 6)
	testIntegerObject(t, array.Elements[2], 9)
}

func TestArrayIndexExpressions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected any
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},

		//
		{"let i = 0; [1][i]", 1},
		{"[1, 2, 3][1 + 1]", 3},

		//
		{"let myArray = [1, 2, 3]; myArray[2];", 3},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i];", 2},

		// off-by-one error
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)

			expectedInteger, ok := tt.expected.(int)
			if ok {
				testIntegerObject(t, evaluated, int64(expectedInteger))
			} else {
				testNullObject(t, evaluated)
			}
		})
	}

}
