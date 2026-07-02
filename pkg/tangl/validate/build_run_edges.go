//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what buildRunEdges — the case->case graph induced by run cascade edges
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// buildRunEdges returns, for every case, the list of target case names
// reached by its `run <case> when <node>` edges (document order).
func buildRunEdges(doc *ast.Document) map[string][]string {
	edges := make(map[string][]string)
	for _, c := range doc.Cases {
		for _, e := range c.Execs {
			if e.Kind != ast.RunExec {
				continue
			}
			edges[c.Name] = append(edges[c.Name], e.Case)
		}
	}
	return edges
}
