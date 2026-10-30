//ff:func feature=engine type=engine control=sequence
//ff:what calc — h-Categoriser lazy recursive verdict computation
package toulmin

import "fmt"

// calc computes the verdict for a node using lazy h-Categoriser.
// Executes func on first visit, caches result. Returns [-1.0, +1.0].
// A node whose rule verdict is -1 (not applicable) returns -1 without
// evaluating its attackers (lazy evaluation).
// Sets ec.err and returns -1.0 if a rule function panics or returns an
// out-of-range verdict.
// Cycle-free graph is guaranteed by detectCycle in newEvalContext.
func (ec *evalContext) calc(id string, ctx Context) float64 {
	if ec.err != nil {
		return -1.0
	}
	if v, ok := ec.verdictCache[id]; ok {
		return v
	}
	fn, ok := ec.fnMap[id]
	if !ok || fn == nil {
		return -1.0
	}
	if !ec.ran[id] {
		ec.ran[id] = true
		v, evidence, err := safeCall(fn, ctx, ec.specsMap[id])
		if err != nil {
			ec.err = fmt.Errorf("rule %q: %w", id, err)
			return -1.0
		}
		ec.setOutcome(id, v, evidence)
	}
	if !ec.active[id] {
		return -1.0
	}
	return ec.combine(id, func(aid string) float64 { return ec.calc(aid, ctx) })
}
