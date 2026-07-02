//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what caseNodeLineIndex — maps each node name in one case to its declaration line
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// caseNodeLineIndex maps every node name registered in c to the line where
// it was declared.
func caseNodeLineIndex(c ast.Case) map[string]int {
	idx := make(map[string]int, len(c.Nodes))
	for _, n := range c.Nodes {
		idx[n.Name] = n.Line
	}
	return idx
}
