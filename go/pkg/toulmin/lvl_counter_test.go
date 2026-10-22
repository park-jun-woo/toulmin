//ff:func feature=engine type=engine control=sequence
//ff:what lvlCounter — test helper: always-active counter rule for Level integration
package toulmin

func lvlCounter(ctx Context, specs Specs) (bool, any) { return true, nil }
