//ff:func feature=tangl type=codegen control=sequence
//ff:what writeAndHelper — writes the shared tanglAnd function
package gen

import "strings"

// writeAndHelper writes the shared tanglAnd function: a short-circuit "and"
// over term closures that stops at the first false term (returning its
// witness) and otherwise returns true with the last term's witness, mirroring
// pkg/tangl/eval.andClosure.
func writeAndHelper(w *strings.Builder) {
	w.WriteString(`func tanglAnd(terms ...func() (bool, any)) (bool, any) {
	var w any
	for _, t := range terms {
		b, ww := t()
		w = ww
		if !b {
			return false, ww
		}
	}
	return true, w
}

`)
}
