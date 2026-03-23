//ff:func feature=engine type=engine control=selection
//ff:what evalRule — dispatches to calc or calcTrace based on option
package toulmin

// evalRule selects the appropriate calculation method based on EvalOption.
func (ctx *evalContext) evalRule(id string, claim, ground any, opt EvalOption) float64 {
	switch {
	case opt.Trace:
		return ctx.calcTrace(id, claim, ground, opt.Duration)
	default:
		return ctx.calc(id, claim, ground)
	}
}
