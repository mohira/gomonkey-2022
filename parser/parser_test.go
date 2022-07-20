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

}
