//ff:type feature=tangl type=model
//ff:what Compare — a field/operator/value comparison expression
package ast

// Compare is a leaf condition: `<field> <op> <value>`.
// Op is one of: "is empty", "is not empty", "equals", "is in",
// "is greater than", "is less than", "is at most", "is at least".
//
// Lhs and Rhs are additive operand slots. Invariants (guaranteed by the
// parser, the single producer): Lhs == nil ⇔ Field is valid (a single-step
// backtick left side always uses Field with Lhs nil, so pre-existing
// documents keep a byte-identical AST), and Rhs != nil ⇒ Value == "".
type Compare struct {
	Field string  `json:"field"`
	Op    string  `json:"op"`
	Value string  `json:"value,omitempty"`
	Lhs   Operand `json:"lhs,omitempty"`
	Rhs   Operand `json:"rhs,omitempty"`
}
