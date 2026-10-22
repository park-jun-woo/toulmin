//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what walkSteps — walks map-key/struct-field steps into a value, returning (value, ok)
package eval

// walkSteps walks the given steps into cur via stepInto. A missing key/field or
// nil intermediate yields ok=false; with no steps it returns cur unchanged.
func walkSteps(cur any, steps []string) (any, bool) {
	for _, step := range steps {
		v, ok := stepInto(cur, step)
		if !ok {
			return nil, false
		}
		cur = v
	}
	return cur, true
}
