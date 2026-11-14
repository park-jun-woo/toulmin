//ff:func feature=engine type=engine control=sequence
//ff:what generalDiscount — test helper: general discount (supersede target), always applies
package toulmin

func generalDiscount(ctx Context, specs Specs) (bool, any) { return true, nil }
