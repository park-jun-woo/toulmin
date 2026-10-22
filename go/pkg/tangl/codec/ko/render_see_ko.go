//ff:func feature=tangl type=parser control=sequence
//ff:what renderSeeKo — render a See declaration (`` `p` 는 `path` 에서 참조한다 ``)
package ko

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// renderSeeKo renders a tangl:참조 entry back to "`<alias>` 는 `<package>` 에서 참조한다",
// the inverse of parseSeeKo. The 는 particle after the alias is chosen by 받침 (은/는);
// 에서 is fixed. Either allomorph re-parses identically (both are in the ignore set).
func renderSeeKo(s ast.See) string {
	return "`" + s.Alias + "` " + josaFor(s.Alias, "은", "는") + " `" + s.Package + "` 에서 참조한다"
}
