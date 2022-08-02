package parser

import (
	"fmt"
	"gomonkey/ast"
	"gomonkey/lexer"
	"testing"
)

func TestParseLetStatements(t *testing.T) {
	// let文の <expression> 以外のところをパースする。<expression>のパースはあとでやるよ。ムズいんでね。

	// 3つの let文 だけが含まれる正しいプログラムですね。
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	// エラーを確かめるための実験(あとで消してね)
	//	input = `
	//let x 5;
	//let = 10;
	//let 838383;
	//`
	l := lexer.New(input)
	p := New(l)

	// それぞれの let文 をパースする前に
	// 1. ast.Programノードが作れているかを確認する
	// 2. このProgram は 3つの(何かしらの)Statement からなることを確認する

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil なのはおかしいよね'")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("3つのStatementからなるProgramじゃないのはおかしいよね。got = %q", len(program.Statements))
	}

	// ここからが let文 のパースって感じです。
	// expressionは後回しにして、 <identifier> がちゃんと解析できているかを、まずは、調べていくぞ。
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]

		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}

}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func testLetStatement(t *testing.T, stmt ast.Statement, expectedName string) bool {
	// <expression> つまり LetStatement.Value の検証はムズいので後回しだよ
	t.Helper()

	// まずは LETトークンかどうかをちゃんと調べる
	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral() が 'let` じゃないよ。 got = %q", stmt.TokenLiteral())
		return false
	}

	// LETトークンをもっていても、LetStatement型じゃない可能性もあるのか
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt が LetStatement じゃないよ。got=%Tl", stmt)
		return false
	}

	// ここまで来たら ast.LetStatement 確定なので、素直に属性としての Name(つまり、<identifier>)を調べる
	if letStmt.Name.Value != expectedName {
		t.Errorf("letStmt.Name.Value が %q じゃないよ。got=%q", letStmt.Name.Value, expectedName)
		return false
	}

	// よくわからんけど、LetStatement.Name.Value だけじゃなくて LetStatement.Name.TokenLiteral()もしらべてる
	if letStmt.Name.TokenLiteral() != expectedName {
		t.Errorf("letStmt.Name.TokenLiteral() が %q じゃないよ。got=%q", letStmt.Name.TokenLiteral(), expectedName)
		return false
	}

	return true
}

func TestParseReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	// それぞれのreturn文をParseする前に Statement が 3つある ことを確認する
	if len(program.Statements) != 3 {
		t.Fatalf("Statementは3つじゃないとおかしいね. got=%q", len(program.Statements))
	}

	// ここから1文ずつ検証する
	for _, stmt := range program.Statements {
		// まず、return文 なのかを確かめる
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", returnStmt)
			continue
		}

		// returnStmtがちゃんと "return"トークンを持っているかを調べる
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral() が 'return' になっていない。got=%q", returnStmt.TokenLiteral())
		}
	}
}

func TestParseIdentifier(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements が 1文のみになってないよ！ got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("その<文> は ast.ExpressionStatement <式文> になってないぞ！ got=%T", program.Statements[0])
	}

	ident, ok := exprStmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("その<式文> は <IDENT> になってないよ！ got=%T", exprStmt.Expression)
	}

	// フィールド検証
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() not %s, got %s", "foobar", ident.TokenLiteral())
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s, got %s", "foobar", ident.Value)
	}
}

func TestParseIntegerLiteral(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements が 1文のみになってないよ！ got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("その<文> は ast.ExpressionStatement <式文> になってないぞ！ got=%T", program.Statements[0])
	}

	intLit, ok := exprStmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("ast.IntegerLiteral に変換できてないよ. got=%T", exprStmt.Expression)
	}

	// フィールドの検証
	if intLit.TokenLiteral() != "5" {
		t.Errorf("intLit.TokenLiteral() not %s, got=%s", "5", intLit.TokenLiteral())
	}

	if intLit.Value != 5 {
		t.Errorf("intLit.Value not %d, got=%d", 5, intLit.Value)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue any // ついになんでも！
	}{
		{"!5;", "!", 5}, // !5 が何を返すかは多分決めてないと思う。あくまで前置演算式って構造だけ。
		{"-15;", "-", 15},

		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("1 Statement になってないのはおかしいよ, got=%d", len(program.Statements))
		}

		exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("ast.ExpressionStatmentに変換できないよ, got=%T", program.Statements[0])
		}

		prefixExpr, ok := exprStmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("ast.PrefixExrepssionに変換できないよ, got=%T", exprStmt.Expression)
		}

		// ast.PrefixExpressionのフィールド検証
		if prefixExpr.Operator != tt.operator {
			t.Errorf("Opeartorが違うよ. got=%s, want=%s", prefixExpr.Operator, tt.operator)
		}

		if testLiteralExpression(t, prefixExpr.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integerLiteral, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLitearl, got=%T", il)
		return false
	}

	if integerLiteral.Value != value {
		t.Errorf("integerLiteral not %d, got=%d", integerLiteral.Value, value)
		return false
	}

	if integerLiteral.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integerLiteral.TokenLiteral not %d, got=%s", value, integerLiteral.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	// <expr> <infix op> <expr>; だけど
	// このテストにおける <expr> は、具体的には、IntegerLiteralのみ
	infixTests := []struct {
		input             string
		LeftIntegerValue  any
		Operator          string
		RightIntegerValue any
	}{
		{"3 + 4;", 3, "+", 4},
		{"3 - 4;", 3, "-", 4},
		{"3 * 4;", 3, "*", 4},
		{"3 / 4;", 3, "/", 4},
		{"3 > 4;", 3, ">", 4},
		{"3 < 4;", 3, "<", 4},
		{"3 == 4;", 3, "==", 4},
		{"3 != 4;", 3, "!=", 4},

		// BooleanLiteralあり
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("1 Statement じゃないよ！ got=%d", len(program.Statements))
		}

		exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("おかしいよ. got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, exprStmt.Expression, tt.LeftIntegerValue, tt.Operator, tt.RightIntegerValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	// ASTの文字列表現が優先順位を含んだカッコの感じになっているかチェック。LISP感。
	tests := []struct {
		input             string
		expectedStringAST string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d", "((a + (b * c)) + d)"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},

		// booleanまじり
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},

		// グループ化された式。つまり、カッコによってAST構築の優先順位が変わることを確かめる
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},

		// 関数の呼び出し式のカッコをうまく処理できるか
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8)) ", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expectedStringAST {
			t.Errorf("expected=%q, got=%q", tt.expectedStringAST, actual)
		}

	}
}

