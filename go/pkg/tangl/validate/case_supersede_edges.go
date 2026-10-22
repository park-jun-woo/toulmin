//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what caseSupersedeEdges — the node->node supersede graph induced by one case's supersede edges
package validate

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// caseSupersedeEdges returns, for every superseder node in c, the list of
// node names it supersedes via "`x` supersedes `y`" edges.
func caseSupersedeEdges(c ast.Case) map[string][]string {
	edges := make(map[string][]string)
	for _, s := range c.Supersedes {
		edges[s.Superseder] = append(edges[s.Superseder], s.Superseded)
	}
	return edges
}
