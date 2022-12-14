package evaluator

import (
	"gomonkey/ast"
	"gomonkey/object"
	"gomonkey/token"
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

func ExpandMacros(program ast.Node, env *object.Environment) ast.Node {
	modifier := func(node ast.Node) ast.Node {
		// CallExprでない時点で「マクロ呼び出し」なわけがない
		callExpr, ok := node.(*ast.CallExpression)
		if !ok {
			return node // 「このノードは変更しないよ」ってこと
		}

		// マクロ呼び出し なのかをチェック
		macroObj, ok := isMacroCall(callExpr, env)
		if !ok {
			return node
		}

		// 実引数 を "評価しないまま" にしたいので、それぞれを *object.Quote で包む
		quotedArgs := quoteArgs(callExpr)

		// マクロ展開時の環境を作る
		// quotedArgs たちが登録された環境をつくります
		macroEnv := extendMarcoEnv(macroObj, quotedArgs)

		// マクロオブジェクトの評価
		evaluated := Eval(macroObj.Body, macroEnv)

		quote, ok := evaluated.(*object.Quote)
		if !ok {
			panic("マクロ呼び出しの結果は Quoteオブジェクト(ASTノードを内包している) じゃないとダメなルールです！ これはルールです。")
		}

		return quote.Node
	}

	return ast.Modify(program, modifier)
}

func extendMarcoEnv(macroObj *object.Macro, quotedArgs []*object.Quote) *object.Environment {
	macroEnv := object.NewEnclosedEnvironment(macroObj.Env)

	for i, param := range macroObj.Parameters {
		macroEnv.Set(param.Value, quotedArgs[i])
	}

	return macroEnv
}

func quoteArgs(callExpr *ast.CallExpression) []*object.Quote {
	// 実引数 を "評価しないまま" にしたいので、それぞれを *object.Quote で包む
	var quotedArgs []*object.Quote

	for _, arg := range callExpr.Arguments {
		quotedArgs = append(quotedArgs, &object.Quote{Node: arg})
	}

	return quotedArgs
}

func isMacroCall(callExpr *ast.CallExpression, env *object.Environment) (*object.Macro, bool) {
	identifier, ok := callExpr.Function.(*ast.Identifier)
	if !ok {
		return nil, false
	}

	// 呼び出し式 の 識別子 が環境に登録されているかチェック
	obj, ok := env.Get(identifier.Value)
	if !ok {
		return nil, false
	}

	macroObj, ok := obj.(*object.Macro)
	if !ok {
		return nil, false
	}

	return macroObj, true
}

func 俺たちでつくったの一応残しておく_ExpandMacros(program *ast.Program, env *object.Environment) ast.Node {
	// program.Statementsから「マクロ呼び出し」を見つけて、そのマクロ呼び出しを 適用して 置き換える。
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

	// ## 最終的には a: Quote(2, +, 2) , b: Quote(10-5) というかたちでenvに登録されればいい
	// 参考としては、CallExpressionのEvalを使うと良い

	// 1. CallExpr.Args(実引数) を Quoteで包む ← うっかり評価しちゃだめだぞ！
	var quotedArgs []*object.Quote
	for _, arg := range callExpr.Arguments {
		quotedArgs = append(quotedArgs, &object.Quote{Node: arg})
	}

	// 2. エラーチェック: 実引数の数 と 仮引数の数 が 不一致だったらダメだぞ！(本編とは関係なので書いてない)

	// 3. "拡張したEnv" に quotedArgs たちを登録する
	macroEnv := object.NewEnclosedEnvironment(macroObj.Env)

	for i, param := range macroObj.Parameters {
		macroEnv.Set(param.Value, quotedArgs[i])
	}

	obj = Eval(macroObj.Body, macroEnv)

	// マクロは object.Quote を返さないとダメだよルールなので、チェックします
	quoteObj, ok := obj.(*object.Quote)
	if !ok {
		panic("ぱにっくだ！！！！")
	}

	// *object.Quote &{(1 + 2)}
	// fmt.Printf("%[1]T %[1]v\n", obj)

	// 「マクロ呼び出しを展開した結果のノード」を 「マクロ呼び出しのノード部分」 にすげ替える

	expr, ok := quoteObj.Node.(ast.Expression)
	if !ok {
		panic("ぱにっくだ！！！！")
	}

	newExprStmt := &ast.ExpressionStatement{
		Token:      token.Token{},
		Expression: expr,
	}

	program.Statements[0] = newExprStmt

	return program
}
