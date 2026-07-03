//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what renderWithKo — render a node's `, 기준 `a` 그리고 `b`` Definitions-term clause
package ko

// renderWithKo renders a node's 기준 clause: the first term follows ", 기준", each later
// term follows " 그리고" (parseWithTermsKo's inverse). An empty list renders to the
// empty string.
func renderWithKo(with []string) string {
	s := ""
	for i, t := range with {
		if i == 0 {
			s += ", 기준 `" + t + "`"
		} else {
			s += " 그리고 `" + t + "`"
		}
	}
	return s
}
