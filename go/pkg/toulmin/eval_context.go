//ff:type feature=engine type=engine
//ff:what evalContext — shared state for lazy graph evaluation
package toulmin

// evalContext holds the shared state for h-Categoriser lazy evaluation.
type evalContext struct {
	fnMap         map[string]func(Context, Specs) (float64, any)
	qualMap       map[string]float64
	strMap        map[string]Strength
	specsMap      map[string]Specs
	edges         map[string][]string
	attackerSet   map[string]bool
	ran           map[string]bool
	active        map[string]bool      // derived: ruleVerdict > -1
	ruleVerdict   map[string]float64   // rule function's own return value v in [-1, 1]
	applicability map[string]float64   // a_eff = (v+1)/2 · Π_except (1 − F·(v_j+1)/2)
	forces        map[string][]float64 // forces[to][i] is the force of edges[to][i]; missing → 1
	evidence      map[string]any
	trace         []TraceEntry
	roleMap       map[string]string
	err           error
	verdictCache  map[string]float64
	supersedeBy   map[string][]string
	levelMap      map[string]Level
	levelSetMap   map[string]bool
}
