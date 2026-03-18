//ff:func feature=rule type=engine control=sequence
//ff:what Register — adds a rule to the engine
package toulmin

// Register appends a rule to the engine's rule set.
func (e *Engine) Register(meta RuleMeta) {
	e.rules = append(e.rules, meta)
}
