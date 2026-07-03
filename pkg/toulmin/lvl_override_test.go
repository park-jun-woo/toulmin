//ff:func feature=engine type=engine control=sequence
//ff:what lvlOverride — test helper: always-active counter rule with explicit Level override
package toulmin

func lvlOverride(ctx Context, specs Specs) (bool, any) { return true, nil }
