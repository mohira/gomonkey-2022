package parser_test

import (
	"fmt"
	"gomonkey/ast"
	"gomonkey/lexer"
	"gomonkey/parser"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
		// セミコロンは省略できる
		{"let foobar = y", "foobar", "y"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("おかしいぞ")
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
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
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},

		// セミコロンは省略できる
		{"return foobar", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		returnStmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("ast.ReturnStatementじゃないよ. got=%T", program.Statements[0])
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral is not 'return'. got=%s", returnStmt.TokenLiteral())
		}

		testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue)

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

		// 2.8.5 CallExpression
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},

		// 4.4.3 添字の計算は最強！ プライマリだ！
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},

		// 添字演算子の優先順位が、呼び出し式のみならず中置式の「 * 」演算子よりも高いことを 期待している。
		//    2 * b[0]
		// o: (2 * (b[0]))
		// x: ((2 * b)[0]) ← ちがうよ！
		{"2 * b[0]", "(2 * (b[0]))"},
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

func TestFunctionLiteral(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("1文じゃないよ. got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	// fn(x, y) { return x + y; }
	functionLit, ok := exprStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("exprStmt.Expression not *ast.FunctionLiteral. got=%T", exprStmt.Expression)
	}

	// 関数リテラルのパラメータ関連
	if len(functionLit.Parameters) != 2 {
		t.Fatalf("function literal Parameters wrong. want 2, got=%d", len(functionLit.Parameters))
	}
	testLiteralExpression(t, functionLit.Parameters[0], "x")
	testLiteralExpression(t, functionLit.Parameters[1], "y")

	// 関数リテラルのBody関連
	if len(functionLit.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d", len(functionLit.Body.Statements))
	}

	bodyStmt, ok := functionLit.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not *ast.ExpressionStatement. got=%T", functionLit.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")

}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		exprStmt := program.Statements[0].(*ast.ExpressionStatement)
		functionLit := exprStmt.Expression.(*ast.FunctionLiteral)

		if len(functionLit.Parameters) != len(tt.expectedParams) {
			fmt.Printf("👺 %[1]T %[1]v\n", functionLit)
			t.Errorf("parameterの数がおかしいよ。 want=%d got=%d", len(tt.expectedParams), len(functionLit.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, functionLit.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("1文じゃないよ")
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("だめだよ")
	}

	callExpr, ok := exprStmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("ast.CallExpressionじゃないよ！ got=%T", exprStmt.Expression)
	}

	// add(1, 2 * 3, 4 + 5);
	if !testIdentifier(t, callExpr.Function, "add") {
		return
	}

	// Argsのチェック
	if len(callExpr.Arguments) != 3 {
		t.Fatalf("Argumentsの数が違うよ. want=3, got=%d", len(callExpr.Arguments))
	}

	testLiteralExpression(t, callExpr.Arguments[0], 1)
	testInfixExpression(t, callExpr.Arguments[1], 2, "*", 3)
	testInfixExpression(t, callExpr.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	exprStmt := program.Statements[0].(*ast.ExpressionStatement)
	strLit, ok := exprStmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("*ast.StringLiteralじゃないよ。got=%T", exprStmt.Expression)
	}

	if strLit.Value != "hello world" {
		t.Errorf("literal.Value not %s, got=%s", "hello world", strLit.Value)
	}

}

func TestParsingArrayLiterals(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("*ast.ExpressionStatementじゃないよ。got=%T", program.Statements[0])
	}

	arrayLiteral, ok := exprStmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("*ast.ArrayLiteralじゃないよ。got=%T", exprStmt.Expression)
	}

	if len(arrayLiteral.Elements) != 3 {
		t.Fatalf("要素数が3じゃないよ！ got=%d", len(arrayLiteral.Elements))
	}
	// [1, 2 * 2, 3 + 3]
	testIntegerLiteral(t, arrayLiteral.Elements[0], 1)
	testInfixExpression(t, arrayLiteral.Elements[1], 2, "*", 2)
	testInfixExpression(t, arrayLiteral.Elements[2], 3, "+", 3)
}

