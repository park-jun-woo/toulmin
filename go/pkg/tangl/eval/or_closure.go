//ff:func feature=tangl type=engine control=sequence
//ff:what orClosure — builds a short-circuit first-true closure over compiled terms
package eval

// orClosure builds a closure that stops at the first true term (returning its
// witness) and otherwise returns false with the last term's witness. It backs
// both Logic "or" and Either.
func orClosure(terms []evalFunc) evalFunc {
	return func(sc *scope) (bool, any) {
		var w any
		for _, ef := range terms {
			b, ww := ef(sc)
			w = ww
			if b {
				return true, ww
			}
		}
		return false, w
	}
}
