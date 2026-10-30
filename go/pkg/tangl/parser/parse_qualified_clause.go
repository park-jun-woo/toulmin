//ff:func feature=tangl type=parser control=sequence
//ff:what parseQualifiedClause — parse the weight value of a "qualified" clause ([0, 1], at most 4 decimals)
package parser

// parseQualifiedClause parses the FLOAT following "qualified": a weight in
// [0, 1] with at most four decimal places (parseUnitDecimal).
func parseQualifiedClause(s string, path string, line int) (float64, string, error) {
	return parseUnitDecimal(s, "qualified", path, line)
}
