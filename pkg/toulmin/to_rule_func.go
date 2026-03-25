//ff:func feature=engine type=util control=selection
//ff:what toRuleFunc — converts fn (any) to 3-arg rule func, supporting both signatures
package toulmin

// toRuleFunc accepts fn as any and returns the 3-arg rule func.
// Supports func(any,any,Backing)(bool,any) and legacy func(any,any)(bool,any).
// Panics if fn is neither.
func toRuleFunc(fn any) func(any, any, Backing) (bool, any) {
	switch f := fn.(type) {
	case func(any, any, Backing) (bool, any):
		return f
	case func(any, any) (bool, any):
		return wrapLegacy(f)
	default:
		panic("toulmin: fn must be func(any,any,Backing)(bool,any) or func(any,any)(bool,any)")
	}
}
