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

	l.skipWhitespace()

	switch l.ch {
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()

			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()

			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if l.isLetter() {
			literal := l.readIdentifier()
			tok.Type = token.LookupIdent(literal)
			tok.Literal = literal
			return tok
		} else if l.isDigit() {
			literal := l.readNumber()
			tok.Type = token.INT
			tok.Literal = literal
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()

	return tok
}

func newToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
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

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) isLetter() bool {
	return 'a' <= l.ch && l.ch <= 'z' || 'A' <= l.ch && l.ch <= 'Z' || l.ch == '_'
}

func (l *Lexer) readIdentifier() string {
	oldPosition := l.position

	for l.isLetter() {
		l.readChar()
	}

	return l.input[oldPosition:l.position]
}

func (l *Lexer) isDigit() bool {
	return '0' <= l.ch && l.ch <= '9'
}

func (l *Lexer) readNumber() string {
	oldPosition := l.position

	for l.isDigit() {
		l.readChar()
	}

	return l.input[oldPosition:l.position]
}
