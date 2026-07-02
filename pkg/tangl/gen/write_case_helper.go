//ff:func feature=tangl type=codegen control=sequence
//ff:what writeCaseHelper — writes the shared tanglCaseActive helper
package gen

import "strings"

// writeCaseHelper writes the shared tanglCaseActive helper used by
// "checking" wrappers and Internal "every ... until" clauses to compose a
// case's Evaluate results (one per active warrant node) into a single
// active/inactive verdict: active if any warrant result is active.
func writeCaseHelper(w *strings.Builder) {
	w.WriteString(`func tanglCaseActive(results []toulmin.EvalResult) bool {
	for _, r := range results {
		if r.Verdict > 0 {
			return true
		}
	}
	return false
}

`)
}
