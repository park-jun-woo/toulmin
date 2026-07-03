//ff:func feature=tangl type=parser control=iteration dimension=1 pattern=recursive
//ff:what renderCondList — render an expression as a nested condition-list (list-marker lines at a depth)
package en

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// renderCondList renders expression e as the lines of a nested condition list at
// the given depth: it flattens e's top-level and/or chain into prefixed terms and
// renders each term (recursing through renderTerm for either/not groups). depth
// drives the 2-space relative indentation the core renderer later shifts to the
// section base (README contract §3).
func renderCondList(e ast.Expr, depth int) []string {
	var lines []string
	for _, ct := range flattenLogic(e) {
		lines = append(lines, renderTerm(ct.Prefix, ct.Term, depth)...)
	}
	return lines
}
