//ff:func feature=engine type=engine control=sequence
//ff:what buildTraceEntry — builds one node's TraceEntry from the evalContext
package toulmin

// buildTraceEntry builds the TraceEntry for a single node from the filled evalContext.
// Role is taken from the graph's role map, or inferred when unset. ctx is the shared Ground.
// excluded is keyed by full id; when true the node is crisply superseded out of the tally.
// RuleVerdict is the rule function's own v, or -1 when the node never ran.
func (g *Graph) buildTraceEntry(ec *evalContext, name string, ctx Context, excluded map[string]bool) TraceEntry {
	role := g.roles[name]
	if role == "" {
		role = inferRole(ec.strMap, ec.attackerSet, name)
	}
	ruleVerdict, ok := ec.ruleVerdict[name]
	if !ok {
		ruleVerdict = -1.0
	}
	return TraceEntry{
		Name:        shortName(name),
		RuleVerdict: ruleVerdict,
		Role:        role,
		Activated:   ec.active[name],
		Qualifier:   ec.qualMap[name],
		Verdict:     ec.verdictCache[name],
		Evidence:    ec.evidence[name],
		Ground:      ctx,
		Specs:       ec.specsMap[name],
		Excluded:    excluded[name],
		Level:       ec.effectiveLevel(name, role),
	}
}
