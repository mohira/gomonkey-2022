package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	EOF     = ""
	ILLEGAL = "ILLEGAL"

	LET      = "let"
	FUNCTION = "fn"

	INT   = "INT"
	IDENT = "IDENT"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	SEMICOLON = ";"

	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	BANG     = "!"

	EQ     = "=="
	NOT_EQ = "!="

	ASSIGN = "="
)
