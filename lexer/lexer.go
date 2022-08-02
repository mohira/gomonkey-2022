package lexer

import (
	"fmt"
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

	// 最初の文字を読んでおく
	l.readChar()

	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	fmt.Println(l.readPosition, string(l.ch))
	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	}

	l.readChar()

	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) { // 終端チェックってこと
		// 次の文字位置が入陸文字数より多いなら終了だよね
		l.ch = 0
	} else {
		// まだ読んでいない文字があるなら、「現在検査中の文字」を次の文字にする
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition // 現在の文字位置 を 次の文字位置 に進める
	l.readPosition += 1         // 「次の文字位置」を更に次の位置に進める
}

func newToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}
