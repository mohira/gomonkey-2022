package parser

import (
	"gomonkey/ast"
	"gomonkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	if program == nil {
		t.Fatalf("👺 ParseProgram() return nil なのはおかしいよね")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("👺 program.Statements は 3つの Statement が含まれているべきだよ。でも、 got=%q", len(program.Statements))
	}

	// こっからParser本編って感じ
	tests := []struct {
		expectedIdentifier string
	}{
		// let文の左辺である <identifier> だけチェックする。 右辺の<expression>はいつかやるんでしょうね。
		// なぜ整数リテラル（ 5 、 10 など）が正しく構文解析されているかを確認し か？ 答えは、「あとでやる」だ。まずはlet文が正しく構文解析できるかを確かめる必要があるので、 Value には目をつぶる。
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]

		// ヘルパーメソッドを使って、今見ている <文> が let文 かどうかをチェックする
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	// let <identifier> = <expression>;

	// type LetStatement struct {
	//	Token token.Token // token.LET っていうトークンが入るだけじゃんね。
	//	Name  *Identifier // 要は、左辺の<識別子>
	//	Value Expression  // こっちは、右辺の<式>
	// }
	t.Helper()

	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
	}

	// LetStatement.Name は <identifier> であり、型でいうと ast.Identifier
	// type Identifier struct {
	//	    Token token.Token // 何かしらのトークンでもあるよね
	//		Value string      // <識別子> の "実際の値" とでも言えばいいかな。
	//	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	// <identifier> チェックね。
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}
