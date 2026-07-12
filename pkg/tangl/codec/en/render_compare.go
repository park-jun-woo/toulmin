//ff:func feature=tangl type=parser control=sequence
//ff:what renderCompare — render a leaf `field` op [value] comparison (the inverse of parseCompare)
package en

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// renderCompare renders a Compare leaf back to "<lhs> <op> [rhs]". The left side
// is the bare backtick Field unless Lhs is a non-nil operand; the right side is a
// rendered Rhs operand when non-nil, otherwise the verbatim Value (empty for the
// value-less "is empty"/"is not empty" forms). It trusts the parser invariants
// (Lhs == nil ⇔ Field valid, Rhs != nil ⇒ Value == ""), so a Compare with both
// Lhs and Rhs nil renders byte-identically to the pre-existing path.
func renderCompare(c ast.Compare) string {
	lhs := "`" + c.Field + "`"
	if c.Lhs != nil {
		lhs = renderOperand(c.Lhs)
	}
	s := lhs + " " + c.Op
	if c.Rhs != nil {
		return s + " " + renderOperand(c.Rhs)
	}
	if c.Value != "" {
		s += " " + c.Value
	}
	return s
}
