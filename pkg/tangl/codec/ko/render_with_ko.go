//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what renderWithKo — render the '`a` 그리고 `b` 을 기준으로 하는' adnominal With modifier
package ko

// renderWithKo renders a node's 기준 modifier: backtick terms chained by 그리고,
// then the 받침-selected object particle of the last term, then '기준으로 하는'
// (parseWithModifierKo's inverse). An empty list renders to the empty string.
func renderWithKo(with []string) string {
	if len(with) == 0 {
		return ""
	}
	s := ""
	for i, t := range with {
		if i > 0 {
			s += " 그리고 "
		}
		s += "`" + t + "`"
	}
	last := with[len(with)-1]
	return s + " " + josaFor(last, "을", "를") + " 기준으로 하는"
}
