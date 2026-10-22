//ff:func feature=tangl type=parser control=sequence
//ff:what renderRequire — render a required-field declaration (`` `f` is required [as <T>]``)
package en

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// renderRequire renders a required context field back to "`<field>` is required",
// with an optional "as <Type>" clause. It serves both Cases and Provides bodies
// (identical surface form), inverting the parser's parseRequireItem.
func renderRequire(r ast.Require) string {
	s := "`" + r.Field + "` is required"
	if r.Type != "" {
		s += " as " + renderType(r.Type)
	}
	return s
}
