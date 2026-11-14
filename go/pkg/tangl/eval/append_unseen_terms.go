//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what appendUnseenTerms — append the terms of one group not yet seen, marking them seen
package eval

// appendUnseenTerms appends every term of group that is not in seen to out, in
// order, and marks it seen.
func appendUnseenTerms(out []string, seen map[string]bool, group []string) []string {
	for _, term := range group {
		if !seen[term] {
			seen[term] = true
			out = append(out, term)
		}
	}
	return out
}
