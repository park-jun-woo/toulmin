//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseUnitDecimal — parse a leading weight literal (qualifier/force) and check it against [0, 1] at 4 decimals
package parser

import (
	"strconv"
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseUnitDecimal parses the weight literal following keyword ("qualified" or
// "force") and returns its value and the unconsumed remainder. The literal must
// parse as a float, lie in [0, 1] and carry at most four decimal places (the
// GEUL SIGNED-DECIMAL precision); anything else is a load error, never a
// silently rounded weight. The authored spelling (e.g. "0.50") is not kept in
// the AST — mode-1 decode replays it from SOURCE, mode-2 renders the shortest
// form.
func parseUnitDecimal(s, keyword, path string, line int) (float64, string, error) {
	s = trimStart(s)
	i := 0
	for i < len(s) && (s[i] == '.' || s[i] == '-' || (s[i] >= '0' && s[i] <= '9')) {
		i++
	}
	if i == 0 {
		return 0, "", errAt(path, line, "expected float after '%s'", keyword)
	}
	lit := s[:i]
	v, err := strconv.ParseFloat(lit, 64)
	if err != nil {
		return 0, "", errAt(path, line, "invalid %s value %q: %v", keyword, lit, err)
	}
	if dot := strings.IndexByte(lit, '.'); dot >= 0 && len(lit)-dot-1 > 4 {
		return 0, "", errAt(path, line, "invalid %s value %q: more than 4 decimal places", keyword, lit)
	}
	if err := ast.CheckUnitInterval(v); err != nil {
		return 0, "", errAt(path, line, "invalid %s value %q: %v", keyword, lit, err)
	}
	return v, s[i:], nil
}
