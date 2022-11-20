package evaluator

import (
	"gomonkey/ast"
	"gomonkey/object"
)

func quote(node ast.Node) object.Object {

	return &object.Quote{Node: node}
}
