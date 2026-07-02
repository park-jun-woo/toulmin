//ff:func feature=tangl type=parser control=iteration dimension=1 pattern=recursive
//ff:what checkOrderedSiblings — verify ordered list numbering is 1-based and contiguous
package parser

// checkOrderedSiblings verifies that, within each sibling list (at every
// nesting level), ordered ("N. ") items are numbered 1, 2, 3, ... in
// appearance order. A gap or wrong starting number is a parser error.
func checkOrderedSiblings(items []item, path string) error {
	expect := 1
	for _, it := range items {
		next, err := checkOrderedSibling(it, expect, path)
		if err != nil {
			return err
		}
		expect = next
		if err := checkOrderedSiblings(it.Children, path); err != nil {
			return err
		}
	}
	return nil
}
