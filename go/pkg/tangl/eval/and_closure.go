//ff:func feature=tangl type=engine control=sequence
//ff:what andClosure — builds a short-circuit "and" closure over compiled terms
package eval

// andClosure builds a closure that stops at the first false term (returning its
// witness) and otherwise returns true with the last term's witness.
func andClosure(terms []evalFunc) evalFunc {
	return func(sc *scope) (bool, any) {
		var w any
		for _, ef := range terms {
			b, ww := ef(sc)
			w = ww
			if !b {
				return false, ww
			}
		}
		return true, w
	}
}
