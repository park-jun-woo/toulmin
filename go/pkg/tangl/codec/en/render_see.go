//ff:func feature=tangl type=parser control=sequence
//ff:what renderSee — render a See declaration (`see `alias` from `package``)
package en

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// renderSee renders a tangl:See entry back to "see `<alias>` from `<package>`",
// the inverse of the parser's parseSeeItem.
func renderSee(s ast.See) string {
	return "see `" + s.Alias + "` from `" + s.Package + "`"
}
