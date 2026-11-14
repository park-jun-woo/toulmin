//ff:func feature=tangl type=parser control=sequence
//ff:what takeKeywordKo — consume a leading keyword at a word boundary
package ko

import "strings"

// takeKeywordKo consumes a leading literal keyword kw (which may contain internal
// spaces, e.g. "다음 중") from s after trimming leading spaces/tabs, returning the
// remaining text. The match must land on a word boundary — the following byte, if
// any, must be a space, tab, or comma — so "그리고" does not spuriously match the
// prefix of "그리고도", while a modifier keyword abutting a separator comma (e.g.
// "하는, …", left by strings.Join(mods, ", ")) still matches; the caller's
// skipParticlesKo then peels the leading comma. It mirrors
// pkg/tangl/parser.takeKeyword; the new package redefines it.
func takeKeywordKo(s string, kw string) (string, bool) {
	s = strings.TrimLeft(s, " \t")
	if !strings.HasPrefix(s, kw) {
		return s, false
	}
	rest := s[len(kw):]
	if len(rest) > 0 && rest[0] != ' ' && rest[0] != '\t' && rest[0] != ',' {
		return s, false
	}
	return rest, true
}
