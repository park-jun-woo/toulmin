//ff:func feature=engine type=engine control=sequence
//ff:what recordTrace — executes rule function and appends trace entry on first visit
package toulmin

// recordTrace runs the rule function for id and records a TraceEntry.
// Called once per node when calcTrace visits it for the first time.
func (ctx *evalContext) recordTrace(id string, claim, ground any) {
	ctx.ran[id] = true
	ctx.active[id], ctx.evidence[id] = ctx.fnMap[id](claim, ground, ctx.backingMap[id])
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
		Backing:   ctx.backingMap[id],
	})
}
