//ff:func feature=tangl type=codegen control=sequence
//ff:what writeNotHelper — writes the shared tanglNot function
package gen

import "strings"

// writeNotHelper writes the shared tanglNot function: it evaluates the inner
// term and inverts the verdict, propagating the inner witness unchanged,
// mirroring pkg/tangl/eval.compileNot.
func writeNotHelper(w *strings.Builder) {
	w.WriteString(`func tanglNot(term func() (bool, any)) (bool, any) {
	b, w := term()
	return !b, w
}

`)
}
