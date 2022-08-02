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

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = "<"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "function"
	LET      = "let"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "return"
	TRUE     = "true"
	FALSE    = "false"

	EQ     = "=="
	NOT_EQ = "!="
)

var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
