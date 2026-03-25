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
	Fn        func(ctx Context, specs Specs) (bool, any)
}
