//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseQualifiedClause — parse the float value of a "qualified" clause
package parser

import "strconv"

// parseQualifiedClause parses the FLOAT following "qualified".
func parseQualifiedClause(s string, path string, line int) (float64, string, error) {
	s = trimStart(s)
	i := 0
	for i < len(s) && (s[i] == '.' || s[i] == '-' || (s[i] >= '0' && s[i] <= '9')) {
		i++
	}
	if i == 0 {
		return 0, "", errAt(path, line, "expected float after 'qualified'")
	}
	q, err := strconv.ParseFloat(s[:i], 64)
	if err != nil {
		return 0, "", errAt(path, line, "invalid qualified value %q: %v", s[:i], err)
	}
	return q, s[i:], nil
}
