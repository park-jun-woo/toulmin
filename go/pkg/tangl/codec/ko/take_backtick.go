//ff:func feature=tangl type=parser control=sequence
//ff:what takeBacktick — consume a leading backtick-quoted name token
package ko

import "strings"

// takeBacktick consumes a leading backtick-quoted name from s after trimming
// leading spaces/tabs, returning the name, the remaining text, and whether a
// well-formed backtick token was found. The backtick delimiters fix the name
// boundary so a Hangul name may contain spaces and ASCII (e.g. `IP 차단`). It
// follows the same contract as pkg/tangl/parser.takeBacktick.
func takeBacktick(s string) (name string, rest string, ok bool) {
	s = strings.TrimLeft(s, " \t")
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
