//ff:type feature=engine type=engine
//ff:what evalPlan — the part of a graph's evaluation state that no evaluation changes, built once per graph, with the per-evaluation states kept for reuse
package toulmin

import "sync"

// static holds the maps newEvalContext derives from the graph's rules and
// edges (functions, qualifiers, strengths, specs, edges, forces, attacker
// set, supersede map, levels). No evaluation writes to them, so every
// evaluation of the graph shares them. err is the cycle error of the graph,
// if any. pool keeps the states whose per-evaluation maps were cleared after
// an evaluation that handed nothing of them out.
type evalPlan struct {
	static *evalContext
	err    error
	pool   sync.Pool
}
