package parser

import (
	"gomonkey/ast"
	"gomonkey/lexer"
	"testing"
)

func TestParseLetStatement(t *testing.T) {
	// let文の <expression> 以外のところをパースする。<expression>のパースはあとでやるよ。ムズいんでね。

	// 3つの let文 だけが含まれる正しいプログラムですね。
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)

	// それぞれの let文 をパースする前に
	// 1. ast.Programノードが作れているかを確認する
	// 2. このProgram は 3つの(何かしらの)Statement からなることを確認する

	program := p.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned nil なのはおかしいよね'")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("3つのStatementからなるProgramじゃないのはおかしいよね。got = %q", len(program.Statements))
	}

	// ここからが let文 のパースって感じです。
	// expressionは後回しにして、 <identifier> がちゃんと解析できているかを、まずは、調べていくぞ。
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]

		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}

}

func testLetStatement(t *testing.T, stmt ast.Statement, expectedName string) bool {
	// <expression> つまり LetStatement.Value の検証はムズいので後回しだよ
	t.Helper()

	// まずは LETトークンかどうかをちゃんと調べる
	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral() が 'let` じゃないよ。 got = %q", stmt.TokenLiteral())
		return false
	}

	// LETトークンをもっていても、LetStatement型じゃない可能性もあるのか
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt が LetStatement じゃないよ。got=%Tl", stmt)
		return false
	}

	// ここまで来たら ast.LetStatement 確定なので、素直に属性としての Name(つまり、<identifier>)を調べる
	if letStmt.Name.Value != expectedName {
		t.Errorf("letStmt.Name.Value が %q じゃないよ。got=%q", letStmt.Name.Value, expectedName)
		return false
	}

	// よくわからんけど、LetStatement.Name.Value だけじゃなくて LetStatement.Name.TokenLiteral()もしらべてる
	if letStmt.Name.TokenLiteral() != expectedName {
		t.Errorf("letStmt.Name.TokenLiteral() が %q じゃないよ。got=%q", letStmt.Name.TokenLiteral(), expectedName)
		return false
	}

	return true
}
