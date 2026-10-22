//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what sanitizeWord — romanizes Hangul then strips non letter/digit runes from a word
package gen

import "strings"

// sanitizeWord first transliterates any Hangul syllables to Revised
// Romanization (so a Korean tangl name like "출금" becomes "chulgeum"
// instead of collapsing to nothing) and then strips every rune that is not
// an ASCII letter or digit, so stray punctuation never leaks into a Go
// identifier. ASCII input passes through the romanize pre-pass unchanged, so
// its output is byte-for-byte identical to before.
func sanitizeWord(w string) string {
	var romanized strings.Builder
	for _, r := range w {
		romanized.WriteString(romanizeRune(r))
	}
	var b []rune
	for _, r := range romanized.String() {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b = append(b, r)
		}
	}
	return string(b)
}
