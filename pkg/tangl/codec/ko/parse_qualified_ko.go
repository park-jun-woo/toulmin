//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseQualifiedKo — parse the float value of a '가중치' clause
package ko

import "strconv"

// parseQualifiedKo parses the FLOAT following 가중치, returning the value and the
// unconsumed remainder. It mirrors pkg/tangl/parser.parseQualifiedClause.
func parseQualifiedKo(s string, line int, path string) (float64, string, error) {
	i := 0
	for i < len(s) && s[i] == ' ' {
		i++
	}
	s = s[i:]
	j := 0
	for j < len(s) && (s[j] == '.' || s[j] == '-' || (s[j] >= '0' && s[j] <= '9')) {
		j++
	}
	if j == 0 {
		return 0, "", errAtKo(path, line, "expected float after '가중치'")
	}
	q, err := strconv.ParseFloat(s[:j], 64)
	if err != nil {
		return 0, "", errAtKo(path, line, "invalid 가중치 value %q: %v", s[:j], err)
	}
	return q, s[j:], nil
}
