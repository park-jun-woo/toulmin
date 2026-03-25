//ff:func feature=engine type=engine control=sequence
//ff:what InactiveR — test helper: always-inactive rule
package toulmin

func InactiveR(ctx Context, backing Backing) (bool, any) { return false, nil }
