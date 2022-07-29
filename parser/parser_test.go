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
		integerValue int64 // ← !5; とか -3; みたいに、<expression>が<integer_literal>限定のテストってこと
	}{
		{"!5;", "!", 5}, // !5 が何を返すかは多分決めてないと思う。あくまで前置演算式って構造だけ。
		{"-15;", "-", 15},
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

		if testIntegerLiteral(t, prefixExpr.Right, tt.integerValue) {
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
		LeftIntegerValue  int64
		Operator          string
		RightIntegerValue int64
	}{
		{"3 + 4;", 3, "+", 4},
		{"3 - 4;", 3, "-", 4},
		{"3 * 4;", 3, "*", 4},
		{"3 / 4;", 3, "/", 4},
		{"3 > 4;", 3, ">", 4},
		{"3 < 4;", 3, "<", 4},
		{"3 == 4;", 3, "==", 4},
		{"3 != 4;", 3, "!=", 4},
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

		infixExpr, ok := exprStmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("おかしいよ. got=%T", exprStmt.Expression)
		}

		// フィールド検証
		if !testIntegerLiteral(t, infixExpr.Left, tt.LeftIntegerValue) {
			return
		}

		if infixExpr.Operator != tt.Operator {
			t.Errorf("おかしいよ")
		}

		if !testIntegerLiteral(t, infixExpr.Right, tt.RightIntegerValue) {
			return
		}
	}
}
