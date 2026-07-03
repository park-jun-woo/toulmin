//ff:func feature=engine type=engine control=sequence
//ff:what lvlOptOut — test helper: always-active counter rule that opts out of Level default
package toulmin

func lvlOptOut(ctx Context, specs Specs) (bool, any) { return true, nil }
