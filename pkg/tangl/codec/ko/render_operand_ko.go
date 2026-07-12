//ff:func feature=tangl type=parser control=selection
//ff:what renderOperandKo — render a value operand (Path / SpecRef / Call) plus its 조사 base token
package ko

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderOperandKo renders a value operand to its KO surface and returns the token that
// drives 받침-based 조사 selection when the operand is a comparison subject or object:
// a Path/SpecRef is the 의-chain (spec adds the 기준 keyword) with the last step as the
// 조사 base, and a Call is the FUNC_REF plus the 의 결과 marker with 결과 as the base. It
// is the inverse of parseOperandKo and is shared by the 서술형 and 관형형 comparison
// renderers and the exists renderers.
func renderOperandKo(o ast.Operand) (text string, josaBase string, err error) {
	switch v := o.(type) {
	case ast.Path:
		return renderPathChainKo(v.Steps), v.Steps[len(v.Steps)-1], nil
	case ast.SpecRef:
		return "기준 " + renderPathChainKo(v.Steps), v.Steps[len(v.Steps)-1], nil
	case ast.Call:
		return renderRefKo(v.Fn) + " 의 결과", "결과", nil
	default:
		return "", "", fmt.Errorf("ko: unsupported operand %T", o)
	}
}
