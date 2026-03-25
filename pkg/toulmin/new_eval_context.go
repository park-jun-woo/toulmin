//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what newEvalContext — creates evalContext from rules and defeat edges
package toulmin

// newEvalContext builds an evalContext from rules and explicit defeat edges.
// If defeatEdges is nil, edges are derived from RuleMeta.Defeats.
// Returns an error if the defeat graph contains a cycle.
func newEvalContext(rules []RuleMeta, defeatEdges []defeatEdge, roleMap map[string]string) (*evalContext, error) {
	ec := &evalContext{
		fnMap:    make(map[string]func(Context, Specs) (bool, any)),
		qualMap:  make(map[string]float64),
		strMap:   make(map[string]Strength),
		specsMap: make(map[string]Specs),
		edges:    make(map[string][]string),
		ran:      make(map[string]bool),
		active:   make(map[string]bool),
		evidence: make(map[string]any),
		roleMap:  roleMap,
	}
	for _, r := range rules {
		ec.fnMap[r.Name] = r.Fn
		ec.qualMap[r.Name] = r.Qualifier
		ec.strMap[r.Name] = r.Strength
		ec.specsMap[r.Name] = r.Specs
	}
	if defeatEdges != nil {
		for _, d := range defeatEdges {
			ec.edges[d.to] = append(ec.edges[d.to], d.from)
		}
	} else {
		buildEdgesFromRules(ec.edges, rules)
	}
	ec.attackerSet = buildAttackerSet(ec.edges)
	if err := DetectCycle(ec.edges); err != nil {
		return nil, err
	}
	return ec, nil
}
