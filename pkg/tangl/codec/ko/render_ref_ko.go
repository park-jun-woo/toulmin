//ff:func feature=tangl type=parser control=sequence
//ff:what renderRefKo — render an ast.Ref as `alias`.`name` or a bare `name` (the inverse of parseRefKo)
package ko

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// renderRefKo renders a FUNC_REF back to its KO surface form: a qualified
// `alias`.`name` when Alias is set, or a bare `name` otherwise. It is the exact
// inverse of parseRefKo, so the output re-parses to the same Ref.
func renderRefKo(r ast.Ref) string {
	if r.Alias != "" {
		return "`" + r.Alias + "`.`" + r.Name + "`"
	}
	return "`" + r.Name + "`"
}
