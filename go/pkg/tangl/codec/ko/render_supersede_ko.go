//ff:func feature=tangl type=parser control=sequence
//ff:what renderSupersedeKo — render a supersede edge (`` `X` 면 `Y` 배제한다 ``, superseder-first)
package ko

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// renderSupersedeKo renders a supersede edge back to "`<superseder>` 면 `<superseded>` 배제한다".
// Unlike renderAttackKo the order is NOT reversed: the first backtick is the superseder
// (mirroring the English "`X` supersedes `Y`"). The 면(면/이면) particle after the
// superseder is 받침-selected.
func renderSupersedeKo(s ast.Supersede) string {
	if len(s.Chain) >= 3 {
		text := ""
		for i, name := range s.Chain {
			if i > 0 {
				text += " > "
			}
			text += "`" + name + "`"
		}
		return text + " 순서로 배제한다"
	}
	return "`" + s.Superseder + "` " + josaFor(s.Superseder, "이면", "면") + " `" + s.Superseded + "` 배제한다"
}
