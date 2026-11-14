//ff:func feature=engine type=engine control=sequence
//ff:what release — clears an evalContext's per-evaluation maps and keeps it for the plan's next evaluation
package toulmin

// release is for a state nothing outside the evaluation holds any more. The
// maps are cleared here, so that the evidence of an evaluation is not kept
// alive by the pool.
func (p *evalPlan) release(ec *evalContext) {
	clear(ec.ran)
	clear(ec.active)
	clear(ec.ruleVerdict)
	clear(ec.applicability)
	clear(ec.evidence)
	clear(ec.verdictCache)
	ec.trace = nil
	ec.err = nil
	p.pool.Put(ec)
}
