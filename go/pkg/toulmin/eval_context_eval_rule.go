//ff:func feature=engine type=engine control=selection
//ff:what evalRule — dispatches to calc or calcTrace based on option
package toulmin

// evalRule selects the appropriate calculation method based on EvalOption.
func (ec *evalContext) evalRule(id string, ctx Context, opt EvalOption) float64 {
	switch {
	case opt.Trace:
		return ec.calcTrace(id, ctx, opt.Duration)
	default:
		return ec.calc(id, ctx)
	}
}
