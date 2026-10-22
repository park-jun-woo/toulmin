//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what buildCheckingEdges — the case->case graph induced by checking clauses
package validate

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// buildCheckingEdges returns, for every case, the list of target case names
// referenced by its nodes' `checking <case>` clauses (document order).
func buildCheckingEdges(doc *ast.Document) map[string][]string {
	edges := make(map[string][]string)
	for _, c := range doc.Cases {
		for _, n := range c.Nodes {
			if n.Checking == nil || n.Checking.Alias != "" {
				continue // alias'd targets are cross-document — out of this single-doc graph
			}
			edges[c.Name] = append(edges[c.Name], n.Checking.Name)
		}
	}
	return edges
}
