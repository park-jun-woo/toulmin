//ff:func feature=tangl type=parser control=sequence
//ff:what parseQualifiedKo — parse the weight value of a '가중치' clause ([0, 1], at most 4 decimals)
package ko

// parseQualifiedKo parses the FLOAT following 가중치 — a weight in [0, 1] with
// at most four decimal places (parseUnitDecimalKo) — returning the value and the
// unconsumed remainder. It mirrors pkg/tangl/parser.parseQualifiedClause.
func parseQualifiedKo(s string, line int, path string) (float64, string, error) {
	return parseUnitDecimalKo(s, "가중치", line, path)
}
