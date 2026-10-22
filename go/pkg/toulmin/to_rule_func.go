//ff:func feature=engine type=util control=selection
//ff:what toRuleFunc — converts fn (any) to rule func
package toulmin

// toRuleFunc accepts fn as any and returns the rule func.
// Supports func(Context, Specs)(bool, any).
// Panics if fn is not the expected signature.
func toRuleFunc(fn any) func(Context, Specs) (bool, any) {
	switch f := fn.(type) {
	case func(Context, Specs) (bool, any):
		return f
	default:
		panic("toulmin: fn must be func(Context, Specs)(bool, any)")
	}
}
