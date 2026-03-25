//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what calcTrace — h-Categoriser lazy recursive verdict computation with trace
package toulmin

// calcTrace computes the verdict like calc, but also collects TraceEntry
// for each executed rule. Role is resolved from roleMap or inferred.
// When duration is true, measures execution time per rule.
// Returns -1.0 immediately if ec.err is set.
// Cycle-free graph is guaranteed by detectCycle in newEvalContext.
func (ec *evalContext) calcTrace(id string, ctx Context, duration bool) float64 {
	if ec.err != nil {
		return -1.0
	}
	fn, ok := ec.fnMap[id]
	if !ok || fn == nil {
		return -1.0
	}
	if !ec.ran[id] {
		ec.recordTrace(id, ctx, duration)
	}
	if ec.err != nil {
		return -1.0
	}
	if !ec.active[id] {
		return -1.0
	}
	sum := 0.0
	if ec.strMap[id] != Strict {
		for _, aid := range ec.edges[id] {
			raw := (ec.calcTrace(aid, ctx, duration) + 1.0) / 2.0
			sum += raw
		}
	}
	raw := ec.qualMap[id] / (1.0 + sum)
	return 2*raw - 1
}
