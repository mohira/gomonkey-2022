package parser

import (
	"fmt"
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
	fmt.Println(p)

}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	program.Statements = []ast.Statement{}

	fmt.Println(p.curToken)

	return program
}
