//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what calc — h-Categoriser lazy recursive verdict computation
package toulmin

import "fmt"

// calc computes the verdict for a node using lazy h-Categoriser.
// Executes func on first visit, caches result. Returns [-1.0, +1.0].
// Sets ctx.err and returns -1.0 if a rule function panics.
// Cycle-free graph is guaranteed by detectCycle in newEvalContext.
func (ec *evalContext) calc(id string, ctx Context) float64 {
	if ec.err != nil {
		return -1.0
	}
	fn, ok := ec.fnMap[id]
	if !ok || fn == nil {
		return -1.0
	}
	if !ec.ran[id] {
		ec.ran[id] = true
		active, evidence, err := safeCall(fn, ctx, ec.backingMap[id])
		if err != nil {
			ec.err = fmt.Errorf("rule %q: %w", id, err)
			return -1.0
		}
		ec.active[id] = active
		ec.evidence[id] = evidence
	}
	if !ec.active[id] {
		return -1.0
	}
	sum := 0.0
	if ec.strMap[id] != Strict {
		for _, aid := range ec.edges[id] {
			raw := (ec.calc(aid, ctx) + 1.0) / 2.0
			sum += raw
		}
	}
	raw := ec.qualMap[id] / (1.0 + sum)
	return 2*raw - 1
}
