//ff:func feature=tangl type=codegen control=sequence
//ff:what writeCtxValHelper — writes the shared tanglCtxVal function
package gen

import "strings"

// writeCtxValHelper writes the shared tanglCtxVal function: it resolves a
// ctx-rooted data path by fetching the first step from ctx and walking the
// remaining steps via tanglWalk, mirroring pkg/tangl/eval.pathVal for the
// non-binder case. A missing ctx key yields (nil, false).
func writeCtxValHelper(w *strings.Builder) {
	w.WriteString(`func tanglCtxVal(ctx toulmin.Context, steps ...string) (any, bool) {
	v, ok := ctx.Get(steps[0])
	if !ok {
		return nil, false
	}
	return tanglWalk(v, steps[1:]...)
}

`)
}
