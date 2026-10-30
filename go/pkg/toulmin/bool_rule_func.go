//ff:func feature=engine type=util control=sequence
//ff:what boolRuleFunc — wraps a legacy bool rule func as a verdict rule func (true → +1, false → -1)
package toulmin

// boolRuleFunc adapts a legacy func(Context, Specs)(bool, any) to the verdict
// signature: true maps to +1 (fully applies) and false maps to -1 (not applicable).
// Evidence is passed through unchanged.
func boolRuleFunc(f func(Context, Specs) (bool, any)) func(Context, Specs) (float64, any) {
	return func(ctx Context, specs Specs) (float64, any) {
		ok, evidence := f(ctx, specs)
		if ok {
			return 1.0, evidence
		}
		return -1.0, evidence
	}
}
