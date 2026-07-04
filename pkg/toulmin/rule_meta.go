//ff:type feature=rule type=model
//ff:what RuleMeta — rule metadata (name, qualifier, strength, defeats, specs, fn)
package toulmin

// RuleMeta holds a rule's metadata and its boolean evaluation function.
type RuleMeta struct {
	Name      string
	Qualifier float64
	Strength  Strength
	Defeats   []string
	Specs     Specs
	Grounds   []string // ground fact names this rule's Fn resolves at eval time (audit metadata only; engine does not resolve/cache)
	Fn        func(ctx Context, specs Specs) (bool, any)

	Level    Level // crisp reporting level (Unset unless set)
	LevelSet bool  // true once Rule.Level was called; explicit value wins over role default

	RunOn NodeHandler // run handler, fires on Active (nil allowed)

	RunGraph *Graph // sub-graph to Run when this node is Active (nil allowed)
}
