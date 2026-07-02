//ff:func feature=tangl type=parser control=sequence
//ff:what takeArticle — consume a leading "a" or "an" article keyword
package parser

// takeArticle consumes a leading "a" or "an" keyword, as used before
// ROLE_EXPR in a node declaration ("is a general rule" / "is an except rule").
func takeArticle(s string) (string, bool) {
	if rest, ok := takeKeyword(s, "a"); ok {
		return rest, true
	}
	if rest, ok := takeKeyword(s, "an"); ok {
		return rest, true
	}
	return s, false
}
