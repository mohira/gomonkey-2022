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

}
