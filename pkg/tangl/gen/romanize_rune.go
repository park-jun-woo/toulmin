//ff:func feature=tangl type=codegen control=sequence
//ff:what romanizeRune — maps one rune toward ASCII: Hangul syllable to RR, everything else passes through
package gen

// romanizeRune transliterates a single Hangul syllable to its Revised
// Romanization spelling and passes every other rune (ASCII, punctuation,
// non-Hangul letters) through unchanged. It only produces the romanized
// form; the caller still applies the ASCII filter afterwards, so this stays
// a pure Hangul-to-ASCII pre-pass with no behavior change for ASCII input.
func romanizeRune(r rune) string {
	if rr := romanizeHangul(r); rr != "" {
		return rr
	}
	return string(r)
}
