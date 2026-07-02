//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what sanitizeWord — strips non letter/digit runes from a word
package gen

// sanitizeWord strips every rune that is not a letter or digit from w, so
// stray punctuation in a tangl name (rare) never leaks into a Go
// identifier.
func sanitizeWord(w string) string {
	var b []rune
	for _, r := range w {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b = append(b, r)
		}
	}
	return string(b)
}
