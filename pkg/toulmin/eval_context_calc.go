//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what calc — h-Categoriser lazy recursive verdict computation
package toulmin

// calc computes the verdict for a node using lazy h-Categoriser.
// Executes func on first visit, caches result. Returns [-1.0, +1.0].
func (ctx *evalContext) calc(id string, claim, ground any, depth int) float64 {
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
	}
	if !ctx.active[id] {
		return -1.0
	}
	sum := 0.0
	if ctx.strMap[id] != Strict {
		for _, aid := range ctx.edges[id] {
			raw := (ctx.calc(aid, claim, ground, depth+1) + 1.0) / 2.0
			sum += raw
		}
	}
	raw := ctx.qualMap[id] / (1.0 + sum)
	return 2*raw - 1
}
