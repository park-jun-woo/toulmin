//ff:func feature=engine type=util control=sequence
//ff:what wrapLegacy — wraps legacy 2-arg rule func to 3-arg backing-aware func
package toulmin

// wrapLegacy converts a legacy func(claim, ground) to func(claim, ground, backing).
func wrapLegacy(fn func(any, any) (bool, any)) func(any, any, any) (bool, any) {
	return func(claim any, ground any, backing any) (bool, any) {
		return fn(claim, ground)
	}
}
