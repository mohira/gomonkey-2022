package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

// トークンタイプは有限 => 定数として定義できる
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"

	ASSIGN = "="
	PLUS   = "+"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "function"
	LET      = "let"
)

var kewwords = map[string]Type{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupIdent(ident string) Type {
	if tok, ok := kewwords[ident]; ok {
		return tok
	}

	return IDENT
}
