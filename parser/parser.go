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

	p.registerPrefixFn(token.TRUE, p.parseBoolean)
	p.registerPrefixFn(token.FALSE, p.parseBoolean)

	p.registerPrefixFn(token.LPAREN, p.parseGroupedExpression)

	p.registerPrefixFn(token.IF, p.parseIfExpression)

	p.registerPrefixFn(token.FUNCTION, p.parseFunctionLiteral)

	// 中置演算式用の構文解析関数の用意
	p.infixFns = make(map[token.TokenType]parseInfixFn)

	p.registerInfixFn(token.EQ, p.parseInfixExpression)
	p.registerInfixFn(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixFn(token.LT, p.parseInfixExpression)
	p.registerInfixFn(token.GT, p.parseInfixExpression)
	p.registerInfixFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixFn(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixFn(token.SLASH, p.parseInfixExpression)

	// 関数呼び出し用？ よくわかってない
	p.registerInfixFn(token.LPAREN, p.parseCallExpression)
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

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	// return文 という Node を構築していくぞ
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	// TODO: <expression>のパースは後回しなので、セミコロンがくるで読み飛ばしている
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

func (p *Parser) registerPrefixFn(tokenType token.TokenType, fn parsePrefixFn) {
	p.prefixFns[tokenType] = fn
}

func (p *Parser) registerInfixFn(tokenType token.TokenType, fn parseInfixFn) {
	p.infixFns[tokenType] = fn
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	// MEMO: traceがちょっとうるさいし、うるさいから見ないのでOFFにしとく。Debugモードが有ると良いっぽい
	// defer untrace(trace("parseExpressionStatement <式文>"))
	exprStmt := &ast.ExpressionStatement{Token: p.curToken}

	exprStmt.Expression = p.parseExpression(LOWEST)

	// トークンの読み進めもお忘れなく ← 構文解析関数では読み進めない仕様にしている！
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return exprStmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// MEMO: traceがちょっとうるさいし、うるさいから見ないのでOFFにしとく。Debugモードが有ると良いっぽい
	// defer untrace(trace("parseExpression <式>"))
	prefixFn := p.prefixFns[p.curToken.Type]

	if prefixFn == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExpr := prefixFn()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// 「今の演算子の優先順位」が、「次の演算子の優先順位」より「小さい」なら、
		//

		infixFn := p.infixFns[p.peekToken.Type]

		if infixFn == nil {
			return leftExpr
		}

		p.nextToken()

		leftExpr = infixFn(leftExpr)
	}

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
	// MEMO: traceがちょっとうるさいし、うるさいから見ないのでOFFにしとく。Debugモードが有ると良いっぽい
	// defer untrace(trace("parseIntegerLiteral <INT>"))
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
	// MEMO: traceがちょっとうるさいし、うるさいから見ないのでOFFにしとく。Debugモードが有ると良いっぽい
	// defer untrace(trace("parsePrefixExpression <前置式>"))
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

// 演算の優先順位テーブル
var precedences = map[token.TokenType]int{
	token.EQ:     EQUALS,
	token.NOT_EQ: EQUALS,

	token.LT: LESSGREATER,
	token.GT: LESSGREATER,

	token.PLUS:  SUM,
	token.MINUS: SUM,

	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,

	token.LPAREN: CALL,
}

func (p *Parser) curPrecedence() int {
	if precedence, ok := precedences[p.curToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) parseInfixExpression(leftExpr ast.Expression) ast.Expression {
	// MEMO: traceがちょっとうるさいし、うるさいから見ないのでOFFにしとく。Debugモードが有ると良いっぽい
	// defer untrace(trace("parseInfixExpression <中置演算式>"))
	infixExpr := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     leftExpr,
		Operator: p.curToken.Literal,
		Right:    nil,
	}

	precedence := p.curPrecedence() // この時点ではcurTokenは何かしらの演算子トークンのはず
	p.nextToken()                   // だから、1つ読み進める

	infixExpr.Right = p.parseExpression(precedence)

	return infixExpr
}

func (p *Parser) parseBoolean() ast.Expression {
	b := &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE), // TRUEトークンかどうか調べればいいのか！ なるほど！
	}

	return b
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	// 「グループ化された式」は  `(<expression>)` という構造。
	// だから、 ( を見つけたらパース開始して、1つ飛ばして、expr解析して、次に)が来るはずってやるだけ。
	// まじで手品すぎる。
	p.nextToken() // `(` だから読み飛ばしていいよね。

	expr := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	// if ( <condition> ) { <consequence> }
	// if ( <condition> ) { <consequence> } else { <alternative> }
	ifExpr := &ast.IfExpression{
		Token:       p.curToken,
		Condition:   nil,
		Consequence: nil,
		Alternative: nil,
	}

	// curTokenが IF のはずだから、次に来るのは ( のはず(LPARENね)
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken() // LPAREN を 読み飛ばして <condition> へ。
	ifExpr.Condition = p.parseExpression(LOWEST)

	// expectPeek は トークン が期待通りなら進めるぞ！
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// expectPeek は トークン が期待通りなら進めるぞ！
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	ifExpr.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

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

	// この時点では、curTokenは LBRACE のはずだからね
	p.nextToken()

	// Blockの終端、つまり、 RBRACE が くるまでStatementを探せばいい(あとはEOF)
	// ここの構造は、 p.ParseProgram() と 一緒！ どちらも、複数のStatementを持つからね。
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
	// fn(x, y) { x + y; }
	// fn ( <parameters> ) { <block statement> }
	fnLit := &ast.FunctionLiteral{
		Token:      p.curToken,
		Parameters: nil,
		Body:       nil,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fnLit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	fnLit.Body = p.parseBlockStatement()

	return fnLit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier

	// 引数がない場合もある
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	// 1個めの引数
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	// 2個目以降があるときは、必ず「,」COMMAがあるはず！
	// 次に「,」が来る限り、それは引数だよね
	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // 今見ている識別子を読み飛ばして、
		p.nextToken() // 「,」も読み飛ばして、ようやく次の識別子だね

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	// 引数は、 ) で終わらないと文法おかしいからね
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	callExpr := &ast.CallExpression{
		Token:     p.curToken,
		Function:  function,
		Arguments: nil,
	}

	callExpr.Arguments = p.parseCallArguments()

	return callExpr
}

func (p *Parser) parseCallArguments() []ast.Expression {
	var args []ast.Expression

	// 引数が0個の場合はこれ
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // exprを読み飛ばし、
		p.nextToken() // カンマを読み飛ばし
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}
