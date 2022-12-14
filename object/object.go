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
	HashObj     = "HASH"

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

	QuoteObj = "QUOTE"
	MacroObj = "MACRO"
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

type Hashable interface {
	HashKey() HashKey
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

// HashPair がわざわざ必要なのなんでなん？ ← Keyを記録したいから
// 後からREPLでMonkeyのハッシュを表示するとき、ハッシュに格納されている値だけでなく、そのキーも表示したいんだ。
type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair // map[HashKey]Object ではだめな理由は？
}

func (h *Hash) Type() Type {
	return HashObj
}

func (h *Hash) Inspect() string {
	// {"name": "Bob",
	//  "age": 25,
	//  "points": {"a": 100, "b: 200},
	// }
	var out strings.Builder

	var pairs []string
	for _, hashPair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", hashPair.Key.Inspect(), hashPair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type Quote struct {
	Node ast.Node
}

func (q *Quote) Type() Type {
	return QuoteObj
}

func (q *Quote) Inspect() string {
	return "QUOTE(" + q.Node.String() + ")"
}

type Macro struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment // ← マクロの中の「環境」これむずくね？
}

func (f *Macro) Type() Type {
	return MacroObj
}

func (f *Macro) Inspect() string {
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
