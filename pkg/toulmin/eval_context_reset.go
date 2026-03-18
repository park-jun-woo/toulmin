//ff:func feature=engine type=engine control=sequence
//ff:what reset — resets per-warrant evaluation state
package toulmin

// reset clears ran, active, evidence, and trace for a fresh warrant evaluation.
func (ctx *evalContext) reset() {
	ctx.ran = make(map[string]bool)
	ctx.active = make(map[string]bool)
	ctx.evidence = make(map[string]any)
	ctx.trace = nil
}
