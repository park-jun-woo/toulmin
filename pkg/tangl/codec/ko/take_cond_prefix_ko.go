//ff:func feature=tangl type=parser control=sequence
//ff:what takeCondPrefixKo — consume a leading 그리고/또는 condition-chaining prefix
package ko

// takeCondPrefixKo consumes a leading 그리고 or 또는 keyword used to chain a
// condition-list item onto the running expression, returning the canonical
// English operator ("and"/"or"). It mirrors pkg/tangl/parser.takeCondPrefix.
func takeCondPrefixKo(s string) (op string, rest string, ok bool) {
	if r, ok := takeKeywordKo(s, "그리고"); ok {
		return "and", r, true
	}
	if r, ok := takeKeywordKo(s, "또는"); ok {
		return "or", r, true
	}
	return "", s, false
}
