//ff:type feature=tangl type=model
//ff:what Call — a host function call escape hatch (FUNC_REF used as value or predicate)
package ast

// Call is a host function call escape hatch: a FUNC_REF used either as a
// value operand (in a "result of" position) or as a predicate term (in a
// "holds" position).
type Call struct {
	Fn Ref `json:"fn"`
}
