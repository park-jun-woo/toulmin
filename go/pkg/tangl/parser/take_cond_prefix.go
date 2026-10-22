//ff:func feature=tangl type=parser control=sequence
//ff:what takeCondPrefix — consume a leading "and"/"or" condition-chaining prefix
package parser

// takeCondPrefix consumes a leading "and" or "or" keyword used to chain a
// condition-list item onto the running expression.
func takeCondPrefix(s string) (op string, rest string, ok bool) {
	if r, ok := takeKeyword(s, "and"); ok {
		return "and", r, true
	}
	if r, ok := takeKeyword(s, "or"); ok {
		return "or", r, true
	}
	return "", s, false
}
