//ff:func feature=tangl type=codegen control=selection
//ff:what renderOperandVal — dispatches an Operand to its per-kind valCode emitter
package gen

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderOperandVal dispatches an Operand to its per-kind valCode emitter
// (pathValCode, specValCode, callValCode). needOK asks for the resolution
// flag expression, needActual (lhs position) additionally asks for the
// witness Fact's Actual expression, which differs from the value only for a
// host Call whose error becomes the Actual.
func renderOperandVal(o ast.Operand, env map[string]binderInfo, prefix string, needOK, needActual bool) (valCode, error) {
	switch v := o.(type) {
	case ast.Path:
		return pathValCode(v, env, prefix, needOK), nil
	case ast.SpecRef:
		return specValCode(v, prefix, needOK), nil
	case ast.Call:
		return callValCode(v, prefix, needOK, needActual), nil
	default:
		return valCode{}, fmt.Errorf("gen: unsupported operand node %T", o)
	}
}
