//ff:type feature=tangl type=codegen
//ff:what reprCode — a Where-notation Go expression, constant-folded when possible
package gen

// reprCode is the Go expression that produces an operand's Where notation in
// a generated witness Fact. When the notation is compile-time constant
// (Const true) Val holds the raw string and Expr its quoted Go literal;
// a binder-indexed path yields a runtime concatenation Expr instead.
type reprCode struct {
	Expr  string
	Const bool
	Val   string
}
