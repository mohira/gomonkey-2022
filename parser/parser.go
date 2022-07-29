package parser

import (
	"fmt"
	"gomonkey/ast"
	"gomonkey/lexer"
	"gomonkey/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // =
	LESSGREATER // > または <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X または !X
	CALL        // myFunction(X)
)

type Parser struct {
	l *lexer.Lexer

	// Lexerのときと同じ作戦。Lexerは「1文字」単位で、Parserは「1トークン」単位ですね。
	curToken  token.Token
	peekToken token.Token

	// デバッグを楽にするため error を記録しておく
	errors []string // stringが楽で良さげ

	// トークンタイプごとの構文解析関数をもっておくmap
	prefixFns map[token.TokenType]parsePrefixFn
	infixFns  map[token.TokenType]parseInfixFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// 初期化時に2回読み進めておく
	p.nextToken()
	p.nextToken()

	// 前置の構文解析関数のmap初期化
	p.prefixFns = make(map[token.TokenType]parsePrefixFn)

	p.registerPrefixFn(token.IDENT, p.parseIdentifier)
	p.registerPrefixFn(token.INT, p.parseIntegerLiteral)
	p.registerPrefixFn(token.BANG, p.parsePrefixExpression)
	p.registerPrefixFn(token.MINUS, p.parsePrefixExpression)

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
		return p.parseExpressionStatement()
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

type (
	parsePrefixFn func() ast.Expression
	// <left> <op> <right> の <left> が引数な
	parseInfixFn func(ast.Expression) ast.Expression
)

func (p Parser) registerPrefixFn(tokenType token.TokenType, fn parsePrefixFn) {
	p.prefixFns[tokenType] = fn
}

func (p Parser) registerInfixFn(tokenType token.TokenType, fn parseInfixFn) {
	p.infixFns[tokenType] = fn
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	exprStmt := &ast.ExpressionStatement{Token: p.curToken}

	exprStmt.Expression = p.parseExpression(LOWEST)

	// トークンの読み進めもお忘れなく ← 構文解析関数では読み進めない仕様にしている！
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return exprStmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFn := p.prefixFns[p.curToken.Type]

	if prefixFn == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExpr := prefixFn()

	return leftExpr
}

func (p *Parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	return ident
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("colud not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	intLit := &ast.IntegerLiteral{
		Token: p.curToken,
		Value: value,
	}

	return intLit

}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	// 適応する前置用の構文解析関数が見つからないときにエラー情報を格納するやつ
	msg := fmt.Sprintf("トークンタイプ '%s' 用のPrefixParseFnがありませんよ？", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	prefixExpr := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	// この関数が呼ばれている次点で、 token.BANG または token.MINUS なので(そう指定したから)
	// そのままトークンを読み進めればおk
	p.nextToken()

	prefixExpr.Right = p.parseExpression(PREFIX)

	return prefixExpr
}
