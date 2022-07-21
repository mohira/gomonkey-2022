package parser

import (
	"gomonkey/ast"
	"gomonkey/lexer"
	"testing"
)

func TestParseLetStatements(t *testing.T) {
	// let文の <expression> 以外のところをパースする。<expression>のパースはあとでやるよ。ムズいんでね。

	// 3つの let文 だけが含まれる正しいプログラムですね。
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	// エラーを確かめるための実験(あとで消してね)
	//	input = `
	//let x 5;
	//let = 10;
	//let 838383;
	//`
	l := lexer.New(input)
	p := New(l)

	// それぞれの let文 をパースする前に
	// 1. ast.Programノードが作れているかを確認する
	// 2. このProgram は 3つの(何かしらの)Statement からなることを確認する

	program := p.ParseProgram()
	checkParserErrors(t, p)

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

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
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

func TestParseReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	// それぞれのreturn文をParseする前に Statement が 3つある ことを確認する
	if len(program.Statements) != 3 {
		t.Fatalf("Statementは3つじゃないとおかしいね. got=%q", len(program.Statements))
	}

	// ここから1文ずつ検証する
	for _, stmt := range program.Statements {
		// まず、return文 なのかを確かめる
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", returnStmt)
			continue
		}
	}
}
