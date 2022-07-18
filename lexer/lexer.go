package lexer

import (
	"gomonkey/token"
)

type Lexer struct {
	input        string
	position     int  // 入力における現在の位置。現在の文字を指し示す。常に最後に読み込んだ場所を表す。
	readPosition int  // これから読み込む位置。現在の文字の"次"の文字を指し示す
	ch           byte // 現在検査中の文字
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // Lexer初期化時に最初の文字にセットしておく
	return l
}

// 文字を1つ進める
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // Q. 「0 は 何を意味する？」 → A. 0 は 終端(ASCIIコードの NUL文字)を表す
	} else {
		l.ch = l.input[l.readPosition] // 検査中の文字 を 次の文字 に移動する
	}
	l.position = l.readPosition // 現在の位置 を 次の位置 にずらす
	l.readPosition += 1         // 次の位置 を その次に ずらす

}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// ホワイトスペースは食べ尽くす
	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
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
	case '!':
		tok = newToken(token.BANG, l.ch)
	case 0:
		// newToken は 第2引数が byte なので使えない
		// l.ch は '\x00' が入っているので tok.Literal に代入してもダメ
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		// 記号トークン(？) 出ない場合は
		// 	文字(letter)か
		//  数値(number)か
		//	それ以外かで区別する
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	// ホワイトスペースに該当する 文字 であるかぎり、読み進める
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	// readIdentifier と同じ発想
	position := l.position // 最初の位置を覚えておく

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readIdentifier() string {
	position := l.position // 最初の位置を覚えておく

	// 現在の検査対象 が 文字 である限り読み込んでいく
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
