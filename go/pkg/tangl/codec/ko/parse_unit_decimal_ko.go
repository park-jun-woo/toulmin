//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseUnitDecimalKo — parse a leading weight literal (가중치/세기) and check it against [0, 1] at 4 decimals
package ko

import (
	"strconv"
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseUnitDecimalKo parses the weight literal following keyword (가중치 or 세기)
// and returns its value and the unconsumed remainder (a trailing '의', attached or
// spaced, is left for the caller's skipParticlesKo). The literal must parse as a
// float, lie in [0, 1] and carry at most four decimal places; anything else is a
// load error. It mirrors pkg/tangl/parser.parseUnitDecimal.
func parseUnitDecimalKo(s, keyword string, line int, path string) (float64, string, error) {
	s = strings.TrimLeft(s, " \t")
	j := 0
	for j < len(s) && (s[j] == '.' || s[j] == '-' || (s[j] >= '0' && s[j] <= '9')) {
		j++
	}
	if j == 0 {
		return 0, "", errAtKo(path, line, "expected float after '%s'", keyword)
	}
	lit := s[:j]
	v, err := strconv.ParseFloat(lit, 64)
	if err != nil {
		return 0, "", errAtKo(path, line, "invalid %s value %q: %v", keyword, lit, err)
	}
	if dot := strings.IndexByte(lit, '.'); dot >= 0 && len(lit)-dot-1 > 4 {
		return 0, "", errAtKo(path, line, "invalid %s value %q: more than 4 decimal places", keyword, lit)
	}
	if err := ast.CheckUnitInterval(v); err != nil {
		return 0, "", errAtKo(path, line, "invalid %s value %q: %v", keyword, lit, err)
	}
	return v, s[j:], nil
}
