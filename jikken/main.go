package main

import (
	"fmt"
	"gomonkey/ast"
)

type MyKey struct {
	Value string
}

func main() {
	// Goのmapで同じ構造体を2つkeyにしたらだめになるの？
	//num := &ast.StringLiteral{Value: "num"}
	one := func() ast.Expression { return &ast.IntegerLiteral{Value: 1} }
	ONE := one()
	hashLiteral := &ast.HashLiteral{
		Pairs: map[ast.Expression]ast.Expression{
			ONE: ONE,
			ONE: ONE,
			ONE: ONE,
			//one(): one(),
			//one(): one(),
			//num:   one(),
		},
	}
	_ = hashLiteral

	//m := map[string]int{
	//	"a": 1,
	//	"a": 1,
	//}
	// map[main.MyKey]int map[{Value:a}:1]
	//m := map[MyKey]int{
	//	MyKey{Value: "a"}: 1,
	//	MyKey{Value: "a"}: 1,
	//}

	// map[*main.MyKey]int map[0x14000104210:1 0x14000104220:1]
	m := map[*MyKey]int{
		&MyKey{Value: "a"}: 1,
		&MyKey{Value: "a"}: 1,
	}

	fmt.Printf("%[1]T %+[1]v\n", m)
	//fmt.Printf("%[1]T %+[1]v\n", hashLiteral)
}
