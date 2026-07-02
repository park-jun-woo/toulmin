//ff:func feature=tangl type=codegen control=iteration dimension=2
//ff:what usedNodeNames — collects every node name referenced by a case's Attacks or Execs
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// usedNodeNames collects every node name a case's Attacks (as target or
// attacker) or Execs (as trigger) actually reference. A node var absent
// from this set has no further statement referencing it and needs a
// blank-identifier discard to avoid a Go "declared and not used" error —
// e.g. a `checking` node with no attached don't/do/run edge.
func usedNodeNames(c ast.Case) map[string]bool {
	used := make(map[string]bool)
	for _, a := range c.Attacks {
		used[a.Target] = true
		used[a.Attacker] = true
	}
	for _, e := range c.Execs {
		used[e.Node] = true
	}
	return used
}