func TestParsingArrayLiterals_空の配列の場合(t *testing.T) {
	input := `[]`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("*ast.ExpressionStatementじゃないよ。got=%T", program.Statements[0])
	}

	arrayLiteral, ok := exprStmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("*ast.ArrayLiteralじゃないよ。got=%T", exprStmt.Expression)
	}

	if len(arrayLiteral.Elements) != 0 {
		t.Fatalf("要素数が0じゃないよ！ got=%d", len(arrayLiteral.Elements))
	}
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 2]"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("だめだよ")
	}

	indexExpr, ok := exprStmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("*ast.IndexExpressionになってないよ！ got=%T", exprStmt.Expression)
	}

	if !testIdentifier(t, indexExpr.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExpr.Index, 1, "+", 2) {
		return
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("*ast.ExpressionStatementじゃないよ。got=%T", program.Statements[0])
	}

	hashLiteral, ok := exprStmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("*ast.HashLiteral じゃないよ。got=%T", exprStmt.Expression)
	}

	if len(hashLiteral.Pairs) != 3 {
		t.Errorf("要素数が3じゃないよ.got=%d", len(hashLiteral.Pairs))
	}

	want := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hashLiteral.Pairs {
		str, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("ハッシュリテラルのKeyが ast.StringLiteral じゃないよ。got=%T", key)
		}

		wantValue := want[str.String()]

		testIntegerLiteral(t, value, wantValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("*ast.ExpressionStatementじゃないよ。got=%T", program.Statements[0])
	}

	hashLiteral, ok := exprStmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("*ast.HashLiteral じゃないよ。got=%T", exprStmt.Expression)
	}

	if len(hashLiteral.Pairs) != 0 {
		t.Errorf("要素数が0じゃないよ.got=%d", len(hashLiteral.Pairs))
	}
}

func TestParsingHashLiteralsWithExpression(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5, }`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("*ast.ExpressionStatementじゃないよ。got=%T", program.Statements[0])
	}

	hashLiteral, ok := exprStmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("*ast.HashLiteral じゃないよ。got=%T", exprStmt.Expression)
	}

	if len(hashLiteral.Pairs) != 3 {
		t.Errorf("要素数が3じゃないよ.got=%d", len(hashLiteral.Pairs))
	}

	// このテスト方式にするなら、infixExpr以外のケースも入れたほうが良くない？
	tests := map[string]func(ast.Expression){
		"one": func(expr ast.Expression) {
			testInfixExpression(t, expr, 0, "+", 1)
		},
		"two": func(expr ast.Expression) {
			testInfixExpression(t, expr, 10, "-", 8)
		},
		"three": func(expr ast.Expression) {
			testInfixExpression(t, expr, 15, "/", 5)
		},
	}

	for key, value := range hashLiteral.Pairs {
		str, ok := key.(*ast.StringLiteral)

		if !ok {
			t.Errorf("ハッシュリテラルのKeyが ast.StringLiteral じゃないよ。got=%T", key)
			continue
		}

		testFunc, ok := tests[str.String()]
		if !ok {
			t.Errorf("そのKey %q に対するテスト関数が見つからないよ！", str.String())
			continue
		}

		testFunc(value)
	}
}

func TestMacroLiteralParsing(t *testing.T) {
	input := `macro(x, y) { x + y; }`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("おかしいぞ")
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("おかしいぞ")
	}

	macroLiteral, ok := exprStmt.Expression.(*ast.MacroLiteral)
	if !ok {
		t.Fatalf("おかしいぞ")
	}

	// macro(x, y) { x + y; }
	if len(macroLiteral.Parameters) != 2 {
		t.Fatalf("おかしいぞ")
	}

	testLiteralExpression(t, macroLiteral.Parameters[0], "x")
	testLiteralExpression(t, macroLiteral.Parameters[1], "y")

	if len(macroLiteral.Body.Statements) != 1 {
		t.Fatalf("macroLiteral.Body.Statements が 1 じゃないよ。 got =%d", len(macroLiteral.Body.Statements))
	}

	bodyStmt, ok := macroLiteral.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("おかしいぞ")
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}
