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
}
