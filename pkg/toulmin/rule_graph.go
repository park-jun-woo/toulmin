//ff:type feature=engine type=model
//ff:what RuleGraph — defeats graph (nodes and attack edges)
package toulmin

// RuleGraph holds nodes and attack edges for verdict computation.
type RuleGraph struct {
	Nodes map[string]*Node
	Edges map[string][]string // target → []attacker names
}
