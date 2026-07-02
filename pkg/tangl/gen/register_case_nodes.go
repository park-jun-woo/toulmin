//ff:func feature=tangl type=codegen control=iteration dimension=2
//ff:what registerCaseNodes — writes one rule registration statement per case node
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// registerCaseNodes writes one g.Rule/Counter/Except(...) registration
// statement per case node, in document order, chaining With/Qualifier
// clauses, and returns the node name -> Go variable map later steps use
// to attach defeat and execution edges.
func registerCaseNodes(w *strings.Builder, gc *genContext, c ast.Case) map[string]nodeInfo {
	nodes := make(map[string]nodeInfo, len(c.Nodes))
	for _, n := range c.Nodes {
		v := goIdent(n.Name)
		fmt.Fprintf(w, "\t%s := g.%s(%s)", v, roleMethod(n.Role), nodeFn(gc, n))
		for _, term := range n.With {
			fmt.Fprintf(w, ".With(%s)", withArg(gc, term))
		}
		if n.Qualified != nil {
			fmt.Fprintf(w, ".Qualifier(%s)", formatFloat(*n.Qualified))
		}
		fmt.Fprintln(w)
		nodes[n.Name] = nodeInfo{Var: v, Node: n}
	}
	return nodes
}
