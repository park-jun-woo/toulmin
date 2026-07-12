//ff:func feature=tangl type=codegen control=sequence
//ff:what writeHoldsHelper — writes the shared tanglHolds function
package gen

import "strings"

// writeHoldsHelper writes the shared tanglHolds function: the predicate-call
// witness rule mirroring pkg/tangl/eval.compileCall — a non-nil host evidence
// passes through unchanged (never overwritten by a machine Fact); otherwise
// the witness is a Fact recording the call and its boolean result.
func writeHoldsHelper(w *strings.Builder) {
	w.WriteString(`func tanglHolds(b bool, ev any, where string) (bool, any) {
	if ev != nil {
		return b, ev
	}
	return b, toulmin.Fact{Where: where, Expected: "holds", Actual: b}
}

`)
}
