package object

import (
	"fmt"
	"gomonkey/ast"
	"hash/fnv"
	"strings"
)

type Type string

const (
	IntegerObj     = "INTEGER"
	StringObj      = "STRING"
	BooleanObj     = "BOOLEAN"
	NullObj        = "NULL"
	ReturnValueObj = "RETURN_VALUE"
	ErrorObj       = "ERROR"

	FunctionObj = "FUNCTION"
	BuiltinObj  = "BUILTIN"

	// 発見: ユーザー定義"関数" と 組み込み"関数" は monkeyのオブジェクトの世界だと全く別物！
	/*
		Python でも 文字列表現としては、組み込み関数とユーザー定義関数は確かに違う扱いだった！
		っていうか、CpythonならprintはCで実装されているので、それはそうというのはあとになって気づきました
		>>> print
		<built-in function print>
		>>> def add():pass
		...
		>>> add
		<function add at 0x1013d3250>
	*/
	ArrayObj = "ARRAY"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type {
	return IntegerObj
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type {
	return BooleanObj
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct{} // 何の値もラップしていないことが「値の不存在」を表現している

func (n *Null) Type() Type {
	return NullObj
}

func (n *Null) Inspect() string {
	return "NULL"
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() Type {
	return ReturnValueObj
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() Type {
	return ErrorObj
}

func (e *Error) Inspect() string {
	return "💥 ERROR:" + e.Message
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() Type {
	return FunctionObj
}

func (f *Function) Inspect() string {
	var out strings.Builder

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() Type {
	return StringObj
}

func (s *String) Inspect() string {
	return s.Value
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() Type {
	return BuiltinObj
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() Type {
	return ArrayObj
}

func (a *Array) Inspect() string {
	// [1, add(2, 3), 4 + 5]
	var out strings.Builder

	var elements []string
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// HashKey ハッシュ値を表現する構造体。別にObjectインタフェースは満足してないよ
type HashKey struct {
	Type  Type
	Value uint64
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func (i *Integer) HashKey() HashKey {
	// ポインタじゃないからこれでおk
	// ハッシュ値の演算は不要
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (b *Boolean) HashKey() HashKey {
	var v uint64
	if b.Value {
		v = 1
	}
	return HashKey{Type: b.Type(), Value: v}
}
