//ff:func feature=tangl type=codegen control=iteration dimension=2
//ff:what registerCaseNodes — writes one rule registration statement per case node
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// registerCaseNodes writes one g.Rule/Counter/Except(...) registration
// statement per case node, in document order, chaining
// With/Qualifier/Grounds/Level clauses, and returns the node name -> Go
// variable map later steps use
// to attach defeat and execution edges.
func registerCaseNodes(w *strings.Builder, gc *genContext, c ast.Case) (map[string]nodeInfo, error) {
	nodes := make(map[string]nodeInfo, len(c.Nodes))
	for _, n := range c.Nodes {
		v := goIdent(n.Name)
		fn, err := nodeFn(gc, n)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(w, "\t%s := g.%s(%s)", v, roleMethod(n.Role), fn)
		for _, term := range n.With {
			fmt.Fprintf(w, ".With(%s)", withArg(gc, term))
		}
		if n.Qualified != nil {
			fmt.Fprintf(w, ".Qualifier(%s)", formatFloat(*n.Qualified))
		}
		if len(n.Grounds) > 0 {
			fmt.Fprintf(w, ".Grounds(%s)", groundsArg(n.Grounds))
		}
		if n.Level != ast.LevelUnset {
			fmt.Fprintf(w, ".Level(%s)", levelArg(n.Level))
		}
		fmt.Fprintln(w)
		nodes[n.Name] = nodeInfo{Var: v, Node: n}
	}
	return nodes, nil
}
