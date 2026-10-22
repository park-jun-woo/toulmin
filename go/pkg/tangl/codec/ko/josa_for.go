//ff:func feature=tangl type=parser control=sequence
//ff:what josaFor — pick a 받침-driven Korean particle allomorph from a name's last character
package ko

import "unicode/utf8"

// josaFor selects a particle allomorph by the 받침 (final consonant) of name's last
// character, returning batchim (받침 있음) or open (받침 없음). The six pairs are
// {은,는} {이,가} {을,를} {과,와} {으로,로} {이면,면} — the caller passes {batchim, open}.
// A Hangul syllable (U+AC00~D7A3) has 받침 when (r-0xAC00)%28 != 0, with one 특칙:
// a ㄹ jongseong (%28==8) takes 로 in the 으로/로 pair (detected by open=="로"). A
// trailing digit uses its 낭독 받침 (0·1·3·6·7·8 받침, 2·4·5·9 무받침); any other
// character (ASCII, symbols, empty) is treated as 받침 없음 (문법 v0.3 §조사 정책).
func josaFor(name, batchim, open string) string {
	r, size := utf8.DecodeLastRuneInString(name)
	if size == 0 {
		return open
	}
	if r >= 0xAC00 && r <= 0xD7A3 {
		jong := (r - 0xAC00) % 28
		if jong == 0 {
			return open
		}
		if jong == 8 && open == "로" {
			return open
		}
		return batchim
	}
	if r >= '0' && r <= '9' {
		switch r {
		case '0', '1', '3', '6', '7', '8':
			return batchim
		default:
			return open
		}
	}
	return open
}
