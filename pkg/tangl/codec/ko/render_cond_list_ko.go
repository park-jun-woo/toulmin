//ff:func feature=tangl type=parser control=iteration dimension=1 pattern=recursive
//ff:what renderCondListKo — render an expression as a nested condition-list (list-marker lines at a depth)
package ko

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// renderCondListKo renders expression e as the lines of a nested condition list at
// the given depth: it flattens e's top-level 그리고/또는 chain into prefixed terms and
// renders each (recursing through renderTermKo for a 다음 중 either group). depth drives
// the 2-space relative indentation the core renderer later shifts to the section base.
func renderCondListKo(e ast.Expr, depth int) ([]string, error) {
	var lines []string
	for _, ct := range flattenLogicKo(e) {
		ls, err := renderTermKo(ct.Prefix, ct.Term, depth)
		if err != nil {
			return nil, err
		}
		lines = append(lines, ls...)
	}
	return lines, nil
}
