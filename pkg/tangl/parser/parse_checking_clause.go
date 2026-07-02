//ff:func feature=tangl type=parser control=sequence
//ff:what parseCheckingClause — parse the target case name of a "checking" clause
package parser

// parseCheckingClause parses the backtick-quoted case name following
// "checking" (verdict composition, compiled to Evaluate).
func parseCheckingClause(s string, path string, line int) (string, string, error) {
	caseName, rest, ok := takeBacktick(s)
	if !ok {
		return "", "", errAt(path, line, "expected backtick-quoted case name after 'checking'")
	}
	return caseName, rest, nil
}
