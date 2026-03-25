//ff:func feature=tangl type=parser control=sequence
//ff:what stripNumberedPrefix — remove numbered list prefix like "1. " from a line
package parser

import "strings"

// stripNumberedPrefix removes a numbered list prefix (e.g., "1. ") from a string.
func stripNumberedPrefix(s string) string {
	dotIdx := findDotAfterDigits(s)
	if dotIdx < 0 {
		return ""
	}
	rest := s[dotIdx+1:]
	if strings.HasPrefix(rest, " ") {
		return strings.TrimSpace(rest)
	}
	return ""
}
