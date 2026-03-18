//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what collectAttackers — returns set of rule names that appear as attackers
package toulmin

// collectAttackers returns a set of rule names that appear as attackers in defeat edges.
func collectAttackers(defeats []defeatEdge) map[string]bool {
	m := make(map[string]bool)
	for _, d := range defeats {
		m[d.from] = true
	}
	return m
}
