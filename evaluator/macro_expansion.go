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
	return nil
}
