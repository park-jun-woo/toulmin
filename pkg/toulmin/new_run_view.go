//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what newRunView — builds the immutable RunView snapshot from a full-pass evalContext
package toulmin

// newRunView classifies every g.rules node from the full-pass evalContext into a
// final NodeEvent and builds the attacker index, producing one immutable snapshot
// shared by all handlers (determinism, regardless of firing order or ctx mutation).
func newRunView(g *Graph, ec *evalContext) *runView {
	v := &runView{
		order:     make([]NodeEvent, 0, len(g.rules)),
		byName:    make(map[string]NodeEvent, len(g.rules)),
		attackers: make(map[string][]NodeEvent, len(ec.edges)),
	}
	for i := range g.rules {
		name := g.rules[i].Name
		ev := NodeEvent{
			Name:     shortName(name),
			Role:     g.roles[name],
			Type:     classifyEvent(ec.active[name], ec.verdictCache[name]),
			Verdict:  ec.verdictCache[name],
			Evidence: ec.evidence[name],
		}
		v.order = append(v.order, ev)
		v.byName[ev.Name] = ev
	}
	for to, attackerIDs := range ec.edges {
		events := make([]NodeEvent, 0, len(attackerIDs))
		for _, aid := range attackerIDs {
			ev := NodeEvent{
				Name:     shortName(aid),
				Role:     g.roles[aid],
				Type:     classifyEvent(ec.active[aid], ec.verdictCache[aid]),
				Verdict:  ec.verdictCache[aid],
				Evidence: ec.evidence[aid],
			}
			events = append(events, ev)
		}
		v.attackers[shortName(to)] = events
	}
	return v
}
