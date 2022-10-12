package ast_test

import (
	"gomonkey/ast"
	"gomonkey/token"
	"testing"
)

// p.54のデモ; ノードに String() が実装されていることで *ast.Program の構造を簡単にテストできる話

func TestString(t *testing.T) {
	// let myVar = anotherVar; というlet文を組み立てる
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},

				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},

				Value: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	want := "let myVar = anotherVar;"
	if program.String() != want {
		t.Errorf("program.String() wrong. want=%s, got=%q", want, program.String())
	}
}