func TestOperatorPrecedenceParsing確認用(t *testing.T) {
	// ASTの文字列表現が優先順位を含んだカッコの感じになっているかチェック。LISP感。
	input := "1 + 2 + 3"
	want := "((1 + 2) + 3)"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	actual := program.String()
	if actual != want {
		t.Errorf("expected=%q, got=%q", want, actual)
	}
}

func testIdentifier(t *testing.T, expr ast.Expression, value string) bool {
	t.Helper()

	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Errorf("expr not *ast.Identifier. got=%T", expr)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, expr ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case int64:
		return testIntegerLiteral(t, expr, v)
	case string:
		return testIdentifier(t, expr, v)
	case bool:
		return testBooleanLiteral(t, expr, v)
	}

	t.Errorf("type of expr not handled. まだ対応できないんだわ。got=%T", expr)

	return true
}

func testInfixExpression(t *testing.T, expr ast.Expression, left any, operator string, right any) bool {
	infixExpr, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expr is not ast.InfixExpression. got=%T(%s)", expr, expr)
		return false
	}

	if !testLiteralExpression(t, infixExpr.Left, left) {
		return false
	}

	if infixExpr.Operator != operator {
		t.Errorf("expr.Operator is not %s. got=%s", operator, infixExpr.Operator)
		return false
	}

	if !testLiteralExpression(t, infixExpr.Right, right) {
		return false
	}

	return true
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("1 Statementじゃないぞ！ got=%d", len(program.Statements))
		}

		exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("*ast.ExpressionStatementに変換できないぞ。got=%T", program.Statements[0])
		}

		boolean, ok := exprStmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("*ast.BooleanLiteralに変換できないぞ。got=%T", exprStmt.Expression)
		}

		// フィールド検証
		if boolean.Value != tt.expectedValue {
			t.Errorf("boolean.Value not %t got=%t", tt.expectedValue, boolean.Value)
		}
	}
}

