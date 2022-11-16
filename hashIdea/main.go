package main

import (
	"fmt"
	"gomonkey/object"
)

// monkey code
// let myHash = {"name": "Monkey"};
// myHash["name"]

// SimpleHash 素朴なアイデア
type SimpleHash struct {
	Pairs map[object.Object]object.Object
}

func main() {
	name1 := &object.String{Value: "name"}
	monkey := &object.String{Value: "Monkey"}

	myHash := SimpleHash{Pairs: map[object.Object]object.Object{}}
	myHash.Pairs[name1] = monkey

	fmt.Printf("%[1]T %+[1]v\n", myHash.Pairs[name1])

	name2 := &object.String{Value: "name"}

	fmt.Printf("%t\n", name1 == name2) // falseになる(ポインタ比較だから)

	fmt.Printf("%[1]T %[1]v\n", myHash.Pairs[name2]) // nil になっちゃう！ 持っているValueは同じなのに！
}
