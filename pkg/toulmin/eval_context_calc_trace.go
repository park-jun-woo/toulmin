//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what calcTrace — h-Categoriser lazy recursive verdict computation with trace
package toulmin

// calcTrace computes the verdict like calc, but also collects TraceEntry
// for each executed rule. Role is resolved from roleMap or inferred.
func (ctx *evalContext) calcTrace(id string, claim, ground any, depth int) float64 {
	if depth >= maxDepth {
		return 0.0
	}
	fn, ok := ctx.fnMap[id]
	if !ok || fn == nil {
		return -1.0
	}
	if !ctx.ran[id] {
		ctx.ran[id] = true
		ctx.active[id], ctx.evidence[id] = fn(claim, ground)
		role := ctx.roleMap[id]
		if role == "" {
			role = inferRole(ctx.strMap, ctx.attackerSet, id)
		}
		ctx.trace = append(ctx.trace, TraceEntry{
			Name:      id,
			Role:      role,
			Activated: ctx.active[id],
			Qualifier: ctx.qualMap[id],
			Evidence:  ctx.evidence[id],
		})
	}
	if !ctx.active[id] {
		return -1.0
	}
	sum := 0.0
	if ctx.strMap[id] != Strict {
		for _, aid := range ctx.edges[id] {
			raw := (ctx.calcTrace(aid, claim, ground, depth+1) + 1.0) / 2.0
			sum += raw
		}
	}
	raw := ctx.qualMap[id] / (1.0 + sum)
	return 2*raw - 1
}
