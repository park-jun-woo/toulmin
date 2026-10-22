//ff:func feature=tangl type=parser control=sequence
//ff:what wordPrefix — report whether text begins with a keyword on a word boundary
package en

import "strings"

// wordPrefix reports whether text begins with kw on a word boundary (the next
// byte, if any, is a space or tab), so "run" does not match the start of "running".
func wordPrefix(text, kw string) bool {
	if !strings.HasPrefix(text, kw) {
		return false
	}
	rest := text[len(kw):]
	return rest == "" || rest[0] == ' ' || rest[0] == '\t'
}
