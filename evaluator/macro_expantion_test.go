package evaluator_test

import (
	"gomonkey/ast"
	"gomonkey/evaluator"
	"gomonkey/lexer"
	"gomonkey/object"
	"gomonkey/parser"
	"testing"
)

func TestDefineMacros(t *testing.T) {
	input := `
let number = 1;
let function = fn(x, y) { x + y }; 
let mymacro = macro(x, y) { x + y };
`
	env := object.NewEnvironment()
	program := testParseProgram(input)

	// マクロ定義をみつけて、そのマクロを環境に登録して、ASTから消し去る。なかなか色んな仕事をするやつ
	evaluator.DefineMacros(program, env)

	if len(program.Statements) != 2 {
		t.Fatalf("mymacroをASTから消し去るので、2文になってないといけませんね？ got=%d", len(program.Statements))
	}

	// Evalしてないので、変数numberや変数functionは環境にはありませんね？
	if _, ok := env.Get("number"); ok {
		t.Fatalf("変数numberが環境に登録されているのはおかしいぜ！？")
	}
	if _, ok := env.Get("function"); ok {
		t.Fatalf("変数functionが環境に登録されているのはおかしいぜ！？")
	}

	obj, ok := env.Get("mymacro")
	if !ok {
		t.Fatalf("mymacro が 環境にないのはダメだぞ！")
	}

	// 環境に登録されているマクロがちゃんと期待通りかチェックする(結構かったるいチェックです)
	macro, ok := obj.(*object.Macro)
	if !ok {
		t.Fatalf("マクロオブジェクトじゃないがな！ got=%[1]T(%+[1]v)", obj)
	}

	if len(macro.Parameters) != 2 {
		t.Fatalf("おかしいぜ")
	}

	if macro.Parameters[0].String() != "x" {
		t.Fatalf("おかしいぜ")
	}
	if macro.Parameters[1].String() != "y" {
		t.Fatalf("おかしいぜ")
	}

	expectedBody := "(x + y)"
	if macro.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, macro.Body.String())
	}
}

func testParseProgram(input string) *ast.Program {
	// 評価はしないやつ！
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
