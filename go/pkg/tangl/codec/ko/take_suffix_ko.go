//ff:func feature=tangl type=parser control=sequence
//ff:what takeSuffixKo — consume a trailing sentence-final keyword at a word boundary
package ko

import "strings"

// takeSuffixKo consumes a trailing keyword kw (which may contain internal spaces,
// e.g. "한 번만 실행한다" or "발생 시") from the end of s, returning the text before it
// with trailing spaces/tabs trimmed. The match must land on a word boundary — the
// byte immediately before kw, if any, must be a space, tab, or backtick — so
// "면" does not match inside "같으면". Longest-match is the caller's duty: try the
// compound keyword before its shorter suffix.
func takeSuffixKo(s string, kw string) (string, bool) {
	if !strings.HasSuffix(s, kw) {
		return s, false
	}
	end := len(s) - len(kw)
	if end > 0 {
		b := s[end-1]
		if b != ' ' && b != '\t' && b != '`' {
			return s, false
		}
	}
	return strings.TrimRight(s[:end], " \t"), true
}
