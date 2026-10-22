//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what newEvalContext — creates evalContext from rules and defeat edges
package toulmin

// newEvalContext builds an evalContext from rules and explicit defeat edges.
// If defeatEdges is nil, edges are derived from RuleMeta.Defeats.
// supersedes populates supersedeBy (target → superseders) for the post-hoc
// exclusion pass; supersede edges never enter the h-Categoriser traversal.
// Returns an error if the defeat graph or the supersede graph contains a cycle.
func newEvalContext(rules []RuleMeta, defeatEdges []defeatEdge, supersedes []defeatEdge, roleMap map[string]string) (*evalContext, error) {
	ec := &evalContext{
		fnMap:        make(map[string]func(Context, Specs) (bool, any)),
		qualMap:      make(map[string]float64),
		strMap:       make(map[string]Strength),
		specsMap:     make(map[string]Specs),
		edges:        make(map[string][]string),
		ran:          make(map[string]bool),
		active:       make(map[string]bool),
		evidence:     make(map[string]any),
		verdictCache: make(map[string]float64),
		roleMap:      roleMap,
		supersedeBy:  make(map[string][]string),
		levelMap:     make(map[string]Level),
		levelSetMap:  make(map[string]bool),
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
