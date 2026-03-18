//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what newEvalContext — creates evalContext from rules and defeat edges
package toulmin

// newEvalContext builds an evalContext from rules and explicit defeat edges.
// If defeatEdges is nil, edges are derived from RuleMeta.Defeats.
func newEvalContext(rules []RuleMeta, defeatEdges []defeatEdge, roleMap map[string]string) *evalContext {
	ctx := &evalContext{
		fnMap:    make(map[string]func(any, any) (bool, any)),
		qualMap:  make(map[string]float64),
		strMap:   make(map[string]Strength),
		edges:    make(map[string][]string),
		ran:      make(map[string]bool),
		active:   make(map[string]bool),
		evidence: make(map[string]any),
		roleMap:  roleMap,
	}
	for _, r := range rules {
		ctx.fnMap[r.Name] = r.Fn
		ctx.qualMap[r.Name] = r.Qualifier
		ctx.strMap[r.Name] = r.Strength
	}
	if defeatEdges != nil {
		for _, d := range defeatEdges {
			ctx.edges[d.to] = append(ctx.edges[d.to], d.from)
		}
	} else {
		buildEdgesFromRules(ctx.edges, rules)
	}
	ctx.attackerSet = buildAttackerSet(ctx.edges)
	return ctx
}
