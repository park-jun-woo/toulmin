//ff:func feature=tangl type=parser control=sequence
//ff:what takeKeyword — consume a leading literal keyword at a word boundary
package parser

import "strings"

// takeKeyword consumes a leading literal keyword (which may contain internal
// spaces, e.g. "in case of") from s (after trimming leading spaces/tabs) and
// returns the remaining text. The match must land on a word boundary (the
// following byte, if any, must be a space, tab, or end of string) so that
// "do" does not spuriously match the prefix of "don't".
func takeKeyword(s string, kw string) (string, bool) {
	s = trimStart(s)
	if !strings.HasPrefix(s, kw) {
		return s, false
	}
	rest := s[len(kw):]
	if len(rest) > 0 && rest[0] != ' ' && rest[0] != '\t' {
		return s, false
	}
	return rest, true
}
