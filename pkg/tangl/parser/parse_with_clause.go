//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseWithClause — parse "with_arg [and with_arg]*" term list
package parser

import "strings"

// parseWithClause parses a "with `term1` and `term2` ..." clause into the
// list of referenced Definitions term names.
func parseWithClause(s string, path string, line int) ([]string, string, error) {
	var terms []string
	rest := s
	for {
		name, tail, ok := takeBacktick(rest)
		if !ok {
			return nil, "", errAt(path, line, "expected backtick-quoted term after 'with'/'and'")
		}
		terms = append(terms, name)
		rest = tail
		trimmed := strings.TrimSpace(rest)
		if r2, ok2 := takeKeyword(trimmed, "and"); ok2 {
			rest = r2
			continue
		}
		break
	}
	return terms, rest, nil
}
