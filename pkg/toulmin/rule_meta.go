//ff:type feature=rule type=model
//ff:what RuleMeta — rule metadata (name, qualifier, strength, defeats, backing, fn)
package toulmin

// RuleMeta holds a rule's metadata and its boolean evaluation function.
type RuleMeta struct {
	Name      string
	Qualifier float64
	Strength  Strength
	Defeats   []string
	Backing   string
	Fn        func(claim any, ground any) bool
}