func testBooleanLiteral(t *testing.T, expr ast.Expression, value bool) bool {
	b, ok := expr.(*ast.Boolean)
	if !ok {
		t.Errorf("expr not *ast.Boolean got=%T", expr)
		return false
	}

	if b.Value != value {
		t.Errorf("b.Value not %t. got=%t", value, b.Value)
		return false
	}

	if b.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("b.TokenLiteral not %t got %s", value, b.TokenLiteral())
		return false
	}

	return true
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("1 Statement じゃないぞ！ got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("*ast.ExpressionStatementに変換できないぞ。got=%T", program.Statements[0])
	}

	ifExpr, ok := exprStmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("*ast.IfExpressonに変換できないぞ。got=%T", exprStmt.Expression)
	}

	// フィールド検証
	// input は `if (x < y) { x }` なので、 <condition> は x < y という中置演算式

	if !testInfixExpression(t, ifExpr.Condition, "x", "<", "y") {
		return
	}

	// 1つ目のBlockの<文>の数も数えないとね
	if len(ifExpr.Consequence.Statements) != 1 {
		t.Errorf("<consequence> が 1 じゃないのはおかしい。got=%d", len(ifExpr.Consequence.Statements))
	}

	//  `if (x < y) { x }` における BlockStmtは 1つであると決めているからね。
	// consequenceExprStmt は 1つ目の<式文>ってこと。BlockStatement型じゃないから注意ね。名前がややこしい。
	consequenceExprStmt, ok := ifExpr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ifExpr.Consequence.Statements[0] not *aset.ExpressionStatement.got=%T", ifExpr.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequenceExprStmt.Expression, "x") {
		return
	}

	if ifExpr.Alternative != nil {
		t.Errorf("elseないんだから ifExpr.Alternative は nil になってもらわんと困る.got=%+v", ifExpr.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("1 Statement じゃないぞ！ got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("*ast.ExpressionStatementに変換できないぞ。got=%T", program.Statements[0])
	}

	ifExpr, ok := exprStmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("*ast.IfExpressonに変換できないぞ。got=%T", exprStmt.Expression)
	}

	// <condition> が x < y という 中置演算式であることを確認する
	if !testInfixExpression(t, ifExpr.Condition, "x", "<", "y") {
		return
	}

	// <consequence>が 1文 で構成されていることを確認する
	if len(ifExpr.Consequence.Statements) != 1 {
		t.Fatalf("ifExpr.Consequence.Statements[0] not *aset.ExpressionStatement.got=%T", ifExpr.Consequence.Statements[0])
	}

	// <consequence>の 1文目 が x という ast.Identifier であることを確認する
	consequenceExprStmt, ok := ifExpr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("")
	}
	if !testIdentifier(t, consequenceExprStmt.Expression, "x") {
		return
	}

	// <alternative> が 1文 で構成されていることを確認する
	if len(ifExpr.Alternative.Statements) != 1 {
		t.Fatalf("")
	}

	// <alternative>の 1文目 が y という ast.Identifier であることを確認する
	alternativeExprStmt, ok := ifExpr.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", ifExpr.Alternative.Statements[0])
	}
	if !testIdentifier(t, alternativeExprStmt.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("1 Statement じゃないのはおかしいね？ got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("*ast.ExpressionStatement じゃないのはおかしいね？ got=%T", program.Statements[0])
	}

	fnLit, ok := exprStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("*ast.FunctionLiteral じゃないのはおかしいね？ got=%T", exprStmt.Expression)
	}

	// フィールド検証
	// Token検証;
	if fnLit.TokenLiteral() != "fn" {
		t.Errorf("")
	}

	// 引数検証; fn(x, y) { x + y; }
	if len(fnLit.Parameters) != 2 {
		t.Fatalf("Parameteters は 2個 じゃないとおかしいよ。got=%d", len(fnLit.Parameters))
	}
	if !testIdentifier(t, fnLit.Parameters[0], "x") {
		t.Errorf("第1引数は x じゃないとおかしいね。got=%v", fnLit.Parameters[0])
	}
	if !testIdentifier(t, fnLit.Parameters[1], "y") {
		t.Errorf("第1引数は y じゃないとおかしいね。got=%v", fnLit.Parameters[1])
	}

	// Body検証; fn(x, y) { x + y; }
	if len(fnLit.Body.Statements) != 1 {
		t.Fatalf("関数のBodyは 1 Statementじゃないとおかしいよ.got=%d", len(fnLit.Body.Statements))
	}

	bodyStmt, ok := fnLit.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("")
	}
	if !testInfixExpression(t, bodyStmt.Expression, "x", "+", "y") {
		t.Errorf("")
	}
}

func TestFunctionParameters(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y) {};", []string{"x", "y"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		exprStmt := program.Statements[0].(*ast.ExpressionStatement)
		fnLit := exprStmt.Expression.(*ast.FunctionLiteral)

		if len(fnLit.Parameters) != len(tt.expectedParams) {
			t.Errorf("FunctionLitearl.Parmetersの数がおかしいよ。got=%d want=%d", len(tt.expectedParams), len(fnLit.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testIdentifier(t, fnLit.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("1 Statement じゃないよ。got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("*ast.ExpressionStatementじゃないよ。got=%T", program.Statements[0])
	}

	callExpr, ok := exprStmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("*ast.CallExpressionじゃないよ. got=%T", exprStmt.Expression)
	}

	// フィールド検証
	// add(1, 2 * 3, 4 + 5);
	if !testIdentifier(t, callExpr.Function, "add") {
		return
	}

	// add(1, 2 * 3, 4 + 5);
	if len(callExpr.Arguments) != 3 {
		t.Errorf("3 Arguments じゃないよ。got=%d", len(callExpr.Arguments))
	}

	testLiteralExpression(t, callExpr.Arguments[0], 1)
	testInfixExpression(t, callExpr.Arguments[1], 2, "*", 3)
	testInfixExpression(t, callExpr.Arguments[2], 4, "+", 5)

}

func TestLetStatements(t *testing.T) {
	// MEMO: 面倒だから新しくテストケースをまるまる追加(2.8.6)
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("1 Statement じゃないのはおかしいよ。got=%d", len(program.Statements))
			}

			stmt := program.Statements[0]
			if !testLetStatement(t, stmt, tt.expectedIdentifier) {
				return
			}

			value := stmt.(*ast.LetStatement).Value
			if !testLiteralExpression(t, value, tt.expectedValue) {
				return
			}

		})
	}

}
