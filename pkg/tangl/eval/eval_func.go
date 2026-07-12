//ff:type feature=tangl type=engine
//ff:what evalFunc — a compiled Expr closure evaluating to (verdict, evidence)
package eval

// evalFunc is a compiled condition-tree closure: given a runtime scope it
// returns the boolean verdict and the witness evidence for that verdict.
type evalFunc func(sc *scope) (bool, any)
