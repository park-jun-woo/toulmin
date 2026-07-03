//ff:func feature=tangl type=util control=sequence
//ff:what DecodeSpan — decode 3 words into a Span (TID, offset, length)
package word

import "errors"

// DecodeSpan decodes a Span from the first three words of ws (TID, offset,
// length). It errors if fewer than three words are present.
func DecodeSpan(ws []Word) (Span, error) {
	if len(ws) < spanWords {
		return Span{}, errors.New("tangeul/word: SPAN requires 3 words")
	}
	return Span{
		TID: uint16(ws[0]),
		Off: uint16(ws[1]),
		Len: uint16(ws[2]),
	}, nil
}
