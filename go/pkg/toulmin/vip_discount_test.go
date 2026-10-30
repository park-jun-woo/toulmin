//ff:func feature=engine type=engine control=sequence
//ff:what vipDiscount — test helper: VIP discount (superseder), always applies
package toulmin

func vipDiscount(ctx Context, specs Specs) (bool, any) { return true, nil }
