//ff:type feature=tangl type=engine
//ff:what valueFunc — a compiled Operand closure resolving to (value, ok, error)
package eval

// valueFunc is a compiled value-operand closure: given a runtime scope it
// returns the resolved value, whether a value was produced, and any error
// raised by a host value Call.
type valueFunc func(sc *scope) (any, bool, error)
