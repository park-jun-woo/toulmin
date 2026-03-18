//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what buildAttackerSet — builds attacker lookup from defeat edges
package toulmin

// buildAttackerSet returns a set of all node IDs that appear as attackers.
func buildAttackerSet(edges map[string][]string) map[string]bool {
	set := make(map[string]bool)
	for _, attackers := range edges {
		for _, aid := range attackers {
			set[aid] = true
		}
	}
	return set
}
