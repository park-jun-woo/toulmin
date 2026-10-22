//ff:func feature=tangl type=parser control=sequence
//ff:what renderRef — render an ast.Ref as `alias`.`name` or a bare `name` (the inverse of parseRef)
package en

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// renderRef renders a FUNC_REF back to its English surface form: a qualified
// `alias`.`name` when Alias is set, or a bare `name` otherwise. It is the exact
// inverse of the parser's parseRef, so the output re-parses to the same Ref.
func renderRef(r ast.Ref) string {
	if r.Alias != "" {
		return "`" + r.Alias + "`.`" + r.Name + "`"
	}
	return "`" + r.Name + "`"
}
