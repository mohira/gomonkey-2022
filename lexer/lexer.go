package lexer

import (
	"gomonkey/token"
)

type Lexer struct {
	input        string
	position     int  // 現在位置
	readPosition int  // 次の位置
	ch           byte // 現在位置の文字
}

func New(input string) *Lexer {
	l := &Lexer{input: input}

	l.readChar()

	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	switch l.ch {
	case '(':
		tok = token.Token{
			Type:    token.LPAREN,
			Literal: string(l.ch),
		}
	case ')':
		tok = token.Token{
			Type:    token.RPAREN,
			Literal: string(l.ch),
		}
	case '{':
		tok = token.Token{
			Type:    token.LBRACE,
			Literal: string(l.ch),
		}
	case '}':
		tok = token.Token{
			Type:    token.RBRACE,
			Literal: string(l.ch),
		}
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	}

	l.readChar()

	return tok
}

// 1文字読み進める
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}
