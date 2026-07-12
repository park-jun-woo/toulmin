//ff:func feature=tangl type=parser control=selection
//ff:what renderCondPredKo — render a single-term condition (Compare / Not{Compare}) as one predicate line
package ko

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// renderCondPredKo renders a condition that fits on one line — a leaf Compare, a
// negated Compare (ast.Not wrapping a Compare), a host predicate Call (“ … 참이면 “),
// or a one-line Exists (whose body is itself a single leaf) — to its predicate text.
// ok is false for a Logic or Either, a nested Exists (a multi-term body), and the
// impossible Not{non-Compare}, which all need the nested 다음이면/다음 중/다음인 multi-line
// form, so callers fall back to the nested renderer.
func renderCondPredKo(e ast.Expr) (string, bool, error) {
	switch t := e.(type) {
	case ast.Compare:
		s, err := renderCompareKo(t, false)
		return s, true, err
	case ast.Not:
		if inner, ok := t.Term.(ast.Compare); ok {
			s, err := renderCompareKo(inner, true)
			return s, true, err
		}
		return "", false, nil
	case ast.Call:
		return renderCallPredKo(t), true, nil
	case ast.Exists:
		return renderExistsPredKo(t)
	default:
		return "", false, nil
	}
}
