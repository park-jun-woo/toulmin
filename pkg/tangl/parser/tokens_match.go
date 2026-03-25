//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what tokensMatch — check if word slice prefix matches token slice
package parser

// tokensMatch checks if the beginning of words matches the given tokens.
func tokensMatch(words []string, tokens []string) bool {
	for i, t := range tokens {
		if words[i] != t {
			return false
		}
	}
	return true
}
