//ff:func feature=engine type=engine control=sequence
//ff:what lvlExcept — test helper: always-active except rule for Level integration
package toulmin

func lvlExcept(ctx Context, specs Specs) (bool, any) { return true, nil }
