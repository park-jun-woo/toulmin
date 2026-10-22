//ff:func feature=tangl type=parser control=sequence
//ff:what lastTokenKo — return the last whitespace-delimited token for suffix dispatch
package ko

import "strings"

// lastTokenKo returns the last whitespace-delimited token of s, or "" if s has
// no tokens. SOV statements end in a keyword, so the last token drives statement
// dispatch. Because the caller consumes the sentence-final keyword before applying
// skipParticlesKo to the head, "`X` 의 경우에" yields "경우에" — the 에 of 경우에 is
// never peeled off as a standalone particle.
func lastTokenKo(s string) string {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return ""
	}
	return fields[len(fields)-1]
}
