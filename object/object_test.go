package object_test

import (
	"gomonkey/object"
	"testing"
)

func TestStringHashKey(t *testing.T) {
	hello1 := &object.String{Value: "Hello World"}
	hello2 := &object.String{Value: "Hello World"}

	diff1 := &object.String{Value: "My name is johnny"}
	diff2 := &object.String{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("同じ値なのに、ハッシュ値が異なるぞ！？")
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("同じ値なのに、ハッシュ値が異なるぞ！？")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("違う値なのに、ハッシュ値が一緒になっているのはおかしいぞ！")
	}
}
func TestIntegerHashKey(t *testing.T) {
	five1 := &object.Integer{Value: 5}
	five2 := &object.Integer{Value: 5}

	ten1 := &object.Integer{Value: 10}
	ten2 := &object.Integer{Value: 10}

	if five1.HashKey() != five2.HashKey() {
		t.Errorf("同じ値なのに、ハッシュ値が異なるぞ！？")
	}
	if ten1.HashKey() != ten2.HashKey() {
		t.Errorf("同じ値なのに、ハッシュ値が異なるぞ！？")
	}

	if five1.HashKey() == ten1.HashKey() {
		t.Errorf("違う値なのに、ハッシュ値が一緒になっているのはおかしいぞ！")
	}
}
func TestBooleanHashKey(t *testing.T) {
	true1 := &object.Boolean{Value: true}
	true2 := &object.Boolean{Value: true}

	false1 := &object.Boolean{Value: false}
	false2 := &object.Boolean{Value: false}

	if true1.HashKey() != true2.HashKey() {
		t.Errorf("同じ値なのに、ハッシュ値が異なるぞ！？")
	}
	if false1.HashKey() != false2.HashKey() {
		t.Errorf("同じ値なのに、ハッシュ値が異なるぞ！？")
	}

	if true1.HashKey() == false1.HashKey() {
		t.Errorf("違う値なのに、ハッシュ値が一緒になっているのはおかしいぞ！")
	}
}
