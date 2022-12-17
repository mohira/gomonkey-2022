package evaluator

import (
	"gomonkey/ast"
	"gomonkey/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	var macroDefinitionIndexes []int

	for idx, stmt := range program.Statements {
		if isMacroDefinition(stmt) {
			// 1) 環境にマクロ定義を登録する
			addMacro(stmt, env)

			// ループ中にASTから消し去るのはご法度なので位置だけ記憶
			macroDefinitionIndexes = append(macroDefinitionIndexes, idx)
		}
	}

	// 2) indexを使ってマクロ定義をASTから消し去る
	// memo: スライスを使った中抜きアルゴリズム
	for i := len(macroDefinitionIndexes) - 1; i >= 0; i = i - 1 {
		definitionIndex := macroDefinitionIndexes[i]
		program.Statements = append(program.Statements[:definitionIndex], program.Statements[definitionIndex+1:]...)
	}

}

func isMacroDefinition(stmt ast.Statement) bool {
	// マクロ定義 := Let文 かつ 右辺がマクロリテラル であること
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		return false
	}

	if _, ok := letStmt.Value.(*ast.MacroLiteral); !ok {
		return false
	}

	return true
}

func addMacro(stmt ast.Statement, env *object.Environment) {
	// エラー処理が冗長なので無視しています。ごめんな。
	letStmt, _ := stmt.(*ast.LetStatement)
	macroLit, _ := letStmt.Value.(*ast.MacroLiteral)

	macroObj := &object.Macro{
		Parameters: macroLit.Parameters,
		Body:       macroLit.Body,
		Env:        env, // よくわからんし今のところ解説もないが、もらったenvをそのまま突っ込んでいます。
	}

	env.Set(letStmt.Name.Value, macroObj)
}

func ExpandMacros(program *ast.Program, env *object.Environment) ast.Node {
	// program.Statementsから「マクロ呼び出し」を見つける
	// マクロ呼び出し := CallExpr かつ  CallExpr.Function の識別子の.Value が envに登録されているやつ

	// TODO: forは後回しにする。1Statementのテストケースだけだから！
	stmt := program.Statements[0]
	exprStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		panic("ぱにっくだ！！！！")
	}

	callExpr, ok := exprStmt.Expression.(*ast.CallExpression)
	if !ok {
		panic("ぱにっくだ！！！！")
	}

	// envから名前探し
	obj, ok := env.Get(callExpr.Function.TokenLiteral())
	if !ok {
		panic("ぱにっくだ！！！！")
	}

	// *object.Macro &{[] quote((1 + 2)) 0x14000052a10}
	//fmt.Printf("%[1]T %[1]v\n", obj)

	macroObj, ok := obj.(*object.Macro)
	if !ok {
		// チェックしなくても良いかも？
		panic("ぱにっくだ！！！！")
	}

	obj = Eval(macroObj.Body, env)

	// マクロは object.Quote を返さないとダメだよルールなので、チェックします
	quoteObj, ok := obj.(*object.Quote)
	if !ok {
		panic("ぱにっくだ！！！！")
	}

	// *object.Quote &{(1 + 2)}
	// fmt.Printf("%[1]T %[1]v\n", obj)

	return quoteObj.Node
}
