//ff:type feature=tangl type=codegen
//ff:what nodeInfo — a registered case node's Go variable identifier plus its AST definition
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// nodeInfo records a registered case node's Go variable identifier
// alongside its AST definition, keyed by node name in the enclosing
// case's node map.
type nodeInfo struct {
	Var  string
	Node ast.Node
}
