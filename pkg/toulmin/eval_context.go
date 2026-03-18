//ff:type feature=engine type=engine
//ff:what evalContext — shared state for lazy graph evaluation
package toulmin

// evalContext holds the shared state for h-Categoriser lazy evaluation.
type evalContext struct {
	fnMap       map[string]func(any, any) (bool, any)
	qualMap     map[string]float64
	strMap      map[string]Strength
	edges       map[string][]string
	attackerSet map[string]bool
	ran         map[string]bool
	active      map[string]bool
	evidence    map[string]any
	trace       []TraceEntry
	roleMap     map[string]string
}
