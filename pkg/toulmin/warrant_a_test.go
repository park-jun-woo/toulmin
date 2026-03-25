//ff:func feature=engine type=engine control=sequence
//ff:what WarrantA — test helper: always-active warrant rule
package toulmin

func WarrantA(ctx Context, backing Backing) (bool, any) { return true, nil }
