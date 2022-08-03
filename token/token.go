package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	EOF = ""

	LPAREN = "("
)
