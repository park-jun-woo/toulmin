//ff:func feature=engine type=engine control=sequence
//ff:what reset — resets per-warrant evaluation state
package toulmin

// reset clears ran, active, rule verdicts, applicability, evidence, and trace for a fresh warrant evaluation.
func (ctx *evalContext) reset() {
	ctx.ran = make(map[string]bool)
	ctx.active = make(map[string]bool)
	ctx.ruleVerdict = make(map[string]float64)
	ctx.applicability = make(map[string]float64)
	ctx.evidence = make(map[string]any)
	ctx.trace = nil
	ctx.err = nil
	ctx.verdictCache = make(map[string]float64)
}
