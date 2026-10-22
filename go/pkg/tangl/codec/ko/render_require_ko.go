//ff:func feature=tangl type=parser control=sequence
//ff:what renderRequireKo — render a required-field declaration (`` `f` 는 필수다[, 타입 T] ``)
package ko

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// renderRequireKo renders a required context field back to "`<field>` 는 필수다", with
// an optional ", 타입 <Type>" clause. It serves both 사안 and 업무 bodies (identical
// surface form), inverting parseRequireKo. The 는(은/는) particle is 받침-selected.
func renderRequireKo(r ast.Require) string {
	s := "`" + r.Field + "` " + josaFor(r.Field, "은", "는") + " 필수다"
	if r.Type != "" {
		s += ", 타입 " + renderTypeKo(r.Type)
	}
	return s
}
