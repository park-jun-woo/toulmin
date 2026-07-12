//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what buildEdges — builds the rule -> referenced sibling-rule dependency graph
package eval

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// buildEdges builds the dependency graph mapping each rule name to the sibling
// document rules it references via a bare predicate Call.
func buildEdges(rules []ast.InlineRule, ruleSet map[string]bool) map[string][]string {
	edges := map[string][]string{}
	for _, r := range rules {
		var deps []string
		collectRuleCalls(r.Cond, ruleSet, &deps)
		edges[r.Name] = deps
	}
	return edges
}
