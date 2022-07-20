package parser

import (
	"gomonkey/lexer"
	"testing"
)

func TestParseLetStatement(t *testing.T) {
	// let文の <expression> 以外のところをパースする。<expression>のパースはあとでやるよ。ムズいんでね。

	// 3つの let文 だけが含まれる正しいプログラムですね。
	input := `
let x = 5;
let y = 10;
let foobar = 838383
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

}
