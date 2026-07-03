//ff:func feature=tangl type=parser control=sequence
//ff:what renderSupersede — render an exclusion edge (`` `X` supersedes `Y` ``)
package en

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// renderSupersede renders a supersede edge back to "`<superseder>` supersedes
// `<superseded>`", the canonical form of the parser's crisp-exclusion edge —
// distinct from renderAttack's defeat edge, which only lowers a verdict float.
func renderSupersede(s ast.Supersede) string {
	return "`" + s.Superseder + "` supersedes `" + s.Superseded + "`"
}
