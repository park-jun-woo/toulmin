//ff:func feature=tangl type=parser control=selection
//ff:what parseRoleKo — map a KO role word (일반/반박/예외) to its ast.Role
package ko

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// parseRoleKo maps the role word that immediately precedes 규칙이다 — 일반, 반박, or
// 예외 — to its canonical ast.Role (general/counter/except). ok is false for any
// other word. It mirrors pkg/tangl/parser.parseRoleExpr's role set.
func parseRoleKo(word string) (ast.Role, bool) {
	switch word {
	case "일반":
		return ast.GeneralRule, true
	case "반박":
		return ast.CounterRule, true
	case "예외":
		return ast.ExceptRule, true
	default:
		return 0, false
	}
}
