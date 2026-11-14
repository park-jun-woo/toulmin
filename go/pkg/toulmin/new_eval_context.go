//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what newEvalContext — creates evalContext from rules and defeat edges
package toulmin

// newEvalContext builds an evalContext from rules and explicit defeat edges.
// If defeatEdges is nil, edges are derived from RuleMeta.Defeats (every such edge has force 1).
// Each explicit defeat edge's force is stored parallel to edges in forces.
// supersedes populates supersedeBy (target → superseders) for the post-hoc
// exclusion pass; supersede edges never enter the h-Categoriser traversal.
// Returns an error if the defeat graph or the supersede graph contains a cycle.
func newEvalContext(rules []RuleMeta, defeatEdges []defeatEdge, supersedes []defeatEdge, roleMap map[string]string) (*evalContext, error) {
	ec := &evalContext{
		fnMap:         make(map[string]func(Context, Specs) (float64, any)),
		qualMap:       make(map[string]float64),
		strMap:        make(map[string]Strength),
		specsMap:      make(map[string]Specs),
		edges:         make(map[string][]string),
		ran:           make(map[string]bool),
		active:        make(map[string]bool),
		ruleVerdict:   make(map[string]float64),
		applicability: make(map[string]float64),
		forces:        make(map[string][]float64),
		evidence:      make(map[string]any),
		verdictCache:  make(map[string]float64),
		roleMap:       roleMap,
		supersedeBy:   make(map[string][]string),
		levelMap:      make(map[string]Level),
		levelSetMap:   make(map[string]bool),
	}
	for _, r := range rules {
		ec.fnMap[r.Name] = r.Fn
		ec.qualMap[r.Name] = r.Qualifier
		ec.strMap[r.Name] = r.Strength
		ec.specsMap[r.Name] = r.Specs
		ec.levelMap[r.Name] = r.Level
		ec.levelSetMap[r.Name] = r.LevelSet
	}
	if defeatEdges != nil {
		for _, d := range defeatEdges {
			ec.edges[d.to] = append(ec.edges[d.to], d.from)
			ec.forces[d.to] = append(ec.forces[d.to], d.force)
		}
	} else {
		buildEdgesFromRules(ec.edges, rules)
	}
	ec.attackerSet = buildAttackerSet(ec.edges)
	if err := DetectCycle(ec.edges); err != nil {
		return nil, err
	}
	for _, s := range supersedes {
		ec.supersedeBy[s.to] = append(ec.supersedeBy[s.to], s.from)
	}
	if err := detectSupersedeCycle(supersedes); err != nil {
		return nil, err
	}
	return ec, nil
}
