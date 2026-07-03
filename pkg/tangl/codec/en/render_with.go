//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what renderWith — render a node's `with `t` and `u`` definitions-term clause
package en

// renderWith renders a node's `with` clause: the first term follows "with", each
// later term follows "and" (parseWithClause's inverse). An empty list renders to
// the empty string.
func renderWith(with []string) string {
	s := ""
	for i, t := range with {
		if i == 0 {
			s += " with `" + t + "`"
		} else {
			s += " and `" + t + "`"
		}
	}
	return s
}
