//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what expandSupersedeChain — desugar a priority chain into every earlier-to-later supersede edge
package parser

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// expandSupersedeChain returns the transitive closure of the chain as pair
// edges in canonical (i, j>i) order. The first edge carries a copy of the full
// chain so renderers can restore the chain surface form.
func expandSupersedeChain(names []string, line int) []ast.Supersede {
	edges := make([]ast.Supersede, 0, len(names)*(len(names)-1)/2)
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			edges = append(edges, ast.Supersede{Superseder: names[i], Superseded: names[j], Line: line})
		}
	}
	if len(edges) > 0 {
		edges[0].Chain = append([]string(nil), names...)
	}
	return edges
}
