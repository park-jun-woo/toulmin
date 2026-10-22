//ff:func feature=tangl type=parser control=sequence
//ff:what parseEveryClause — parse an "<interval> [until `case`]" clause
package parser

import "strings"

// parseEveryClause parses the interval and optional "until `case`" clause
// following "every".
func parseEveryClause(s string, path string, line int) (string, string, error) {
	idx := strings.Index(s, " until ")
	if idx < 0 {
		interval := strings.TrimSpace(s)
		if interval == "" {
			return "", "", errAt(path, line, "expected interval after 'every'")
		}
		return interval, "", nil
	}
	interval := strings.TrimSpace(s[:idx])
	if interval == "" {
		return "", "", errAt(path, line, "expected interval after 'every'")
	}
	rest := strings.TrimSpace(s[idx+len(" until "):])
	caseName, rest, ok := takeBacktick(rest)
	if !ok {
		return "", "", errAt(path, line, "expected backtick-quoted case name after 'until'")
	}
	if strings.TrimSpace(rest) != "" {
		return "", "", errAt(path, line, "unexpected trailing text: %q", rest)
	}
	return interval, caseName, nil
}
