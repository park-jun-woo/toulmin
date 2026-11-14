//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseGroundedClause — parse "on" ground_arg [("and") ground_arg]* term list
package parser

import "strings"

// parseGroundedClause parses a "grounded on `g1` and `g2` ..." clause into the
// list of external ground fact names. s begins immediately after the
// "grounded" keyword (already consumed by the caller); this function first
// consumes the required "on" keyword, then the backtick+"and" list.
func parseGroundedClause(s string, path string, line int) ([]string, string, error) {
	rest, ok := takeKeyword(s, "on")
	if !ok {
		return nil, s, errAt(path, line, "expected 'on' after 'grounded'")
	}
	var names []string
	for {
		name, tail, ok := takeBacktick(rest)
		if !ok {
			return nil, "", errAt(path, line, "expected backtick-quoted name after 'grounded on'/'and'")
		}
		names = append(names, name)
		rest = tail
		trimmed := strings.TrimSpace(rest)
		if r2, ok2 := takeKeyword(trimmed, "and"); ok2 {
			rest = r2
			continue
		}
		break
	}
	return names, rest, nil
}
