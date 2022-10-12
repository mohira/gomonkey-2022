package parser_test

import (
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
