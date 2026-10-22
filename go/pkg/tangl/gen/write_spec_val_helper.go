//ff:func feature=tangl type=codegen control=sequence
//ff:what writeSpecValHelper — writes the shared tanglSpecVal function
package gen

import "strings"

// writeSpecValHelper writes the shared tanglSpecVal function: it resolves a
// spec reference by matching the first step against Specs.Find (SpecName) and
// walking any remaining steps via tanglWalk, mirroring
// pkg/tangl/eval.specVal. No matching Spec yields (nil, false).
func writeSpecValHelper(w *strings.Builder) {
	w.WriteString(`func tanglSpecVal(specs toulmin.Specs, steps ...string) (any, bool) {
	sp := specs.Find(steps[0])
	if sp == nil {
		return nil, false
	}
	if len(steps) == 1 {
		return sp, true
	}
	return tanglWalk(sp, steps[1:]...)
}

`)
}
