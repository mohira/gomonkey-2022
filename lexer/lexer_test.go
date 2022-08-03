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
		t.Fatalf("TokenLiteralが違うよ。got=%s, want=%s", tok.Literal, expectedLiteral)
	}
}

func Test1文字トークンの字句解析(t *testing.T) {
	input := `(){}=+-*/!,;[]:<>`

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
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.LBRACKET, "["},
		{token.RBRACKET, "]"},
		{token.COLON, ":"},
		{token.LT, "<"},
		{token.GT, ">"},
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

func Test空白文字対処(t *testing.T) {
	input := `
( 
	) ! `

	tests := []struct {
		expectedToken   token.Type
		expectedLiteral string
	}{
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.BANG, "!"},
		{token.EOF, ""},
	}
	l := New(input)

	for _, tt := range tests {
		tok := l.NextToken()

		testToken(t, tok, tt.expectedToken, tt.expectedLiteral)
	}
}

func TestILLEGALなトークン(t *testing.T) {
	input := `@#$%`

	tests := []struct {
		expectedToken   token.Type
		expectedLiteral string
	}{
		{token.ILLEGAL, "@"},
		{token.ILLEGAL, "#"},
		{token.ILLEGAL, "$"},
		{token.ILLEGAL, "%"},
		{token.EOF, ""},
	}
	l := New(input)

	for _, tt := range tests {
		tok := l.NextToken()

		testToken(t, tok, tt.expectedToken, tt.expectedLiteral)
	}
}

func Test可変長文字数(t *testing.T) {
	input := `let x = 1;
let total = 234;
let add = fn(x, y) { x + y };

if (5 < 10) {
	return true;
} else {
	return false;
};
`

	tests := []struct {
		expectedToken   token.Type
		expectedLiteral string
	}{
		// let x = 1;
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},

		// let total = 234;
		{token.LET, "let"},
		{token.IDENT, "total"},
		{token.ASSIGN, "="},
		{token.INT, "234"},
		{token.SEMICOLON, ";"},

		// let add = fn(x, y) { x + y };
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		/*
			if (5 < 10) {
				return true;
			} else {
				return false;
			};
		*/
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}
	l := New(input)

	for _, tt := range tests {
		tok := l.NextToken()

		testToken(t, tok, tt.expectedToken, tt.expectedLiteral)
	}
}

func Test文字列(t *testing.T) {
	input := `
"Alice";
"hello world";
"";
" ";
"  ";
`
	//"123"
	//""←カラの文字列は！？
	tests := []struct {
		expectedToken   token.Type
		expectedLiteral string
	}{
		{token.STRING, "Alice"},
		{token.SEMICOLON, ";"},

		{token.STRING, "hello world"},
		{token.SEMICOLON, ";"},

		{token.STRING, ""},
		{token.SEMICOLON, ";"},

		{token.STRING, " "},
		{token.SEMICOLON, ";"},

		{token.STRING, "  "},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}
	l := New(input)

	for _, tt := range tests {
		tok := l.NextToken()

		testToken(t, tok, tt.expectedToken, tt.expectedLiteral)
	}
}
