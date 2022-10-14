package parser

import (
	"fmt"
	"gomonkey/ast"
	"gomonkey/lexer"
	"gomonkey/token"
	"strconv"
)

const (
	_           int = iota // 0
	LOWEST                 // 1
	EQUALS                 // 2
	LESSGREATER            // 3
	SUM                    // 4
	PRODUCT                // 5
	PREFIX                 // 6
	CALL                   // 7 //myFunction(X)
)

var precedences = map[token.Type]int{
	token.EQ:     EQUALS,
	token.NOT_EQ: EQUALS,

	token.LT: LESSGREATER,
	token.GT: LESSGREATER,

	token.PLUS:  SUM,
	token.MINUS: SUM,

	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := Parser{
		l:      l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)

	// 前置演算式は ! と - の2種類だけです
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	// 2.8.2 グループ化された式に対応していく
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)

	return &p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("😢 次のトークンは %s になってほしいけど、 %s が来ちゃってる！", t, p.peekToken.Type)

	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
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
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	letStmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: セミコロンに遭遇するまで式を読み飛ばしている
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return letStmt
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken() // 1個進めている！ ← 期待通りなら1つすすめるのはよさそう(自然っぽい)
		return true
	} else {
		p.peekError(t) // 期待にそぐわなかったらエラーとして追加する
		return false
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStmt := &ast.ReturnStatement{
		Token:       p.curToken,
		ReturnValue: nil,
	}

	p.nextToken()

	// TODO: セミコロンに遭遇するまで読み飛ばしている
	if !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return returnStmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	defer untrace(trace("parseExpressionStatement()"))
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(curPrecedence int) ast.Expression {
	defer untrace(trace(fmt.Sprintf("parseExression() precedence=%d", curPrecedence)))
	// 現在のトークンに応じて、前置演算式のパース用の関数を探しにいく
	prefixFn := p.prefixParseFns[p.curToken.Type]

	if prefixFn == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefixFn() // ast.IntegerLiteral{3} になっている

	// 条件1: 次がセミコロンだったら、 `3;`みたいなやつだったってこと！
	// 条件2: この条件に当てはまる
	// => 次のトークンの処理を今よりも優先しろ！ ex: `3 + 4` だったら、+を先に処理する => 中置演算式としてパースしなさい
	for !p.peekTokenIs(token.SEMICOLON) && curPrecedence < p.peekPrecedence() {
		// 次のトークンのパースを優先しろって話なので、次のトークン用のパース関数を探しにいく。中置演算式のパースですね。
		infixFn := p.infixParseFns[p.peekToken.Type]

		if infixFn == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infixFn(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	defer untrace(trace("parseIntegerLiteral()"))
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) noPrefixParseFnError(t token.Type) {
	//msg := fmt.Sprintf("no prefix parse function for %s found", t)
	msg := fmt.Sprintf("👺 %s に対する前置演算のパースの関数がないよ！ マジで！", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	defer untrace(trace("parsePrefixExpression()"))
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX) // ここまじで意味わからん

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	defer untrace(trace(fmt.Sprintf("parseInfixExpression() left=%[1]T %[1]q", left)))
	// `3 + 4`
	infixExpr := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	curPrecedence := p.curPrecedence()

	p.nextToken()

	infixExpr.Right = p.parseExpression(curPrecedence)

	return infixExpr
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	boolean := &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}

	return boolean
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	// カッコがあるASTノードなんてものは、不要！ そんなのいらない！
	// 単純にこれでいい！
	// 	- カッコを剥がした式でパースする ← 閉じカッコは優先順位が最低値だから、どのみちカッコ内の演算子が優先されるので心配いらない
	// 	- 閉じカッコがちゃんとあるかを確かめる
	p.nextToken()

	expr := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expr
}
