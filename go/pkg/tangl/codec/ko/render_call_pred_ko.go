//ff:func feature=tangl type=parser control=sequence
//ff:what renderCallPredKo — render a host predicate call term (`` `alias`.`name` 가 참이면 ``)
package ko

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// renderCallPredKo renders a Call used as a predicate condition term back to its KO
// surface `<FUNC_REF> 가 참이면`, the inverse of parseCallPredKo. The subject 조사 이/가
// is 받침-selected from the reference's name (the last backtick token before 참이면). A
// Call reaches this renderer both as a rule head and as a condition-list item, so it is
// wired into renderCondPredKo alongside Compare.
func renderCallPredKo(c ast.Call) string {
	return renderRefKo(c.Fn) + " " + josaFor(c.Fn.Name, "이", "가") + " 참이면"
}
