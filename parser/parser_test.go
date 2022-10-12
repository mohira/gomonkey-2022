package parser_test

import (
	"fmt"
	"gomonkey/ast"
	"gomonkey/lexer"
	"gomonkey/parser"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	// わざとエラー起こすための入力( = がない)
	//	input = `
	//let x 5;
	//let = 10;
	//let 838383;
	//`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements が 3文 じゃないよ. got=%d", len(program.Statements))
	}

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

func checkParseErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	// エラーが起きている時点で処理をとめちゃうべき
	// エラーに気づいている状態で移行のテストしても無駄だからね
	t.FailNow()
}

func testLetStatement(t *testing.T, stmt ast.Statement, expectedName string) bool {
	t.Helper()
	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt が *ast.LetStatement じゃないよ。got=%T", letStmt)
		return false
	}
	// MEMO: LetStatementのValueのテストは後回し(<expression>だから大変なので)
	if letStmt.Name.Value != expectedName {
		t.Errorf("letStmt.Name.Value が '%s' じゃないよ。got=%s", expectedName, letStmt.Name.Value)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements が 3文 になっていないんだよね。got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt が *ast.ReturnStatementじゃないよ! got=%T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got=%q", returnStmt.TokenLiteral())
		}

	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program の文が足りないよ.got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] が *ast.ExpressionStatementじゃないよ！got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expressionが、 *ast.Identiferじゃないよ！ got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value が %s じゃないよ。got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() が %s じゃないよ. got=%s", "foobar", ident.TokenLiteral())
	}

}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program の文が足りないよ.got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] が *ast.ExpressionStatementじゃないよ.got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp が *ast.IntegerLiteral じゃないよ。got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value が %d じゃないよ。got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral() が %s じゃないよ。got=%s", "5", literal.TokenLiteral())
	}

}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements が %d 文じゃないよ！got=%d\b", 1, len(program.Statements))
		}

		exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("だめでした。got=%T", program.Statements[0])
		}

		prefixExpr, ok := exprStmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("型アサーション失敗！ got=%T", exprStmt.Expression)
		}

		if prefixExpr.Operator != tt.operator {
			t.Fatalf("prefixExpr.Operator が '%s' じゃないぞ！ got=%s", tt.operator, prefixExpr.Operator)
		}

		if !testIntegerLiteral(t, prefixExpr.Right, tt.integerValue) {
			return
		}
	}

}

func testIntegerLiteral(t *testing.T, expr ast.Expression, expectedValue int64) bool {
	t.Helper()

	integerLiteral, ok := expr.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("与えられた expr が *ast.IntegerLiteralじゃないよ！！！ got=%T", integerLiteral)
		return false
	}

	if integerLiteral.Value != expectedValue {
		t.Errorf("integerLiteral.Value が %d じゃないぞ！ got=%d", expectedValue, integerLiteral.Value)
		return false
	}

	if integerLiteral.TokenLiteral() != fmt.Sprintf("%d", expectedValue) {
		t.Errorf("integerLiteral.TokenLiteral が %d じゃないよ. got=%s", expectedValue, integerLiteral.TokenLiteral())
		return false
	}

	return true
}
