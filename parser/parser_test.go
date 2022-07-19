package parser

import (
	"gomonkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	// あとで追加

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
