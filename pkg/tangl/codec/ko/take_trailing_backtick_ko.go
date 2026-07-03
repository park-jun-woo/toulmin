//ff:func feature=tangl type=parser control=sequence
//ff:what takeTrailingBacktick — split off the last backtick-quoted name at the end of a statement
package ko

import "strings"

// takeTrailingBacktick splits s into the prefix and the trailing backtick-quoted
// name. In the KO 규칙 정규형 the rule name is the last backtick group ("<조건> `n` 이다",
// with 이다 already consumed), and the condition — which itself contains backticks —
// is everything before it. It returns the name, the leading prefix, and ok.
func takeTrailingBacktick(s string) (name string, prefix string, ok bool) {
	s2 := strings.TrimRight(s, " \t")
	if !strings.HasSuffix(s2, "`") {
		return "", s, false
	}
	inner := s2[:len(s2)-1]
	i := strings.LastIndex(inner, "`")
	if i < 0 {
		return "", s, false
	}
	return inner[i+1:], s2[:i], true
}
