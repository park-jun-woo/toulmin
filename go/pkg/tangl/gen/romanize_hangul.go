//ff:func feature=tangl type=codegen control=sequence
//ff:what romanizeHangul — transliterates one precomposed Hangul syllable to Revised Romanization jamo
package gen

// Precomposed Hangul syllable block: U+AC00 (가) .. U+D7A3 (힣). Each
// syllable index decomposes into initial (choseong, 19), medial
// (jungseong, 21), and final (jongseong, 28 including "none").
const (
	hangulSyllableBase = 0xAC00
	hangulSyllableLast = 0xD7A3
	jungseongCount     = 21
	jongseongCount     = 28
)

// Revised Romanization jamo tables (2000 MOCT notation). Phonological
// variation (liaison, assimilation) is deliberately NOT applied: a given
// syllable always maps to the same spelling so identifiers stay stable.
var (
	rrChoseong  = [...]string{"g", "kk", "n", "d", "tt", "r", "m", "b", "pp", "s", "ss", "", "j", "jj", "ch", "k", "t", "p", "h"}
	rrJungseong = [...]string{"a", "ae", "ya", "yae", "eo", "e", "yeo", "ye", "o", "wa", "wae", "oe", "yo", "u", "wo", "we", "wi", "yu", "eu", "ui", "i"}
	rrJongseong = [...]string{"", "k", "k", "k", "n", "n", "n", "t", "l", "k", "m", "l", "l", "l", "p", "l", "m", "p", "p", "t", "t", "ng", "t", "t", "k", "t", "p", "t"}
)

// romanizeHangul returns the Revised Romanization of r when r is a
// precomposed Hangul syllable, and "" otherwise. "출" (U+CD9C) becomes
// "chul" (ch + u + l); the silent initial ㅇ contributes nothing so "아"
// becomes "a".
func romanizeHangul(r rune) string {
	if r < hangulSyllableBase || r > hangulSyllableLast {
		return ""
	}
	idx := int(r - hangulSyllableBase)
	cho := idx / (jungseongCount * jongseongCount)
	jung := (idx / jongseongCount) % jungseongCount
	jong := idx % jongseongCount
	return rrChoseong[cho] + rrJungseong[jung] + rrJongseong[jong]
}
