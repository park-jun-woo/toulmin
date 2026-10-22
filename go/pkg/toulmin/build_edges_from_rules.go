//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what buildEdgesFromRules — builds defeat edges from RuleMeta.Defeats
package toulmin

// buildEdgesFromRules populates edges map from RuleMeta slice.
func buildEdgesFromRules(edges map[string][]string, rules []RuleMeta) {
	for _, r := range rules {
		for _, target := range r.Defeats {
			edges[target] = append(edges[target], r.Name)
		}
	}
}
