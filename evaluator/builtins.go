package evaluator

import (
	"gomonkey/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			// 配列とかハッシュマップに対してのlenはまたあとでな！
			// わざわざビルトイン関数の実装の中で引数の数チェックをするのは分かる。
			// なぜなら、len関数は引数が1個です」というドキュメント的な作用があるから
			if len(args) != 1 {
				return newError("argument error: wrong number of arguments (given %d, expected %d)", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			// 引数は1個でないとだめ
			if len(args) != 1 {
				return newError("argument error: wrong number of arguments (given %d, expected %d)", len(args), 1)
			}

			// 引数のデータ型が配列じゃないとだめ
			switch arg := args[0].(type) {
			case *object.Array:
				// MEMO: 要素数が0なら、どうするか？ → myArray[0]と同じ挙動にするなら NULL を返す
				if len(arg.Elements) == 0 {
					return NULL
				}
				return arg.Elements[0]
			default:
				return newError("argument to `first` not supported, got %s", arg.Type())
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			// 引数は1個でないとだめ
			if len(args) != 1 {
				return newError("argument error: wrong number of arguments (given %d, expected %d)", len(args), 1)
			}

			// 引数のデータ型が配列じゃないとだめ
			switch arg := args[0].(type) {
			case *object.Array:
				// MEMO: 要素数が0なら、どうするか？ NULL にしました
				if len(arg.Elements) == 0 {
					return NULL
				}
				return arg.Elements[len(arg.Elements)-1]
			default:
				return newError("argument to `last` not supported, got %s", arg.Type())
			}
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			// 引数は1個でないとだめ
			if len(args) != 1 {
				return newError("argument error: wrong number of arguments (given %d, expected %d)", len(args), 1)
			}

			// 引数のデータ型が配列じゃないとだめ
			switch arg := args[0].(type) {
			case *object.Array:
				// MEMO: 要素数が0の場合の rest は [] という設計もある。
				length := len(arg.Elements)
				if length == 0 {
					return NULL
				}

				rest := make([]object.Object, length-1)
				copy(rest, arg.Elements[1:length])

				return &object.Array{Elements: rest}
			default:
				return newError("argument to `rest` not supported, got %s", arg.Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			// 引数は2個でないとダメ
			if len(args) != 2 {
				return newError("argument error: wrong number of arguments (given %d, expected %d)", len(args), 2)
			}

			arg1 := args[0]
			arg2 := args[1]

			// 第1引数のデータ型が配列じゃないとだめ
			if arg1.Type() != object.ArrayObj {
				return newError("first argument to `push` not supported, got %s", arg1.Type())
			}

			arrayOrg := arg1.(*object.Array)
			newElements := append(arrayOrg.Elements, arg2)
			// make()で確保したほうが伸長しなくなるっぽいのでメモリ的に有利っぽい。
			// そのときにappendじゃなくて、arr[length] = arg2 みたいな代入になって
			// append感がちょっと減るじゃん？。
			// それが嫌だなーっておもったので、append関数にしたんだと思います(後付)

			return &object.Array{Elements: newElements}
		},
	},
}
