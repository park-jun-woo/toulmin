//ff:type feature=engine type=engine
//ff:what Engine — rule engine (registers rules, evaluates claims)
package toulmin

// Engine accumulates rules and evaluates claims against them.
type Engine struct {
	rules []RuleMeta
}
