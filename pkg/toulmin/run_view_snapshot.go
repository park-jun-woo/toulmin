//ff:type feature=engine type=engine
//ff:what runView — immutable RunView implementation backed by snapshot maps
package toulmin

// runView is the immutable snapshot implementation of RunView.
type runView struct {
	order     []NodeEvent            // registration order
	byName    map[string]NodeEvent   // shortName → event
	attackers map[string][]NodeEvent // shortName → attacker events
}
