//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what uniqueTerms — concatenate term groups keeping each term's first occurrence
package eval

// uniqueTerms concatenates groups in order, keeping only the first occurrence
// of each term.
func uniqueTerms(groups ...[]string) []string {
	seen := map[string]bool{}
	var out []string
	for _, group := range groups {
		out = appendUnseenTerms(out, seen, group)
	}
	return out
}
