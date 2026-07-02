//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what internalReachableCases — the run-cascade closure of every case an Internal trigger runs directly
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// internalReachableCases returns the set of case names reachable from
// tangl:Internal's direct `run` targets, closed over the run-cascade edges
// (Internal `check` targets are Evaluate-only and are not included).
func internalReachableCases(doc *ast.Document) map[string]bool {
	edges := buildRunEdges(doc)
	reached := make(map[string]bool)
	var queue []string
	for _, in := range doc.Internals {
		for _, name := range in.Runs {
			if reached[name] {
				continue
			}
			reached[name] = true
			queue = append(queue, name)
		}
	}
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		for _, next := range edges[name] {
			if reached[next] {
				continue
			}
			reached[next] = true
			queue = append(queue, next)
		}
	}
	return reached
}
