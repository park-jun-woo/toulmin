//ff:func feature=tangl type=parser control=sequence
//ff:what renderCompare — render a leaf `field` op [value] comparison (the inverse of parseCompare)
package en

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// renderCompare renders a Compare leaf back to "`field` <op> [value]". Op is the
// canonical English operator string, and Value is appended verbatim only when
// non-empty (the value-less "is empty"/"is not empty" forms carry no Value).
func renderCompare(c ast.Compare) string {
	s := "`" + c.Field + "` " + c.Op
	if c.Value != "" {
		s += " " + c.Value
	}
	return s
}
