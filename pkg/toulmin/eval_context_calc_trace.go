//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what calcTrace — h-Categoriser lazy recursive verdict computation with trace
package toulmin

// calcTrace computes the verdict like calc, but also collects TraceEntry
// for each executed rule. Role is resolved from roleMap or inferred.
// When duration is true, measures execution time per rule.
// Returns -1.0 immediately if ctx.err is set.
// Cycle-free graph is guaranteed by detectCycle in newEvalContext.
func (ctx *evalContext) calcTrace(id string, claim, ground any, duration bool) float64 {
	if ctx.err != nil {
		return -1.0
	}
	fn, ok := ctx.fnMap[id]
	if !ok || fn == nil {
		return -1.0
	}
	if !ctx.ran[id] {
		ctx.recordTrace(id, claim, ground, duration)
	}
	if ctx.err != nil {
		return -1.0
	}
	if !ctx.active[id] {
		return -1.0
	}
	sum := 0.0
	if ctx.strMap[id] != Strict {
		for _, aid := range ctx.edges[id] {
			raw := (ctx.calcTrace(aid, claim, ground, duration) + 1.0) / 2.0
			sum += raw
		}
	}
	raw := ctx.qualMap[id] / (1.0 + sum)
	return 2*raw - 1
}
