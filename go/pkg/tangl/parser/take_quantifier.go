//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what takeQuantifier — consume a leading quantifier keyword (any / all / count / sum)
package parser

// quantifierKeywords lists the recognized quantifier keywords. Only "any" is
// implemented; "all", "count", and "sum" are reserved and rejected with a
// clear error by parseExistsTerm.
var quantifierKeywords = []string{"any", "all", "count", "sum"}

// takeQuantifier consumes a leading quantifier keyword from s and returns the
// matched keyword, the remaining text, and whether a quantifier was found.
func takeQuantifier(s string) (kw string, rest string, ok bool) {
	for _, q := range quantifierKeywords {
		if r, matched := takeKeyword(s, q); matched {
			return q, r, true
		}
	}
	return "", s, false
}
