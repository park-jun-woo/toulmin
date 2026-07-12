//ff:func feature=tangl type=engine control=selection
//ff:what reprOperand — renders an Operand as a stable Where string for evidence
package eval

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// reprOperand renders an Operand as the human-readable Where string used in
// evidence Facts: a Path via pathRepr, a SpecRef as its dotted steps, and a
// Call as its reference.
func reprOperand(o ast.Operand, sc *scope) string {
	switch v := o.(type) {
	case ast.Path:
		return pathRepr(v, sc)
	case ast.SpecRef:
		return strings.Join(v.Steps, ".")
	case ast.Call:
		return callRepr(v)
	default:
		return ""
	}
}
