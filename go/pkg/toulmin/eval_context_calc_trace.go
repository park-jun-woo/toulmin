//ff:func feature=engine type=engine control=sequence
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
	if v, ok := ec.verdictCache[id]; ok {
		return v
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
	return ec.combine(id, func(aid string) float64 { return ec.calcTrace(aid, ctx, duration) })
}
