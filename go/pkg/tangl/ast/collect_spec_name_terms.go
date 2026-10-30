//ff:func feature=tangl type=model control=iteration dimension=1 pattern=recursive
//ff:what collectSpecNameTerms — collect SpecRef names from each term of a Logic/Either group in order
package ast

// collectSpecNameTerms applies collectSpecNames to every term in order.
func collectSpecNameTerms(terms []Expr, seen map[string]bool, names []string) []string {
	for _, term := range terms {
		names = collectSpecNames(term, seen, names)
	}
	return names
}
