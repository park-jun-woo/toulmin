//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what topoOrder — dependency-first topological order of rule names over sibling edges
package eval

// topoOrder returns names in dependency-first order: for edges (rule name ->
// referenced sibling rule names), every referenced rule appears before the rule
// that references it. It assumes the graph is acyclic (verified beforehand) and
// visits names in the given order for determinism.
func topoOrder(names []string, edges map[string][]string) []string {
	visited := map[string]bool{}
	order := make([]string, 0, len(names))
	var visit func(n string)
	visit = func(n string) {
		if visited[n] {
			return
		}
		visited[n] = true
		for _, dep := range edges[n] {
			visit(dep)
		}
		order = append(order, n)
	}
	for _, n := range names {
		visit(n)
	}
	return order
}
