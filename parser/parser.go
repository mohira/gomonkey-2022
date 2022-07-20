package parser

import (
	"gomonkey/ast"
	"gomonkey/lexer"
	"gomonkey/token"
)

type Parser struct {
	l *lexer.Lexer

	// Lexerのときと同じ作戦。Lexerは「1文字」単位で、Parserは「1トークン」単位ですね。
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// 初期化時に2回読み進めておく
	p.nextToken()
	p.nextToken()
	return p
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
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
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

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		// 期待通りならトークンを1つ読みすすめる
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}
