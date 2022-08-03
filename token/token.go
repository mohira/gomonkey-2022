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
	RETURN   = "return"

	INT   = "INT"
	IDENT = "IDENT"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	SEMICOLON = ";"
	COMMA     = ","

	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	BANG     = "!"

	EQ     = "=="
	NOT_EQ = "!="

	ASSIGN = "="
)

var keywords = map[string]Type{
	"let":    LET,
	"fn":     FUNCTION,
	"return": RETURN,
}

func LookupIdent(key string) Type {
	if v, ok := keywords[key]; ok {
		return v
	}

	return IDENT
}
