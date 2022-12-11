package evaluator

import (
	"gomonkey/ast"
	"gomonkey/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	var macro定義のインデックス集合 []int

	for idx, stmt := range program.Statements {
		letStmt, ok := stmt.(*ast.LetStatement)
		if !ok {
			continue
		}

		macroLit, ok := letStmt.Value.(*ast.MacroLiteral)
		if !ok {
			continue
		}

		// マクロリテラルだったので
		//	1) 環境に保存する
		// TODO: シグネチャ怪しいけどテスト通してから直すので安心してください
		registerMarco(env, letStmt.Name.Value, macroLit)

		//	ループ中にASTから消し去るのはご法度なので位置だけ記憶
		macro定義のインデックス集合 = append(macro定義のインデックス集合, idx)
	}

	// 2) indexを使ってマクロ定義をASTから消し去る
	for i := len(macro定義のインデックス集合) - 1; i >= 0; i = i - 1 {
		targetIndex := macro定義のインデックス集合[i]
		program.Statements = append(program.Statements[:targetIndex], program.Statements[targetIndex+1:]...)
	}

}

func registerMarco(env *object.Environment, name string, macroLit *ast.MacroLiteral) {
	macroObj := &object.Macro{
		Parameters: macroLit.Parameters,
		Body:       macroLit.Body,
		Env:        nil, // ??? あとで
	}

	env.Set(name, macroObj)
}
