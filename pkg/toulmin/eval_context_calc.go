//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what calc — h-Categoriser lazy recursive verdict computation
package toulmin

import "fmt"

// calc computes the verdict for a node using lazy h-Categoriser.
// Executes func on first visit, caches result. Returns [-1.0, +1.0].
// Sets ctx.err and returns -1.0 if a rule function panics.
// Cycle-free graph is guaranteed by detectCycle in newEvalContext.
func (ctx *evalContext) calc(id string, claim, ground any) float64 {
	if ctx.err != nil {
		return -1.0
	}
	fn, ok := ctx.fnMap[id]
	if !ok || fn == nil {
		return -1.0
	}
	if !ctx.ran[id] {
		ctx.ran[id] = true
		active, evidence, err := safeCall(fn, claim, ground, ctx.backingMap[id])
		if err != nil {
			ctx.err = fmt.Errorf("rule %q: %w", id, err)
			return -1.0
		}
		ctx.active[id] = active
		ctx.evidence[id] = evidence
	}
	if !ctx.active[id] {
		return -1.0
	}
	sum := 0.0
	if ctx.strMap[id] != Strict {
		for _, aid := range ctx.edges[id] {
			raw := (ctx.calc(aid, claim, ground) + 1.0) / 2.0
			sum += raw
		}
	}
	raw := ctx.qualMap[id] / (1.0 + sum)
	return 2*raw - 1
}
