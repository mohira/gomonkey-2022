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
		input         string
		operator      string
		expectedValue any
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},

		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements が %d 文じゃないよ！got=%d\n", 1, len(program.Statements))
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

		if !testLiteralExpression(t, prefixExpr.Right, tt.expectedValue) {
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

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"3 + 4;", 3, "+", 4},
		{"3 - 4;", 3, "-", 4},
		{"3 * 4;", 3, "*", 4},
		{"3 / 4;", 3, "/", 4},
		{"3 > 4;", 3, ">", 4},
		{"3 < 4;", 3, "<", 4},
		{"3 == 4;", 3, "==", 4},
		{"3 != 4;", 3, "!=", 4},

		// BooleanLiteral系
		{"true == true", true, "==", true},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("1文じゃないよ！ got=%d\n", len(program.Statements))
		}

		exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatementでっせ！ got=%T", program.Statements[0])
		}

		infixExpr, ok := exprStmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("infixExpr が *ast.InfixExpression じゃないぞ！ got=%T", exprStmt.Expression)
		}

		if !testInfixExpression(t, infixExpr, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}

	}

}

func Testえび実験_3項とかになっても大丈夫かな(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"3 + 4 * 5;", "(3 + (4 * 5))"},
		{"3 * 4 + 5;", "((3 * 4) + 5)"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if program.String() != tt.want {
			t.Errorf("got=%s, want=%s", program.String(), tt.want)
		}
	}

}

func TestOperatorPrecedenceParsing(t *testing.T) {
	// 異なる優先順位を持っているもっと複雑なパターンの検証
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5", "(3 + (4 * 5))"},
		{"3 * 1 + 4 * 5", "((3 * 1) + (4 * 5))"},

		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"3 + 4 * 5 == 6 * 7", "((3 + (4 * 5)) == (6 * 7))"},

		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},

		// 2.8.2 グループ化された式
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(3 + 4) * 2", "((3 + 4) * 2)"},
		{"2 / (3 + 4)", "(2 / (3 + 4))"},
		{"-(3 + 4)", "(-(3 + 4))"},
		{"!(true == true)", "(!(true == true))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		got := program.String()
		if got != tt.expected {
			t.Errorf("got=%s, want=%s", got, tt.expected)
		}
	}

}

func TestPratt構文解析の仕組みの実験(t *testing.T) {
	// p.76あたりからの説明
	input := `1 + 2 + 3`
	want := `((1 + 2) + 3)`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	got := program.String()
	if got != want {
		t.Errorf("おかしいよ.got=%s want=%s", got, want)
	}

}

func testIdentifier(t *testing.T, expr ast.Expression, expectedValue string) bool {
	t.Helper()

	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Errorf("expr not *ast.Identifier. got=%T", expr)
		return false
	}

	if ident.Value != expectedValue {
		t.Errorf("ident.Value not %s. got=%s", expectedValue, ident.Value)
		return false
	}

	if ident.TokenLiteral() != expectedValue {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", expectedValue, ident.TokenLiteral())
		return false
	}

	return true
}

// Expressionを検証する より一般的なヘルパーテスト関数 リテラル編
func testLiteralExpression(t *testing.T, expr ast.Expression, expected any) bool {
	t.Helper()
	// 期待値によって検証項目を切り替えるというシンプルな発想
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case string:
		return testIdentifier(t, expr, v)
	case bool:
		return testBooleanLiteral(t, expr, v)
	}

	t.Errorf("😢 そのtypeのリテラルのテストヘルパーはまだないんだわこりゃ. got=%T", expr)
	return false
}

func testInfixExpression(t *testing.T, expr ast.Expression, left any, operator string, right any) bool {
	infixExpr, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expr is not *ast.InfixExpression. got=%T(%s)", expr, expr)
		return false
	}

	if !testLiteralExpression(t, infixExpr.Left, left) {
		return false
	}

	if infixExpr.Operator != operator {
		t.Errorf("Operater が '%s' じゃないぞ！ got='%s'", operator, infixExpr.Operator)
		return false
	}

	if !testLiteralExpression(t, infixExpr.Right, right) {
		return false
	}

	return true
}

func Testテストヘルパーを使っていい感じにテストコードが書けることを試すやつ(t *testing.T) {
	input := `5 + 10;`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("だめ〜")
	}

	if !testInfixExpression(t, exprStmt.Expression, 5, "+", 10) {
		return
	}

}

func TestBooleanLiteral(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("文の数がおかしいね？ want 1, got=%d", len(program.Statements))
		}

		exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("*ast.ExpressionStatement じゃないけど？ got=%T", program.Statements[0])
		}

		booleanLiteral, ok := exprStmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("*ast.Boolean じゃないよ？ got=%T", exprStmt.Expression)
		}

		if booleanLiteral.Value != tt.expectedBoolean {
			t.Errorf("got=%t want=%t", booleanLiteral.Value, tt.expectedBoolean)
		}

	}
}

func testBooleanLiteral(t *testing.T, expr ast.Expression, value bool) bool {
	t.Helper()

	booleanLiteral, ok := expr.(*ast.Boolean)
	if !ok {
		t.Errorf("*ast.Boolean じゃないよ。got=%T", expr)
		return false
	}

	if booleanLiteral.Value != value {
		t.Errorf("booleanLiteral.Value got=%t, want=%t", booleanLiteral.Value, value)
		return false
	}

	if booleanLiteral.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("booleanLiteral.TokenLiteral() not '%t' got '%s'", value, booleanLiteral.TokenLiteral())
	}

	return true
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statemnts not contain %d statements. got=%d", 1, len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ifExpr, ok := exprStmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("not *ast.IfExpression. got=%T", exprStmt.Expression)
	}

	// if (x < y) { x }
	if !testInfixExpression(t, ifExpr.Condition, "x", "<", "y") {
		return
	}

	// consequence
	if len(ifExpr.Consequence.Statements) != 1 {
		t.Errorf("consequnce is not 1 statement. got=%d", len(ifExpr.Consequence.Statements))
	}

	consequence, ok := ifExpr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not *ast.ExpressionStatement. got=%T", ifExpr.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	// elseブロックは無いよ
	if ifExpr.Alternative != nil {
		t.Errorf("ifExpr.Alternative が nil じゃないよ. got=%+v", ifExpr.Alternative)
	}

}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statemnts not contain %d statements. got=%d", 1, len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ifExpr, ok := exprStmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("not *ast.IfExpression. got=%T", exprStmt.Expression)
	}

	// if (x < y) { x } else { y }
	if !testInfixExpression(t, ifExpr.Condition, "x", "<", "y") {
		return
	}

	// consequence
	if len(ifExpr.Consequence.Statements) != 1 {
		t.Errorf("consequnce is not 1 statement. got=%d", len(ifExpr.Consequence.Statements))
	}

	consequence, ok := ifExpr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not *ast.ExpressionStatement. got=%T", ifExpr.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	// if (x < y) { x } else { y }
	if len(ifExpr.Alternative.Statements) != 1 {
		t.Errorf("ifExpr.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(ifExpr.Alternative.Statements))
	}

	alternative, ok := ifExpr.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			ifExpr.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}
