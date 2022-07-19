package parser

import (
	"gomonkey/ast"
	"gomonkey/lexer"
	"gomonkey/token"
)

type Parser struct {
	l *lexer.Lexer

	// Lexer  は 1文字ずつ    見ていく
	// Parser は 1トークンずつ 見ていく。Parserはトークン単位で見ている。
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// 2つトークンを読み込む。トークンの読み込みの雰囲気はLexerと同じだね。
	// 2つ進めているので
	//
	// curToken
	// |
	// let age = 5;
	//     |
	//     peekToken
	// ↓↓ こういう感じ
	// 0) nil, nil
	// 1) nil, LET
	// 2) LET, IDENTIFIER
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken // Lexer.position = Lexer.read_position と同じ関係
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	// expectPeekは <IDENT> なら 読み進めちゃっているよ！
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// let <identifier> = <expression>
	// なわけだから、次に <ASSIGN> じゃないのはおかしいよね
	// → おかしな構文 → nil
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: セミコロンに遭遇するまで式を読み飛ばしてしまっている
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken() // ← 期待するTokenなら、読み進めちゃうぞ！！！！
		return true
	} else {
		return false
	}
}
