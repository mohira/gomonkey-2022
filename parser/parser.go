package parser

import (
	"fmt"
	"gomonkey/ast"
	"gomonkey/lexer"
	"gomonkey/token"
)

type Parser struct {
	l *lexer.Lexer

	// Lexerのときと同じ作戦。Lexerは「1文字」単位で、Parserは「1トークン」単位ですね。
	curToken  token.Token
	peekToken token.Token

	// デバッグを楽にするため error を記録しておく
	errors []string // stringが楽で良さげ
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// 初期化時に2回読み進めておく
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("👺 次のトークンは %s になってほしいけど、いまんところ %s になっちゃってるよ", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// ここらへん、擬似コードまんま！
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()

		// 文が来るとは限らないからね。式だったり、ILLEGALだったりするからね
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	// Statement の Parse の エントリーポイント ってかんじだね
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()

	default:
		return nil
	}
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	// チラ見して、期待通りならトークンを1つ読みすすめる。そうでなければ、何もしない。
	// 構文解析器のよくある動作らしいね
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	// let文 という Node を構築していく感じ。
	stmt := &ast.LetStatement{Token: p.curToken}

	// 次のトークンは <identifier> であってほしいからね。
	if !p.expectPeek(token.IDENT) {
		return nil // いまんところErrorじゃなくてnilで
	}

	// <identifier> を 登録(？)すればいい
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// let文なら、つぎのトークンは 「=」 だよね
	if !p.expectPeek(token.ASSIGN) {
		return nil // いまんところErrorじゃなくてnilで
	}

	// TOのDO: <expression>のパースを後回しにするので、セミコロンがくるまで読み飛ばしている。
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	// return文 という Node を構築していくぞ
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TOOD: <expression>のパースは後回しなので、セミコロンがくるで読み飛ばしている
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
