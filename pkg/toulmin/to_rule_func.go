//ff:func feature=engine type=util control=selection
//ff:what toRuleFunc — converts fn (any) to rule func
package toulmin

// toRuleFunc accepts fn as any and returns the rule func.
// Supports func(Context, Backing)(bool, any).
// Panics if fn is not the expected signature.
func toRuleFunc(fn any) func(Context, Backing) (bool, any) {
	switch f := fn.(type) {
	case func(Context, Backing) (bool, any):
		return f
	default:
		panic("toulmin: fn must be func(Context, Backing)(bool, any)")
	}
}
