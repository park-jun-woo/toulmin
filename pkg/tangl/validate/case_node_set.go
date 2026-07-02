//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what caseNodeSet — the set of node names registered inside one case
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// caseNodeSet returns the set of every node name registered in c.
func caseNodeSet(c ast.Case) map[string]bool {
	set := make(map[string]bool, len(c.Nodes))
	for _, n := range c.Nodes {
		set[n.Name] = true
	}
	return set
}
