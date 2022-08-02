package lexer

import (
	"gomonkey/token"
)

type Lexer struct {
	input        string
	position     int  // 現在の文字位置を指し示す
	readPosition int  // 次の文字位置を指し示す
	ch           byte // 現在検査中の文字
}

func New(input string) *Lexer {
	l := &Lexer{input: input}

	return l
}

func (l Lexer) NextToken() token.Token {
	return token.Token{}
}
