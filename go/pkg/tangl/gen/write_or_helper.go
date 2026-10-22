//ff:func feature=tangl type=codegen control=sequence
//ff:what writeOrHelper — writes the shared tanglOr function
package gen

import "strings"

// writeOrHelper writes the shared tanglOr function: a short-circuit
// first-true "or" over term closures that stops at the first true term
// (returning its witness) and otherwise returns false with the last term's
// witness, mirroring pkg/tangl/eval.orClosure. It backs both Logic "or" and
// Either groups.
func writeOrHelper(w *strings.Builder) {
	w.WriteString(`func tanglOr(terms ...func() (bool, any)) (bool, any) {
	var w any
	for _, t := range terms {
		b, ww := t()
		w = ww
		if b {
			return true, ww
		}
	}
	return false, w
}

`)
}
