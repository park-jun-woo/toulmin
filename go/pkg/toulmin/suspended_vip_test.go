//ff:func feature=engine type=engine control=sequence
//ff:what suspendedVIP — test helper: suspended-VIP exception against vipDiscount, always applies
package toulmin

func suspendedVIP(ctx Context, specs Specs) (bool, any) { return true, nil }
