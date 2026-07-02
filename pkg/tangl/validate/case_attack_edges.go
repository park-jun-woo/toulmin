//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what caseAttackEdges — the node->node defeat graph induced by one case's don't edges
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// caseAttackEdges returns, for every attacker node in c, the list of node
// names it attacks via `don't <target> when <attacker>` edges.
func caseAttackEdges(c ast.Case) map[string][]string {
	edges := make(map[string][]string)
	for _, a := range c.Attacks {
		edges[a.Attacker] = append(edges[a.Attacker], a.Target)
	}
	return edges
}
