//ff:func feature=graph type=util control=iteration dimension=1
//ff:what buildEdgesFromDef — converts EdgeDef slice to defeat edges map
package graphdef

// buildEdgesFromDef converts a slice of EdgeDef into a map of target → attackers.
func buildEdgesFromDef(defeats []EdgeDef) map[string][]string {
	edges := make(map[string][]string)
	for _, d := range defeats {
		edges[d.To] = append(edges[d.To], d.From)
	}
	return edges
}
