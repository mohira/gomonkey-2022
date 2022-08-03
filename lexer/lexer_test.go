package lexer

import (
	"gomonkey/token"
	"testing"
)

func testToken(t *testing.T, tok token.Token, expectedToken token.Type, expectedLiteral string) {
	t.Helper()

	if tok.Type != expectedToken {
		t.Fatalf("TokenTypeが違うよ。got=%q, want=%q", tok.Type, expectedToken)
	}
	if tok.Literal != expectedLiteral {
		t.Fatalf("TokenLiteralが違うよ。got=%s, want=%s", tok.Type, expectedToken)
	}
}

func Test1文字トークンの字句解析(t *testing.T) {
	input := `(){}=+-*/!;`

	tests := []struct {
		expectedToken   token.Type
		expectedLiteral string
	}{
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.BANG, "!"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := New(input)

	for _, tt := range tests {
		tok := l.NextToken()

		testToken(t, tok, tt.expectedToken, tt.expectedLiteral)
	}
}

func Test2文字トークンの字句解析(t *testing.T) {
	input := `==!=`

	tests := []struct {
		expectedToken   token.Type
		expectedLiteral string
	}{
		{token.EQ, "=="},
		{token.NOT_EQ, "!="},
		{token.EOF, ""},
	}
	l := New(input)

	for _, tt := range tests {
		tok := l.NextToken()

		testToken(t, tok, tt.expectedToken, tt.expectedLiteral)
	}
}
