//ff:func feature=engine type=engine control=sequence
//ff:what buildTraceEntry — builds one node's TraceEntry from the evalContext
package toulmin

// buildTraceEntry builds the TraceEntry for a single node from the filled evalContext.
// Role is taken from the graph's role map, or inferred when unset. ctx is the shared Ground.
func (g *Graph) buildTraceEntry(ec *evalContext, name string, ctx Context) TraceEntry {
	role := g.roles[name]
	if role == "" {
		role = inferRole(ec.strMap, ec.attackerSet, name)
	}
	return TraceEntry{
		Name:      shortName(name),
		Role:      role,
		Activated: ec.active[name],
		Qualifier: ec.qualMap[name],
		Verdict:   ec.verdictCache[name],
		Evidence:  ec.evidence[name],
		Ground:    ctx,
		Specs:     ec.specsMap[name],
	}
}
