//ff:func feature=tangl type=parser control=sequence
//ff:what checkOrderedSibling — verify one sibling's ordered-list number matches the running counter
package parser

// checkOrderedSibling checks a single sibling item's ordered-list number (if
// any) against the running expect counter, returning the counter to use for
// the next ordered sibling.
func checkOrderedSibling(it item, expect int, path string) (int, error) {
	if !it.Ordered {
		return expect, nil
	}
	if it.Number != expect {
		return expect, errAt(path, it.Line, "ordered list item numbered %d, expected %d", it.Number, expect)
	}
	return expect + 1, nil
}
