//ff:func feature=tangl type=parser control=sequence
//ff:what takeBacktick — consume a leading backtick-quoted name token
package parser

import "strings"

// takeBacktick consumes a leading backtick-quoted name from s (after
// trimming leading spaces/tabs) and returns the name, the remaining text,
// and whether a well-formed backtick token was found.
func takeBacktick(s string) (name string, rest string, ok bool) {
	s = trimStart(s)
	if !strings.HasPrefix(s, "`") {
		return "", s, false
	}
	end := strings.Index(s[1:], "`")
	if end < 0 {
		return "", s, false
	}
	name = s[1 : 1+end]
	rest = s[1+end+1:]
	return name, rest, true
}
