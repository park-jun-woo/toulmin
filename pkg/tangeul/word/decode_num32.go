//ff:func feature=tangl type=util control=sequence
//ff:what DecodeNum32 — decode a 2-word NUM32 into a uint32 (high word first)
package word

import "errors"

// DecodeNum32 decodes a NUM32 from the first two words of ws (high word first)
// and returns the uint32 value. It errors if fewer than two words are present.
func DecodeNum32(ws []Word) (uint32, error) {
	if len(ws) < num32Words {
		return 0, errors.New("tangeul/word: NUM32 requires 2 words")
	}
	return uint32(ws[0])<<16 | uint32(ws[1]), nil
}
