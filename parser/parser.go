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
	INDEX                  // 8
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

	token.LPAREN: CALL,

	token.LBRACKET: INDEX,
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

	// 2.8.3 if式
	p.registerPrefix(token.IF, p.parseIfExpression)

	// 2.8.2 グループ化された式に対応していく
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	// 2.8.4 関数リテラル
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	// 4.2.2 文字列リテラル
	p.registerPrefix(token.STRING, p.parseStringLiteral)

	// 4.4.2 配列リテラル
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)

	// 2.8.5 関数の呼び出し式
	// `(` を 中置演算式におけるOperatorだと思うってこと！
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	// 4.4.3 添字演算子式 ← 配列アクセスの式
	// myArray[0]
	// `[` を 中置演算式におけるOperatorだと思うってこと！
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
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

	// let x = 5;
	p.nextToken()
	letStmt.Value = p.parseExpression(LOWEST)

	// let文のセミコロンは省略できる！
	if p.peekTokenIs(token.SEMICOLON) {
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

	returnStmt.ReturnValue = p.parseExpression(LOWEST)

	// return文のセミコロンは省略できる！
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return returnStmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	//defer untrace(trace("parseExpressionStatement()"))
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(curPrecedence int) ast.Expression {
	//defer untrace(trace(fmt.Sprintf("parseExression() precedence=%d", curPrecedence)))

	// myArray[1 + 2]
	//   ↑

	// 現在のトークンに応じて、前置演算式のパース用の関数を探しにいく
	prefixFn := p.prefixParseFns[p.curToken.Type]

	if prefixFn == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefixFn() // *ast.Identifier{myArray}

	// 条件1: 次がセミコロンだったら、 `3;`みたいなやつだったってこと！
	// 条件2: この条件に当てはまる
	// => 次のトークンの処理を今よりも優先しろ！ ex: `3 + 4` だったら、+を先に処理する => 中置演算式としてパースしなさい
	for !p.peekTokenIs(token.SEMICOLON) && curPrecedence < p.peekPrecedence() {
		// 次のトークンのパースを優先しろって話なので、次のトークン用のパース関数を探しにいく。中置演算式のパースですね。
		infixFn := p.infixParseFns[p.peekToken.Type]

		// myArray[1 + 2]
		//    ↑
		if infixFn == nil {
			return leftExp
		}

		p.nextToken()
		// myArray[1 + 2]
		//        ↑
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
	//defer untrace(trace("parseIntegerLiteral()"))
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
	//defer untrace(trace("parsePrefixExpression()"))
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX) // ここまじで意味わからん

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	//defer untrace(trace(fmt.Sprintf("parseInfixExpression() left=%[1]T %[1]q", left)))
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

func (p *Parser) parseIfExpression() ast.Expression {
	// if (x < y) { x } else { y }
	ifExpr := &ast.IfExpression{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()

	ifExpr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	ifExpr.Consequence = p.parseBlockStatement()

	// elseがある場合の処理
	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // ELSEなう

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		ifExpr.Alternative = p.parseBlockStatement()
	}
	return ifExpr
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	blockStmt := &ast.BlockStatement{
		Token: p.curToken,
	}
	blockStmt.Statements = []ast.Statement{}

	p.nextToken()

	// } は ブロック終端ってことね
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			blockStmt.Statements = append(blockStmt.Statements, stmt)
		}
		p.nextToken()
	}

	return blockStmt

}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	// 	fn(x, y) { x + y; }
	functionLit := &ast.FunctionLiteral{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	functionLit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	functionLit.Body = p.parseBlockStatement()

	return functionLit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // token.Identifier -> token.COMMA
		p.nextToken() // token.COMMA      -> token.Identifier

		ident := &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return identifiers

	//// 自力で書いたパターンも残しておく
	//p.nextToken() // (の次に進んだ
	//
	//// ( または x に遭遇している
	//if p.curTokenIs(token.RPAREN) {
	//	return identifiers
	//}
	//
	//for !p.curTokenIs(token.RPAREN) {
	//	fmt.Printf("👉 %[1]T %[1]v\n", p.curToken)
	//	param := p.parseIdentifier()
	//	ident, ok := param.(*ast.Identifier)
	//	if !ok {
	//		return nil
	//	}
	//	identifiers = append(identifiers, ident)
	//
	//	if p.peekTokenIs(token.COMMA) {
	//		p.nextToken()
	//	}
	//	p.nextToken()
	//}
	//return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	callExpr := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}

	callExpr.Arguments = p.parseExpressionList(token.RPAREN)

	return callExpr
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arrayLiteral := &ast.ArrayLiteral{
		Token: p.curToken,
	}

	arrayLiteral.Elements = p.parseExpressionList(token.RBRACKET)

	return arrayLiteral
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	var args []ast.Expression

	// 引数がない場合
	if p.peekTokenIs(end) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	//           ↓
	// add(2, 3, 4)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return args

}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	// myArray[1 + 2]
	//        ↑
	indexExpr := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()

	// myArray[1 + 2]
	//         ↑
	indexExpr.Index = p.parseExpression(LOWEST)

	// myArray[1 + 2]
	//             ↑
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return indexExpr
}
