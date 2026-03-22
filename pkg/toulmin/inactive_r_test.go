//ff:func feature=engine type=engine control=sequence
//ff:what InactiveR — test helper: always-inactive rule
package toulmin

func InactiveR(claim any, ground any, backing any) (bool, any) { return false, nil }
