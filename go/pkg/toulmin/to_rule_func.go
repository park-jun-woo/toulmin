//ff:func feature=engine type=util control=selection
//ff:what toRuleFunc — converts fn (any) to a verdict-returning rule func
package toulmin

// toRuleFunc accepts fn as any and returns the rule func.
// Supports func(Context, Specs)(float64, any) as-is, and the legacy
// func(Context, Specs)(bool, any) wrapped so that true → +1 and false → -1.
// Panics if fn is neither signature.
func toRuleFunc(fn any) func(Context, Specs) (float64, any) {
	switch f := fn.(type) {
	case func(Context, Specs) (float64, any):
		return f
	case func(Context, Specs) (bool, any):
		return boolRuleFunc(f)
	default:
		panic("toulmin: fn must be func(Context, Specs)(bool, any) or func(Context, Specs)(float64, any)")
	}
}
