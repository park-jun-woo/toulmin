//ff:func feature=tangl type=parser control=selection
//ff:what renderOperand — render a value-producing operand (path / spec ref / host value call)
package en

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// renderOperand renders an ast.Operand back to its English surface form: a bare
// backtick dot-chain for a Path, "spec " plus the chain for a SpecRef, and
// "result of " plus the FUNC_REF for a value-position Call. It is the exact
// inverse of the parser's parseOperand, so the output re-parses to the same
// operand.
func renderOperand(o ast.Operand) string {
	switch v := o.(type) {
	case ast.Path:
		return renderPath(v.Steps)
	case ast.SpecRef:
		return "spec " + renderPath(v.Steps)
	case ast.Call:
		return "result of " + renderRef(v.Fn)
	default:
		return "`?`"
	}
}
