//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what ReadWords — parse big-endian bytes into words, requiring even length
package word

import "errors"

// ReadWords parses b into words, reading each word big-endian (high byte
// first). Words are aligned to 2 bytes, so an odd-length input is an error.
func ReadWords(b []byte) ([]Word, error) {
	if len(b)%2 != 0 {
		return nil, errors.New("tangeul/word: byte length must be even (word aligned)")
	}
	out := make([]Word, 0, len(b)/2)
	for i := 0; i < len(b); i += 2 {
		out = append(out, Word(b[i])<<8|Word(b[i+1]))
	}
	return out, nil
}
