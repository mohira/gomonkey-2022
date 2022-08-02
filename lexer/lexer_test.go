package lexer

import (
	"gomonkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
let five = 5;
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token.Type が違うよ。got=%q, want=%q", i, tok.Type, tt.expectedType)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - token.Literal が違うよ。got=%q, want=%q", i, tok.Literal, tt.expectedLiteral)
		}
	}
}
