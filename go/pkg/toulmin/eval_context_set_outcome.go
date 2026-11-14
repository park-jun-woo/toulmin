//ff:func feature=engine type=engine control=sequence
//ff:what setOutcome — stores a rule's own verdict, derived active flag, and evidence
package toulmin

// setOutcome records the rule function's result for id: its own verdict v,
// active (v > -1, i.e. the rule applies at least partially), and evidence.
func (ec *evalContext) setOutcome(id string, v float64, evidence any) {
	ec.ruleVerdict[id] = v
	ec.active[id] = v > -1.0
	ec.evidence[id] = evidence
}
