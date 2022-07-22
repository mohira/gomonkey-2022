package ast

import (
	"gomonkey/token"
	"testing"
)

func TestString(t *testing.T) {
	// 手動でASTを組み立てるのまじで良い！ 型やインタフェース、ASTの確認になる！
	// これが式の構文解析を行う場合に特に便利
	program := &Program{Statements: []Statement{
		&LetStatement{
			Token: token.Token{
				Type:    token.LET,
				Literal: "let",
			},
			Name: &Identifier{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "myVar",
				},
				Value: "myVar",
			},
			Value: &Identifier{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "anotherVar",
				},
				Value: "anotherVar",
			},
		},
	},
	}

	want := "let myVar = anotherVar;"
	if program.String() != want {
		t.Errorf("program.String() wrong. got=%q, want=%q ", program.String(), want)
	}
}

func TestString2(t *testing.T) {
	// 自力でASTを組み立てる練習
	// return result;
	program := &Program{Statements: []Statement{
		&ReturnStatement{
			Token: token.Token{
				Type:    token.RETURN,
				Literal: "return",
			},
			ReturnValue: &Identifier{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "result",
				},
				Value: "result",
			},
		},
	},
	}

	want := "return result;"
	if program.String() != want {
		t.Errorf("program.String()！ got=%q, want=%q", program.String(), want)
	}
}
